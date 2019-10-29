# shortLink URL转短地址服务
## 什么是短地址服务？
```
将长地址缩短到一个很短的地址，用户访问这个短地址可以重定向到原本的长地址。
```
## 开发运用知识点
- HTTP Router和Handler设计
- HTTP处理流程中的Middleware设计
- 利用Golang接口实现可扩展设计
- 利用Redis自增长序列生成短地址

## API接口说明
- POST /api/shorten 生成短地址
```
{
	"url": "https://www.tracy.com",          // 将被转换的长地址
	"expiration_in_minutes": 1    // 设置过期时间（单位分钟）
}
```
响应
```
{
	"shortlink": "Aflb1"
}
```
- GET /api/info?shortlink=shortlink 获取短地址详细信息
```
例如：GET /api/info?shortlink=Aflb1
结果：
{
	"url":"https://www.tracy.com",
	"created":"2019-10-29 11:10:46.196556 +0800 CST m=+1033.841130980",
	"expiration_in_minutes":10
}
```
- GET /:shortlink 重定向到长地址（return 302 code）
```
例如：请求GET /Aflb1
结果：重定向到https://www.tracy.com
```
## 项目模块架构
- 主服务模块（app.go）
- 中间件模块(middleware.go)
- 存储模块(storage.go)
