package data

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

const (
	FROM     = "guinyote.manyosoft@gmail.com"
	PASSWORD = "WeAreManyos123"
)

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func send(toEmail string, toUsername string) {
	// Receiver email address.
	to := []string{
		toEmail,
	}
	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	// Message.
	body := "Hola " + toUsername + ",\nHa creado una cuenta en la aplicación de Guiñote.\n\n" +
		"Atentamente,\nEl equipo de Mañosoft <3"
	msg := composeMimeMail(toEmail, FROM, "Bienvenido a GuiñoteApp", body)

	// Authentication.
	auth := smtp.PlainAuth("", FROM, PASSWORD, smtpServer.host)
	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), auth, FROM, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

// Never fails, tries to format the address if possible
func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
