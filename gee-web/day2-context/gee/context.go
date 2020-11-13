package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	//原始对象
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	//请求信息
	Path   string //路径
	Method string //方法
	// response info
	//响应信息
	StatusCode int //响应码
}

//Context 构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		//URL指定要请求的URI(对于服务器请求)或要访问的URL(用于客户请求)
		//对于服务器请求,URL是从URI解析的在Request-Line上提供,存储在RequestURI中
		//对于大多数请求.除了Path和RawQuery以外的其他字段为空

		//对于客户端请求.URL的主机服务器指定位连接.而"请求的主机"字段则可选
		//指定要在HTTP中发送的Host表头值请求.
		//Path路径(相对路径可能会省略前斜杠)
		Path: req.URL.Path,
		//方法指定HTTP方法(get,post,put等). 空字符串表示get
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	//表单包含已解析的表单数据. 包括URL字段的查询参数以及
	//PATCH,POST,PUT表单数据.
	//此字段仅在调用ParseForm之后调用.
	//HTTP客户端会忽略Form并改用Body
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	//查询应为键=值设置列表.并用与号或分号.
	//没有等号的设置是解析为设置为空的键
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	//WriteHeader发送带有提供的HTTP响应标头状态代码
	//如果未显示调用WriteHeader,则首次调用Write将触发一个隐式
	//WriteHeader(http.StatusOk). 因此,对WriteHeader的显示调用
	//主要用于发送错误代码

	//提供的代码必须是有效的HTTP 1xx-5xx状态代码
	//只能写入一个标头,当前不支持发送用户自定义的1xx信息标题,
	//除了100连续响应标头,读取Request.Body后.服务器会自动发送
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	//设置响应头
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
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
