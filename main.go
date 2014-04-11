package main

import (
	"flag"

	"github.com/jmoiron/sqlx"
	"github.com/thoj/go-ircevent"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"log"
	"time"
)

var dsn = flag.String("dsn", ":memory:", "The dsn for the database")
var nick = flag.String("nick", "irclogger", "User to log as")
var channel = flag.String("chan", "", "channel to log")
var server = flag.String("server", "chat.freenode.net:6667", "irc server")
var now = time.Now

func main() {
	flag.Parse()
	conn := irc.IRC(*nick, *nick)
	//	conn.UseTLS = true
	conn.Connect(*server)
	conn.Join(*channel)
	db, err := sqlx.Open("sqlite3", *dsn)
	if err != nil {
		panic(err)
	}
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		_, err := db.NamedExec(`
	INSERT INTO gnotty_ircmessage(
		nickname,message,server,channel,
		message_time,join_or_leave)
	    VALUES(:nickname, :message, :server, :channel, :message_time, :join_or_leave)`,
			map[string]interface{}{
				"nickname":      e.Nick,
				"message":       e.Message(),
				"server":        e.Source,
				"channel":       e.Arguments[0],
				"message_time":  time.Now().Format("2006-01-02 15:04:05.000000"),
				"join_or_leave": "0",
			})
		if err != nil {
			log.Println(err)
		}
	})
	conn.AddCallback("JOIN", func(e *irc.Event) {
		_, err := db.NamedExec(`
	INSERT INTO gnotty_ircmessage(
		nickname,message,server,channel,
		message_time,join_or_leave)
	    VALUES(:nickname, :message, :server, :channel, :message_time, :join_or_leave)`,
			map[string]interface{}{
				"nickname":      e.Nick,
				"message":       e.Message(),
				"server":        e.Source,
				"channel":       e.Arguments[0],
				"message_time":  time.Now().Format("2006-01-02 15:04:05.000000"),
				"join_or_leave": "1",
			})
		if err != nil {
			log.Println(err)
		}
	})
	log.Println("here we go")
	conn.Loop()
}
