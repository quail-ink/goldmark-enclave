# goldmark-enclave

This [goldmark](http://github.com/yuin/goldmark) extension extend commonmark syntax:

- uses Markdown's image syntax `![]()` to support other objects.
- adds highlight syntax for inline text.
- adds title to links
- supports pandoc-style [fenced](https://pandoc.org/MANUAL.html#divs-and-spans) divs in goldmark.
- supports github-style [callouts](https://github.com/orgs/community/discussions/16925) in goldmark.

## Full Demo

[Live Demo](https://quaily.com/blog/p/extended-markdown-syntax)

## Embeded Objects

### Supported Objects

- [x] [YouTube](https://youtube.com) Video
- [x] [Bilibili](https://bilibili.com) Video
- [x] X(formly Twitter)'s Tweet Oembed Widget
- [x] [TradingView](https://tradingview.com) Chart
- [x] [Quaily](https://quaily.com) List and Article
- [x] Add options to images
- [x] [dify.ai](https://dify.ai) Widget
- [x] [Spotify](https://spotify.com) Embed
- [x] html5 audio

### Usage

```go
import (
  enclave "github.com/quailyquaily/goldmark-enclave"
	"github.com/yuin/goldmark"
)
// ...
markdown := goldmark.New(
  goldmark.WithExtensions(
    enclave.New(),
  ),
)
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

![](https://quaily.com/blog)

![](https://quaily.com/blog/p/extended-markdown-syntax?theme=dark)

Image with caption and giving it a width:

![](https://your-image.com/image.png?w=100px "This is a caption")

Dify Widget

![](https://udify.app/chatbot/1NaVTsaJ1t54UrNE)

Spotify Embed

![](https://open.spotify.com/track/5vdp5UmvTsnMEMESIF2Ym7?si=d4ee09bfd0e941c5)

HTML5 Audio

![](https://cdn1.suno.ai/fc991b95-e4e9-4c8f-87e8-e5e4560755e7.mp3)
```

### Some Options

Some objects support options:

- `theme`: The theme of the TradingView chart, twitter tweet and quaily widget. Default: `light`
  - e.g. `![](https://twitter.com/NASA/status/1704954156149084293?theme=dark)`
- `width` / `w` and `height` / `h`: The width and height of images. Default: `auto`
  - e.g. `![](https://your-image.com/image.png?w=100px)`

## Other features

#### Highlight Text

```go
import (
  enclaveMark "github.com/quailyquaily/goldmark-enclave/mark"
	"github.com/yuin/goldmark"
)
// ...
markdown := goldmark.New(
  goldmark.WithExtensions(
    enclaveMark.New(),
  ),
)
```

```md
This is a ==highlighted text==.
```

will be rendered as:

```html
<p>This is a <mark>highlighted text</mark>.</p>
```

### Title to Links

```go
import (
  enclaveHref "github.com/quailyquaily/goldmark-enclave/href"
	"github.com/yuin/goldmark"
)
// ...
markdown := goldmark.New(
  goldmark.WithExtensions(
    enclaveHref.New(&enclaveHref.Config{}),
  ),
)
```

```md
[Quail](/blog "Quail Blog")
```

will be rendered as:

```html
<a href="https://quaily.com/blog" title="Quail Blog">Quail</a>
```

### Pandoc-style Fenced Divs

```go
import (
  enclaveFence "github.com/quailyquaily/goldmark-enclave/fence"
	"github.com/yuin/goldmark"
)
// ...
markdown := goldmark.New(
  goldmark.WithExtensions(
    enclaveFence.New(),
  ),
)
```

### Github-style Callouts

```go
import (
  enclaveCallout "github.com/quailyquaily/goldmark-enclave/callout"
	"github.com/yuin/goldmark"
)
// ...
markdown := goldmark.New(
  goldmark.WithExtensions(
    enclaveCallout.New(),
  ),
)
```

## Installation

```bash
go get github.com/quailyquaily/goldmark-enclave
```
