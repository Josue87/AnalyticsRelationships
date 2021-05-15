import urllib.request
import requests
import re
import argparse


def banner():
    print("""
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

>> Get related domains / subdomains by looking at Google Analytics IDs
>> By @JosueEncinar

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
    pattern2 = "UA-\d+-\d+"
    response = requests.get(url, headers={'User-agent': 'Mozilla/5.0 (Linux; Android 10; SM-A205U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.86 Mobile Safari/537.36'})
    if response.status_code == 200:
        text = response.text
        match = re.findall(pattern, text)
        if match:
            return True, f"https://{match[0].replace('ns.html', 'gtm.js')}"
        else:
            return False, re.findall(pattern2, text)
    return False, None

def get_domains(id):
    pattern = "/relationships/[a-z0-9\-\_]+\.[a-z]+"
    url = f"https://builtwith.com/relationships/tag/{id}"
    try:
        u = urllib.request.urlopen(url)
        data = u.read().decode(errors="ignore")
        return re.findall(pattern, data)
    except:
        return []

def show_data(data):
    if data:
        for u in data:
            analytics_id = u.split('-')
            analytics_id = "-".join(analytics_id[0:2])
            print(f"\n[+] Analytics ID: {analytics_id}")
            domains = get_domains(analytics_id)
            if domains:
                for domain in get_domains(analytics_id):
                    print(f"|__ {domain.replace('/relationships/','')}")
            else:
                print("|__ NOT FOUND")
    else:
        print("[-] Analytics ID not found...")

if __name__ == "__main__":
    banner()
    parser = argparse.ArgumentParser()
    parser.add_argument('-u','--url', help="URL to extract Google Analytics ID",required=True)
    args = parser.parse_args()
    url =  args.url
    print(f"[+] Analyzing url: {url}")
    tagmanager, data = get_googletagmanager(url)
    if tagmanager and data:
        print(f"[+] URL with UA >> {data}")
        uas = get_UA(data)
        show_data(uas)
    elif data:
        show_data(data)
    else:
        print("[-] Tagmanager URL not fount")
