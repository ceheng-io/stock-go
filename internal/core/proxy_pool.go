package core

import (
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
)

// ProxyPoolConfig contains round-robin proxy pool settings.
type ProxyPoolConfig struct {
	URLs []string
}

// ProxyPool selects proxy URLs in round-robin order.
type ProxyPool struct {
	urls []*url.URL
	next atomic.Uint64
}

// NewProxyPool creates a round-robin proxy pool from valid proxy URLs.
func NewProxyPool(rawURLs []string) *ProxyPool {
	if len(rawURLs) == 0 {
		return nil
	}
	urls := make([]*url.URL, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		proxyURL, err := url.Parse(rawURL)
		if err != nil || proxyURL.Scheme == "" || proxyURL.Host == "" {
			continue
		}
		urls = append(urls, proxyURL)
	}
	if len(urls) == 0 {
		return nil
	}
	return &ProxyPool{urls: urls}
}

// Proxy returns the next proxy URL for an HTTP request.
func (p *ProxyPool) Proxy(_ *http.Request) (*url.URL, error) {
	if p == nil || len(p.urls) == 0 {
		return nil, nil
	}
	index := p.next.Add(1) - 1
	return p.urls[index%uint64(len(p.urls))], nil
}
