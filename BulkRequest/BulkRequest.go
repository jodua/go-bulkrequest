package bulkrequest

import (
	utils "github.com/jodua/go-bulkrequest/Utils"
	"log"
	"sync"
	"time"

	jsonparser "github.com/jodua/go-bulkrequest/JSONParser"
	"github.com/jodua/go-bulkrequest/client"
	"github.com/jodua/go-bulkrequest/datatypes"
)

type BulkRequest struct {
	baseUrl         string
	urls            []string
	proxyConfig     *datatypes.ProxyConfig
	userAgentConfig *datatypes.UserAgentConfig
	delayConfig     *datatypes.DelayConfig
	parser          *jsonparser.JSONParser
	timeout         time.Duration
	cookies         map[string]string
	headers         map[string]string
}

func (b *BulkRequest) Fetch() ([]any, error) {

	// Check if proxyConfig is set
	// If proxyConfig is not set, use direct connection
	if b.proxyConfig == nil {
		// Create result array
		var results []any
		// Warn user about missing proxy config
		log.Println("Warning: No proxy config provided. Using direct connection.")
		// Create requestClient
		requestClient := client.NewClient().
			SetBaseUrl(b.baseUrl).
			SetTimeout(b.timeout).
			SetCookies(b.cookies).
			SetHeaders(b.headers).
			SetParser(b.parser).
			SetUserAgent(b.userAgentConfig.GetRandomUserAgent()).
			Build()

		// Iterate over urls
		log.Printf("Fetching %d urls", len(b.urls))
		for _, url := range b.urls {

			// Fetch
			result, err := requestClient.Fetch(url)
			if err != nil {
				return nil, err
			}

			// Add result to results
			results = append(results, result)

			// Wait between requests
			time.Sleep(b.delayConfig.GetRandomDelay())
		}

		// Return results
		log.Printf("Done fetching %d urls", len(b.urls))
		return results, nil
	} else {
		// Create waitGroup and mutex in order to synchronize goroutines
		var wg sync.WaitGroup
		var mutex = &sync.Mutex{}

		// Create output array to store results
		var results [][]any

		// Shuffle proxy list and urls
		utils.ShuffleList(b.proxyConfig.ProxyList)
		utils.ShuffleList(b.urls)

		// Take n proxies from proxy list
		proxyAmountNeeded := len(b.urls) / b.proxyConfig.RequestsPerProxy
		if len(b.urls)%b.proxyConfig.RequestsPerProxy != 0 {
			proxyAmountNeeded++
		}
		proxyList := utils.Take(b.proxyConfig.ProxyList, 0, proxyAmountNeeded)

		// Iterate over proxies
		for i, proxy := range proxyList {
			// Create requestClient
			requestClient := client.NewClient().
				SetBaseUrl(b.baseUrl).
				SetTimeout(b.timeout).
				SetCookies(b.cookies).
				SetHeaders(b.headers).
				SetParser(b.parser).
				SetProxy(proxy).
				SetUserAgent(b.userAgentConfig.GetRandomUserAgent()).
				Build()

			// Get requestsPerProxy urls
			requestPerProxyUrls := utils.Take(b.urls, i*b.proxyConfig.RequestsPerProxy, b.proxyConfig.RequestsPerProxy)

			wg.Add(1)
			log.Printf("Using proxy %s", proxy)

			// Fetch requests
			go FetchUsingProxy(requestClient, &wg, mutex, proxy, requestPerProxyUrls, b.delayConfig, &results)
		}
		// Wait for all goroutines to finish
		wg.Wait()
		// Return results
		return utils.Flatten(results), nil
	}
}

