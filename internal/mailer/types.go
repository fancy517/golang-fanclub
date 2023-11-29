package mailer

type BaseMailType struct {
	Recipient string `json:"recipient"` // email
	Name      string `json:"name"`
}

type MailPackage struct {
	Recipient string
	Template  string
	Data      any
}

type MailUserActivation struct {
	BaseMailType
	Code string `json:"code"`
}

type MailUserPasswordReset struct {
	MailUserActivation
}

type MailBetWin struct {
	BaseMailType
	Amount float64 `json:"amount"`
}

type MailBetLose struct {
	BaseMailType
	Amount float64 `json:"amount"`
}

type MailMatchCancel struct {
	BaseMailType
	MatchID int `json:"match_id"`
}

type MailNewsletter struct {
	BaseMailType
	Title   string `json:"title"`
	Content string `json:"content"`
}
