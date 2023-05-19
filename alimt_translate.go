// alimt 翻译
package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	Endpoint = "mt.aliyuncs.com"
)

var (
	convertLanguageCode = map[string]string{
		"zh-CN": "zh",
		"zh-cn": "zh",
		"hans":  "zh",
		"hant":  "zh-tw",
	}
)

// alimt 返回结构
type ALIMTResponse struct {
	Data  ALIMTResponseData `json:"data"`
	Error ALIMTError        `json:"error"`
}
type ALIMTResponseTranslations struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}
type ALIMTResponseData struct {
	Translations []ALIMTResponseTranslations `json:"translations"`
}

// ALIMT 翻译失败返回结构
type ALIMTErrors struct {
	Message string `json:"message"`
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
}
type ALIMTError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ALIMTErrors `json:"errors"`
}

type ALIMTTranslateService struct {
	Client *Client
}

// 使用AK&SK初始化账号Client
func (gt *ALIMTTranslateService) CreateClient(accessKeyId *string, accessKeySecret *string) (_result *alimt20181012.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String(Endpoint)
	_result = &alimt20181012.Client{}
	_result, _err = alimt20181012.NewClient(config)
	return _result, _err
}

// 翻译
// alimt 翻译一次性最大只能翻译128个文本，调用前请做好 Contents 切片长度控制
func (gt *ALIMTTranslateService) Trans(c context.Context, contents []SourceText) ([]*Result, error) {
	keys := strings.Split(gt.Client.Key, "/")
	if len(keys) != 2 {
		return nil, fmt.Errorf("AccessKeyId Or AccessKeySecret failed")
	}

	client, err := gt.CreateClient(tea.String(keys[0]), tea.String(keys[1]))
	if err != nil {
		return nil, err
	}

	mapContents := make(map[string]string, len(contents))
	for index, con := range contents {
		indexStr := strconv.Itoa(index)
		mapContents[indexStr] = string(con)
	}
	mapContentsStr, err := json.Marshal(mapContents)
	if err != nil {
		return nil, err
	}

	if code, ok := convertLanguageCode[string(gt.Client.Source)]; ok {
		gt.Client.Source = SourceLanguage(code)
	}
	if code, ok := convertLanguageCode[string(gt.Client.Target)]; ok {
		gt.Client.Target = TargetLanguage(code)
	}

	request := &alimt20181012.GetBatchTranslateRequest{
		FormatType:     tea.String(gt.Client.Format),
		SourceLanguage: tea.String(string(gt.Client.Source)),
		TargetLanguage: tea.String(string(gt.Client.Target)),
		SourceText:     tea.String(string(mapContentsStr)),
		Scene:          tea.String(gt.Client.Scene),
		ApiType:        tea.String(gt.Client.ApiType),
	}

	response, err := client.GetBatchTranslate(request)
	if err != nil {
		return nil, err
	}
	if int(*response.Body.Code) != 200 {
		return nil, fmt.Errorf(string(*response.Body.Message))
	}
	var result []*Result
	// 返回翻译结果
	for index, cont := range contents {
		indexStr := strconv.Itoa(index)
		var tempResult *Result
		for _, trans := range response.Body.TranslatedList {
			if trans["index"] == indexStr {
				tempResult = &Result{
					SourceText:     cont,
					SourceLanguage: gt.Client.Source,
				}
				if trans["code"] == "200" {
					tempResult.Text = trans["translated"].(string)
					wordNumsStr := trans["wordCount"].(string)
					wordNumsUint32, _ := strconv.Atoi(wordNumsStr)
					tempResult.WordNums = uint32(wordNumsUint32)
				} else {
					tempResult.ErrMsg = trans["errorMsg"].(string)
				}
				continue
			}
		}
		result = append(result, tempResult)
	}

	return result, nil

}
