package main

import (
    irc "github.com/fluffle/goirc/client"
    "github.com/moopie/gobot/message"
    "flag"
    "strings"
    "fmt"
    // modules
    "github.com/moopie/gobot/modules/hello"
)

var (
    server = flag.String("server", "irc.quakenet.org", "Location of the server")
    port = flag.Int("port", 6667, "Port of the server")
    nick = flag.String("nick", "Archer", "Which nick to use")
    user = flag.String("user", "gobot", "Username")
    name = flag.String("name", "gobot", "Ident")
    pass = flag.String("pass", "", "Server password, not nickserv")
    chans = flag.String("channels", "#redditeu", "Join channels on connect")
    connection = new(irc.Conn)
    listener = make(chan message.Message)
    responder = make(chan message.Message)
)

func main() {
    flag.Parse()
    connection := irc.SimpleClient(*nick, *user, *name)
    connection.AddHandler(irc.CONNECTED, connect) // Join channels when you connect
    connection.AddHandler("PRIVMSG", recieve) // Do stuff when you recieve a PRIVMSG
    quit := make(chan bool)
    connection.AddHandler(irc.DISCONNECTED, func (conn *irc.Conn, line *irc.Line) {quit <- true})
    hello.Register(listener, responder)

    go hello.Start()

    // No port yet, TODO: find out how to append an int to a string (yes, really)
    if err := connection.Connect(*server); err != nil {
        fmt.Println("Connection error: %s", err.Error())
    }

    for {
        select {
            case msg := <-responder:
            connection.Privmsg(msg.Channel, msg.Message)
        }
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
    listener <- *message.Line(line)

    fmt.Println("[", line.Args[0], "]", line.Nick, ":", line.Args[1])
}
