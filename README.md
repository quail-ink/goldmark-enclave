# goldmark-enclave

This [goldmark](http://github.com/yuin/goldmark) extension uses Markdown's image syntax  `![]()` to support other objects.

## Supported Objects

- [x] [YouTube](https://youtube.com) Video
- [x] [Bilibili](https://bilibili.com) Video
- [x] X(formly Twitter)'s Tweet Oembed Widget
- [x] [TradingView](https://tradingview.com) Chart

## Planned Objects

- [ ] [Quail](https://quail.ink) List and Article
- [ ] [Discord](https://discord.com) Server Widget
- [ ] Add options to resize and position images

## Usage

```go
  markdown := goldmark.New(
    goldmark.WithExtensions(
      enclave.New(),
    ),
  )
  var buf bytes.Buffer
  if err := markdown.Convert([]byte(source), &buf); err != nil {
    panic(err)
  }
  fmt.Print(buf)
}
```

And then you can use it like this:

```md
![](https://youtu.be/dQw4w9WgXcQ?si=0kalBBLQpIXT1Wcd)
```

```md
![](https://www.bilibili.com/video/BV1uT4y1P7CX)
```

```md
![](https://twitter.com/NASA/status/1704954156149084293)
```

### Installation

```bash
go get github.com/quail.ink/goldmark-enclave
```