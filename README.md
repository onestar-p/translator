# translator
Google Translation &amp; Bing Translation


### Install
```
go get github.com/onestar-p/translator
```

### e.g.
```
c := context.Background()
var (
  key    = "key"
  target = "en"
  text   = []translator.SourceText{
    "翻译测试",
    "翻译测试",
  }
)
res, err := translator.NewTranslator(key, target, text).GoogleTranslate().Do(c)
if err != nil {
  panic(err)
}
for _, r := range res {
  fmt.Println(r.Text)
}
```
