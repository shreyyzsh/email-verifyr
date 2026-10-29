package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"

	"github.com/shreyyzsh/email-verifyr.git/internal/verify"

	"github.com/joho/godotenv"
)

func main() {
	// env var
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	sender := os.Getenv("SMTP_SENDER")
	if sender == "" {
		log.Fatal("SMTP_SENDER is not set in .env")
	}

	fmt.Printf("Welcome to Email Verifyr...\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter the email to verify: ")

	for scanner.Scan() {
		// split domain and mailbox
		email := scanner.Text()
		parts := strings.SplitN(email, "@", 2)

		if len(parts) != 2 || parts[0] == "" {
			fmt.Printf("Invalid email format\n")
			fmt.Print("Enter the email to verify: ")
			continue
		}

		mailbox := parts[0]
		domainName := parts[1]

		resDomain := verify.DomainCheck(domainName)

		if !resDomain.IsValid {
			fmt.Println("Domain is invalid.\n")
			fmt.Print("Enter the email to verify: ")
			continue
		} 

		fmt.Println("Domain is Valid.\n")
		
		status, err := verify.VerifyMailbox(
			resDomain.MXRecord,
			sender,
			email,
		)
		if err != nil {
			fmt.Println("Mailbox verification inconclusive: ", err)
		} else {
			fmt.Printf("Mailbox status: %s\n", status)
		}

		_ = mailbox 
		break
	}

	if err := scanner.Err(); err!=nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}

}