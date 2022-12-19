package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"csdn2wordpress/log"
	ids "csdn2wordpress/utils"
)

func init() {
	viper.SetConfigName("base")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		err := errors.Wrap(err, "read config")
		Logger.Err(err).Stack().Msg("")
	}
}

var Logger = log.Logger
var CsdnDB *sql.DB

func main() {
	initMySQL()
	crawlingCSDN()
}

func initMySQL() {
	mysqlConf := viper.GetStringMapString("mysql")
	host := mysqlConf["host"]
	port := mysqlConf["port"]
	dbname := mysqlConf["dbname"]
	username := mysqlConf["username"]
	pwd := mysqlConf["pwd"]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, pwd, host, port, dbname)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		err := errors.Wrap(err, "sql open")
		Logger.Err(err).Stack().Msg("")
	}
	err = conn.Ping()
	if err != nil {
		err := errors.Wrap(err, "conn ping")
		Logger.Err(err).Stack().Msg("")
	}
	CsdnDB = conn
}

type CSDNArticle struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId"`
	Data    struct {
		List []struct {
			Type         string `json:"type"`
			FormatTime   string `json:"formatTime"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			HasOriginal  bool   `json:"hasOriginal"`
			DiggCount    int    `json:"diggCount"`
			CommentCount int    `json:"commentCount"`
			PostTime     int64  `json:"postTime"`
			CreateTime   int64  `json:"createTime"`
			URL          string `json:"url"`
			ArticleType  int    `json:"articleType"`
			ViewCount    int    `json:"viewCount"`
			Rtype        string `json:"rtype"`
		} `json:"list"`
		Total interface{} `json:"total"`
	} `json:"data"`
}

// CrawlingCSDN 爬取指定的CSDN博客
func crawlingCSDN() {
	userName := viper.GetString("csdnUserName")
	articleTypes := viper.GetString("articleTypes")
	Logger.Info().Msg(fmt.Sprintf("开始爬取指定博客:%s...", userName))
	url := "https://blog.csdn.net/community/home-api/v1/get-business-list?"
	hasMore := true
	page := 0
	for hasMore {
		page += 1
		params := fmt.Sprintf("page=%d&size=20&businessType=lately&username=%s", page, userName)
		r, err := http.Get(url + params)
		if err != nil {
			err := errors.Wrap(err, "http get")
			Logger.Err(err).Stack().Msg("")
			break
		}
		body, _ := ioutil.ReadAll(r.Body)
		var art CSDNArticle
		err = json.Unmarshal(body, &art)
		if err != nil {
			err := errors.Wrap(err, "json unmarshal")
			Logger.Err(err).Stack().Msg("")
			break
		}
		// 没有文章则不再获取
		if art.Data.List == nil || len(art.Data.List) < 1 {
			hasMore = false
			break
		}
		for i := 0; i < len(art.Data.List); i++ {
			item := art.Data.List[i]
			if strings.Contains(articleTypes, strconv.Itoa(item.ArticleType)) { // 1 原创 2 转载 4 翻译
				articleId := ids.GenerateID()
				articleLink := item.URL
				updateTime := time.Unix(item.PostTime/1000, 0).String()[0:19]
				publishTime := time.Unix(item.CreateTime/1000, 0).String()[0:19]
				articleType := item.ArticleType
				createTime := publishTime
				articleName := item.Title
				articleDesc := item.Description
				resp, err := soup.Get(articleLink)
				if err != nil {
					err := errors.Wrap(err, "soup get")
					Logger.Err(err).Stack().Msg("")
					break
				}
				doc := soup.HTMLParse(resp)
				tagBox := doc.Find("div", "class", "tags-box")
				aTags := tagBox.FindAll("a", "class", "tag-link")
				categorys := ""
				tags := ""
				for _, tag := range aTags {
					link := tag.Attrs()["href"]
					if has := strings.Contains(link, "category"); has {
						// link中包含category则属于分类
						categorys += tag.Text() + ";"
					} else {
						// 否则属于tag
						tags += tag.Text() + ";"
					}
				}
				if len(categorys) > 0 {
					categorys = categorys[:len(categorys)-1]
				}
				if len(tags) > 0 {
					tags = tags[:len(tags)-1]
				}
				contentViews := doc.Find("div", "id", "content_views")
				childs := contentViews.Children()
				content := ""
				for _, child := range childs {
					if child.NodeValue == "svg" || child.NodeValue == "\n                    " {
						continue
					}
					content += child.HTML()
					imgs := child.FindAll("img")
					if len(imgs) > 0 {
						// 转化成 base64 存储
						for _, img := range imgs {
							src := img.Attrs()["src"]
							alt := img.Attrs()["alt"]
							if src == "" {
								continue
							}
							r, err := http.Get(src)
							if err != nil {
								err := errors.Wrap(err, "get img src")
								Logger.Err(err).Stack().Msg("")
								break
							}
							data, err := ioutil.ReadAll(r.Body)
							if err != nil {
								err := errors.Wrap(err, "read img")
								Logger.Err(err).Stack().Msg("")
								break
							}
							mimeType := http.DetectContentType(data)
							var imgBase64Str string
							switch mimeType {
							case "image/jpeg":
								imgBase64Str += "data:image/jpeg;base64,"
							case "image/png":
								imgBase64Str += "data:image/png;base64,"
							}
							imgBase64Str += base64.StdEncoding.EncodeToString(data)
							// imgBase64Str存储到数据库
							picId := insertImg2DB(imgBase64Str, articleId, alt)
							replaceCont := fmt.Sprintf("{picId:%d}", picId)
							content = strings.Replace(content, img.HTML(), replaceCont, -1)
						}
					}
				}
				sqlStr := "INSERT INTO csdn_article" +
					"(article_id, author, article_desc, article_link, article_name, article_type, classification, article_tags, article_content, publish_time, create_time, update_time) " +
					"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
				_, err = CsdnDB.Exec(sqlStr, articleId, userName, articleDesc, articleLink, articleName, articleType, categorys,
					tags, content, publishTime, createTime, updateTime)
				if err != nil {
					err := errors.Wrap(err, "insert article sql")
					Logger.Err(err).Stack().Msg("")
				}
				Logger.Info().Msg(fmt.Sprintf("获取文章【%s】完成。", articleName))
			}
		}
		err = r.Body.Close()
		if err != nil {
			err := errors.Wrap(err, "close body")
			Logger.Err(err).Stack().Msg("")
			break
		}
		// defer func() { _ = r.Body.Close() }()
	}
	Logger.Info().Msg(fmt.Sprintf("爬取博客%s完成...", userName))
}

// insertImg2DB 将imgBase64Str存入数据库
func insertImg2DB(imgBase64Str string, articleId uint64, alt string) (picId uint64) {
	picId = ids.GenerateID()
	sqlStr := "INSERT INTO csdn_pic(pic_id, article_id, pic_content, pic_alt) VALUES (?, ?, ?, ?);"
	_, err := CsdnDB.Exec(sqlStr, picId, articleId, imgBase64Str, alt)
	if err != nil {
		err := errors.Wrap(err, "insert img sql")
		Logger.Err(err).Stack().Msg("")
	}
	return
}
