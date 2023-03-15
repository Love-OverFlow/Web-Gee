package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{} // 方便构建 JSON

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

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm 从 Request 传来的 Form 表单中, 根据 key 找对应的 value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 从形如: /hello?name=gee 的 URL 中获取参数(gee)
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Param 当路由为 /hello/:name 时, 从 /hello/gee 中获取参数(gee)
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 方便设置消息头的格式
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj H) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		fmt.Println("json 编码出错")
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
