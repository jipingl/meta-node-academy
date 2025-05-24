### 题目

Go 语言结合 Gin 框架和 GORM 库开发一个个人博客系统的后端，实现博客文章的基本管理功能，包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能

### 运行环境

GO + Gin + GORM + MySQL

版本参考 go.mod 文件

### 依赖安装

1. 安装好 mysql 数据库服务器，进入配置文件配置好数据库连接信息

```shell
cd basic_1/task_4/blog/config.yaml
```

2. 进入项目根目录，执行以下命令安装依赖

```shell
cd basic_1/task_4/blog
go mod download
```

### 运行

项目默认运行在 8080 端口

进入项目主目录，执行以下命令运行项目

```shell
cd basic_1/task_4/blog
go run .
```

### 接口文档

1. 注册

```shell
POST /register HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 54

{
    "username": "gene",
    "password": "121212"
}
```

2. 登录

```shell
POST /login HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 54

{
    "username": "gene",
    "password": "121212"
}
```

3. 创建文档

```shell
POST /api/post HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
Content-Type: application/json
Content-Length: 70

{
    "title": "hello 2025",
    "content": "i am coming"
}
```

4. 获取文档

```shell
GET /api/post/1 HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
```

5. 更新文档

```shell
PUT /api/post/2 HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
Content-Type: application/json
Content-Length: 59

{
    "title": "happy 2025",
    "content": "i am coming"
}
```

6. 删除文档

```shell
DELETE /api/post/1 HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
Content-Type: application/json
Content-Length: 70
```

7. 评论

```shell
POST /api/comment HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
Content-Type: application/json
Content-Length: 43

{
    "post_id": 1,
    "content": "nice"
}
```

8. 获取评论

```shell
GET /api/comment?post_id=1 HTTP/1.1
Host: localhost:8080
Authorization: [jwt-token]
Content-Type: application/json
Content-Length: 43
```
