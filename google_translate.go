// 谷歌翻译
package translator

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	GoogleTransUrl = "translation.googleapis.com/language/translate/v2"
)

// google 返回结构
type GoogleResponse struct {
	Data  GoogleResponseData `json:"data"`
	Error GoogleError        `json:"error"`
}
type GoogleResponseTranslations struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}
type GoogleResponseData struct {
	Translations []GoogleResponseTranslations `json:"translations"`
}

// Google 翻译失败返回结构
type GoogleErrors struct {
	Message string `json:"message"`
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
}
type GoogleError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Errors  []GoogleErrors `json:"errors"`
}

type GoogleTranslateService struct {
	Client *Client
}

type GoogleParameters struct {
	Q      []SourceText `json:"q"`      // 翻译内容
	Target string       `json:"target"` // 目标语言
	Format string       `json:"format"` // 源文本的格式，HTML（默认）或纯文本格式。值html表示 HTML，值text 表示纯文本。
	Source string       `json:"source"` // 源语言，如未设置，Google将自动识别
	Model  string       `json:"model"`
	Key    string       `json:"key"` // API 密钥

}

// 谷歌翻译URL
func (gt *GoogleTranslateService) GoogleBaseUrl() (string, error) {
	if gt.Client.Key == "" {
		return "", fmt.Errorf("key is null")
	}
	baseUrl := fmt.Sprintf("https://%s?key=%s", GoogleTransUrl, gt.Client.Key)
	return baseUrl, nil
}

// 翻译
// 谷歌翻译一次性最大只能翻译128个文本，调用前请做好 Contents 切片长度控制
func (gt *GoogleTranslateService) Trans(c context.Context, contents []SourceText) ([]*Result, error) {
	params := GoogleParameters{
		Q:      contents,
		Target: string(gt.Client.Target),
		Format: gt.Client.Format,
	}
	googleURL, err := gt.GoogleBaseUrl()
	if err != nil {
		return nil, err
	}

	response := new(GoogleResponse)
	resByte, err := gt.Client.Request.Post(googleURL, params)
	if err != nil {
		if err2 := json.Unmarshal([]byte(err.Error()), &response); err2 == nil {
			return nil, fmt.Errorf("translate fail. code: %d message: %s", response.Error.Code, response.Error.Message)
		} else {
			return nil, err
		}
	}

	if err := json.Unmarshal(resByte, &response); err != nil {
		return nil, err
	}
	var result []*Result
	for k, data := range response.Data.Translations {
		result = append(result, &Result{
			SourceText:     contents[k],
			Text:           data.TranslatedText,
			SourceLanguage: SourceLanguage(data.DetectedSourceLanguage),
		})
	}
	return result, nil

}
