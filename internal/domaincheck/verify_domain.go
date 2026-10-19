package domaincheck

import (
	"log"
	"strings"
	"net"
)

type DomainCheckResult struct {
	domain string
	hasMX bool
	mxRecord []string
	hasSPF bool
	spfRecord string
	hasDMARC bool
	dmarcRecord string
	isValid bool
}

func CheckDomain(domain string) DomainCheckResult {

	var result DomainCheckResult
	//MX, SPF, DMARC lookup 

	mxRecords, err := net.LookupMX(domain)
	if err!=nil {
		log.Printf("MX lookup failed: %v\n", err)
	} else if len(mxRecords) > 0 {
		result.hasMX = true
		for _, record := range mxRecords {
			result.mxRecord = append(
				result.mxRecord, 
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
			result.hasSPF = true
			result.spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err!=nil {
		log.Printf("DMARC lookup failed: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			result.hasDMARC = true
			result.dmarcRecord = record
			break
		}
	}

	if result.hasMX == true {
		result.isValid = true
	} else {
		result.isValid = false
	}

	return result
}