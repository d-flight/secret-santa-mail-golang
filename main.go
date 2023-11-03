package main

import (
	"fmt"
	"regexp"
	"strings"

	"math/rand"

	"gopkg.in/gomail.v2"
)

// smtp settings
var host = "smtp.example.com"
var port = 1337
var user = "user@example.com"
var pass = "password"

// mail settings
var sender = "user@example.com"
var subject = "Secret Santa"

var participants = []string{
	// "<Mason Johnston> m.johnston@gmail.com",
	// "<Cody Dawson> cody.dawson@hotmail.com",
	// "<Bethany Davies> daviesbeth@icloud.com"
}

const template = `
Hi %name%,

you're secret santa for: %target%
`

func main() {
	dialer := gomail.NewDialer(host, port, user, pass)

	for recipient, target := range generateRandomPairs(participants) {
		fmt.Printf("sending message to %s\n", recipient)
		sendMail(recipient, target, dialer)
	}

	fmt.Println("Done.")
}

func sendMail(recipient string, target string, dialer *gomail.Dialer) {
	recipientName, recipientAddress := splitAddress(recipient)
	targetName, targetAddress := splitAddress(target)

	body := strings.Replace(template, "%name%", recipientName, -1)
	body = strings.Replace(body, "%target%", fmt.Sprintf("%s, %s", targetName, targetAddress), -1)

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetAddressHeader("To", recipientAddress, recipientName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := dialer.DialAndSend(m); err != nil {
		panic(err)
	}
}

func splitAddress(address string) (string, string) {
	mailRegexp := regexp.MustCompile(`<([^>]+)>\s*(.*)`)
	result := mailRegexp.FindAllStringSubmatch(address, -1)

	matches := result[0]

	return matches[1], matches[2]
}

func generateRandomPairs(addresses []string) map[string]string {
	pairs := map[string]string{}
	var targets = make([]string, len(addresses))
	copy(targets, addresses)

	for _, recipient := range addresses {
		var target = ""
		var targetIndex = 0

		for ok := true; ok; ok = (target == "" || target == recipient) {
			targetIndex = rand.Intn(len(targets))
			target = targets[targetIndex]

			if target == recipient && len(targets) == 1 {
				return generateRandomPairs(addresses)
			}
		}

		pairs[recipient] = target
		targets = append(targets[:targetIndex], targets[targetIndex+1:]...)
	}

	return pairs
}
