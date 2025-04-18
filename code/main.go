package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xulinus/certwatch/pkg/email"
	"github.com/xulinus/certwatch/pkg/models"
)

const DOMAINS_FILE = "./domains.txt"

func main() {
	log.Println("Program up and running...")

	ticker := time.NewTicker(1 * time.Minute)

	for range ticker.C {
		if isItMondayThreeAM() {
			log.Println("Checking domains...")

			domains, err := readDomainFile(DOMAINS_FILE)
			if err != nil {
				log.Fatal(fmt.Sprintf("Could not read %s\n%s", DOMAINS_FILE, err))
			}

			certs := getCertData(domains)
			checkCerts(certs)
		}
	}
}

func isItMondayThreeAM() bool {
	now := time.Now().UTC()
	return now.Weekday() == time.Monday && now.Hour() == 3 && now.Minute() == 0
}

func getCertData(domains []string) []models.Cert {
	var data []models.Cert

	for _, d := range domains {
		cert, err := getCert(d)
		if err != nil {
			log.Println(err)
		} else {
			data = append(data, cert)
		}
	}

	return data
}

func getCert(domainname string) (models.Cert, error) {
	host := domainname + ":443"

	conn, err := tls.Dial("tcp", host, &tls.Config{})
	if err != nil {
		return models.Cert{}, err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return models.Cert{}, fmt.Errorf("No certifcates found for: %s", domainname)
	}

	cert := certs[0]
	return models.Cert{
		NotAfter:   cert.NotAfter,
		CommonName: cert.Subject.CommonName,
		DNSNames:   cert.DNSNames,
	}, nil
}

func checkCerts(certs []models.Cert) {
	for _, c := range certs {
		checkCert(c)
	}
}

func checkCert(cert models.Cert) {
	now := time.Now().UTC()

	certTimeLeft := cert.NotAfter.Sub(now)

	if certTimeLeft < 24*20*time.Hour {
		email.SendReminder(cert)
	} else {
		log.Printf("Certifcate for %s will expire in %s", cert.CommonName, certTimeLeft)
	}
}

func readDomainFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var domains []string
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}
