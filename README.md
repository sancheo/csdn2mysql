A go project for crawl csdn blog to local database.  

想把csdn上的博客搬迁到其他博客里，于是写了这个小工具。  
# 使用前的准备工作
1. 确保本地装了Mysql并新建存储博客文章的DB。
2. 把数据库信息（如dbname、账号密码等）填到`conf\base.yml`文件中，文件里的`csdnUserName`字段填入要抓取的博客的用户名。
3. 运行`create_tables.sql`里的sql语句，创建出`存储博客表`和`存储图片表`。
4. 把`main.go`跑起来，运行项目。

# 处理逻辑说明
抓取时存储的是html格式的文章，遇到图片时会将图片base64编码存储到数据表里，同时将原来的`img标签`替换为`{picId:xxxxxx}`这种格式的占位符，再将文章发到其他博客平台时要自己处理图片的逻辑。