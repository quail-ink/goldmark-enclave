package enclave

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type astTransformer struct{}

var defaultASTTransformer = &astTransformer{}

func (a *astTransformer) InsertFailedHint(n ast.Node, msg string) {
	msgNode := ast.NewString([]byte(fmt.Sprintf("\n<!-- goldmark-enclave: %s -->\n", msg)))
	msgNode.SetCode(true)
	n.Parent().InsertAfter(n.Parent(), n, msgNode)
}

func (a *astTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	replaceImages := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		img := n.(*ast.Image)
		u, err := url.Parse(string(img.Destination))
		if err != nil {
			a.InsertFailedHint(n, fmt.Sprintf("failed to parse url: %s, %s", img.Destination, err))
			return ast.WalkContinue, nil
		}

		oid := ""
		theme := "light"
		provider := ""
		if u.Host == "www.youtube.com" && u.Path == "/watch" {
			// this is a youtube video: https://www.youtube.com/watch?v={vid}
			provider = EnclaveProviderYouTube
			oid = u.Query().Get("v")
		} else if u.Host == "youtu.be" {
			// this is also a youtube video: https://youtu.be/{vid}
			provider = EnclaveProviderYouTube
			oid = u.Path[1:]
			oid = strings.Trim(oid, "/")

		} else if u.Host == "www.bilibili.com" && strings.HasPrefix(u.Path, "/video/") {
			// this is a bilibili video: https://www.bilibili.com/video/{vid}
			provider = EnclaveProviderBilibili
			oid = u.Path[7:]
			oid = strings.Trim(oid, "/")

		} else if u.Host == "twitter.com" || u.Host == "m.twitter.com" || u.Host == "x.com" {
			// https://twitter.com/{username}/status/{id number}?theme=dark
			provider = EnclaveProviderTwitter
			oid = string(img.Destination)
			if u.Host == "x.com" {
				// replace x.com with twitter.com, because x.com doesn't support using x.com as the source host, what a shame
				oid = strings.Replace(oid, "x.com", "twitter.com", 1)
			}
			theme = u.Query().Get("theme")

		} else if u.Host == "tradingview.com" || u.Host == "www.tradingview.com" {
			// https://www.tradingview.com/chart/UC0wWW9o/?symbol=BITFINEX%3ABTCUSD
			provider = EnclaveProviderTradingView
			oid = u.Query().Get("symbol")
			theme = u.Query().Get("theme")

		} else {
			a.InsertFailedHint(n, fmt.Sprintf("unsupported object: %s", img.Destination))
			return ast.WalkContinue, nil
		}

		if oid != "" {
			ev := NewEnclave(
				&Enclave{
					Image:    *img,
					Provider: provider,
					ObjectID: oid,
					Theme:    theme,
				})
			n.Parent().ReplaceChild(n.Parent(), n, ev)
		}

		return ast.WalkContinue, nil
	}

	ast.Walk(node, replaceImages)
}
