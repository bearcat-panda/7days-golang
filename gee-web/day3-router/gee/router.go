package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //存储每种请求的树根节点
	handlers map[string]HandlerFunc //存储每种请求方式的HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
//解析路由  仅仅允许一个*号
func parsePattern(pattern string) []string {
	//按/进行切分
	vs := strings.Split(pattern, "/")

	//所有的子项
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { //如果为*表示匹配所有. 这一定是结尾
				break
			}
		}
	}
	return parts
}

//@param method 请求方法
//@param pattern 路由
//@param handler 处理方法
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//获取所有的子项
	parts := parsePattern(pattern)

	//键
	key := method + "-" + pattern
	_, ok := r.roots[method] //如果存在.不做处理.   不存在.就添加一个
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//获取所有的子项
	searchParts := parsePattern(path)
	params := make(map[string]string)
	//获取方法对应的前缀树. 节点
	root, ok := r.roots[method]

	//不存在. 返回空
	if !ok {
		return nil, nil
	}
	//查找对应的 最终节点
	n := root.search(searchParts, 0)

	if n != nil {
		//获取所有子项
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			//获取字符串的第一个字节
			if part[0] == ':' {
				//params[name] = test
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//替换*好
				//a/*/b   a/c/b
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
