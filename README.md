# goldmark-enclave

This [goldmark](http://github.com/yuin/goldmark) extension uses Markdown's image syntax  `![]()` to support other objects.

## Supported Objects

- [x] [YouTube](https://youtube.com) Video
- [x] [Bilibili](https://bilibili.com) Video
- [x] X(formly Twitter)'s Tweet Oembed Widget
- [x] [TradingView](https://tradingview.com) Chart
- [x] [Quail](https://quail.ink) List and Article

## Planned Objects

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
Youtube Video:

![](https://youtu.be/dQw4w9WgXcQ?si=0kalBBLQpIXT1Wcd)

Bilibili Video:

![](https://www.bilibili.com/video/BV1uT4y1P7CX)

Twitter Tweet:

![](https://twitter.com/NASA/status/1704954156149084293)

TradingView Chart:

![](https://www.tradingview.com/chart/AA0aBB8c/?symbol=BITFINEX%3ABTCUSD)

Quail List and Post

![](https://quail.ink/blog)

![](https://quail.ink/blog/p/extended-markdown-syntax?theme=dark)
```

## Demo

[Live Demo](https://quail.ink/blog/p/extended-markdown-syntax)

## Options

- `theme`: The theme of the TradingView chart, twitter tweet and quail widget. Default: `light`
  - e.g. `![](https://twitter.com/NASA/status/1704954156149084293?theme=dark)`

## Installation

```bash
go get github.com/quail.ink/goldmark-enclave
```
