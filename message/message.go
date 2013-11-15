package message

import (
    irc "github.com/fluffle/goirc/client"
    // "github.com/moopie/gobot/main"
    "time"
)

type Message struct {
    Timestamp time.Time
    Channel   string
    Sender    string
    Message   string
}

func Line(line *irc.Line) *Message {
    m := new(Message)

    m.Timestamp = line.Time
    m.Sender = line.Nick
    m.Channel = line.Args[0]
    m.Message = line.Args[1]

    return m
}

func Response(channel, message string) *Message {
    m := new(Message)

    m.Timestamp = time.Now()
    // m.Sender = main.Nick
    m.Channel = channel
    m.Message = message

    return m
}
