package main

import (
	"flag"
	"fmt"
	"github.com/joshbetz/config"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	cfg          *config.Config
	ipv4Endpoint string
	ipv6Endpoint string

	logger *log.Logger
)

// Preparations
func init() {
	// Logging
	logOpts := log.Ldate | log.Ltime | log.LUTC | log.Lshortfile
	logger = log.New(os.Stdout, "", logOpts)

	logger.Println("PINGING dyndnsd ENDPOINT")

	// Configuration
	var configPath string
	flag.StringVar(&configPath, "config", "", "Config file path")
	flag.Parse()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config path does not exist. Aborting!\n")
		flag.Usage()
		os.Exit(3)
	}

	cfg = config.New(configPath)
	cfg.Get("ipv4_endpoint", &ipv4Endpoint)
	cfg.Get("ipv6_endpoint", &ipv6Endpoint)
}

// Do the update!
func main() {
	// Base URL
	endpoint, err := buildEndpointUrl()
	if err != nil {
		logger.Println("Couldn't build endpoint URL")
		logger.Printf("%s", err)
		return
	}

	// IPv4 is required
	if ipv4, err := getUrl(ipv4Endpoint); err == nil {
		if ipv4Valid := net.ParseIP(ipv4); ipv4Valid == nil {
			logger.Println("Invalid IPv4 address returned by endpoint")
			logger.Printf("%s", err)
			return
		} else {
			query := endpoint.Query()
			query.Set("myip", ipv4Valid.String())
			endpoint.RawQuery = query.Encode()
		}
	} else {
		logger.Println("Couldn't retrieve IPv4 address")
		logger.Printf("%s", err)
		return
	}

	// IPv6 is optional
	// Leave empty to skip
	if len(ipv6Endpoint) > 0 {
		if ipv6, err := getUrl(ipv6Endpoint); err == nil {
			if ipv6Valid := net.ParseIP(ipv6); ipv6Valid == nil {
				logger.Println("Invalid IPv6 address returned by endpoint")
				logger.Printf("%s", err)
			} else {
				var ipv6String string
				var usePrefix bool
				cfg.Get("ipv6_use_prefix", &usePrefix)
				ipMask := fmt.Sprintf("%s/%d", ipv6Valid.String(), 64)

				if ipv6, network, err := net.ParseCIDR(ipMask); usePrefix && err == nil {
					ipv6String = network.String()
					ipv6String = strings.Replace(ipv6String, "/64", "1", 1)
				} else {
					ipv6String = ipv6.String()
				}

				query := endpoint.Query()
				query.Set("myip6", ipv6String)
				endpoint.RawQuery = query.Encode()
			}
		} else {
			logger.Println("Couldn't retrieve IPv6 address")
			logger.Printf("%s", err)
		}
	}

	// Send the update
	dyndns, err := getUrl(endpoint.String())
	if err != nil {
		logger.Println("Couldn't update dyndnsd endpoint")
		logger.Printf("%s", err)
		return
	}

	logger.Println("SUCCESS! Ping sent.")
	logger.Printf("%s", dyndns)
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
