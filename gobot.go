package main

import (
    irc "github.com/fluffle/goirc/client"
    "github.com/moopie/gobot/message"
    "flag"
    "strings"
    "fmt"
    "net/http"
    // modules
    "github.com/moopie/gobot/modules/hello"
    "github.com/moopie/gobot/modules/pms"
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
    // Hack to make the bot stay alive on heroku
    // http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html
    fmt.Println("Listening...")
    go fakeHttp()
    fmt.Println("Continuing with the rest of app")
    // Create connection stuff
    connection := irc.SimpleClient(*nick, *user, *name)
    connection.AddHandler(irc.CONNECTED, connect) // Join channels when you connect
    connection.AddHandler("PRIVMSG", recieve) // Do stuff when you recieve a PRIVMSG
    quit := make(chan bool)
    connection.AddHandler(irc.DISCONNECTED, func (conn *irc.Conn, line *irc.Line) {quit <- true})
    hello.Register(listener, responder)
    pms.Register(listener, responder)

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

// Said heroku hack in action
// makes heroku think we are running web app
type httpHandler struct {}

func (h httpHandler) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request) {
        fmt.Println(w, "Hello, $user!")
}

func fakeHttp() {
    // http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    //     fmt.Fprintln(res, "Hello, $user!")
    // })
    http.ListenAndServe(":5000", httpHandler{})
}
