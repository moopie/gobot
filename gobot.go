package main

import (
    "github.com/moopie/gobot/bot"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
)

var (
    server = flag.String("server", "irc.quakenet.org", "Location of the server")
    port   = flag.Int("port", 6667, "Port of the server")
    nick   = flag.String("nick", "Archer", "Which nick to use")
    user   = flag.String("user", "gobot", "Username")
    name   = flag.String("name", "gobot", "Ident")
    pass   = flag.String("pass", "", "Server password, not nickserv")
    chans  = flag.String("join", "#redditeu", "Join channels on connect")
)

func main() {
    flag.Parse()
    // Hack to make the bot stay alive on heroku
    // http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html
    fmt.Println("Listening...")
    go fakeHttp()
    fmt.Println("Continuing with the rest of the app")
    // Create connection stuff
    join := strings.Split(*chans, ",")
    bot.Connect(*nick, *user, *name, *server, *port, join)

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
