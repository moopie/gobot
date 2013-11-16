package razwork

import (
    "github.com/moopie/gobot/message"
    "strings"
)

type Razwork struct{}

func (r *Razwork) Respond(msg *message.Message) *message.Message {
    m := message.Empty()
    msglc := strings.ToLower(msg.Message)
    nicklc := strings.ToLower(msg.Sender)

    if strings.Contains(nicklc, "raziel") {
        if strings.Contains(msglc, "work") {
            m = message.Response(msg.Channel, "why are you not at work?")
        } else if strings.Contains(msglc, "home") {
            m = message.Response(msg.Channel, "how was work? :)")
        }
    }

    return m
}
