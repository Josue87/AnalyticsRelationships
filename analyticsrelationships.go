package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
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

`
	data += "\033[32m> \033[0mGet related domains / subdomains by looking at Google Analytics IDs\n"
	data += "\033[32m> \033[0mGO Version\n"
	data += "\033[32m> \033[0mBy @JosueEncinar\n"

	println(data)
}

func getURLResponse(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3}
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

func getGoogleTagManager(targetURL string) (bool, []string) {
	var resultTagManager []string
	response := getURLResponse(targetURL)
	if response != "" {
		pattern := regexp.MustCompile(`www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+`)
		data := pattern.FindStringSubmatch(response)
		if len(data) > 0 {
			resultTagManager = append(resultTagManager, "https://"+strings.Replace(data[0], "ns.html", "gtm.js", -1))
		} else {
			pattern = regexp.MustCompile("GTM-[A-Z0-9]+")
			data = pattern.FindStringSubmatch(response)
			if len(data) > 0 {
				resultTagManager = append(resultTagManager, "https://www.googletagmanager.com/gtm.js?id="+data[0])
			} else {
				pattern = regexp.MustCompile(`UA-\d+-\d+`)
				aux := pattern.FindAllStringSubmatch(response, -1)
				var result []string
				for _, r := range aux {
					result = append(result, r[0])
				}
				return true, result
			}
		}
	}
	return false, resultTagManager
}

func getUA(url string) []string {
	pattern := regexp.MustCompile("UA-[0-9]+-[0-9]+")
	response := getURLResponse(url)
	var result []string
	if response != "" {
		aux := pattern.FindAllStringSubmatch(response, -1)
		for _, r := range aux {
			result = append(result, r[0])
		}
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
	pattern := regexp.MustCompile(`/relationships/[a-z0-9\-\_\.]+\.[a-z]+`)
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
	if response != "" && !strings.Contains(response, "API count exceeded") {
		allDomains = strings.Split(response, "\n")
	}
	return allDomains
}

func getDomains(id string) []string {
	var allDomains []string = getDomainsFromBuiltWith(id)
	domains2 := getDomainsFromHackerTarget(id)
	if len(domains2) != 0 {
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

func showDomains(ua string) {
	fmt.Println(">> " + ua)
	allDomains := getDomains(ua)
	if len(allDomains) == 0 {
		fmt.Println("|__ NOT FOUND")
	}
	for _, domain := range allDomains {
		fmt.Println("|__ " + domain)
	}
	fmt.Println("")
}

func start(url string) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	println("[+] Analyzing url: " + url)
	uaResult, resultTagManager := getGoogleTagManager(url)
	if len(resultTagManager) > 0 {
		var visited = []string{}
		var allUAs []string
		if !uaResult {
			urlGoogleTagManager := resultTagManager[0]
			println("[+] URL with UA: " + urlGoogleTagManager)
			allUAs = getUA(urlGoogleTagManager)
		} else {
			println("[+] Found UA directly")
			allUAs = resultTagManager
		}
		println("[+] Obtaining information from builtwith and hackertarget\n")
		for _, ua := range allUAs {
			baseUA := strings.Join(strings.Split(ua, "-")[0:2], "-")
			if !contains(visited, baseUA) {
				visited = append(visited, baseUA)
				showDomains(baseUA)
			}
		}
		println("\n[+] Done!")
	} else {
		println("[-] Tagmanager URL not found")
	}
}

func main() {
	url := flag.String("url", "", "URL to extract Google Analytics ID")
	flag.Parse()
	banner()
	if url != "" {
		//start main processing
		start(url)
	} else {
		//read from standard input (stdin)

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			//data is being piped to stdin
			//read stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					crash("bufio couldn't read stdin correctly.", err)
				} else {
					start(scanner.Text())
				}
			}

		} //else { //stdin is from a terminal }

	}
}
