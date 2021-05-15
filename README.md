![Supported Python versions](https://img.shields.io/badge/python-3.6+-blue.svg?style=flat-square&logo=python)
![License](https://img.shields.io/badge/license-GNU-green.svg?style=flat-square&logo=gnu)

# DomainRelationShips

This script try to get related domains / subdomains by looking at Google Analytics IDs from a URL. First search for ID of Google Analytics in the webpage and then request to **builtwith** with the ID.

**Note**: It does not work with all websites.It is searched by the following expressions: 

```
->  "www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+"
->  "UA-\d+-\d+"
```

## Installation:

```
> sudo pip3 install -r requirements.txt
```

## Usage

```
> python3 analyticsrelationships.py -u https://www.example.com
```

# Author

This project has been developed by:

* **Josué Encinar García** -- [@JosueEncinar](https://twitter.com/JosueEncinar)


# Disclaimer!

This is a PoC. The author is not responsible for any illegitimate use.
