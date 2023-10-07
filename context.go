package goo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//原对象
	Writer http.ResponseWriter
	Req    *http.Request

	//请求信息
	Path   string
	Methon string
	Params map[string]string

	//响应消息
	StatusCode int

	//中间件
	handlers []HandlerFunc
	index    int

	//Engine指针
	engine *Engine
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Methon: r.Method,
		index:  -1,
	}
}

func (context *Context) Next() {
	context.index++
	Len := len(context.handlers)
	for ; context.index < Len; context.index++ {
		context.handlers[context.index](context)
	}
}

func (context *Context) Status(code int) {
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(code int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) HTML(code int, html string, data interface{}) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplates.ExecuteTemplate(context.Writer, html, data); err != nil {
		context.Fail(500, err.Error())
	}
}

func (context *Context) Data(code int, data []byte) {
	context.Status(code)
	context.Writer.Write(data)
}

func (context *Context) Param(key string) string {
	value, _ := context.Params[key]
	return value
}

func (context *Context) Fail(code int, err string) {
	context.index = len(context.handlers)
	context.JSON(code, H{"message": err})
}
