# DomainRelationShips

This script try to get related domains / subdomains by looking at Google Analytics IDs from a URL. First search for ID of Google Analytics in the webpage and then request to **builtwith** with the ID.

**Note**: It does not work with all websites.It is searched by the following expressions: 

```
->  "www\.googletagmanager\.com/ns\.html\?id=[A-Z0-9\-]+"
-> GTM-[A-Z0-9]+
->  "UA-\d+-\d+"
```

## 2 Versions:

* [Python](Python)
* [GO](GO)

# Author

This project has been developed by:

* **Josué Encinar García** -- [@JosueEncinar](https://twitter.com/JosueEncinar)


# Disclaimer!

This is a PoC. The author is not responsible for any illegitimate use.
