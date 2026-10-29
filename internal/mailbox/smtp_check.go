package mailbox

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

func checkSMTP(mxHosts []string) error {
	if len(mxHosts) == 0 {
		return fmt.Errorf("no MX hosts provided")
	}

	var lastErr error 

	for _, mxHost := range mxHosts {
		host := strings.TrimSuffix(
			strings.TrimSuffix(mxHost),
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

		// handshake succeeded 
		_ = client.Quit()
		conn.Close()
		return nil
	}

	if lastErr == nil {
		return fmt.Errorf("no usable MX hostnames provided")
	}

	return fmt.Errorf(
		"all MX servers failed: %w",
		lastErr,
	)
}

