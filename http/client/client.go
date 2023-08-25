package client

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/multipart"
)

type Client struct {
	*http.Client
	header  http.Header
	useGzip bool
}

func TimeoutDialer(cTimeout, rwTimeout time.Duration) func(netw, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}

		if rwTimeout > 0 {
			conn.SetDeadline(time.Now().Add(rwTimeout))
		}

		return conn, nil
	}
}

func NewClientWithTransport(transport *http.Transport) *Client {
	c := &Client{Client: &http.Client{}}
	c.Transport = transport
	return c
}

func NewClientWhthTransportAndHeader(transport *http.Transport, header http.Header) *Client {
	c := &Client{
		Client: &http.Client{},
		header: header,
	}

	c.Transport = transport
	return c
}

func TransportWithProxyAndTimeOut(proxy func(*http.Request) (*url.URL, error), tlsClientConfig *tls.Config, headerTimeout, connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy:                 proxy,
		Dial:                  TimeoutDialer(connectTimeout, readWriteTimeout),
		TLSClientConfig:       tlsClientConfig,
		ResponseHeaderTimeout: headerTimeout,
	}
}

func TransportWithTimeOut(tlsClientConfig *tls.Config, headerTimeout, connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	return TransportWithProxyAndTimeOut(http.ProxyFromEnvironment, tlsClientConfig, headerTimeout, connectTimeout, readWriteTimeout)
}

func TransportWithProxyAndSkipVerifyAndTimeOut(proxy func(*http.Request) (*url.URL, error), headerTimeout, connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	return TransportWithProxyAndTimeOut(proxy, tlsClientConfig, headerTimeout, connectTimeout, readWriteTimeout)
}

func TransportWithSkipVerifyAndTimeOut(headerTimeout, connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	return TransportWithTimeOut(tlsClientConfig, headerTimeout, connectTimeout, readWriteTimeout)
}

func DefaultTransport() *http.Transport {
	return TransportWithSkipVerifyAndTimeOut(60*time.Second, 60*time.Second, 2*time.Minute)
}

func DefaultTransportWithProxy(proxy func(*http.Request) (*url.URL, error)) *http.Transport {
	return TransportWithProxyAndSkipVerifyAndTimeOut(proxy, 60*time.Second, 60*time.Second, 2*time.Minute)
}

func DefaultTransportWithTLS(tlsClientConfig *tls.Config) *http.Transport {
	return TransportWithTimeOut(tlsClientConfig, 60*time.Second, 60*time.Second, 2*time.Minute)
}

func DefaultTransportWithTLSAndProxy(proxy func(*http.Request) (*url.URL, error), tlsClientConfig *tls.Config) *http.Transport {
	return TransportWithProxyAndTimeOut(proxy, tlsClientConfig, 60*time.Second, 60*time.Second, 2*time.Minute)
}

func NewClient() *Client {
	transport := DefaultTransport()
	return NewClientWithTransport(transport)
}

func NewClientWithHeader(header http.Header) *Client {
	transport := DefaultTransport()
	return NewClientWhthTransportAndHeader(transport, header)
}

func NewClientWithProxy(proxy func(*http.Request) (*url.URL, error)) *Client {
	transport := DefaultTransportWithProxy(proxy)
	return NewClientWithTransport(transport)
}

func NewClientWithProxyAndHeader(proxy func(*http.Request) (*url.URL, error), header http.Header) *Client {
	transport := DefaultTransportWithProxy(proxy)
	return NewClientWhthTransportAndHeader(transport, header)
}

func NewClientWithTLS(tlsClientConfig *tls.Config) *Client {
	transport := DefaultTransportWithTLS(tlsClientConfig)
	return NewClientWithTransport(transport)
}

func NewClientWhthTLSAndHeader(tlsClientConfig *tls.Config, header http.Header) *Client {
	transport := DefaultTransportWithTLS(tlsClientConfig)
	return NewClientWhthTransportAndHeader(transport, header)
}

func NewClientWithTLSAndProxy(tlsClientConfig *tls.Config, proxy func(*http.Request) (*url.URL, error)) *Client {
	transport := DefaultTransportWithTLSAndProxy(proxy, tlsClientConfig)
	return NewClientWithTransport(transport)
}

func NewClientWhthTLSAndProxyAndHeader(tlsClientConfig *tls.Config, proxy func(*http.Request) (*url.URL, error), header http.Header) *Client {
	transport := DefaultTransportWithTLSAndProxy(proxy, tlsClientConfig)
	return NewClientWhthTransportAndHeader(transport, header)
}

func (c *Client) Gzip() *Client {
	c.useGzip = true
	return c
}

func (c *Client) GzipDisable() *Client {
	c.useGzip = false
	return c
}

func (c *Client) doMethod(method, rawurl string, header http.Header, body io.Reader) (*http.Response, error) {
	if header == nil {
		header = http.Header{}
	}

	if c.header != nil && len(c.header) > 0 {
		for k, v := range c.header {
			if _, found := header[k]; !found {
				header[k] = v
			}
		}
	}

	var ir io.Reader
	bu := new(bytes.Buffer)
	if _, err := bu.ReadFrom(body); err != nil {
		return nil, err
	}

	if c.useGzip {
		header.Add("Content-Encoding", "gzip")

		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write(bu.Bytes())
		gzw.Close()
		ir = bufw
	} else {
		ir = bu
	}

	req, err := http.NewRequest(method, rawurl, ir)
	if err != nil {
		log.Error("NewRequest error:", err)
		return nil, err
	}

	if len(header) > 0 {
		req.Header = header
	}

	resp, respErr := c.Do(req)
	return resp, respErr
}

