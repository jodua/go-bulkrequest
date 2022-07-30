package main

import (
	"bufio"
	"github.com/jodua/go-bulkrequest/JSONParser/schemas"
	"github.com/jodua/go-bulkrequest/bulkrequest"
	"github.com/jodua/go-bulkrequest/datatypes"
	"log"
	"os"
	"time"
)

func main() {
	TestJSONPlaceholder()
	TestJSONPlaceholderProxy()
}

func TestJSONPlaceholder() {
	var parser = schemas.JSONPlaceholderTodoParser
	var baseUrl = "https://jsonplaceholder.typicode.com/todos/"
	var urls = []string{"1", "2", "3", "4", "5"}
	var UAConfig = datatypes.UserAgentConfig{
		UserAgentList: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
			"go-bulkrequest",
		},
	}

	var delayConfig = datatypes.DelayConfig{
		DelayMin: 500,
		DelayMax: 1000,
	}

	var bulkRequest = bulkrequest.NewBulkRequest().
		SetBaseUrl(baseUrl).
		SetUrls(urls).
		SetParser(&parser).
		SetTimeout(time.Second * 10).
		SetUserAgentConfig(&UAConfig).
		SetDelayConfig(&delayConfig).
		Build()

	fetch, err := bulkRequest.Fetch()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("%v", fetch)
}

func TestJSONPlaceholderProxy() {
	proxyList, _ := LoadFile("examples/proxies.txt")
	urlList, _ := LoadFile("examples/urls.txt")

	var parser = schemas.JSONPlaceholderTodoParser
	var baseUrl = "https://jsonplaceholder.typicode.com/todos/"
	var proxyConfig = datatypes.ProxyConfig{
		ProxyList:        proxyList,
		RequestsPerProxy: 5,
	}
	var UAConfig = datatypes.UserAgentConfig{
		UserAgentList: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
			"go-bulkrequest",
		},
	}

	var delayConfig = datatypes.DelayConfig{
		DelayMin: 500,
		DelayMax: 1000,
	}

	var bulkRequest = bulkrequest.NewBulkRequest().
		SetBaseUrl(baseUrl).
		SetUrls(urlList).
		SetParser(&parser).
		SetTimeout(time.Second * 10).
		SetProxyConfig(&proxyConfig).
		SetUserAgentConfig(&UAConfig).
		SetDelayConfig(&delayConfig).
		Build()

	fetch, err := bulkRequest.Fetch()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("%v", fetch)
}

func LoadFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
