package pms

import (
    "github.com/moopie/gobot/message"
    "math/rand"
    "strings"
)

type Pms struct{}

func (p *Pms) Respond(msg *message.Message) *message.Message {
    m := message.Empty()

    lowered := strings.ToLower(msg.Message)
    if strings.Contains(lowered, "i love you") {
        m = message.Response(msg.Channel, getMood())
    }

    return m
}

func random() int {
    return rand.Intn(3)
}

func getMood() string {
    mood := random()
    if mood == 1 {
        return "I love you too <3"
    } else if mood == 2 {
        return "meh"
    } else {
        return "fuck off"
    }
}
