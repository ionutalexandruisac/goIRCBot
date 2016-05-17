package main

import (
    "irc"
    "fmt"
)

func main() {
    // IRC Variables to be defined (required for connecting the bot)
    var botVersion = "v0.1"
    var botNickname = "goBot"
    var botUsername = "goBot"
    var roomName = "#bx"
    var ircServer = "irc.freenode.net:6667"
    var botMaster = "darks1de"

    // Trying to connect to server
    goBot := irc.IRC(botNickname, botUsername)
    err := goBot.Connect(ircServer)
    if err != nil {
        fmt.Println("Connection failed.")
        return
    }

    // Message on join function
    goBot.Join(roomName)
    goBot.AddCallback("JOIN", func (e *irc.Event) {
      if e.Nick == botNickname {
			     goBot.Privmsg(roomName, "Hi. I am a NON friendly goBot undergoing development.")
		   } else {
          goBot.Privmsg(roomName, "Hello " + e.Nick + "! Welcome to " + roomName + ".")
      }
    })

    // Annoyingly repeat whatever the bot master types (test function - to be removed! ;-)
    goBot.AddCallback("PRIVMSG", func (e *irc.Event) {
         if e.Nick == botMaster {
              goBot.Privmsg(roomName, e.Message())
         }
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