func FetchUsingProxy(client *client.Client, wg *sync.WaitGroup, mutex *sync.Mutex, proxy string, requestPerProxyUrls []string, delayConfig *datatypes.DelayConfig, results *[][]any) {
	var localResults []any
	for _, url := range requestPerProxyUrls {
		// Fetch
		result, err := client.Fetch(url)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		// Add result to results
		localResults = append(localResults, result)
		// Delay
		time.Sleep(delayConfig.GetRandomDelay())
	}
	// Add results to global results
	mutex.Lock()
	*results = append(*results, localResults)
	mutex.Unlock()
	// Done
	defer log.Printf("Done with proxy %s", proxy)
	defer wg.Done()
}

type BulkRequestBuilder struct {
	bulkRequest *BulkRequest
}

func NewBulkRequest() *BulkRequestBuilder {
	return &BulkRequestBuilder{
		bulkRequest: &BulkRequest{},
	}
}

func (b *BulkRequestBuilder) SetBaseUrl(baseUrl string) *BulkRequestBuilder {
	// Check if url is valid
	if !utils.IsValidUrl(baseUrl) {
		panic("Invalid base url")
	}
	b.bulkRequest.baseUrl = baseUrl
	return b
}

func (b *BulkRequestBuilder) SetUrls(urls []string) *BulkRequestBuilder {
	b.bulkRequest.urls = urls
	return b
}

func (b *BulkRequestBuilder) SetDelayConfig(delayConfig *datatypes.DelayConfig) *BulkRequestBuilder {
	err := delayConfig.Validate()
	if err != nil {
		panic(err)
	}
	b.bulkRequest.delayConfig = delayConfig
	return b
}

func (b *BulkRequestBuilder) SetTimeout(timeout time.Duration) *BulkRequestBuilder {
	b.bulkRequest.timeout = timeout
	return b
}

func (b *BulkRequestBuilder) SetProxyConfig(proxyConfig *datatypes.ProxyConfig) *BulkRequestBuilder {
	err := proxyConfig.Validate()
	if err != nil {
		panic(err)
	}
	b.bulkRequest.proxyConfig = proxyConfig
	return b
}

func (b *BulkRequestBuilder) SetUserAgentConfig(userAgentConfig *datatypes.UserAgentConfig) *BulkRequestBuilder {
	err := userAgentConfig.Validate()
	if err != nil {
		panic(err)
	}
	b.bulkRequest.userAgentConfig = userAgentConfig
	return b
}

func (b *BulkRequestBuilder) AddCookie(name, value string) *BulkRequestBuilder {
	if b.bulkRequest.cookies == nil {
		b.bulkRequest.cookies = make(map[string]string)
	}
	b.bulkRequest.cookies[name] = value
	return b
}

func (b *BulkRequestBuilder) SetParser(parser *jsonparser.JSONParser) *BulkRequestBuilder {
	b.bulkRequest.parser = parser
	return b
}

func (b *BulkRequestBuilder) AddHeader(key, value string) *BulkRequestBuilder {
	if b.bulkRequest.headers == nil {
		b.bulkRequest.headers = make(map[string]string)
	}
	b.bulkRequest.headers[key] = value
	return b
}

func (b *BulkRequestBuilder) Build() *BulkRequest {
	log.Printf("BaseUrl: %s", b.bulkRequest.baseUrl)
	log.Printf("Request amount: %d", len(b.bulkRequest.urls))
	log.Printf("Timeout: %s", b.bulkRequest.timeout)
	log.Printf("Delay between: %d and %d", b.bulkRequest.delayConfig.DelayMin, b.bulkRequest.delayConfig.DelayMax)
	// if proxy is set
	if b.bulkRequest.proxyConfig != nil {
		log.Printf("Proxy amount: %d", len(b.bulkRequest.proxyConfig.ProxyList))
		log.Printf("Requests per proxy: %d", b.bulkRequest.proxyConfig.RequestsPerProxy)
	}
	log.Printf("User agent amount: %d", len(b.bulkRequest.userAgentConfig.UserAgentList))
	log.Printf("Cookies: %v", b.bulkRequest.cookies)
	log.Printf("Headers: %v", b.bulkRequest.headers)
	log.Printf("Parser: %v", b.bulkRequest.parser.Name)
	return b.bulkRequest
}
