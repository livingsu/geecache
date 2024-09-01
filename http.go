package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/geecache/"

type HTTPNode struct {
	url      string
	basePath string
}

func NewHTTPNode(url string) *HTTPNode {
	return &HTTPNode{
		url:      url,
		basePath: defaultBasePath,
	}
}

func (node *HTTPNode) Log(format string, args ...interface{}) {
	log.Printf("[%s] %s\n", node.url, fmt.Sprintf(format, args...))
}

// 路由格式为/geecache/{group}/{key}
func (node *HTTPNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, node.basePath) {
		http.Error(w, "unsupported path", http.StatusBadRequest)
		return
	}
	params := strings.SplitN(r.URL.Path[len(node.basePath):], "/", 2)
	if len(params) != 2 {
		http.Error(w, "unsupported path", http.StatusBadRequest)
		return
	}
	group := params[0]
	key := params[1]

	g := GetGroup(group)
	if g == nil {
		http.Error(w, "unknown group", http.StatusBadRequest)
		return
	}
	v, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(v.Bytes())
}
