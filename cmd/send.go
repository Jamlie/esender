/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Jamlie/esender/types"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

type tomlConfig struct {
	Smtp struct {
		Email    string `toml:"email"`
		Password string `toml:"password"`
	}
}

type EmailService struct {
	From      string
	Password  string
	To        []string
	Subject   string
	Body      string
	EmailType types.EmailType
	Port      int16
}

type response struct {
	err error
}

func (e *EmailService) SendEmail() {
	to := strings.Join(e.To, ",")
	message := "To: " + to + "\r\n" +
		"Subject: " + e.Subject + "\r\n" +
		"\r\n" + e.Body

	auth := smtp.PlainAuth("", e.From, e.Password, string(e.EmailType))

	err := smtp.SendMail(string(e.EmailType)+fmt.Sprintf(":%d", e.Port), auth, e.From, e.To, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
}

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Opens a form to fill in the email details and send it.",
	Long:  `This command opens a form to fill in the email details and send it.`,
	Run: func(cmd *cobra.Command, args []string) {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile(path.Join(homedir, ".esender.toml"), os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var tomlContent tomlConfig
		if _, err := toml.NewDecoder(file).Decode(&tomlContent); err != nil {
			log.Fatal(err)
		}

		var emailService EmailService
		emailService.Password = tomlContent.Smtp.Password
		emailService.Port = 587
		emailService.From = tomlContent.Smtp.Email

		doesSend := EmailSenderForm(&emailService)
		if !doesSend {
			return
		}

		err = spinner.New().
			Title("Sending Email").
			Action(func() {
				emailService.SendEmail()
			}).
			Run()

		if err != nil {
			log.Fatal(err)
		}
	},
}

func EmailSenderForm(emailInfo *EmailService) bool {
	var emails string

	confirmed := false

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[types.EmailType]().
				Title("What's your email service?").
				Options(
					huh.NewOption("Gmail", types.Gmail),
					huh.NewOption("Outlook", types.Outlook),
					huh.NewOption("Yahoo", types.Yahoo),
				).
				Value(&emailInfo.EmailType),
			huh.NewInput().
				Title("What's your email?").
				Value(&emailInfo.From).
				Validate(func(s string) error {
					_, err := mail.ParseAddress(s)
					if err != nil {
						return errors.New("The string you entered is not a valid email.")
					}

					return nil
				}),

			huh.NewInput().
				Title("To? (Separate by (\",\"))").
				Value(&emails).
				Validate(func(s string) error {
					emailAddrs := strings.Split(s, ",")
					for _, email := range emailAddrs {
						if _, err := mail.ParseAddress(email); err != nil {
							return errors.New(fmt.Sprintf("%s is not a valid email address", email))
						}
					}

					return nil
				}),

			huh.NewInput().
				Title("What is the subject?").
				Value(&emailInfo.Subject),

			huh.NewText().
				Title("What is the Body?").
				Value(&emailInfo.Body),

			huh.NewConfirm().
				Title("Do you want to send the email?").
				Value(&confirmed),
		),
	)

	form.Run()

	if !confirmed {
		fmt.Println("Exiting...")
		return confirmed
	}

	splitAddrs := strings.Split(emails, ",")
	for _, s := range splitAddrs {
		emailInfo.To = append(emailInfo.To, s)
	}

	return confirmed
}

func init() {
	rootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
