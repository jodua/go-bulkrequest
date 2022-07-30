package client

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	jsonparser "github.com/jodua/go-bulkrequest/JSONParser"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	headers    map[string]string
	userAgent  string
	parser     *jsonparser.JSONParser
}

type ClientBuilder struct {
	httpClient *http.Client
	cookies    map[string]string
	timeout    time.Duration
	transport  *http.Transport
	headers    map[string]string
	baseUrl    string
	userAgent  string
	proxy      string
	parser     *jsonparser.JSONParser
}

func NewClient() *ClientBuilder {
	return &ClientBuilder{
		httpClient: &http.Client{},
		transport:  &http.Transport{},
		timeout:    10 * time.Second,
	}
}

func (c *ClientBuilder) SetBaseUrl(baseUrl string) *ClientBuilder {
	c.baseUrl = baseUrl
	return c
}

func (c *ClientBuilder) SetProxy(proxyUrl string) *ClientBuilder {
	c.proxy = proxyUrl
	return c
}

func (c *ClientBuilder) SetParser(parser *jsonparser.JSONParser) *ClientBuilder {
	c.parser = parser
	return c
}

func (c *ClientBuilder) SetTimeout(timeout time.Duration) *ClientBuilder {
	c.timeout = timeout
	return c
}

func (c *ClientBuilder) SetCookies(cookies map[string]string) *ClientBuilder {
	c.cookies = cookies
	return c
}

func (c *ClientBuilder) SetHeaders(headers map[string]string) *ClientBuilder {
	c.headers = headers
	return c
}

func (c *ClientBuilder) SetUserAgent(userAgent string) *ClientBuilder {
	c.userAgent = userAgent
	return c
}

func (c *ClientBuilder) Build() *Client {
	// Create new client
	client := &Client{
		httpClient: &http.Client{},
		baseUrl:    c.baseUrl,
		headers:    c.headers,
	}

	// Set proxy
	if c.proxy != "" {
		proxyUrl, err := url.Parse(c.proxy)
		if err != nil {
			log.Fatalf("Invalid proxy url: %s", err)
		}
		c.transport.Proxy = http.ProxyURL(proxyUrl)
	}

	// Set timeout
	client.httpClient.Timeout = c.timeout

	// Set transport
	client.httpClient.Transport = c.transport

	// Set cookies
	// Create new cookiejar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookiejar: %v", err)
	}
	// Convert cookies map to []http.Cookie
	var cookies []*http.Cookie
	for name, value := range c.cookies {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}
	// Parse url
	cookiesUrl, err := url.Parse(c.baseUrl)
	if err != nil {
		log.Fatalf("Error parsing url: %v", err)
	}
	// Set cookies to cookiejar
	cookieJar.SetCookies(cookiesUrl, cookies)
	// Set cookiejar to client
	client.httpClient.Jar = cookieJar
	// Set parser
	client.parser = c.parser

	return client
}

func (c *Client) Fetch(url string) (any, error) {
	// Create request
	request, err := http.NewRequest("GET", c.baseUrl+url, nil)
	if err != nil {
		return nil, err
	}
	// Set headers
	for name, value := range c.headers {
		request.Header.Set(name, value)
	}
	// Set user agent
	if c.userAgent != "" {
		request.Header.Set("User-Agent", c.userAgent)
	}
	// Fetch
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	// Read body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// Close body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(response.Body)

	// Parse and return
	return c.parser.Parse(body, url)
}
