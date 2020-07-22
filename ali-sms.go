package ginger

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

const (
	//ErrNoPhoneParam 没有设置电话参数
	ErrNoPhoneParam = 1 << iota
	//ErrNoRegionParam 没有设置区域参数
	ErrNoRegionParam
	//ErrNoKeyParam 没有设置key参数
	ErrNoKeyParam
	//ErrNoSecretParam 没有设置secret参数
	ErrNoSecretParam
	//ErrNoSignNameParam 没有设置签名名称参数
	ErrNoSignNameParam
	//ErrNoTemplateCodeParam 没有设置模板代码参数
	ErrNoTemplateCodeParam
)

//AliSMS 阿里SMS类
type AliSMS struct {
	region        string
	accessKey     string
	accessSecret  string
	phone         string
	signName      string
	templateCode  string
	templateParam interface{}
}

//NewAliSMS 初始化类
func NewAliSMS() *AliSMS {
	return &AliSMS{}
}

//SetTemplateCode 设置模板代码
func (a *AliSMS) SetTemplateCode(templateCode string) *AliSMS {
	a.templateCode = templateCode
	return a
}

//SetSignName 设置短信签名名称
func (a *AliSMS) SetSignName(signName string) *AliSMS {
	a.signName = signName
	return a
}

//SetPhone 设置发送到哪个手机
func (a *AliSMS) SetPhone(phone string) *AliSMS {
	a.phone = phone
	return a
}

//SetAccessSecret 设置secret
func (a *AliSMS) SetAccessSecret(accessSecret string) *AliSMS {
	a.accessSecret = accessSecret
	return a
}

//SetAccessKey 设置访问key
func (a *AliSMS) SetAccessKey(accessKey string) *AliSMS {
	a.accessKey = accessKey
	return a
}

//SetRegion 设置区域
func (a *AliSMS) SetRegion(region string) *AliSMS {
	a.region = region
	return a
}

//GetErr 获取错误信息
func (a *AliSMS) GetErr(errCode int) error {
	var result = "未知错误"
	switch errCode {
	case ErrNoPhoneParam:
		result = "没有设置电话参数"
		break
	case ErrNoRegionParam:
		result = "没有设置区域参数"
		break
	case ErrNoKeyParam:
		result = "没有设置访问Key参数"
		break
	case ErrNoSecretParam:
		result = "没有设置访问secret参数"
		break
	case ErrNoSignNameParam:
		result = "没有设置签名名称参数"
		break
	case ErrNoTemplateCodeParam:
		result = "没有设置模板代码参数"
		break
	}

	return errors.New(result)
}

//CheckParam 检查内部参数
func (a *AliSMS) CheckParam() (err error) {
	if a.region == "" {
		err = a.GetErr(ErrNoRegionParam)
	}
	if a.accessKey == "" {
		err = a.GetErr(ErrNoKeyParam)
	}
	if a.accessSecret == "" {
		err = a.GetErr(ErrNoSecretParam)
	}
	if a.phone == "" {
		err = a.GetErr(ErrNoPhoneParam)
	}
	if a.signName == "" {
		err = a.GetErr(ErrNoSignNameParam)
	}
	if a.templateCode == "" {
		err = a.GetErr(ErrNoTemplateCodeParam)
	}
	return err
}

//Send 发送短信
func (a *AliSMS) Send(templateParam interface{}) (response *dysmsapi.SendSmsResponse, err error) {
	//设置templateParam
	a.templateParam = templateParam
	//检查参数
	if err = a.CheckParam(); err != nil {
		return nil, err
	}
	//将pram转换为json string
	var param string
	var paramB []byte
	if a.templateParam != nil {
		if paramB, err = json.Marshal(a.templateParam); err == nil {
			param = string(paramB)
		}
	}
	//发送短信
	var client *dysmsapi.Client
	client, err = dysmsapi.NewClientWithAccessKey(a.region, a.accessKey, a.accessSecret)
	if err != nil {
		return nil, err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = a.phone
	request.SignName = a.signName
	request.TemplateCode = a.templateCode
	request.TemplateParam = param

	response, err = client.SendSms(request)
	if err != nil {
		return nil, err
	}

	return response, err
}