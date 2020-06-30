package ginger

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
)

type Curl struct {
	uri         string
	userName    string
	passWord    string
	tokenEncode string
	contentType string
	headers     []Headers
	method      string
	requestData interface{}
	result      []byte
}

type Headers struct {
	headerKey   string
	headerValue string
}

//新的CURL对象
func NewCurl(maybeUrl ...interface{}) *Curl {
	//maybeUrl可以在初始化时就添加要访问的URL，也可以为空，后续调用Url方法来设定
	var url string
	if len(maybeUrl) == 1 {
		url = fmt.Sprintf("%v", maybeUrl[0])
	}
	return &Curl{contentType: "application/json", uri: url}
}

//设定要请求的URI地址
func (h *Curl) Url(uri string) *Curl {
	h.uri = uri
	return h
}

//basic auth
func (h *Curl) Auth(username string, password string) *Curl {
	h.userName = username
	h.passWord = password
	return h
}

//basic auth
func (h *Curl) ContentType(contentType string) *Curl {
	h.contentType = contentType
	return h
}

//basic auth
func (h *Curl) Header(headerKey string, headerValue string) *Curl {
	h.headers = append(h.headers, Headers{
		headerKey:   headerKey,
		headerValue: headerValue,
	})
	return h
}

//token 加密后的 和basic auth不能同时存在
func (h *Curl) Token(token string) *Curl {
	h.tokenEncode = token
	return h
}

//POST请求
func (h *Curl) Post(data ...interface{}) *Curl {
	h.method = "POST"
	if len(data) == 1 {
		h.requestData = data[0]
	}
	return h
}

//GET请求
func (h *Curl) Get(data ...interface{}) *Curl {
	h.method = "GET"
	if len(data) == 1 {
		h.requestData = data[0]
	}
	return h
}

//PUT请求
func (h *Curl) Put(data ...interface{}) *Curl {
	h.method = "PUT"
	if len(data) == 1 {
		h.requestData = data[0]
	}
	return h
}

//HEAD请求
func (h *Curl) Head(data ...interface{}) *Curl {
	h.method = "HEAD"
	if len(data) == 1 {
		h.requestData = data[0]
	}
	return h
}

//PATCH请求
func (h *Curl) Patch(data ...interface{}) *Curl {
	h.method = "PATCH"
	if len(data) == 1 {
		h.requestData = data[0]
	}
	return h
}

//发送请求 此方法作为内部协调方法，不允许外部使用
func (h *Curl) request() error {
	var err error
	var mJson []byte
	if mJson, err = json.Marshal(h.requestData); err != nil {
		return err
	}
	var mIo = bytes.NewReader(mJson)
	var req *http.Request
	if req, err = http.NewRequest(h.method, h.uri, mIo); err != nil {
		return err
	}

	//授权处理，目前支持bearer和basic auth,二者不能同时存在
	var authorize string
	if h.tokenEncode != "" {
		authorize = fmt.Sprintf("Bearer %s", h.tokenEncode)
	} else {
		var eType = fmt.Sprintf("%s:%s", h.userName, h.passWord)
		var esEncode = base64.StdEncoding.EncodeToString([]byte(eType))
		authorize = fmt.Sprintf("Basic %s", esEncode)
	}
	req.Header.Add("Content-Type", h.contentType)
	req.Header.Add("Authorization", authorize)

	//如果有额外的header要添加就在这里添加
	if len(h.headers) > 1 {
		for _, header := range h.headers {
			req.Header.Add(header.headerKey, header.headerValue)
		}
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	defer res.Body.Close()

	//读取结果，这里用io.copy防止内存爆掉
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, res.Body); err != nil {
		return err
	}
	h.result = buffer.Bytes()
	return nil
}

//绑定参数
func (h *Curl) Bind(data interface{}) error {
	//发送请求
	if err := h.request(); err != nil {
		return err
	}

	//断言 确保data的传入值为结构体地址 没写完
	switch data.(type) {
	case string:
	case int:
	case bool:
	case struct{}:
		return errors.New("只能传递结构体的地址")
		break
	}

	//如果是地址，获取真实的指针
	var dataElem = reflect.ValueOf(data).Elem()

	//获取真实的interface,而不是*interface
	realData := dataElem.Interface()
	if err := json.Unmarshal(h.result, realData); err != nil {
		log.Println(err)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(realData); err != nil {
		return err
	}
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(data); err != nil {
		return err
	}
	return nil
}

//绑定参数，这里的参数必须是一个字符串的地址
func (h *Curl) BindString(data interface{}) error {
	var err error

	//发送请求
	if err = h.request(); err != nil {
		return err
	}

	//断言 确保data的传入值为字符串地址 没写完
	switch data.(type) {
	case *string:
		err = nil
	default:
		err = errors.New("只能字符串地址")
	}
	if err != nil {
		return err
	}

	//如果是地址，获取真实的指针
	var dataElem = reflect.ValueOf(data).Elem()
	//获取真实的string,而不是*interface
	var realData = dataElem.String()
	realData = string(h.result)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(realData); err != nil {
		return err
	}
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(data); err != nil {
		return err
	}
	return nil
}
