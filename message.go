package main

import (
	"regexp"
	"strings"

	"github.com/go-mail/mail"
)

var rxEmail = regexp.MustCompile(".+@.+\\..+")

type Message struct {
	Email   string
	Content string
	Errors  map[string]string
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(msg.Email))
	if match == false {
		msg.Errors["Email"] = "Please enter a valid email address."
	}

	if strings.TrimSpace(msg.Content) == "" {
		msg.Errors["Content"] = "Please enter a message"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) DeliverAsync() error {
	username := goDotEnvVariable("MAILTRAP_UN")
	password := goDotEnvVariable("MAILTRAP_PW")

	email := mail.NewMessage()
	email.SetHeader("To", "carlos@csaenz.dev")
	email.SetHeader("From", "server@csaenz.dev")
	email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "New message via Contact Form")
	email.SetBody("text/plain", msg.Content)
	email.AddAlternative("text/html",
		`
		<!DOCTYPE html>
		<html>
		<head>
		<style>
		body {
		font-family: Arial, sans-serif;
		color: #333333;
		background-color: #f9f9f9;
		margin: 0;
		padding: 0;
		}
		.email-container {
		width: 100%;
		max-width: 600px;
		margin: 20px auto;
		background-color: #ffffff;
		border: 1px solid #dddddd;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
		padding: 20px;
		}
		h1 {
		color: #0073e6;
		font-size: 24px;
		text-align: center;
		}
		p {
		font-size: 16px;
		line-height: 1.5;
		}
		.content {
		padding: 10px;
		background-color: #f1f1f1;
		border-radius: 5px;
		margin: 15px 0;
		font-style: italic;
		}
		.footer {
		text-align: center;
		font-size: 14px;
		color: #888888;
		margin-top: 20px;
		}
		</style>
		</head>
		<body>
		<div class="email-container">
		<h1>New Message via Contact Form</h1>
		<p><strong>Email:</strong> `+msg.Email+`</p>
		<div class="content">"`+msg.Content+`"</div>
		<p class="footer">Made with &#128140;<br>Mailtrap</p>
		</div>
		</body>
		</html>
		`)

	return mail.NewDialer("live.smtp.mailtrap.io", 587, username, password).DialAndSend(email)
}
