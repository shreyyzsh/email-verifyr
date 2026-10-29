package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"

	"github.com/shreyyzsh/email-verifyr.git/internal/domain"
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

		res_domain := domain.DomainCheck(parts[1])
		if res_domain.IsValid == true {
			fmt.Printf("\nDomain is Valid.\n")
			break
		}
	}

	if err := scanner.Err(); err!=nil {
		log.Fatal("Error: could not read from input: %v\n", err)
	}

}