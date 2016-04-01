## Usage

Use `signature.Signature` to vertify every request signature against a pre-defined accessKeySecret:

客户端签名方式请参考：[签名机制](https://help.aliyun.com/document_detail/ecs/open-api/requestmethod/signature.html?spm=5176.product8314827_ecs.6.197.HMccmd)

~~~ go
import (
  "github.com/go-martini/martini"
  "github.com/venusource/martini-api-signature"
)

func main() {
  m := martini.Classic()
  // signature vertify every request
  m.Use(signature.Signature("0SB35kS87NNDD8"))
  m.Run()
}
~~~

## Authors
* [Changjun Zhao](http://github.com/ChangjunZhao)
