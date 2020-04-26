package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math/rand"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/balerter/balerter/internal/alert/message"
	"go.uber.org/zap"
)

func (e *Email) Send(message *message.Message) error {
	var (
		c    *smtp.Client
		conn net.Conn
		err  error
	)
	if e.conf.ServerPort == "465" {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}

		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = e.conf.ServerName
		}

		conn, err = tls.Dial("tcp", e.conf.ServerName, tlsConfig)
		if err != nil {
			return fmt.Errorf("establish TLS connection to server: %w", err)
		}
	} else {
		var (
			d   = net.Dialer{}
			err error
		)
		conn, err = d.Dial("tcp", net.JoinHostPort(e.conf.ServerName, e.conf.ServerPort))
		if err != nil {
			return fmt.Errorf("establish connection to server: %w", err)
		}
	}
	c, err = smtp.NewClient(conn, e.conf.ServerName)
	if err != nil {
		conn.Close()
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer func() {
		if err := c.Quit(); err != nil {
			e.logger.Warn("email client", zap.String("smtp", "failed to close SMTP connection"), zap.Error(err))
		}
	}()

	if e.conf.RequireTLS {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return fmt.Errorf("'requireTLS' is true but %q does not advertise the STARTTLS extension",
				net.JoinHostPort(e.conf.ServerName, e.conf.ServerPort))
		}

		tlsConf := &tls.Config{InsecureSkipVerify: true}

		if tlsConf.ServerName == "" {
			tlsConf.ServerName = e.conf.ServerName
		}

		if err := c.StartTLS(tlsConf); err != nil {
			return fmt.Errorf("send STARTTLS command: %w", err)
		}
	}

	if ok, mech := c.Extension("AUTH"); ok {
		auth, err := e.auth(mech)
		if err != nil {
			return fmt.Errorf("find auth mechanism: %w", err)
		}
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("%T auth: %w", auth, err)
			}
		}
	}

	addrs, err := mail.ParseAddressList(e.conf.From)
	if err != nil {
		return fmt.Errorf("parse 'from' addresses: %w", err)
	}
	if len(addrs) != 1 {
		return fmt.Errorf("must be exactly one 'from' address (got: %d)", len(addrs))
	}
	if err = c.Mail(addrs[0].Address); err != nil {
		return fmt.Errorf("send MAIL command: %w", err)
	}
	addrs, err = mail.ParseAddressList(e.conf.To)
	if err != nil {
		return fmt.Errorf("parse 'to' addresses: %w", err)
	}
	for _, addr := range addrs {
		if err = c.Rcpt(addr.Address); err != nil {
			return fmt.Errorf("send RCPT command: %w", err)
		}
	}

	msg, err := c.Data()
	if err != nil {
		return fmt.Errorf("send DATA command: %w", err)
	}
	defer msg.Close()

	buffer := &bytes.Buffer{}
	multipartBuffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(multipartBuffer)

	fmt.Fprintf(buffer, "%s: %s\r\n", "From", mime.QEncoding.Encode("utf-8", e.conf.From))
	fmt.Fprintf(buffer, "%s: %s\r\n", "To", mime.QEncoding.Encode("utf-8", e.conf.To))
	fmt.Fprintf(buffer, "%s: [%s/%s] %s\r\n", "Subject",
		mime.QEncoding.Encode("utf-8", message.AlertName),
		mime.QEncoding.Encode("utf-8", message.Level),
		mime.QEncoding.Encode("utf-8", strings.Join(message.Fields, ",")))

	fmt.Fprintf(buffer, "Message-Id: %s\r\n", fmt.Sprintf("<%d.%d@%s>", time.Now().UnixNano(), rand.Uint64(), e.hostname))
	fmt.Fprintf(buffer, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(buffer, "Content-Type: multipart/alternative;  boundary=%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(buffer, "MIME-Version: 1.0\r\n\r\n")

	_, err = msg.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("write headers: %w", err)
	}

	if len(message.Text) > 0 {
		w, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-Transfer-Encoding": {"quoted-printable"},
			"Content-Type":              {"text/plain; charset=UTF-8"},
		})
		if err != nil {
			return fmt.Errorf("create part for text body: %w", err)
		}

		qw := quotedprintable.NewWriter(w)
		_, err = qw.Write([]byte(message.Text))
		if err != nil {
			return fmt.Errorf("write text part: %w", err)
		}
		err = qw.Close()
		if err != nil {
			return fmt.Errorf("close text part: %w", err)
		}
	}

	err = multipartWriter.Close()
	if err != nil {
		return fmt.Errorf("close multipartWriter: %w", err)
	}

	_, err = msg.Write(multipartBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("write body buffer: %w", err)
	}

	return nil
}
