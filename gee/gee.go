package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct { // 分组结构体
	prefix      string        // 分组对应的前缀
	middlewares []HandlerFunc // 该分组注册的中间件
	parent      *RouterGroup  // 当前分组的父亲分组（子分组可以使用父分组的中间件）
	engine      *Engine       // 所有分组共享同一个 Engine 实例（帮助 group 访问 router ）
}

type Engine struct { // 实现 ServeHTTP 的接口
	*RouterGroup // 匿名内嵌, 算是 Engine 继承 RouterGroup , 可以使用其拥有的相关方法
	router       *router
	groups       []*RouterGroup // 存放所有分组
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}  // RouterGroup 本身持有一个全局的 engine , 初始化时要赋值
	engine.groups = []*RouterGroup{engine.RouterGroup} // engine 本身的 RouterGroup 也算一个分组
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, // 分组前缀为父分组的前缀+当前的前缀
		parent: group,
		engine: engine, // 与夫分组公用同一个 engine
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 注册中间件的方法
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 添加GET请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 添加POST请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run 设置监听在addr端口并启动
func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

// serveHTTP 处理每个http请求 [http监听要求实现的接口函数]
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	e.router.handle(c)
}
