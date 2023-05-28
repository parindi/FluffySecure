package notification

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/mail"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	gomail "github.com/wneessen/go-mail"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/logging"
	"github.com/authelia/authelia/v4/internal/random"
	"github.com/authelia/authelia/v4/internal/templates"
	"github.com/authelia/authelia/v4/internal/utils"
)

// NewSMTPNotifier creates a SMTPNotifier using the notifier configuration.
func NewSMTPNotifier(config *schema.SMTPNotifierConfiguration, certPool *x509.CertPool) *SMTPNotifier {
	var tlsconfig *tls.Config

	if config.TLS != nil {
		tlsconfig = utils.NewTLSConfig(config.TLS, certPool)
	}

	opts := []gomail.Option{
		gomail.WithPort(config.Address.Port()),
		gomail.WithTLSConfig(tlsconfig),
		gomail.WithHELO(config.Identifier),
		gomail.WithTimeout(config.Timeout),
		gomail.WithoutNoop(),
	}

	ssl := config.Address.IsExplicitlySecure()

	if ssl {
		opts = append(opts, gomail.WithSSL())
	}

	switch {
	case ssl:
		break
	case config.DisableStartTLS:
		opts = append(opts, gomail.WithTLSPolicy(gomail.NoTLS))
	case config.DisableRequireTLS:
		opts = append(opts, gomail.WithTLSPolicy(gomail.TLSOpportunistic))
	default:
		opts = append(opts, gomail.WithTLSPolicy(gomail.TLSMandatory))
	}

	var domain string

	at := strings.LastIndex(config.Sender.Address, "@")

	if at >= 0 {
		domain = config.Sender.Address[at+1:]
	} else {
		domain = "localhost.localdomain"
	}

	return &SMTPNotifier{
		config: config,
		domain: domain,
		random: &random.Cryptographical{},
		tls:    utils.NewTLSConfig(config.TLS, certPool),
		log:    logging.Logger(),
		opts:   opts,
	}
}

// SMTPNotifier a notifier to send emails to SMTP servers.
type SMTPNotifier struct {
	config *schema.SMTPNotifierConfiguration
	domain string
	random random.Provider
	tls    *tls.Config
	log    *logrus.Logger
	opts   []gomail.Option
}

// StartupCheck implements model.StartupCheck to perform startup check operations.
func (n *SMTPNotifier) StartupCheck() (err error) {
	var client *gomail.Client

	if client, err = gomail.NewClient(n.config.Address.Hostname(), n.opts...); err != nil {
		return fmt.Errorf("failed to establish client: %w", err)
	}

	ctx := context.Background()

	if err = client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("failed to dial connection: %w", err)
	}

	if err = client.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

// Send a notification via the SMTPNotifier.
func (n *SMTPNotifier) Send(ctx context.Context, recipient mail.Address, subject string, et *templates.EmailTemplate, data any) (err error) {
	msg := gomail.NewMsg(
		gomail.WithMIMEVersion(gomail.Mime10),
		gomail.WithBoundary(n.random.StringCustom(30, random.CharSetAlphaNumeric)),
	)

	n.setMessageID(msg, n.domain)

	if err = msg.From(n.config.Sender.String()); err != nil {
		return fmt.Errorf("notifier: smtp: failed to set from address: %w", err)
	}

	if err = msg.AddTo(recipient.String()); err != nil {
		return fmt.Errorf("notifier: smtp: failed to set to address: %w", err)
	}

	msg.Subject(strings.ReplaceAll(n.config.Subject, "{title}", subject))

	switch {
	case n.config.DisableHTMLEmails:
		if err = msg.SetBodyTextTemplate(et.Text, data); err != nil {
			return fmt.Errorf("notifier: smtp: failed to set body: text template errored: %w", err)
		}
	default:
		if err = msg.AddAlternativeHTMLTemplate(et.HTML, data); err != nil {
			return fmt.Errorf("notifier: smtp: failed to set body: html template errored: %w", err)
		}

		if err = msg.AddAlternativeTextTemplate(et.Text, data); err != nil {
			return fmt.Errorf("notifier: smtp: failed to set body: text template errored: %w", err)
		}
	}

	var client *gomail.Client

	if client, err = gomail.NewClient(n.config.Address.Hostname(), n.opts...); err != nil {
		return fmt.Errorf("notifier: smtp: failed to establish client: %w", err)
	}

	if auth := NewOpportunisticSMTPAuth(n.config); auth != nil {
		client.SetSMTPAuthCustom(auth)
	}

	if err = client.DialWithContext(ctx); err != nil {
		return fmt.Errorf("notifier: smtp: failed to dial connection: %w", err)
	}

	if err = client.Send(msg); err != nil {
		return fmt.Errorf("notifier: smtp: failed to send message: %w", err)
	}

	if err = client.Close(); err != nil {
		return fmt.Errorf("notifier: smtp: failed to close connection: %w", err)
	}

	return nil
}

func (n *SMTPNotifier) setMessageID(msg *gomail.Msg, domain string) {
	rn := n.random.Intn(100000000)
	rm := n.random.Intn(10000)
	rs := n.random.StringCustom(17, random.CharSetAlphaNumeric)
	pid := os.Getpid() + rm

	msg.SetMessageIDWithValue(fmt.Sprintf("%d.%d%d.%s@%s", pid, rn, rm, rs, domain))
}
