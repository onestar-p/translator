package translator

import (
	"context"
	"testing"
)

func TestTranslate(t *testing.T) {
	c := context.Background()

	cases := []struct {
		Name   string
		Text   []SourceText
		Key    string
		Target string
		Func   func(key, target string, text []SourceText) *Client
	}{
		{
			Name: "Bing_Translate",
			Text: []SourceText{
				"必应翻译，用户反馈",
				"必应翻译，用户反馈2",
			},
			Key:    "key",
			Target: "en",
			Func: func(key, target string, text []SourceText) *Client {
				return NewTranslator(key, target, text).BingTranslate()
			},
		},
		{
			Name: "Google_Translate",
			Text: []SourceText{
				"谷歌翻译，用户反馈",
				"谷歌翻译，用户反馈2",
			},
			Key:    "key",
			Target: "fr",
			Func: func(key, targe string, text []SourceText) *Client {
				return NewTranslator(key, targe, text).GoogleTranslate()
			},
		},
		{
			Name: "ALIMT_Translate",
			Text: []SourceText{
				"阿里机器翻译，用户反馈",
				"阿里机器翻译，用户反馈222",
			},
			Key:    "accessKeyID/accessKeySecret", // accessKeyID/accessKeySecret
			Target: "hant",
			Func: func(key, targe string, text []SourceText) *Client {

				return NewTranslator(key, targe, text,
					WithSource("hans"),
					WithScene("general"),
					WithApiType("translate_standard"),
					WithMaxTransTextNums(50),
				).ALIMTTranslate()
			},
		},
	}

	for _, cc := range cases {
		t.Run(cc.Name, func(t *testing.T) {
			_, err := cc.Func(cc.Key, cc.Target, cc.Text).Do(c)
			if err != nil {
				t.Errorf("cannot translate err:%s", err)
			}
		})

	}
}
