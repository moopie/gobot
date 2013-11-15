package hello

import (
    "github.com/moopie/gobot/message"
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
            if msg.Message == "!hello" {
                respond <- *message.Response(msg.Channel, "Hi!")
            }
        }
    }
}
