import urllib.request
import requests
import re
import argparse
from sys import stderr
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

def get_UA(link):
    pattern = "UA-\d+-\d+"
    try:
        u = urllib.request.urlopen(link)
        data = u.read().decode(errors="ignore")
        match = re.findall(pattern, data)
        unique = set()
        unique = unique.union(match)
        return list(unique)
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
                    verify=False)

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
    except Exception as e:
        print(e)
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
    url = f"https://api.hackertarget.com/analyticslookup/?q={id}"
    try:
        response = requests.get(url)
        if response.status_code == 200:
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
    all_uas = [] # avoid duplicates
    if data:
        print("")
        for u in data:
            analytics_id = u.split('-')
            analytics_id = "-".join(analytics_id[0:2])
            if analytics_id not in all_uas:
                all_uas.append(analytics_id)
                print(f">> {analytics_id}")
                domains = get_domains(analytics_id)
                if domains:
                    for domain in get_domains(analytics_id):
                        print(f"|__ {domain}")
                    
                else:
                    print("|__ NOT FOUND")
                print("")
        stderr.writelines("\n[+] Done! \n")
    else:
        stderr.writelines("[-] Analytics ID not found...\n")

if __name__ == "__main__":
    banner()
    parser = argparse.ArgumentParser()
    parser.add_argument('-u','--url', help="URL to extract Google Analytics ID",required=True)
    args = parser.parse_args()
    url =  args.url
    if not url.startswith("http"):
        url = "https://" + url
    stderr.writelines(f"[+] Analyzing url: {url}\n")
    tagmanager, data = get_googletagmanager(url)
    if tagmanager and data:
        stderr.writelines(f"[+] URL with UA: {data}\n")
        stderr.writelines("[+] Obtaining information from builtwith and hackertarget\n")
        uas = get_UA(data)
        show_data(uas)
    elif data:
        show_data(data)
    else:
        stderr.writelines("[-] Tagmanager URL not fount\n")
