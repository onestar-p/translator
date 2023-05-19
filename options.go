package translator

// 翻译客户端配置
type Option func(c *Client)

// 设置翻译源语言
func WithSource(s SourceLanguage) Option {
	return func(c *Client) {
		c.Source = s
	}
}

// 设置源文本的格式
func WithFormat(f string) Option {
	return func(c *Client) {
		if f != "" {
			c.Format = f
		}
	}
}

// 设置每次翻译最多翻译量
func WithMaxTransTextNums(n uint) Option {
	return func(c *Client) {
		if n != 0 {
			c.MaxTransTextNums = n
		}
	}
}

// 场景
// 专业版本支持的场景：商品标题（title），商品描述（description），商品沟通（communication），医疗（medical），社交（social)
// 通用版本支持的场景：general
func WithScene(scene string) Option {
	return func(c *Client) {
		c.Scene = scene
	}
}

// 版本类型
// 通用版本：translate_standard
// 专业版本：translate_ecommerce
func WithApiType(t string) Option {
	return func(c *Client) {
		c.ApiType = t
	}
}
