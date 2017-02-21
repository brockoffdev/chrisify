// https://www.opsdash.com/blog/slack-bot-in-golang.html

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var token string
var chrisify string
var haar string
var faces string
var base_path = "/app/www/"
var base_url string

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "usage: slackbot slack-bot-token /path/to/chrisify /path/to/haar /path/to/faces https://base_url/\n")
		os.Exit(1)
	}

	token = os.Args[1]
	chrisify = os.Args[2]
	haar = os.Args[3]
        faces = os.Args[4]
        base_url = os.Args[5]

	// start a websocket-based Real Time API session
	ws, id := slackConnect(token)
	fmt.Println("slackbot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && m.SubType == "file_share" && strings.Contains(m.Text, "<@"+id+">") {
			go func(m Message) {
				var channel string
				json.Unmarshal(m.Channel, &channel)
				file := SaveTempFile(GetFile(m.File))
				chrisd := Chrisify(file)
				// log.Printf("Uploading to %s", channel)
				// Upload(chrisd, channel)
				url := SaveFile(chrisd)
				postMessage(ws, map[string]string{
					"type":    "message",
					"text":    url,
					"channel": channel,
				})

				defer os.Remove(file)
			}(m)
		}
	}
}
