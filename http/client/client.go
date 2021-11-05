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

func SkipVerifyAndTimeOutTransport(insecureSkipVerify bool, headerTimeout, connectTimeout, readWriteTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  TimeoutDialer(connectTimeout, readWriteTimeout),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
		ResponseHeaderTimeout: headerTimeout,
	}
}

func NewClient() *Client {
	c := &Client{Client: &http.Client{}}
	c.Transport = SkipVerifyAndTimeOutTransport(true, 60*time.Second, 60*time.Second, 2*time.Minute)
	return c
}

func NewClientWithHeader(header http.Header) *Client {
	c := &Client{
		Client: &http.Client{},
		header: header,
	}

	c.Transport = SkipVerifyAndTimeOutTransport(true, 60*time.Second, 60*time.Second, 2*time.Minute)
	return c
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
	req, err := http.NewRequest(method, rawurl, body)
	if err != nil {
		log.Error("NewRequest error:", err)
		return nil, err
	}

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

	if c.useGzip {
		header.Add("Content-Encoding", "gzip")
	}

	if len(header) > 0 {
		req.Header = header
	}

	resp, respErr := c.Do(req)
	return resp, respErr
}

func (c *Client) get(rawurl string, header http.Header, body io.Reader) (*http.Response, error) {
	return c.doMethod("GET", rawurl, header, body)
}

func (c *Client) Get(rawurl string) (*http.Response, error) {
	return c.get(rawurl, nil, nil)
}

func (c *Client) GetWithHeader(rawurl string, header http.Header) (*http.Response, error) {
	return c.get(rawurl, header, nil)
}

func (c *Client) GetForm(rawurl string, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.get(rawurl, header, ir)
}

func (c *Client) GetFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.get(rawurl, header, ir)
}

func (c *Client) PostWithHeader(rawurl string, header http.Header, body io.Reader) (*http.Response, error) {
	return c.doMethod("POST", rawurl, header, body)
}

func (c *Client) Post(rawurl string, body io.Reader) (*http.Response, error) {
	return c.PostWithHeader(rawurl, nil, body)
}

func (c *Client) PostForm(rawurl string, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.PostWithHeader(rawurl, header, ir)
}

func (c *Client) PostFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.PostWithHeader(rawurl, header, ir)
}

func (c *Client) PostMultipart(rawurl string, data map[string][]string) (resp *http.Response, err error) {
	body, ct, err := multipart.Open(data)
	if err != nil {
		return
	}

	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write(body.Bytes())
		gzw.Close()
		body = bufw
	}

	header := http.Header{}
	header.Set("Content-Type", ct)

	return c.PostWithHeader(rawurl, header, body)
}

func (c *Client) PostMultipartWithHeader(rawurl string, header http.Header, data map[string][]string) (resp *http.Response, err error) {
	body, ct, err := multipart.Open(data)
	if err != nil {
		return
	}

	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write(body.Bytes())
		gzw.Close()
		body = bufw
	}

	header.Add("Content-Type", ct)

	return c.PostWithHeader(rawurl, header, body)
}

func (c *Client) PutWithHeader(rawurl string, header http.Header, body io.Reader) (*http.Response, error) {
	return c.doMethod("PUT", rawurl, header, body)
}

func (c *Client) Put(rawurl string, body io.Reader) (*http.Response, error) {
	return c.PutWithHeader(rawurl, nil, body)
}

func (c *Client) PutForm(rawurl string, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.PutWithHeader(rawurl, header, ir)
}

func (c *Client) PutFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.PutWithHeader(rawurl, header, ir)
}

func (c *Client) delete(rawurl string, header http.Header, body io.Reader) (*http.Response, error) {
	return c.doMethod("DELETE", rawurl, header, body)
}

func (c *Client) Delete(rawurl string) (*http.Response, error) {
	return c.delete(rawurl, nil, nil)
}

func (c *Client) DeleteWithHeader(rawurl string, header http.Header) (*http.Response, error) {
	return c.delete(rawurl, header, nil)
}

func (c *Client) DeleteForm(rawurl string, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.delete(rawurl, header, ir)
}

func (c *Client) DeleteFormWithHeader(rawurl string, header http.Header, data url.Values) (*http.Response, error) {
	var ir io.Reader
	if c.useGzip {
		bufw := new(bytes.Buffer)
		gzw := gzip.NewWriter(bufw)
		gzw.Write([]byte(data.Encode()))
		gzw.Close()
		ir = bufw
	} else {
		ir = strings.NewReader(data.Encode())
	}

	header.Add("Content-Type", "application/x-www-form-urlencoded")

	return c.delete(rawurl, header, ir)
}
