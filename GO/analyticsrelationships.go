package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func banner() {
	data := `
██╗   ██╗ █████╗       ██╗██████╗                        
██║   ██║██╔══██╗      ██║██╔══██╗                       
██║   ██║███████║█████╗██║██║  ██║                       
██║   ██║██╔══██║╚════╝██║██║  ██║                       
╚██████╔╝██║  ██║      ██║██████╔╝                       
 ╚═════╝ ╚═╝  ╚═╝      ╚═╝╚═════╝                        
                                                         
██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔════╝
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║███████╗
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║╚════██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║███████║
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝

> Get related domains / subdomains by looking at Google Analytics IDs
> GO Version
> By @JosueEncinar

`
	println(data)
}

func getURLResponse(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func getGoogleTagManager(targetURL string) string {
	url := ""
	response := getURLResponse(targetURL)
	if response != "" {
		pattern := regexp.MustCompile("www\\.googletagmanager\\.com/ns\\.html\\?id=[A-Z0-9\\-]+")
		data := pattern.FindStringSubmatch(response)
		if len(data) > 0 {
			url = "https://" + strings.Replace(data[0], "ns.html", "gtm.js", -1)
		} else {
			pattern = regexp.MustCompile("GTM-[A-Z0-9]+")
			data = pattern.FindStringSubmatch(response)
			if len(data) > 0 {
				url = "https://www.googletagmanager.com/gtm.js?id=" + data[0]
			}
		}
	}
	return url
}

func getUA(url string) [][]string {
	pattern := regexp.MustCompile("UA-[0-9]+-[0-9]+")
	response := getURLResponse(url)
	var result = [][]string{}
	if response != "" {
		result = pattern.FindAllStringSubmatch(response, -1)
	} else {
		result = nil
	}
	return result
}

func cleanRelationShips(domains [][]string) []string {
	var allDomains []string
	for _, domain := range domains {
		allDomains = append(allDomains, strings.Replace(domain[0], "/relationships/", "", -1))
	}
	return allDomains
}

func getDomainsFromBuiltWith(id string) []string {
	pattern := regexp.MustCompile("/relationships/[a-z0-9\\-\\_\\.]+\\.[a-z]+")
	url := "https://builtwith.com/relationships/tag/" + id
	response := getURLResponse(url)
	var allDomains []string = nil
	if response != "" {
		allDomains = cleanRelationShips(pattern.FindAllStringSubmatch(response, -1))
	}
	return allDomains
}

func getDomainsFromHackerTarget(id string) []string {
	url := "https://api.hackertarget.com/analyticslookup/?q=" + id
	response := getURLResponse(url)
	var allDomains []string = nil
	if response != "" {
		allDomains = strings.Split(response, "\n")
	}
	return allDomains
}

func getDomains(id string) []string {
	var allDomains []string = getDomainsFromBuiltWith(id)
	domains2 := getDomainsFromHackerTarget(id)
	if domains2 != nil {
		for _, domain := range domains2 {
			if !contains(allDomains, domain) {
				allDomains = append(allDomains, domain)
			}
		}
	}
	return allDomains
}

func contains(data []string, value string) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func main() {
	url := flag.String("url", "", "URL to extract Google Analytics ID")
	flag.Parse()
	banner()
	if *url == "" {
		println("Usage: ./analyticsrelationships --url https://www.example.com")
		return
	}
	if !strings.HasPrefix(*url, "http") {
		*url = "https://" + *url
	}
	println("[+] Analyzing url: " + *url)
	urlGoogleTagManager := getGoogleTagManager(*url)
	if urlGoogleTagManager != "" {
		println("[+] URL with UA: " + urlGoogleTagManager)
		println("[+] Obtaining information from builtwith hackertarget\n")
		var visited = []string{}
		for _, ua := range getUA(urlGoogleTagManager) {
			baseUA := strings.Join(strings.Split(ua[0], "-")[0:2], "-")
			if !contains(visited, baseUA) {
				visited = append(visited, baseUA)
				fmt.Println(">> " + baseUA)
				allDomains := getDomains(baseUA)
				if len(allDomains) == 0 {
					fmt.Println("|__ NOT FOUND")
				}
				for _, domain := range allDomains {
					fmt.Println("|__ " + domain)
				}
				fmt.Println("")
			}
		}
		println("\n[+] Done!")
	} else {
		println("[-] Tagmanager URL not fount")
	}
}
