package hello

import (
    "github.com/moopie/gobot/message"
)

type Hello struct{}

func (h *Hello) Respond(msg *message.Message) *message.Message {
    m := message.Empty()

    if msg.Message == "!hello" {
        m = message.Response(msg.Channel, "Hi!")
    }

    return m
}
