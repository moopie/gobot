package module

import (
    "github.com/moopie/gobot/message"
)

type Module interface {
    Register(chan message.Message, chan message.Message)
}
