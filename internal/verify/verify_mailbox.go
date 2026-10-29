package verify

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
	"net/textproto"
)

type MailboxStatus string

const (
	MailboxAccepted MailboxStatus = "accepted"
	MailboxRejected MailboxStatus = "rejected"
	MailboxUnknown MailboxStatus = "unknown"
)

func VerifyMailbox(
	mxHosts []string,
	sender string,
	recipient string,
	) (MailboxStatus, error) {

	if len(mxHosts) == 0 {
		return MailboxUnknown, fmt.Errorf("no MX hosts provided")
	}

	permanentRejections := 0
	var lastErr error 

	for _, mxHost := range mxHosts {
		host := strings.TrimSuffix(
			mxHost,
			".",
		)
		
		if host == "" {
			continue
		}

		address := net.JoinHostPort(host, "25")

		conn, err := net.DialTimeout(
			"tcp",
			address,
			5*time.Second,
		)
		if err != nil {
			lastErr = fmt.Errorf(
				"%s: TCP connection failed: %w",
				host, err,
			)
			continue
		}

		_ = conn.SetDeadline(
			time.Now().Add(5 * time.Second),
		)

		// reading server's greeting
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			lastErr = fmt.Errorf(
				"%s: SMTP greeting failed: %w",
				host, err,
			)
			continue
		}

		// EHLO
		err = client.Hello("verifyr.local")
		if err != nil {
			client.Close()
			conn.Close()

			lastErr = fmt.Errorf(
				"%s: SMTP handshake failed: %w",
				host, err,
			)
			continue
		}

		// handshake succeeded, 
		// specify sender
		err = client.Mail(sender)
		if err != nil {
			lastErr = fmt.Errorf(
				"%s: MAIL FROM failed: %w", 
				host, err,
			)

			client.Close()
			conn.Close()
			continue
		}

		// checking recipient
		err = client.Rcpt(recipient)
		if err == nil {
			client.Close()
			conn.Close()
			return MailboxAccepted, nil
		}

		lastErr = err

 		if smtpErr, ok := err.(*textproto.Error); ok {
			switch smtpErr.Code {
			case 550, 551, 553:
				permanentRejections++
			}
		}

		client.Close()
		conn.Close()
	}

	if permanentRejections == len(mxHosts) {
		return MailboxRejected, nil
	}

	if lastErr != nil {
		return MailboxUnknown, fmt.Errorf(
			"mailbox could not be verified: %w", lastErr,
		)
	}

	return MailboxUnknown, nil
}

