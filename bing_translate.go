package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	BingTransUrl   = "api.cognitive.microsofttranslator.com/translate"
	BingApiVersion = "3.0"
)

type BingText struct {
	Text SourceText `json:"Text"` // html text
}

type BingTranslateService struct {
	Client *Client
	Key    string
}

type BingResponse struct {
	DetectedLanguage DetectedLanguage `json:"detectedLanguage"`
	Translations     []Translations   `json:"translations"`
	Error            BingError        `json:"error"`
}
type SentLen struct {
	SrcSentLen   []int `json:"srcSentLen"`
	TransSentLen []int `json:"transSentLen"`
}

type DetectedLanguage struct {
	Language string  `json:"language"`
	Score    float64 `json:"score"`
}

type Translations struct {
	Text    string  `json:"text"`
	To      string  `json:"to"`
	SentLen SentLen `json:"sentLen"`
}

type BingError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 必应翻译URL
func (bt *BingTranslateService) BingBaseUrl(params map[string]string) (string, error) {

	v := url.Values{}
	for key, val := range params {
		v.Add(key, val)
	}
	urlPar := v.Encode()
	urlPar, err := url.QueryUnescape(urlPar)
	if err != nil {
		return "", err
	}

	baseUrl := fmt.Sprintf("https://%s?api-version=%s&%s", BingTransUrl, BingApiVersion, urlPar)
	return baseUrl, nil
}

func (bt *BingTranslateService) Trans(c context.Context, contents []SourceText) ([]*Result, error) {
	if bt.Client.Key == "" {
		return nil, fmt.Errorf("key is null")
	}
	baseUrlPar := map[string]string{
		"to":                    string(bt.Client.Target),
		"textType":              bt.Client.Format,
		"includeAlignment":      "true",
		"includeSentenceLength": "true",
	}
	if bt.Client.Source != "" {
		baseUrlPar["from"] = string(bt.Client.Source)
	}
	bingUrl, err := bt.BingBaseUrl(baseUrlPar)
	if err != nil {
		return nil, err
	}
	bt.Client.Request.SetHeader("Ocp-Apim-Subscription-Key", bt.Client.Key)
	var params []BingText
	for _, val := range contents {
		params = append(params, BingText{Text: val})
	}
	resByte, err := bt.Client.Request.Post(bingUrl, params)
	bingResponse := new(BingResponse)
	if err != nil {
		if err2 := json.Unmarshal([]byte(err.Error()), &bingResponse); err2 == nil {
			return nil, fmt.Errorf("translate fail. code: %d message: %s", bingResponse.Error.Code, bingResponse.Error.Message)
		} else {
			return nil, err2
		}
	}

	var response []BingResponse
	if err := json.Unmarshal(resByte, &response); err != nil {
		return nil, fmt.Errorf("cannot Unmarshal err: %s", err)
	}
	var result []*Result
	for k, val := range response {
		var sourceLanguage SourceLanguage
		if lang := val.DetectedLanguage.Language; lang != "" {
			sourceLanguage = SourceLanguage(lang)
		} else {
			sourceLanguage = bt.Client.Source
		}
		result = append(result, &Result{
			Text:           val.Translations[0].Text,
			SourceText:     contents[k],
			SourceLanguage: sourceLanguage,
		})
	}
	return result, nil
}
