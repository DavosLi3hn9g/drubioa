package aichat

import (
	"VGO/pkg/curl"
	"encoding/json"
)

var (
	aiUrl = "https://openapi.tuling123.com/openapi/api/v2"
)

type Turing struct {
}

type RespData struct {
	Intent struct {
		Code       int               `json:"code"`
		IntentName string            `json:"intentName"`
		ActionName string            `json:"actionName"`
		Parameters map[string]string `json:"parameters"`
	} `json:"intent"`
	Results []struct {
		GroupType  int               `json:"groupType"`
		ResultType string            `json:"resultType"`
		Values     map[string]string `json:"values"`
	} `json:"results"`
}

func (tr Turing) Chat(in string) RespData {

	var resp RespData
	if in == "" {
		return resp
	}

	params := `{
		"reqType":0,
		"perception": {
			"inputText": {
				"text": "` + in + `"
			},			
		},
		"userInfo": {
			"apiKey": "",
			"userId": ""
		}
	}`
	bytesData := []byte(params)
	h := curl.Config{}

	data, _ := h.POSTJSON(aiUrl, bytesData)

	err := json.Unmarshal(data, &resp)
	if err != nil {
		panic(err)
	}

	//core.Print("%+v", resp)
	return resp
}
