package aliyun

import (
	"VGO/pi/internal/cache"
	config2 "VGO/pi/internal/config"
	"VGO/pi/internal/pkg/file"
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var logIO *file.Log
var config = new(Config)
var setting *cache.Settings
var configENV = config2.ENV

type Config struct {
	AccessKeyId     string `form:"aliyun_ak_id" xml:"aliyun_ak_id" json:"aliyun_ak_id"`
	AccessKeySecret string `form:"aliyun_ak_secret" xml:"aliyun_ak_secret" json:"aliyun_ak_secret"`
}

type ClientOSS struct {
	*oss.Client
	RegionId string
}
type ClientISI struct {
	*sdk.Client
	token     string
	RegionId  string
	ISIappKey string
	TTSConfig TTSConfig
}
type ClientSMS struct {
	*sdk.Client
	token *string
}

type RespAuth struct {
	NlsRequestId string `json:"NlsRequestId"`
	RequestId    string `json:"RequestId"`
	Token        struct {
		ExpireTime int    `json:"ExpireTime"`
		Id         string `json:"Id"`
		UserId     string `json:"UserId"`
	} `json:"Token"`
}
type RespDataISI struct {
	TaskId  string `json:"task_id"`
	Result  string `json:"result"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func initConfig() *Config {
	setting = cache.SettingsCache.Read()
	config = &Config{
		setting.AliyunAkId,
		setting.AliyunAkSecret,
	}
	return config
}

func (os *ClientOSS) NewClient(c *Config) *ClientOSS {
	// 创建OSSClient实例
	var err error
	if c != nil {
		config = c
	} else {
		config = initConfig()
	}
	if os.RegionId == "" {
		os.RegionId = setting.AliyunBucket.Location
	}
	if os.RegionId == "" {
		// Endpoint 杭州
		os.RegionId = "oss-cn-hangzhou"
	}
	endpoint := "https://" + os.RegionId + ".aliyuncs.com"
	os.Client, err = oss.New(endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		logIO.Error(err)
	}

	return os
}
func (isi *ClientISI) NewClient(c *Config) *sdk.Client {
	var err error
	if c != nil {
		config = c
	} else {
		config = initConfig()
		isi.ISIappKey = setting.AliyunIsiAppkey
	}
	if isi.RegionId == "" {
		isi.RegionId = "default"
	}
	client, err := sdk.NewClientWithAccessKey(isi.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		logIO.Error(err)
	}
	isi.Client = client
	return client
}
func (isi *ClientISI) CreateToken(c *Config) error {
	var resp RespAuth
	var err error
	isi.RegionId = "cn-shanghai"
	isi.NewClient(c)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = "nls-meta.cn-shanghai.aliyuncs.com"
	request.ApiName = "CreateToken"
	request.Version = "2019-02-28"
	response, err := isi.Client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &resp)
	if err != nil {
		return err
	}
	//log.Print(resp.Token.Id)
	if resp.Token.Id == "" {
		return errors.New("验证失败！未知原因。")
	} else {
		isi.token = resp.Token.Id
	}
	return nil
}
