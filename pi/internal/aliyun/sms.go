package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type data struct {
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Code      string `json:"code"`
}

func (c *ClientSMS) SendOne() *data {
	var data *data
	smsNumber := 1000
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["PhoneNumbers"] = ""
	request.QueryParams["SignName"] = ""
	request.QueryParams["TemplateCode"] = ""
	request.QueryParams["TemplateParam"] = fmt.Sprintf(`{"code":"%d"}`, smsNumber)
	response, err := c.ProcessCommonRequest(request)
	if err != nil {
		logIO.Error(err)
	}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &data)
	if err != nil {
		logIO.Error(err)
	}
	return data
}
