import urllib.request
import requests
import re
import argparse
from sys import stderr
from queue import Queue
import urllib3
urllib3.disable_warnings()


def banner():
    stderr.writelines("""
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
> Python version
> By @JosueEncinar

""")

def clean_ua(uas):
    result = []
    if uas is not None:
        for ua in uas:
            analytics_id = ua.split('-')
            analytics_id = "-".join(analytics_id[0:2])
            if analytics_id not in result:
                result.append(analytics_id)
    return result

def get_UA(link):
    pattern = "UA-\d+-\d+"
    try:
        u = urllib.request.urlopen(link)
        data = u.read().decode(errors="ignore")
        match = re.findall(pattern, data)
        unique = set()
        unique = unique.union(match)
        return clean_ua(list(unique)) 
    except Exception as e:
        print(e)
    return None

def get_googletagmanager(url):
    pattern = "(www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+)"
    pattern2 = "GTM-[A-Z0-9]+"
    pattern3 = "UA-\d+-\d+"
    try:
        response = requests.get(url, 
                    headers={
                        'User-agent': 'Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.86 Mobile Safari/537.36'
                        }, 
                    verify=False,
                    timeout=3)

        if response.status_code == 200:
            text = response.text
            match = re.findall(pattern, text)
            if match:
                return True, f"https://{match[0].replace('ns.html', 'gtm.js')}"
            else:
                match = re.findall(pattern2, text)
                if match:
                    return True, f"https://www.googletagmanager.com/gtm.js?id={match[0]}"
                return False, re.findall(pattern3, text)
    except:
        pass
    return False, None

def clean_relationships(domains):
    all_domains = []
    for domain in domains:
        all_domains.append(domain.replace('/relationships/',''))
    return all_domains

def get_domains_from_builtwith(id):
    pattern = "/relationships/[a-z0-9\-\_\.]+\.[a-z]+"
    url = f"https://builtwith.com/relationships/tag/{id}"
    try:
        u = urllib.request.urlopen(url)
        data = u.read().decode(errors="ignore")
        return clean_relationships(re.findall(pattern, data))
    except:
        pass
    return []

def get_domains_from_hackertarget(id):
    url = f"https://api.hackertarget.com/analyticslookup/?q={id}" # Limited requests!
    try:
        response = requests.get(url)
        if response.status_code == 200 and "API count exceeded" not in response.text: 
            return response.text.split("\n")
    except:
        pass
    return []

def get_domains(id):
    all_domains = set()
    all_domains = all_domains.union(get_domains_from_builtwith(id))
    all_domains = all_domains.union(get_domains_from_hackertarget(id))
    return list(all_domains)
    
def show_data(data):
    for ua, domains in data.items():
        print(f">> {ua}")
        if domains:
            for domain in domains:
                print(f"|__ {domain}")
        else:
            print("|__ NOT FOUND")
        print("")

def get_results(uas):
    data = {}
    if uas is not None:
        for ua in uas:
            data[ua] = get_domains(ua) 
    return data

def get_news_uas(uas, visited_uas):
    news = []
    if uas is not None:
        for ua in uas:
            if ua not in visited_uas and ua not in news:
                news.append(ua)
    return news


if __name__ == "__main__":
    banner()
    parser = argparse.ArgumentParser()
    parser.add_argument('-u','--url', help="URL to extract Google Analytics ID",required=True)
    args = parser.parse_args()
    u =  args.url
    if not u.startswith("http"):
        u = "https://" + u
    urls = Queue()
    urls.put(u)
    visited_uas = []
    visited_urls = []
    while not urls.empty():
        ok = False
        uas = None
        url = urls.get()
        if url in visited_urls:
            continue
        visited_urls.append(url)
        stderr.writelines(f"\n[+] Analyzing url: {url}\n")
        tagmanager, data = get_googletagmanager(url)
        if tagmanager and data:
            uas_aux = get_UA(data)
            ok = True
        elif data:
            uas_aux = clean_ua(data)
            ok = True
        uas = get_news_uas(uas_aux, visited_uas)
        visited_uas.extend(uas)
        if uas:
            stderr.writelines(f"[+] URL with UA: {data}\n")
            stderr.writelines("[+] Obtaining information from builtwith and hackertarget\n\n")
            results = get_results(uas)
            show_data(results)
            for values in results.values():
                for v in values:
                    urls.put("https://" + v)
        elif ok:
            stderr.writelines("[!] No news Analytics IDs found...\n")
        else:
            stderr.writelines("[-] Tagmanager URL not found\n")
    stderr.writelines("\n[+] Done! \n")
