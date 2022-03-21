package translator

import (
	"context"
	"testing"
)

func TestTranslate(t *testing.T) {
	c := context.Background()

	cases := []struct {
		Name string
		Text []SourceText
		Key  string
		Func func(key string, text []SourceText) *Client
	}{
		{
			Name: "Bing_Translate",
			Text: []SourceText{
				"必应翻译，用户反馈",
				"必应翻译，用户反馈2",
			},
			Key: "key",
			Func: func(key string, text []SourceText) *Client {
				return NewTranslator(key, text, WithTarget("en")).BingTranslate()
			},
		},
		{
			Name: "Google_Translate",
			Text: []SourceText{
				"谷歌翻译，用户反馈",
				"谷歌翻译，用户反馈2",
			},
			Key: "key",
			Func: func(key string, text []SourceText) *Client {
				return NewTranslator(key, text, WithTarget("en")).GoogleTranslate()
			},
		},
	}

	for _, cc := range cases {
		t.Run(cc.Name, func(t *testing.T) {
			_, err := cc.Func(cc.Key, cc.Text).Do(c)
			if err != nil {
				t.Errorf("cannot translate err:%s", err)
			}
		})

	}
}
