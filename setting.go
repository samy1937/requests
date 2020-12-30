package req

import (
	"crypto/tls"
	"errors"
	"github.com/projectdiscovery/retryablehttp-go"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	two                   = 2
	ten                   = 10
	defaultMaxWorkers     = 150
	defaultMaxHistorydata = 150
)

var transport = http.Transport{

	DisableKeepAlives: true, // 关闭长长连接
	Proxy:             http.ProxyFromEnvironment,
	//	ResponseHeaderTimeout: time.Duration(2) * time.Second,
	//DisableCompression:    true,
	MaxIdleConnsPerHost: -1,
	MaxIdleConns:        -1,
	//IdleConnTimeout:       1 * time.Second,
	//TLSHandshakeTimeout:   5 * time.Second,
	//ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig: &tls.Config{
		Renegotiation:      tls.RenegotiateOnceAsClient,
		InsecureSkipVerify: true,
	},
	Dial: func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, time.Second*8) //设置建立连接超时
		if err != nil {
			return nil, err
		}
		_ = c.SetDeadline(time.Now().Add(8 * time.Second)) //设置发送接收数据超时
		return c, nil
	},
}

type checkRedirectFunc func(_ *http.Request, requests []*http.Request) error

func makeCheckRedirectFunc(followRedirects bool, maxRedirects int) checkRedirectFunc {
	return func(_ *http.Request, requests []*http.Request) error {
		if !followRedirects {
			return http.ErrUseLastResponse
		}

		if maxRedirects == 0 {
			if len(requests) > ten {
				return http.ErrUseLastResponse
			}

			return nil
		}

		if len(requests) > maxRedirects {
			return http.ErrUseLastResponse
		}

		return nil
	}
}

// create a default client
func newClient() *retryablehttp.Client {

	jar, _ := cookiejar.New(nil)
	retryablehttpOptions := retryablehttp.DefaultOptionsSpraying
	disableKeepAlives := true
	maxIdleConns := 0
	maxConnsPerHost := 0
	maxIdleConnsPerHost := -1

	retryablehttpOptions.RetryWaitMax = 10 * time.Second
	retryablehttpOptions.RetryMax = 1
	followRedirects := true
	maxRedirects := 0

	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		MaxConnsPerHost:     maxConnsPerHost,
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: disableKeepAlives,
	}

	return retryablehttp.NewWithHTTPClient(&http.Client{
		Transport:     transport,
		Jar:           jar,
		Timeout:       time.Duration(10) * time.Second,
		CheckRedirect: makeCheckRedirectFunc(followRedirects, maxRedirects),
	}, retryablehttpOptions)
}

