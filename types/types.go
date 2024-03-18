package types

type EmailType string

const (
	Gmail   EmailType = "smtp.gmail.com"
	Yahoo   EmailType = "smtp.mail.yahoo.com"
	Outlook EmailType = "smtp-mail.outlook.com"
)
