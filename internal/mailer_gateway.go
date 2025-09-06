package internal

import "fmt"

type MailerGateway interface {
	send(recipient string, subject string, message string)
}

type MailerGatewayMemory struct {
}

func NewMailerGatewayMemory() *MailerGatewayMemory {
	return &MailerGatewayMemory{}
}

func (m *MailerGatewayMemory) send(recipient string, subject string, message string) {
	fmt.Println("send", recipient, subject, message)
}
