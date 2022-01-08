package emailService

import (
	"context"
	"fmt"
	"log"
	"os"
	"task-scheduler/internal/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	Subject     string `json:"subject,omitempty"`
	To          string `json:"to,omitempty"`
	HtmlContent string `json:"content,omitempty"`
}

type Mailer struct {
	logHandler logger.Logger
}

func (es *Mailer) SendEmail(ctx context.Context, email Email) error {
	from := mail.NewEmail("Hunain Mehmood", os.Getenv("SENDER_EMAIL_ADDRESS"))
	to := mail.NewEmail("External User", email.To)
	message := mail.NewSingleEmail(from, email.Subject, to, "", email.HtmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return nil
	}
}

// NewService initializes the Users struct with all its dependencies and returns a new instance
// all dependencies of Users should be sent as arguments of NewService
func NewService(l logger.Logger, pqdriver *pgxpool.Pool) (*Mailer, error) {
	return &Mailer{
		logHandler: l,
	}, nil
}
