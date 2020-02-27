package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// type requestPayloadStruct struct {
// 	ProxyCondition string `json:"proxy_condition"`
// }

// Log the typeform payload and redirect url
// func logRequestPayload(requestionPayload requestPayloadStruct, proxyUrl string) {
// 	log.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyUrl)
// }

// Get the url for a given proxy condition
// func getProxyUrl(proxyConditionRaw string) string {
// 	proxyCondition := strings.ToUpper(proxyConditionRaw)
// 	URL := os.Getenv(proxyCondition)
// 	if URL == "" {
// 		return os.Getenv("URL")
// 	}
// 	return URL
// }

// Get a json decoder for a given requests body
// func requestBodyDecoder(request *http.Request) *json.Decoder {
// 	// Read body to buffer
// 	body, err := ioutil.ReadAll(request.Body)
// 	if err != nil {
// 		log.Printf("Error reading body: %v\n", err)
// 	}

// 	// Because go lang is a pain in the ass if you read the body then any susequent calls
// 	// are unable to read the body again....
// 	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

// 	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
// }

// Parse the requests body
// func parseRequestBody(request *http.Request) requestPayloadStruct {
// 	decoder := requestBodyDecoder(request)

// 	var requestPayload requestPayloadStruct
// 	err := decoder.Decode(&requestPayload)

// 	if err != nil {
// 		log.Println("body error", err.Error())
// 	}

// 	return requestPayload
// }

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	// requestPayload := parseRequestBody(req)
	// log.Println("proxy", requestPayload.ProxyCondition)
	log.Println("from ", req.RemoteAddr)
	url := os.Getenv("URL")

	// logRequestPayload(requestPayload, url)

	serveReverseProxy(url, res, req)
}

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getListenAddress() string {
	port := getEnv("PORT", "3000")
	return ":" + port
}

func logSetup() {
	url := os.Getenv("URL")
	log.Printf("Server will run on: %s\n", getListenAddress())
	log.Printf("Redirecting to a url: %s\n", url)
}

func main() {
	logSetup()

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		log.Panicln("server error", err.Error())
	}
}
