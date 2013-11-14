package main

import (
    irc "github.com/fluffle/goirc/client"
    "flag"
    "strings"
    "fmt"
    "time"
)

var (
    server = flag.String("server", "irc.freenode.net", "Location of the server")
    port = flag.Int("port", 6667, "Port of the server")
    nick = flag.String("nick", "BigFriendlyRobot", "Which nick to use")
    user = flag.String("user", "vrobot", "Username")
    name = flag.String("name", "vrobot", "Ident")
    pass = flag.String("pass", "", "Server password, not nickserv")
    chans = flag.String("channels", "", "Join channels on connect")
    listener = make(chan *Message, 20)
    responder = make(chan *Message, 20)
    connection = new(irc.Conn)
)

type Message struct {
    Timestamp time.Time
    Target string
    Sender string
    Message string
}

type Module interface {
    Register(in chan *Message, out chan *Message)
}

func NewMessage(line *irc.Line) *Message {
    m := new(Message)

    if (line.Cmd == "PRIVMSG") {
        m.Target = line.Args[0] // Channel
        m.Message = line.Args[1] // Message
    }

    return m
}

func main() {
    flag.Parse()
    connection := irc.SimpleClient(*nick, *user, *name)
    connection.AddHandler(irc.CONNECTED, connect) // Join channels when you connect
    connection.AddHandler("PRIVMSG", recieve) // Do stuff when you recieve a PRIVMSG
    quit := make(chan bool)
    connection.AddHandler(irc.DISCONNECTED, func (conn *irc.Conn, line *irc.Line) {quit <- true})

    // No port yet, find out how to append an int to string
    if err := connection.Connect(*server); err != nil {
        fmt.Println("Connection error: %s", err.Error())
    }

    <-quit
}

func connect(conn *irc.Conn, line *irc.Line) {
    channels := strings.Split(*chans, ",")
    for i := range(channels) {
        conn.Join(channels[i])
    }
}

func recieve(conn *irc.Conn, line *irc.Line) {
    listener <- NewMessage(line)

    fmt.Println("[", line.Args[0], "]", line.Nick, ":", line.Args[1])
}

func send(msg *Message) {
    connection.Privmsg(msg.Target, msg.Message)
}
