package domaincheck

import (
	"log"
	"strings"
	"net"
)

type DomainCheckResult struct {
	HasMX bool
	MXRecord []string
	HasSPF bool
	SPFRecord string
	HasDMARC bool
	DMARCRecord string
	IsValid bool
}

func CheckDomain(domain string) DomainCheckResult {

	result := DomainCheckResult{}
	//MX, SPF, DMARC lookup 

	mxRecords, err := net.LookupMX(domain)
	if err!=nil {
		log.Printf("MX lookup failed: %v\n", err)
	} else if len(mxRecords) > 0 {
		result.HasMX = true
		for _, record := range mxRecords {
			result.MXRecord = append(
				result.MXRecord, 
				strings.TrimSuffix(record.Host, "."),
			)
		}
	}
	
	txtRecords, err := net.LookupTXT(domain)
	if err!=nil {
		log.Printf("SPF lookup failed: %v\n", err)
	} 
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			result.HasSPF = true
			result.SPFRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err!=nil {
		log.Printf("DMARC lookup failed: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			result.HasDMARC = true
			result.DMARCRecord = record
			break
		}
	}

	if result.HasMX == true {
		result.IsValid = true
	} else {
		result.IsValid = false
	}

	return result
}