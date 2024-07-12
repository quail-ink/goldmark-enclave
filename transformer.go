package enclave

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type astTransformer struct{}

var defaultASTTransformer = &astTransformer{}

var imgLeftNode ast.Node
var imgLeftParentNode ast.Node

func (a *astTransformer) InsertFailedHint(n ast.Node, msg string) {
	msgNode := ast.NewString([]byte(fmt.Sprintf("\n<!-- goldmark-enclave: %s -->\n", msg)))
	msgNode.SetCode(true)
	n.Parent().InsertAfter(n.Parent(), n, msgNode)
}

func (a *astTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	replaceImages := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		fmt.Printf("1 n.Kind(): %v\n", n.Kind())
		if !entering {
			return ast.WalkContinue, nil
		}

		fmt.Printf("2 n.Kind(): %v\n", n.Kind())

		if n.Kind() != ast.KindImage {
			if n.Kind() == ast.KindParagraph {
				imgLeftNode = n
				imgLeftParentNode = n.Parent()
			}
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
		params := map[string]string{}
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

		} else if u.Host == "quail.ink" || u.Host == "dev.quail.ink" {
			// https://quail.ink/{list_slug} or https://quail.ink/{list_slug}/p/{post_slug}
			const re1 = `^([a-zA-Z0-9_-]+)$`
			const re2 = `^([a-zA-Z0-9_-]+)/p/([a-zA-Z0-9_-]+)$`
			if len(u.Path) > 1 {
				p := strings.Trim(u.Path[1:], "/")
				ok1, _ := regexp.MatchString(re1, p)
				ok2, _ := regexp.MatchString(re2, p)
				if ok1 || ok2 {
					provider = EnclaveProviderQuailWidget
					oid = string(img.Destination)
					theme = u.Query().Get("theme")
					params["layout"] = u.Query().Get("layout")
				}
			}
		} else {
			title := string(img.Title)
			w := u.Query().Get("w")
			if w == "" {
				w = u.Query().Get("width")
			}
			h := u.Query().Get("h")
			if h == "" {
				h = u.Query().Get("height")
			}
			if len(title) != 0 || w != "" || h != "" {
				// this is a normal image, but it has a title, so we add a caption
				provider = EnclaveProviderQuailImage
				oid = string(img.Destination)
				if title != "" {
					params["title"] = string(img.Title)
				}
				if w != "" {
					params["width"] = w
				}
				if h != "" {
					params["height"] = h
				}
			} else {
				provider = EnclaveRegularImage
				oid = string(img.Destination)
			}
		}

		if oid != "" {
			ev := NewEnclave(
				&Enclave{
					Image:    *img,
					URL:      u,
					Provider: provider,
					ObjectID: oid,
					Theme:    theme,
					Params:   params,
				})

			// if the outter node is a paragraph node, replace the whole paragraph.
			// because we can not put div in a p tag
			// if imgLeftNode != nil && imgLeftNode.Kind() == ast.KindParagraph && imgLeftParentNode != nil {
			// 	imgLeftParentNode.ReplaceChild(imgLeftParentNode, n, ev)
			// 	imgLeftNode = nil
			// 	imgLeftParentNode = nil
			// } else {
			n.Parent().ReplaceChild(n.Parent(), n, ev)
			// }
		}

		return ast.WalkContinue, nil
	}

	ast.Walk(node, replaceImages)
}
