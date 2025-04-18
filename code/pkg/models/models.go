package models

import "time"

type Cert struct {
	NotAfter   time.Time
	CommonName string
	DNSNames   []string
}
