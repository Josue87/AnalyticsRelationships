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

## Example

An example of use is shown in the following image (just a little bit obfuscated):

![image](https://user-images.githubusercontent.com/16885065/118356444-84475300-b575-11eb-9b1f-bc5c587d620f.png)


# Author

This project has been developed by:

* **Josué Encinar García** -- [@JosueEncinar](https://twitter.com/JosueEncinar)


# Disclaimer!

This is a PoC. The author is not responsible for any illegitimate use.
