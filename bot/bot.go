package bot

import (
    "fmt"
    irc "github.com/fluffle/goirc/client"
    "github.com/moopie/gobot/message"
    "github.com/moopie/gobot/module"
    // Modules
    "github.com/moopie/gobot/modules/hello"
    "github.com/moopie/gobot/modules/pms"
    "github.com/moopie/gobot/modules/razwork"
    "github.com/moopie/gobot/modules/grnaer"
)

var (
    Nick string
    User string
    Name string
    chans []string
    ircClient *irc.Conn
    subscribers = make([]module.Module, 0, 10)
)

func Connect(nick, user, name, server string, port int, join []string) {
    Nick = nick
    User = user
    Name = name
    chans = join
    ircClient := irc.SimpleClient(Nick, User, Name)
    ircClient.AddHandler(irc.CONNECTED, connect) // Join channels when you connect
    ircClient.AddHandler("PRIVMSG", recieve) // Do stuff when you recieve a PRIVMSG
    quit := make(chan bool)
    ircClient.AddHandler(irc.DISCONNECTED, func (conn *irc.Conn, line *irc.Line) {quit <- true})

    register(new(hello.Hello))
    register(new(pms.Pms))
    register(new(razwork.Razwork))
    register(grnaer.Create(Nick))

    // No port yet, TODO: find out how to append an int to a string (yes, really)
    if err := ircClient.Connect(server); err != nil {
        fmt.Println("conn error: %s", err.Error())
    }

    <-quit
}

func register(mod module.Module) {
    subscribers = append(subscribers, mod)
}

func connect(conn *irc.Conn, line *irc.Line) {
    fmt.Println("Connecting to the IRC server")
    for i := range(chans) {
        fmt.Println("Joining channel " + chans[i])
        conn.Join(chans[i])
    }
}

func recieve(conn *irc.Conn, line *irc.Line) {
    for _,mod := range(subscribers) {
        msg := message.Line(line)
        response := mod.Respond(msg)
        if(response.Message != "") {
            // When channel is nick, this means it was a private message
            // so send the message to the sender instead
            // Send it to the full raw nick+hostmask name
            if (response.Channel == Nick) {
                response.Channel = line.Src
            }
            conn.Privmsg(response.Channel, response.Message)
        }
    }

    fmt.Println("["+line.Args[0]+"] "+line.Nick+": "+line.Args[1])
}
