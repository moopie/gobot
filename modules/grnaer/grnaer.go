package grnaer

import (
    "github.com/moopie/gobot/message"
    "github.com/moopie/gobot/main"
    "strings"
    "math/rand"
    "time"
)

type Grnaer struct{}

func shuffle(s string) string {
    rand.Seed(time.Now().UnixNano())
    t := make([]byte, len(s))
    for i, r := range rand.Perm(len(s)) {
        t[i] = s[r]
    }
    return string(t)
}

func (p *Grnaer) Respond(msg *message.Message) *message.Message {
    m := message.Empty()
    msglc := strings.ToLower(msg.Message)
    nicklc := strings.ToLower(msg.Sender)

    if strings.Contains(nicklc, "garner") && strings.HasPrefix(msglc, main.Nick) {
        s := "Archer wtf are you doing"
        words := strings.Split(s, " ")[1:]
        randwords := msg.Sender + ": "
        for _,word := range words {
            randwords += shuffle(word) + " "
        }
        m = randwords
    }

    return m
}