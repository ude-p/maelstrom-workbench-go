package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Message struct {
	Src  string  `json:"src"`
	Dest string  `json:"dest"`
	Body Payload `json:"body"`
}

type Payload struct {
	Type      string `json:"type"`
	MsgID     int32  `json:"msg_id"`
	InReplyTo int32  `json:"in_reply_to"`
	Echo      string `json:"echo"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			handleMessage(scanner.Bytes())
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error occured: %v", err)
		}

	}
}

func handleMessage(input []byte) {
	var msg Message

	if err := json.Unmarshal(input, &msg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshall: %v\n", err)
		return
	}

	response := Message{
		Src:  msg.Dest,
		Dest: msg.Src,
		Body: Payload{
			MsgID:     msg.Body.MsgID,
			InReplyTo: msg.Body.MsgID,
			Echo:      msg.Body.Echo,
		},
	}

	switch msg.Body.Type {
	case "init":
		response.Body.Type = "init_ok"
	case "echo":
		response.Body.Type = "echo_ok"
	default:
		fmt.Fprintf(os.Stderr, "Unfamiliar body type: %v", msg.Body.Type)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshall: %v\n", err)
	}
}
