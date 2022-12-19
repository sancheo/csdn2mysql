# 创建存储博客表 csdn_article
CREATE TABLE `article`.`csdn_article`  (
                                       `article_id` bigint NOT NULL COMMENT '文章id',
                                       `author` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '作者',
                                       `article_desc` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文章简介',
                                       `article_link` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文章链接',
                                       `article_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文章名称',
                                       `article_type` int NULL DEFAULT NULL COMMENT '1 原创 2 转载 4 翻译',
                                       `classification` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '文章文类',
                                       `article_tags` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '文章标签',
                                       `article_content` blob NOT NULL COMMENT '文章内容',
                                       `publish_time` datetime NOT NULL COMMENT '发布时间',
                                       `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
                                       `update_time` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
                                       `is_delete` bit(1) NULL DEFAULT b'0',
                                       PRIMARY KEY (`article_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'CSDN爬取文章存储表' ROW_FORMAT = Dynamic;

# 创建存储图片表 csdn_pic
CREATE TABLE `article`.`csdn_pic`  (
                                       `pic_id` bigint NOT NULL,
                                       `article_id` bigint NULL DEFAULT NULL COMMENT '图片id',
                                       `pic_content` longblob NOT NULL COMMENT '图片内容blob',
                                       `pic_alt` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '图片介绍',
                                       `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
                                       `update_time` datetime NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
                                       `is_delete` bit(1) NULL DEFAULT b'0',
                                       PRIMARY KEY (`pic_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'CSDN文章图片' ROW_FORMAT = Dynamic;

