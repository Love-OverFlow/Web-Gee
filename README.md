# Web-Gee

参考：[从零实现Web框架Gee教程 | 极客兔兔](https://geektutu.com/post/gee.html)

在使用 Go 进行网页后端 api 接口开发时，常使用 Gin 之类的 Web 框架，这是因为一些 Web 开发中的一些基础需求`net/http`并不原生支持，需要手动实现。

- 动态路由：例如`/:username/detail`
- 中间件：一些统一的处理流程只能在每个路由映射的 handler 中实现
- 对特定响应比如 JSON ，没有相应的调用方法，需要手动序列化

- ...

----

## 主要工作

### 封装上下文，支持中间件

参考`Gin`，将以下信息封装在上下文 `Context` 中，方便处理 Web 请求：

- 获取请求的参数，包括 URL ，方法（GET/POST/PUT...）等
- 路由对应的方法，包括中间件

-  提供 JSON 等常用响应的快捷方法，避免手动设置消息头和手动 JSON 序列化

```go
type Context struct {
   Writer     http.ResponseWriter
   Req        *http.Request
   StatusCode int
   Path       string
   Method     string
   Params     map[string]string
   handlers   []HandlerFunc // 包括中间件的所有待执行函数
   index      int           // 执行到第几个函数了
}
```

### 动态路由

- 使用前缀树实现对 URL 中通配符冒号「`:`」和星号「`*`」的动态路径匹配

- 同时将 URL 中对应的参数保存到上下文 Context 中

```go
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang  只有根节点不为空
	part     string  // 路由中的一部分，例如 :lang  保存路径, 这里和经典前缀树稍有不同, 路径是保存在节点中的
	children []*node // 子节点，例如 [doc, tutorial, intro]  因为路径上不确定有多少种字符, 所以是切片
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为 true
}
```
