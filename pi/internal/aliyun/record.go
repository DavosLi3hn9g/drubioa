package aliyun

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"os"
	"time"
)

type Record struct {
	Text      string `json:"text"`      //内容Text
	Content   string `json:"content"`   //内容Json
	Recording string `json:"recording"` //录音文件地址
}

func (isi *ClientISI) InitRecord() {
	isi.ISIappKey = setting.AliyunIsiAppkey
}

func (isi *ClientISI) Record(fileLink, outJson string) *Record {
	var rec = new(Record)
	// 地域ID，常量内容，请勿改变

	const PRODUCT string = "nls-filetrans"
	const DOMAIN string = "filetrans.cn-shanghai.aliyuncs.com"
	const API_VERSION string = "2018-08-17"
	const POST_REQUEST_ACTION string = "SubmitTask"
	const GET_REQUEST_ACTION string = "GetTaskResult"
	// 请求参数key
	const KEY_APP_KEY string = "appkey"
	const KEY_FILE_LINK string = "file_link"
	const KEY_VERSION string = "version"
	const KEY_ENABLE_WORDS string = "enable_words"
	// 响应参数key
	const KEY_TASK string = "Task"
	const KEY_TASK_ID string = "TaskId"
	const KEY_STATUS_TEXT string = "StatusText"
	const KEY_RESULT string = "Result"
	// 状态值
	const STATUS_SUCCESS string = "SUCCESS"
	const STATUS_RUNNING string = "RUNNING"
	const STATUS_QUEUEING string = "QUEUEING"
	isi.RegionId = "cn-shanghai"
	isi.NewClient(nil)
	client := isi.Client
	postRequest := requests.NewCommonRequest()
	postRequest.Domain = DOMAIN
	postRequest.Version = API_VERSION
	postRequest.Product = PRODUCT
	postRequest.ApiName = POST_REQUEST_ACTION
	postRequest.Method = "POST"
	mapTask := make(map[string]string)
	mapTask[KEY_APP_KEY] = isi.ISIappKey
	mapTask[KEY_FILE_LINK] = fileLink
	// 新接入请使用4.0版本，已接入(默认2.0)如需维持现状，请注释掉该参数设置
	mapTask[KEY_VERSION] = "4.0"
	// 设置是否输出词信息，默认为false，开启时需要设置version为4.0
	mapTask[KEY_ENABLE_WORDS] = "false"
	task, err := json.Marshal(mapTask)
	if err != nil {
		logIO.Error(err)
		return rec
	}
	postRequest.FormParams[KEY_TASK] = string(task)
	postResponse, err := client.ProcessCommonRequest(postRequest)
	if err != nil {
		logIO.Error(err)
		return rec
	}
	postResponseContent := postResponse.GetHttpContentString()
	//fmt.Println(postResponseContent)
	if postResponse.GetHttpStatus() != 200 {
		logIO.Error("录音文件识别请求失败，Http错误码: ", postResponse.GetHttpStatus())
		return rec
	}
	var postMapResult map[string]interface{}
	err = json.Unmarshal([]byte(postResponseContent), &postMapResult)
	if err != nil {
		logIO.Error(err)
		return rec
	}
	var taskId string = ""
	var statusText string = ""
	statusText = postMapResult[KEY_STATUS_TEXT].(string)
	if statusText == STATUS_SUCCESS {
		logIO.Println("录音文件识别请求成功!")
		taskId = postMapResult[KEY_TASK_ID].(string)
	} else {
		logIO.Error("录音文件识别请求失败!")
		return rec
	}
	getRequest := requests.NewCommonRequest()
	getRequest.Domain = DOMAIN
	getRequest.Version = API_VERSION
	getRequest.Product = PRODUCT
	getRequest.ApiName = GET_REQUEST_ACTION
	getRequest.Method = "GET"
	getRequest.QueryParams[KEY_TASK_ID] = taskId
	statusText = ""
	var getMapResult map[string]interface{}

	for true {
		getResponse, err := client.ProcessCommonRequest(getRequest)
		if err != nil {
			logIO.Error(err)
			break
		}
		getResponseContent := getResponse.GetHttpContentString()
		//fmt.Println("识别查询结果：", getResponseContent)
		if getResponse.GetHttpStatus() != 200 {
			logIO.Error("识别结果查询请求失败，Http错误码：", getResponse.GetHttpStatus())
			break
		}
		err = json.Unmarshal([]byte(getResponseContent), &getMapResult)
		if err != nil {
			logIO.Error(err)
			break
		}
		statusText = getMapResult[KEY_STATUS_TEXT].(string)
		if statusText == STATUS_RUNNING || statusText == STATUS_QUEUEING {
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	if statusText == STATUS_SUCCESS {
		Result := getMapResult[KEY_RESULT]
		if Result != nil {
			Sentences := Result.(map[string]interface{})["Sentences"]
			jsonFlie, err := os.Create(outJson)
			if err != nil {
				logIO.Error(err)
			} else {
				sen, _ := json.Marshal(Sentences.(interface{}))
				_, err = jsonFlie.Write(sen)
				var talkList string
				for _, v := range Sentences.([]interface{}) {
					talk := v.(map[string]interface{})
					talkList += talk["Text"].(string)
				}
				text := string(talkList)
				return &Record{
					Text:      text,
					Content:   string(sen),
					Recording: fileLink,
				}

			}
			_ = jsonFlie.Close()
			logIO.Println("录音文件识别成功！")
		} else {
			logIO.Warning("录音文件已识别，但没有语音信息！")
		}
	} else {
		logIO.Error("录音文件识别失败！")
	}
	return rec
}
