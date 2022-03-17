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

const colorReset = "\033[0m"
const colorYellow = "\033[33m"
const colorRed = "\033[31m"

func crash(message string, err error) {
	fmt.Print(string(colorRed) + "[ERROR] " + message + string(colorReset) + "\n")
	panic(err)
}

func warning(message string) {
	fmt.Print(string(colorYellow) + "[WARNING] " + message + string(colorReset) + "\n")
}

func info(message string) {
	fmt.Print("[-] " + message + "\n")
}

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

func showDomains(ua string, chainMode bool) {
	allDomains := getDomains(ua)
	if !chainMode {
		fmt.Println(">> " + ua)
		if len(allDomains) == 0 {
			fmt.Println("|__ NOT FOUND")
		}
		for _, domain := range allDomains {
			fmt.Println("|__ " + domain)
		}
		fmt.Println("")
	} else {
		if len(allDomains) == 0 {
			warning("NOT FOUND")
		}
		for _, domain := range allDomains {
			if domain == "error getting results" {
				var err error
				crash("Server-side error on builtwith.com: error getting results", err)
			}
			fmt.Println(domain)
		}
	}

}

func start(url string, chainMode bool) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	if !chainMode {
		info("Analyzing url: " + url)
	}
	uaResult, resultTagManager := getGoogleTagManager(url)
	if len(resultTagManager) > 0 {
		var visited = []string{}
		var allUAs []string
		if !uaResult {
			urlGoogleTagManager := resultTagManager[0]
			if !chainMode {
				info("URL with UA: " + urlGoogleTagManager)
			}
			allUAs = getUA(urlGoogleTagManager)
		} else {
			if !chainMode {
				info("Found UA directly")
			}
			allUAs = resultTagManager
		}
		if !chainMode {
			info("Obtaining information from builtwith and hackertarget\n")
		}
		for _, ua := range allUAs {
			baseUA := strings.Join(strings.Split(ua, "-")[0:2], "-")
			if !contains(visited, baseUA) {
				visited = append(visited, baseUA)
				showDomains(baseUA, chainMode)
			}
		}
		if !chainMode {
			info("Done!")
		}
	} else {
		warning("Tagmanager URL not found")
		//Now, the program exits
	}
}

//needs to be a global variable
var chainMode bool

func main() {
	var url string

	flag.StringVar(&url, "u", "", "URL to extract Google Analytics ID")
	flag.StringVar(&url, "url", "", "URL to extract Google Analytics ID")
	flag.BoolVar(&chainMode, "ch", false, "In \"chain-mode\" we only output the important information. No decorations.")
	flag.BoolVar(&chainMode, "chain-mode", false, "In \"chain-mode\" we only output the important information. No decorations.")

	const usage = `Usage: ./analyticsrelationships -u URL [--chain-mode]
  -u, --url string
      URL to extract Google Analytics ID
  -ch, --chain-mode
      In "chain-mode" we only output the important information. No decorations.
      Default: false
`

	//https://www.antoniojgutierrez.com/posts/2021-05-14-short-and-long-options-in-go-flags-pkg/
	flag.Usage = func() { fmt.Print(usage) }
	//parse CLI arguments
	flag.Parse()
	if !chainMode {
		//display banner
		banner()
	}
	if url != "" {
		//start main processing
		start(url, chainMode)
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
					start(scanner.Text(), chainMode)
				}
			}

		} //else { //stdin is from a terminal }

	}
}
