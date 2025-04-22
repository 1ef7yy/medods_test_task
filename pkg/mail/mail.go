package mail

import (
	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

type SMTPService struct {
	log  logger.Logger
	addr string
}

func NewSMTP(log logger.Logger, addr string) SMTPService {
	return SMTPService{
		log:  log,
		addr: addr,
	}
}

func (s SMTPService) SendMail(to, text string) error {
	// mock smtp send
	s.log.Warnf("sending message to %s: '%s'", to, text)
	return nil
}
