# README

A simple golang program in a docker container for checking if certificates for any domain is about to expire. The check will run Mondays at 03:00 UTC and send an email for any domain that has a certificate that will expire in less than 20 days.

Add a `domains.txt` to the root dir (the one where this file is located). Add all the domains you want to watch for expiring certificates, one per line.

E.g.

```
1.example.com
2.example.com
www.example.com
```
Configure your email settings by copying `template.env` to `.env` and filling out your information.


