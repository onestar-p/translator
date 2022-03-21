// 翻译
// 一次请求第三方翻译接口的Text数量是有限制的，由”Client.MaxTransTextNums“决定每次翻译的数据量，
// 如果翻译的总数大于MaxTransTextNums，将会并发进行请求第三方进行翻译，
// MaxTransTextNums的默认值为100，NewTranslator()翻译对象是可使用WithMaxTransTextNums()配置项进行调节，建议设置值不要大于默认值

package translator

import (
	"context"
	"fmt"
	"math"
	"time"
	"unicode/utf8"
)

const (
	DefaultFormat           = "html"
	DefaultMaxTransTextNums = 100 // 默认一次性最多翻译数量
	DefaultRequestTimeout   = 30  // 请求超时时间（秒）
)

type TranslatorInterfaced interface {
	Trans(c context.Context, contents []SourceText) ([]*Result, error) // 单个翻译
}
type SourceText string     // 翻译源文本
type TargetLanguage string // 目标语言
type SourceLanguage string // 源语言

// 翻译结果
type Result struct {
	SourceText     SourceText     // 翻译原文本，与翻译结果不一定一致，需要测试
	Text           string         // 翻译结果
	SourceLanguage SourceLanguage // 源语言
	WordNums       uint32         // 翻译单词数
}

type Client struct {
	Key              string         // API 密钥
	Target           TargetLanguage // 目标语言
	Source           SourceLanguage // 源语言
	Format           string         // 源文本的格式，HTML（默认）或纯文本格式。值html表示 HTML，值text 表示纯文本。
	Contents         []SourceText   // 翻译源文本
	ContenTotal      int            // 翻译源文本总数
	MaxTransTextNums uint           // 每次翻译最多翻译量
	Request          *Request
	translate        TranslatorInterfaced
}

// 返回翻译客户端
// key 翻译API秘钥
// content 翻译内容
// target 目标语言
// opts 配置项 可查看option.go
func NewTranslator(key, target string, content []SourceText, opts ...Option) *Client {
	c := &Client{
		Key:              key,
		Contents:         content,
		ContenTotal:      len(content),
		Format:           DefaultFormat,
		MaxTransTextNums: DefaultMaxTransTextNums,
		Target:           TargetLanguage(target),
	}

	// 设置超时时间
	c.Request = NewRequest(&RequestOption{
		Timeout: time.Second * DefaultRequestTimeout,
	})

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// run translate
//  e.g. translate.GoogleTranslate().Do(ctx)
func (c *Client) Do(ctx context.Context) ([]*Result, error) {
	if c.translate == nil {
		return nil, fmt.Errorf("c.translate is nil. e.g. ")
	}

	// 数据分组
	contents := arrayChunk(c.Contents, int(c.MaxTransTextNums))
	type results struct {
		Result []*Result
		Index  int // 用于存放分组后的切片的下标
		Err    error
	}
	chNums := len(contents)
	chResults := make(chan results, chNums)
	for i, content := range contents {
		go func(chResults chan results, content []SourceText, index int) {
			var r results
			r.Index = index
			r.Result, r.Err = c.translate.Trans(ctx, content)
			chResults <- r
		}(chResults, content, i)
	}

	// 分组翻译后还原原数组格式返回
	tranRes := []*Result{}
	var err error
	tranResChunk := make([][]*Result, chNums)
	for i := 0; i < chNums; i++ {
		res := <-chResults
		if res.Err == nil {
			tranResChunk[res.Index] = res.Result // 保存到切片中使之排序保存
		} else {
			err = res.Err
		}
	}
	for _, tr := range tranResChunk {
		tranRes = append(tranRes, tr...)
	}

	if len(tranRes) == 0 && err != nil {
		return nil, err
	}

	// 计算源语言单词数
	for _, val := range tranRes {
		val.WordNums = uint32(utf8.RuneCountInString(string(val.SourceText)))
	}
	return tranRes, nil
}

// 切片分组
func arrayChunk(contents []SourceText, nums int) [][]SourceText {
	total := len(contents)
	pageNum := int(math.Ceil(float64(total) / float64(nums)))
	var resArr [][]SourceText
	for i := 0; i < pageNum; i++ {
		offset := int(i * nums)
		end := int((i + 1) * nums)
		end2 := end - total
		var limit int
		if i+1 == pageNum {
			limit = end - end2
		} else {
			limit = end
		}
		temp := contents[offset:limit]
		resArr = append(resArr, temp)
	}
	return resArr
}

// Google翻译
func (c *Client) GoogleTranslate() *Client {
	c.translate = &GoogleTranslateService{Client: c}
	return c
}

// Bing翻译
func (c *Client) BingTranslate() *Client {
	c.translate = &BingTranslateService{Client: c}
	return c
}