func (c *Client) get(rawurl string, header http.Header, data io.Reader) (*http.Response, error) {
	return c.doMethod("GET", rawurl, header, data)
}

func (c *Client) Get(rawurl string) (*http.Response, error) {
	return c.get(rawurl, nil, nil)
}

func (c *Client) GetWithHeader(rawurl string, header http.Header) (*http.Response, error) {
	return c.get(rawurl, header, nil)
}

func (c *Client) GetFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.get(rawurl, header, strings.NewReader(data.Encode()))
}

func (c *Client) GetForm(rawurl string, data url.Values) (*http.Response, error) {
	return c.GetFormWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) GetJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/json")

	return c.get(rawurl, header, bytes.NewReader(data))
}

func (c *Client) GetJson(rawurl string, data []byte) (*http.Response, error) {
	return c.GetJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) GetFormJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Content-Type", "application/json")

	return c.get(rawurl, header, bytes.NewReader(data))
}

func (c *Client) GetFormJson(rawurl string, data []byte) (*http.Response, error) {
	return c.GetFormJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PostWithHeader(rawurl string, header http.Header, data io.Reader) (*http.Response, error) {
	return c.doMethod("POST", rawurl, header, data)
}

func (c *Client) Post(rawurl string, data io.Reader) (*http.Response, error) {
	return c.PostWithHeader(rawurl, nil, data)
}

func (c *Client) PostFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.PostWithHeader(rawurl, header, strings.NewReader(data.Encode()))
}

func (c *Client) PostForm(rawurl string, data url.Values) (*http.Response, error) {
	return c.PostFormWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PostJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/json")

	return c.PostWithHeader(rawurl, header, bytes.NewReader(data))
}

func (c *Client) PostJson(rawurl string, data []byte) (*http.Response, error) {
	return c.PostJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PostFormJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Content-Type", "application/json")

	return c.PostWithHeader(rawurl, header, bytes.NewReader(data))
}

func (c *Client) PostFormJson(rawurl string, data []byte) (*http.Response, error) {
	return c.PostFormJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PostMultipartWithHeader(rawurl string, header http.Header, data map[string][]string) (resp *http.Response, err error) {
	body, ct, err := multipart.Open(data)
	if err != nil {
		return
	}

	header.Add("Content-Type", ct)

	return c.PostWithHeader(rawurl, header, body)
}

func (c *Client) PostMultipart(rawurl string, data map[string][]string) (resp *http.Response, err error) {
	return c.PostMultipartWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PutWithHeader(rawurl string, header http.Header, data io.Reader) (*http.Response, error) {
	return c.doMethod("PUT", rawurl, header, data)
}

func (c *Client) Put(rawurl string, data io.Reader) (*http.Response, error) {
	return c.PutWithHeader(rawurl, nil, data)
}

func (c *Client) PutFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.PutWithHeader(rawurl, header, strings.NewReader(data.Encode()))
}

func (c *Client) PutForm(rawurl string, data url.Values) (*http.Response, error) {
	return c.PutFormWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PutJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/json")
	return c.PutWithHeader(rawurl, header, bytes.NewReader(data))
}

func (c *Client) PutJson(rawurl string, data []byte) (*http.Response, error) {
	return c.PutJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) PutFormJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Content-Type", "application/json")
	return c.PutWithHeader(rawurl, header, bytes.NewReader(data))
}

func (c *Client) PutFormJson(rawurl string, data []byte) (*http.Response, error) {
	return c.PutFormJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) delete(rawurl string, header http.Header, data io.Reader) (*http.Response, error) {
	return c.doMethod("DELETE", rawurl, header, data)
}

func (c *Client) Delete(rawurl string) (*http.Response, error) {
	return c.delete(rawurl, nil, nil)
}

func (c *Client) DeleteWithHeader(rawurl string, header http.Header) (*http.Response, error) {
	return c.delete(rawurl, header, nil)
}

func (c *Client) DeleteFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c.delete(rawurl, header, strings.NewReader(data.Encode()))
}

func (c *Client) DeleteForm(rawurl string, data url.Values) (*http.Response, error) {
	return c.DeleteFormWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) DeleteJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/json")
	return c.delete(rawurl, header, bytes.NewReader(data))
}

func (c *Client) DeleteJson(rawurl string, data []byte) (*http.Response, error) {
	return c.DeleteJsonWithHeader(rawurl, http.Header{}, data)
}

func (c *Client) DeleteFormJsonWithHeader(rawurl string, header http.Header, data []byte) (*http.Response, error) {
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Content-Type", "application/json")
	return c.delete(rawurl, header, bytes.NewReader(data))
}

func (c *Client) DeleteFormJson(rawurl string, data []byte) (*http.Response, error) {
	return c.DeleteFormJsonWithHeader(rawurl, http.Header{}, data)
}
