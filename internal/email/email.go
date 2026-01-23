package email

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type EmailHandler struct {
	mailer *Mailer
}

type Mailer struct {
	Host     string
	Port     string
	User     string
	Password string
}

func New() (*EmailHandler, error) {
	host := os.Getenv("MAIL_SERVICE")
	port := os.Getenv("MAIL_SERVICE_PORT")
	user := os.Getenv("MAIL_SERVICE_USER")
	password := os.Getenv("MAIL_SERVICE_PASSWORD")

	if host == "" || user == "" || password == "" {
		return nil, fmt.Errorf("missing MAIL_SERVICE config")
	}
	if port == "" {
		port = "587"
	}

	client := &Mailer{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	return &EmailHandler{mailer: client}, nil
}

type Email struct {
	From     Target
	To       []Target
	CC       []Target
	BCC      []Target
	Subject  string
	Body     string
	AltBody  string
	IsHTML   bool
	Priority Priority
}

type Target struct {
	Address string
	Name    string
}

type Priority int

const (
	PriorityNormal Priority = 3
	PriorityHigh   Priority = 1
	PriorityLow    Priority = 5
)

// * Pardn Chiu <dev@pardn.io>
func (r Target) toText() string {
	if r.Name == "" {
		return r.Address
	}
	return fmt.Sprintf("%s <%s>", encodeRFC2047(r.Name), r.Address)
}

func (c *Mailer) Send(msg *Email) error {
	auth := smtp.PlainAuth("", c.User, c.Password, c.Host)

	from := msg.From.Address
	if from == "" {
		from = c.User
	}

	recipients := getRecipients(msg)
	body := c.buildMessage(msg)

	addr := c.Host + ":" + c.Port

	if c.Port == "465" {
		return c.sendTLS(addr, auth, from, recipients, body)
	}
	return smtp.SendMail(addr, auth, from, recipients, []byte(body))
}

func (c *Mailer) sendTLS(addr string, auth smtp.Auth, from string, to []string, body string) error {
	tlsConfig := &tls.Config{
		ServerName: c.Host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, c.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(body)); err != nil {
		return err
	}
	return w.Close()
}

func (c *Mailer) buildMessage(msg *Email) string {
	var builder strings.Builder

	from := msg.From
	if from.Address == "" {
		from.Address = c.User
	}
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from.toText()))

	var to []string
	for _, e := range msg.To {
		to = append(to, e.toText())
	}
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))

	if len(msg.CC) > 0 {
		var cc []string
		for _, e := range msg.CC {
			cc = append(cc, e.toText())
		}
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}

	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", encodeRFC2047(msg.Subject)))
	builder.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	builder.WriteString("MIME-Version: 1.0\r\n")

	if msg.Priority != 0 && msg.Priority != PriorityNormal {
		builder.WriteString(fmt.Sprintf("X-Priority: %d\r\n", msg.Priority))
	}

	if msg.IsHTML {
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	builder.WriteString("\r\n")
	builder.WriteString(msg.Body)

	return builder.String()
}

func getRecipients(msg *Email) []string {
	var addrs []string
	for _, r := range msg.To {
		addrs = append(addrs, r.Address)
	}
	for _, r := range msg.CC {
		addrs = append(addrs, r.Address)
	}
	for _, r := range msg.BCC {
		addrs = append(addrs, r.Address)
	}
	return addrs
}

func encodeRFC2047(s string) string {
	for _, r := range s {
		if r > 127 {
			return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(s)))
		}
	}
	return s
}
