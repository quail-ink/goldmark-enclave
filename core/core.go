package core

import (
	"net/url"

	"github.com/yuin/goldmark/ast"
)

const (
	EnclaveProviderYouTube     = "youtube"
	EnclaveProviderBilibili    = "bilibili"
	EnclaveProviderTwitter     = "twitter"
	EnclaveProviderTradingView = "tradingview"
	EnclaveProviderDifyWidget  = "dify-widget"
	EnclaveProviderQuailWidget = "quail-widget"
	EnclaveProviderQuailImage  = "quail-image"
	EnclaveRegularImage        = "regular-image"
)

const (
	IframeDisabledPlaceholderURL = "https://static.quail.ink/assets/not-available-in-email.png"
)

type (
	Config struct {
		DefaultImageAltPrefix string
		IframeDisabled        bool
		VideoDisabled         bool
		TwitterDisabled       bool
		TradingViewDisabled   bool
		DifyWidgetDisabled    bool
		QuailWidgetDisabled   bool
	}

	Enclave struct {
		ast.Image
		URL            *url.URL
		IframeDisabled bool
		Provider       string
		ObjectID       string
		Theme          string
		Params         map[string]string
	}
)

var KindEnclave = ast.NewNodeKind("Enclave")

func (n *Enclave) Kind() ast.NodeKind {
	return KindEnclave
}
