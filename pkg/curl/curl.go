package curl

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Headers map[string]string
}

var defaultConfig = &Config{
	Headers: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
}

var POST = defaultConfig.POST
var GET = defaultConfig.GET

func (c *Config) SetHeader(k, v string) {
	c.Headers[k] = v
}
func (c *Config) GET(httpUrl string, param map[string]string) []byte {

	//  待合成文本
	var data = url.Values{}
	for k, v := range param {
		data.Add(k, v)
	}

	reqBody := data.Encode()
	var netTransport = &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	//  组装http请求头
	req, _ := http.NewRequest("GET", httpUrl+"?"+reqBody, nil)
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody
}
func (c Config) POST(httpUrl string, param map[string]string) ([]byte, http.Header) {

	//  待合成文本
	var data = url.Values{}
	for k, v := range param {
		data.Add(k, v)
	}

	reqBody := data.Encode()
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // disable verify
	}
	client := &http.Client{Transport: transCfg}
	req, _ := http.NewRequest("POST", httpUrl, strings.NewReader(reqBody))
	//  组装http请求头
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, nil
	}

	//core.Print("%+v",resp.Body)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	respHeader := resp.Header
	return respBody, respHeader
}

func (c *Config) POSTFILE(httpUrl string, param map[string]string, file []byte) []byte {

	//  待合成文本
	var data = url.Values{}
	for k, v := range param {
		data.Add(k, v)
	}

	reqBody := data.Encode()

	//fmt.Printf("参数：%+v",reqBody)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", httpUrl+"?"+reqBody, bytes.NewBuffer(file))
	//  组装http请求头
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody
}

func (c *Config) POSTJSON(httpUrl string, params []byte) ([]byte, http.Header) {

	var jsonStr = []byte(params)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // disable verify
	}
	client := &http.Client{Transport: transCfg}
	req, _ := http.NewRequest("POST", httpUrl, bytes.NewBuffer(jsonStr))
	//  组装http请求头
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, nil
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	respHeader := resp.Header
	return respBody, respHeader
}
