package main

import (
    irc "github.com/fluffle/goirc/client"
    "github.com/moopie/gobot/message"
    "github.com/moopie/gobot/module"
    "flag"
    "strings"
    "fmt"
    "os"
    "net/http"
    // modules
    "github.com/moopie/gobot/modules/hello"
    "github.com/moopie/gobot/modules/pms"
    "github.com/moopie/gobot/modules/razwork"
)

var (
    server = flag.String("server", "irc.quakenet.org", "Location of the server")
    port = flag.Int("port", 6667, "Port of the server")
    nick = flag.String("nick", "Archer", "Which nick to use")
    user = flag.String("user", "gobot", "Username")
    name = flag.String("name", "gobot", "Ident")
    pass = flag.String("pass", "", "Server password, not nickserv")
    chans = flag.String("join", "#redditeu", "Join channels on connect")
    connection = new(irc.Conn)
    subscribers = make([]module.Module, 0, 10)
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

    register(new(hello.Hello))
    register(new(pms.Pms))
    register(new(razwork.Razwork))

    // No port yet, TODO: find out how to append an int to a string (yes, really)
    if err := connection.Connect(*server); err != nil {
        fmt.Println("Connection error: %s", err.Error())
    }

    <-quit
}

func register(mod module.Module) {
    subscribers = append(subscribers, mod)
}

func connect(conn *irc.Conn, line *irc.Line) {
    channels := strings.Split(*chans, ",")
    for i := range(channels) {
        conn.Join(channels[i])
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
            if (response.Channel == *nick) {
                response.Channel = line.Src
            }
            conn.Privmsg(response.Channel, response.Message)
        }
    }

    fmt.Println("["+line.Args[0]+"] "+line.Nick+": "+line.Args[1])
}

// Said heroku hack in action
// makes heroku think we are running web app
func fakeHttp() {
    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(res, `<html><title>Gobot</title><body>Ping this place every 40-60 mins<br />
            <a href="http://stackoverflow.com/questions/5480337/easy-way-to-prevent-heroku-idling">
            http://stackoverflow.com/questions/5480337/easy-way-to-prevent-heroku-idling</a>
            </body></html>`)
    })
    http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
