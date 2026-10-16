package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"github.com/shreyyzsh/email-verifyr/internal/domaincheck"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) //creates a scanner to read terminal input
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	// reads input and calls checkDomain 
	for scanner.Scan() { 
		domaincheck.CheckDomain(scanner.Text()) //text entered in terminal
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}
}

