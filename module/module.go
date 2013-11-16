package module

import (
    "github.com/moopie/gobot/message"
)

type Module interface {
    Respond(msg *message.Message) *message.Message
}
