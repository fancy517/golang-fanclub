package mailer

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	workingInterval = time.Millisecond * 500
	idleInterval    = time.Second * 5
)

const (
	TemplateBetWin            = "bet_win.tmpl"
	TemplateBetLose           = "bet_lose.tmpl"
	TemplateBetCanceled       = "match_cancel.tmpl"
	TemplateActivation        = "user_activation.tmpl"
	TemplateSettingActivation = "setting_activation.tmpl"
	TemplatePasswordReset     = "password_reset.tmpl"
	TemplateNewsletter        = "newsletter.tmpl"
	TemplateCongrats          = "congrats.tmpl"
	TemplateTicket            = "ticket.tmpl"
	TemplateReferral          = "referral.tmpl"
	TemplateReferralReward    = "referral_reward.tmpl"
)

func (m mailerImpl) Run() {
	go func() {
		for {
			select {
			case p := <-m.mailQueue:
				if err := m.Send(p.Recipient, p.Template, p.Data); err != nil {
					log.Println(fmt.Errorf("Send mail error, data=%v; %w", p, err))
				}
				time.Sleep(workingInterval)
			default:
				time.Sleep(idleInterval)
			}
		}
	}()
	// return
	// go func() {
	// 	for {
	// 		data, ok := <-m.mailQueue
	// 		if ok {
	// 			switch data := data.(type) {
	// 			case MailUserActivation:
	// 				log.Println("Processing user activation mail")
	// 				err := m.Send(data.Recipient, "user_activation.tmpl", gin.H{
	// 					"code": data.Code,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			case MailUserPasswordReset:
	// 				log.Println("Processing user activation mail")
	// 				err := m.Send(data.Recipient, "password_reset.tmpl", gin.H{
	// 					"code": data.Code,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			case MailBetWin:
	// 				log.Println("Processing bet win")
	// 				err := m.Send(data.Recipient, "bet_win.tmpl", gin.H{
	// 					"amount": data.Amount,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			case MailBetLose:
	// 				log.Println("Processing bet lose")
	// 				err := m.Send(data.Recipient, "bet_lose.tmpl", gin.H{
	// 					"amount": data.Amount,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			case MailMatchCancel:
	// 				log.Println("Processing match cancel")
	// 				err := m.Send(data.Recipient, "bet_lose.tmpl", gin.H{
	// 					"matchID": data.MatchID,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			case MailNewsletter:
	// 				log.Println("Processing newsletter")
	// 				err := m.Send(data.Recipient, "newsletter.tmpl", gin.H{
	// 					"title":   data.Title,
	// 					"content": data.Content,
	// 				})
	// 				if err != nil {
	// 					log.Println(fmt.Errorf("failed to send newsletter mail, data=%#v; %w", data, err))
	// 				} else {
	// 					log.Printf("sent mail, data=%#v\n", data)
	// 				}
	// 			default:
	// 				panic("unknown mail type")
	// 			}
	// 			time.Sleep(time.Second * 1) // sleep for 1 seconds
	// 		} else {
	// 			time.Sleep(time.Second * 5)
	// 		}
	// 	}
	// }()
}

func (m mailerImpl) addQueue(pack MailPackage) {
	go func() {
		m.mailQueue <- pack
	}()
}

func (m mailerImpl) SendActivationCode(recipient, username, code string) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateActivation,
		Data: gin.H{
			"name": username,
			"code": code,
		},
	})
}

func (m mailerImpl) SendSettingCode(recipient, username, code string) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateSettingActivation,
		Data: gin.H{
			"name": username,
			"code": code,
		},
	})
}

func (m mailerImpl) SendResetPasswordCode(recipient, username, code string) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplatePasswordReset,
		Data: gin.H{
			"name": username,
			"code": code,
		},
	})
}

type NewsletterPage struct {
	Title   string
	Content template.HTML
}

func (m mailerImpl) SendNewsletter(recipient, username string, title, content string) {
	html := template.HTML(markdown2HTML([]byte(content)))
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateNewsletter,
		Data: NewsletterPage{
			Title:   title,
			Content: html,
		},
	})
}

func (m mailerImpl) SendCongrats(recipient, name string) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateCongrats,
		Data: gin.H{
			"name": name,
		},
	})
}

func (m mailerImpl) SendTicket(name, email, subject, message string) {
	m.addQueue(MailPackage{
		Recipient: m.adminEmail,
		Template:  TemplateTicket,
		Data: gin.H{
			"name":    name,
			"email":   email,
			"subject": subject,
			"message": message,
		},
	})
}

func (m mailerImpl) SendReferral(recipient, name, refereeEmail, refereeName string, amount float64) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateReferral,
		Data: gin.H{
			"name":          name,
			"referee_name":  refereeName,
			"referee_email": refereeEmail,
			"amount":        amount,
		},
	})
}

func (m mailerImpl) SendRewardsAlertOnReferral(recipient, name, refereeEmail, refereeName string, amount float64) {
	m.addQueue(MailPackage{
		Recipient: recipient,
		Template:  TemplateReferralReward,
		Data: gin.H{
			"name":          name,
			"referee_name":  refereeName,
			"referee_email": refereeEmail,
			"amount":        amount,
		},
	})
}
