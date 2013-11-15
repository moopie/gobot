package pms

import (
    "github.com/moopie/gobot/message"
    "strings"
    "math/rand"
)

var (
    listen  chan message.Message
    respond chan message.Message
)

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

func Register(listener, responder chan message.Message) {
    listen = listener
    respond = responder
}

func Start() {
    for {
        select {
        case msg := <-listen:
            msglc = strings.ToLower(msg)
            if strings.Contains(msglc, "i love you") {
                respond <- *message.Response(msg.Channel, getMood())
            }
        }
    }
}
