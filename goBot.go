package main

import (
    "irc"
    "fmt"
)

func main() {
    // Variables to be defined
    var botVersion = "v0.1"
    var botNickname = "goBot"
    var botUsername = "goBot"
    var roomName = "#bx"
    var ircServer = "irc.freenode.net:6667"

    // Trying to connect to server
    goBot := irc.IRC(botNickname, botUsername)
    err := goBot.Connect(ircServer)
    if err != nil {
        fmt.Println("Connection failed.")
        return
    }

    // Message on join
    goBot.Join(roomName)
    goBot.AddCallback("JOIN", func (e *irc.Event) {
      if e.Nick == botNickname {
			     goBot.Privmsg(roomName, "Hi. I am a NON friendly goBot undergoing development.")
		   } else {
          goBot.Privmsg(roomName, "Hello " + e.Nick + "! Welcome to " + roomName + ".")
      }
    })

    goBot.AddCallback("PRIVMSG", func (e *irc.Event) {
       goBot.Privmsg(roomName, e.Message())
   })

   // Handle VERSION & PING events
   goBot.AddCallback("CTCP_VERSION", func(e *irc.Event) {
	     goBot.Notice(e.Nick, "Currently running GoBot version: " + botVersion)
	 })
   goBot.AddCallback("CTCP_PING", func(e *irc.Event) {
       goBot.SendRawf("NOTICE %s :\x01%s\x01", e.Nick, e.Message())
   })

   goBot.Loop();

}
