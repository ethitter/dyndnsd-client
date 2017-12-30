package main

import (
	"fmt"
	"github.com/joshbetz/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	cfg          *config.Config
	ipv4Endpoint string
	ipv6Endpoint string
)

// Parse config
func init() {
	cfg = config.New("config.json")

	cfg.Get("ipv4_endpoint", &ipv4Endpoint)
	cfg.Get("ipv6_endpoint", &ipv6Endpoint)
}

// Do the update!
func main() {
	// Base URL
	endpoint, err := buildEndpointUrl()
	if err != nil {
		return
	}

	// IPv4 is required
	ipv4, err := getUrl(ipv4Endpoint)
	if err == nil {
		query := endpoint.Query()
		query.Set("myip", ipv4)
		endpoint.RawQuery = query.Encode()
	} else {
		return
	}

	fmt.Println(endpoint)

	// IPv6 is optional
	ipv6, err := getUrl(ipv6Endpoint)
	if err == nil {
		query := endpoint.Query()
		query.Set("myip6", ipv6)
		endpoint.RawQuery = query.Encode()
	}

	// Send the update
	dyndns, err := getUrl(endpoint.String())
	if err != nil {
		return
	}

	// TODO: better formatting
	fmt.Println(dyndns)
}

// Build endpoint URL from configuration
func buildEndpointUrl() (*url.URL, error) {
	var username string
	var password string
	var protocol string
	var host string
	var port int
	var path string
	var hostname string

	cfg.Get("username", &username)
	cfg.Get("password", &password)
	cfg.Get("protocol", &protocol)
	cfg.Get("host", &host)
	cfg.Get("port", &port)
	cfg.Get("path", &path)
	cfg.Get("dns_hostname", &hostname)

	daemonUrl, err := url.Parse("")
	if err != nil {
		return nil, err
	}

	daemonUrl.Scheme = protocol
	daemonUrl.Host = fmt.Sprintf("%s:%d", host, port)
	daemonUrl.Path = path

	userInfo := url.UserPassword(username, password)
	daemonUrl.User = userInfo

	query := daemonUrl.Query()
	query.Set("hostname", hostname)
	daemonUrl.RawQuery = query.Encode()

	return daemonUrl, nil
}

// Retrieve given URL
func getUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	respBodyString := string(respBody)
	return respBodyString, nil
}
