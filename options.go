package translator

// 翻译客户端配置
type Option func(c *Client)

// 设置翻译源语言
func WithScopes(s SourceLanguage) Option {
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
