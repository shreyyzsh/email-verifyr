package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"

	"github.com/shreyyzsh/email-verifyr.git/internal/domaincheck"
)

func main() {
	fmt.Printf("Welcome to Email Verifyr...\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter the email to verify: ")

	for scanner.Scan() {
		// split domain and mailbox
		email := scanner.Text()
		parts := strings.SplitN(email, "@", 2)

		if len(parts) != 2 || parts[0] == "" {
			fmt.Printf("Invalid email format")
			continue
		}

		res_domain := domaincheck.CheckDomain(parts[1])
		if res_domain.IsValid == true {
			fmt.Printf("\nDomain is Valid.")
			break
		}
	}

/* 	fmt.Printf("Fetch records? [y/n]")
	for scanner.Scan() {
		inp := scanner.Text()

		if inp == "y" {
			fmt.Printf("hasMX: %v\nhasSPF: %v\nhasDMARC: %v\n\nMXRecords: %v\nSPFRecords: %v\nDMARCRecords: %v", res_domain.HasMX, res_domain.HasSPF, res_domain.HasDMARC, res_domain.MXRecord, res_domain.SPFRecord, res_domain.DMARCRecord)
		} else {
			break
		}
	} */

	if err := scanner.Err(); err!=nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}

}