package razwork

import (
    "github.com/moopie/gobot/message"
    "strings"
)

var (
    listen  chan message.Message
    respond chan message.Message
)

func Register(listener, responder chan message.Message) {
    listen = listener
    respond = responder
}

func Start() {
    for {
        select {
        case msg := <-listen:
            msglc := strings.ToLower(msg.Message)
            nicklc := strings.ToLower(msg.Sender)
            if (strings.Contains(nicklc, "raziel")) {
                if (strings.Contains(msglc, "work")) {
                    respond <- *message.Response(msg.Channel, "why are you not at work?")
                } else if (strings.Contains(msglc, "home")) {
                    respond <- *message.Response(msg.Channel, "how was work? :)")
                }
            }
        }
    }
}
