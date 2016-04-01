package signature

import (
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var BasicRealm = "Signature Required"

func Signature(accessKeySecret string) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
		//log.Println(req.URL.Query())
		req.ParseForm()
		//log.Println(req.PostForm)
		paramsToSign := NewOrderedParams()
		signature := ""
		for key, value := range req.URL.Query() {
			if key == "Signature" {
				signature = value[0]
				continue
			}
			paramsToSign.Add(key, value[0])
		}

		for key, value := range req.PostForm {
			if key == "Signature" {
				signature = value[0]
				continue
			}
			paramsToSign.Add(key, value[0])
		}
		reqString := requestString(req.Method, "/", paramsToSign)
		realSignature, _ := NewSigner(accessKeySecret).Sign(reqString)
		log.Println("signature:" + signature)
		log.Println("realSignature:" + realSignature)
		if signature != realSignature {
			unauthorized(res)
			return
		}
	}
}

func unauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\""+BasicRealm+"\"")
	http.Error(res, "api signature fail", http.StatusUnauthorized)
}

func requestString(method string, urlPath string, params *OrderedParams) string {
	result := method + "&" + Escape(urlPath)
	for pos, key := range params.Keys() {
		if pos == 0 {
			result += "&"
		} else {
			result += Escape("&")
		}
		result += Escape(fmt.Sprintf("%s=%s", key, params.Get(key)))
	}
	return result
}

//排序后的参数列表
type OrderedParams struct {
	allParams   map[string]string
	keyOrdering []string
}

func NewOrderedParams() *OrderedParams {
	return &OrderedParams{
		allParams:   make(map[string]string),
		keyOrdering: make([]string, 0),
	}
}

func (o *OrderedParams) Get(key string) string {
	return o.allParams[key]
}

func (o *OrderedParams) Keys() []string {
	sort.Sort(o)
	return o.keyOrdering
}

func (o *OrderedParams) Add(key, value string) {
	o.AddUnescaped(key, Escape(value))
}

func (o *OrderedParams) AddUnescaped(key, value string) {
	o.allParams[key] = value
	o.keyOrdering = append(o.keyOrdering, key)
}

func (o *OrderedParams) Len() int {
	return len(o.keyOrdering)
}

func (o *OrderedParams) Less(i int, j int) bool {
	return o.keyOrdering[i] < o.keyOrdering[j]
}

func (o *OrderedParams) Swap(i int, j int) {
	o.keyOrdering[i], o.keyOrdering[j] = o.keyOrdering[j], o.keyOrdering[i]
}

func (o *OrderedParams) Clone() *OrderedParams {
	clone := NewOrderedParams()
	for _, key := range o.Keys() {
		clone.AddUnescaped(key, o.Get(key))
	}
	return clone
}

func Escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}
