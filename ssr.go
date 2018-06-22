package ssr // import "github.com/captaincodeman/appengine-ssr"

import (
	"context"
	"time"

	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ua-parser/uap-go/uaparser"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/captaincodeman/appengine-context"
)

var defaultUserAgents = []string{
	"W3C_Validator",
	"baiduspider",
	"bingbot",
	"facebookexternalhit",
	"LinkedInBot",
	"Pinterest",
	"Slackbot-LinkExpanding",
	"TwitterBot",
	"Googlebot",
	"Mediapartners-Google",
}

type (
	// Option is a configuration function
	Option func(*SSR)

	// SSR holds the configuration for the
	// Server Side Rendering middleware
	SSR struct {
		Puppeteer     string
		Parser        *uaparser.Parser
		UserAgents    []string
		HeadlessParam string
		OverrideParam string
		Cache         CacheProvider
		Timeout       time.Duration
		Expiration    time.Duration
		Verbose       bool
	}
)

// NewSSR creates a new Server Side Rendering middleware
func NewSSR(puppeteer string, opts ...Option) *SSR {
	ssr := &SSR{
		Puppeteer:     puppeteer,
		Parser:        uaparser.NewFromSaved(),
		UserAgents:    defaultUserAgents,
		HeadlessParam: "headless",
		OverrideParam: "ssr",
		Cache:         &memcacheProvider{"ssr:"},
		Timeout:       time.Second * 30,
		Expiration:    time.Hour,
	}

	for _, opt := range opts {
		opt(ssr)
	}

	return ssr
}

// UserAgents configures the UserAgent strings to check
// when determining whether to Server Side Render or not
func UserAgents(userAgents []string) Option {
	return func(ssr *SSR) {
		ssr.UserAgents = userAgents
	}
}

// UserAgentParser allows a non-default User Agent parser
func UserAgentParser(parser *uaparser.Parser) Option {
	return func(ssr *SSR) {
		ssr.Parser = parser
	}
}

// HeadlessParam specifies the name of the querystring
// parameter used to indicate a 'headless' request
func HeadlessParam(name string) Option {
	return func(ssr *SSR) {
		ssr.HeadlessParam = name
	}
}

// OverrideParam specifies the name of the querystring
// parameter used to override SSR rendering (set to
// empty to disable)
func OverrideParam(name string) Option {
	return func(ssr *SSR) {
		ssr.OverrideParam = name
	}
}

// Verbose specifies whether to output verbose logging
func Verbose(verbose bool) Option {
	return func(ssr *SSR) {
		ssr.Verbose = verbose
	}
}

// Timeout specifies the SSR request timeout
func Timeout(timeout time.Duration) Option {
	return func(ssr *SSR) {
		ssr.Timeout = timeout
	}
}

// Middleware ...
func (ssr *SSR) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !ssr.IsHeadless(r) && ssr.IsBot(r) {
			ssr.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	// because we need the appengine context for memcache etc.
	return gaecontext.Middleware(http.HandlerFunc(fn))
}

// IsHeadless returns true if this is a headless request
func (ssr *SSR) IsHeadless(r *http.Request) bool {
	q := r.URL.Query()
	_, headless := q[ssr.HeadlessParam]
	return headless
}

// IsBot returns true if the User Agent provided is a bot
func (ssr *SSR) IsBot(r *http.Request) bool {
	if ssr.OverrideParam != "" {
		q := r.URL.Query()
		if _, ok := q[ssr.OverrideParam]; ok {
			return true
		}
	}

	ua := ssr.Parser.ParseUserAgent(r.UserAgent())

	for _, bot := range ssr.UserAgents {
		if bot == ua.Family {
			return true
		}
	}

	return false
}

// ServeHTTP responds to the request with the Server Side Rendered
// content, caching the results in memcache
func (ssr *SSR) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := gaecontext.Context(r)

	// TODO: provide a transform function to sanitize the URL,
	// remove query strings that shouldn't change the SSR etc...
	key := r.URL.String()

	// TODO: timeout for cache requests
	if data, err := ssr.Cache.Get(ctx, key); err == nil {
		// TODO: cache a GOB struct with bytes + response code, headers etc...
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}

	// add headless query parameter to URL for SSR client use
	q := r.URL.Query()
	q.Set(ssr.HeadlessParam, "")
	if ssr.OverrideParam != "" {
		// remove the querystring override if set
		q.Del(ssr.OverrideParam)
	}
	r.URL.RawQuery = q.Encode()

	u := ssr.Puppeteer + "/ssr?url=" + url.QueryEscape(r.URL.String())

	if ssr.Verbose {
		log.Debugf(ctx, "ssr req %s", u)
	}

	ctx, cancel := context.WithTimeout(ctx, ssr.Timeout)
	defer cancel()

	// TODO: add auth token to request so the
	// Puppeteer endpoint can be locked down
	client := urlfetch.Client(ctx)
	resp, err := client.Get(u)

	if err != nil {
		if ssr.Verbose {
			log.Errorf(ctx, "ssr error %v", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if ssr.Verbose {
			log.Errorf(ctx, "read err %v", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: timeout for cache call in case it's unavailable / slow
	ssr.Cache.Put(ctx, key, ssr.Expiration, data)

	w.WriteHeader(resp.StatusCode)
	w.Write(data)
}