func (r *Req) GetTransport() *http.Transport {
	if r.transport == nil {
		r.transport = &http.Transport{

			DisableKeepAlives: true, // 关闭长长连接
			Proxy:             http.ProxyFromEnvironment,
			//	ResponseHeaderTimeout: time.Duration(2) * time.Second,
			//DisableCompression:    true,
			MaxIdleConnsPerHost: -1,
			MaxIdleConns:        -1,
			//IdleConnTimeout:       1 * time.Second,
			//TLSHandshakeTimeout:   5 * time.Second,
			//ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*8) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				_ = c.SetDeadline(time.Now().Add(8 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		}

		return r.transport
	}
	return r.transport
}

// Client return the default underlying http client
func (r *Req) Client() *retryablehttp.Client {
	if r.client == nil {
		r.client = newClient()
	}
	return r.client
}

// Client return the default underlying http client
func Client() *retryablehttp.Client {
	return std.Client()
}

// SetClient sets the underlying http.Client.
func (r *Req) SetClient(client *retryablehttp.Client) {
	r.client = client // use default if client == nil
}

// SetClient sets the default http.Client for requests.
func SetClient(client *http.Client) {
	c := newClient()
	c.HTTPClient = client
	std.SetClient(c)
}

// SetFlags control display format of *Resp
func (r *Req) SetFlags(flags int) {
	r.flag = flags
}

// SetFlags control display format of *Resp
func SetFlags(flags int) {
	std.SetFlags(flags)
}

// Flags return output format for the *Resp
func (r *Req) Flags() int {
	return r.flag
}

// Flags return output format for the *Resp
func Flags() int {
	return std.Flags()
}

func (r *Req) getTransport() *http.Transport {
	trans, _ := r.Client().HTTPClient.Transport.(*http.Transport)
	return trans
}

// EnableInsecureTLS allows insecure https
func (r *Req) EnableInsecureTLS(enable bool) {
	trans := r.getTransport()
	if trans == nil {
		return
	}
	if trans.TLSClientConfig == nil {
		trans.TLSClientConfig = &tls.Config{}
	}
	trans.TLSClientConfig.InsecureSkipVerify = enable
}

func EnableInsecureTLS(enable bool) {
	std.EnableInsecureTLS(enable)
}

// EnableCookieenable or disable cookie manager
func (r *Req) EnableCookie(enable bool) {
	if enable {
		jar, _ := cookiejar.New(nil)
		r.Client().HTTPClient.Jar = jar
		r.Client().HTTPClient.Jar = jar
	} else {
		r.Client().HTTPClient.Jar = nil
	}
}

// EnableCookieenable or disable cookie manager
func EnableCookie(enable bool) {
	std.EnableCookie(enable)
}

// SetTimeout sets the timeout for every request
func (r *Req) SetTimeout(d time.Duration) {
	r.Client().HTTPClient.Timeout = d
}

// SetTimeout sets the timeout for every request
func SetTimeout(d time.Duration) {
	std.SetTimeout(d)
}

// SetProxyUrl set the simple proxy with fixed proxy url
func (r *Req) SetProxyUrl(rawurl string) error {
	trans := r.getTransport()
	if trans == nil {
		return errors.New("req: no transport")
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	trans.Proxy = http.ProxyURL(u)
	return nil
}

// SetProxyUrl set the simple proxy with fixed proxy url
func SetProxyUrl(rawurl string) error {
	return std.SetProxyUrl(rawurl)
}

// SetProxy sets the proxy for every request
func (r *Req) SetProxy(proxy func(*http.Request) (*url.URL, error)) error {
	trans := r.getTransport()
	if trans == nil {
		return errors.New("req: no transport")
	}
	trans.Proxy = proxy
	return nil
}

// SetProxy sets the proxy for every request
func SetProxy(proxy func(*http.Request) (*url.URL, error)) error {
	return std.SetProxy(proxy)
}

type jsonEncOpts struct {
	indentPrefix string
	indentValue  string
	escapeHTML   bool
}

func (r *Req) getJSONEncOpts() *jsonEncOpts {
	if r.jsonEncOpts == nil {
		r.jsonEncOpts = &jsonEncOpts{escapeHTML: true}
	}
	return r.jsonEncOpts
}

// SetJSONEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func (r *Req) SetJSONEscapeHTML(escape bool) {
	opts := r.getJSONEncOpts()
	opts.escapeHTML = escape
}

// SetJSONEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func SetJSONEscapeHTML(escape bool) {
	std.SetJSONEscapeHTML(escape)
}

// SetJSONIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func (r *Req) SetJSONIndent(prefix, indent string) {
	opts := r.getJSONEncOpts()
	opts.indentPrefix = prefix
	opts.indentValue = indent
}

// SetJSONIndent instructs the encoder to format each subsequent encoded
// value as if indented by the package-level function Indent(dst, src, prefix, indent).
// Calling SetIndent("", "") disables indentation.
func SetJSONIndent(prefix, indent string) {
	std.SetJSONIndent(prefix, indent)
}

type xmlEncOpts struct {
	prefix string
	indent string
}

func (r *Req) getXMLEncOpts() *xmlEncOpts {
	if r.xmlEncOpts == nil {
		r.xmlEncOpts = &xmlEncOpts{}
	}
	return r.xmlEncOpts
}

// SetXMLIndent sets the encoder to generate XML in which each element
// begins on a new indented line that starts with prefix and is followed by
// one or more copies of indent according to the nesting depth.
func (r *Req) SetXMLIndent(prefix, indent string) {
	opts := r.getXMLEncOpts()
	opts.prefix = prefix
	opts.indent = indent
}

// SetXMLIndent sets the encoder to generate XML in which each element
// begins on a new indented line that starts with prefix and is followed by
// one or more copies of indent according to the nesting depth.
func SetXMLIndent(prefix, indent string) {
	std.SetXMLIndent(prefix, indent)
}

// SetProgressInterval sets the progress reporting interval of both
// UploadProgress and DownloadProgress handler
func (r *Req) SetProgressInterval(interval time.Duration) {
	r.progressInterval = interval
}

// SetProgressInterval sets the progress reporting interval of both
// UploadProgress and DownloadProgress handler for the default client
func SetProgressInterval(interval time.Duration) {
	std.SetProgressInterval(interval)
}
