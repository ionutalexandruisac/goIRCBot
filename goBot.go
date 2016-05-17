package main

import (
    "irc"
    "fmt"
    "strings"
)

func main() {
    // IRC Variables to be defined (required for connecting the bot to the IRC server)
    var botVersion = "v0.1"
    var botNickname = "goBot"
    var botUsername = "goBot"
    var ircServer = "irc.freenode.net:6667"
    var botMaster = "darks1de"
    channels := []string{"#bx3", "#bx4"}

    // Trying to connect to server
    goBot := irc.IRC(botNickname, botUsername)
    err := goBot.Connect(ircServer)
    if err != nil {
        fmt.Println("Connection failed.")
        return
    }

    // On connect: Join channels in the list
    for i := 0; i < len(channels); i++ {
		    goBot.Join(channels[i])
	     }

    // Message on join function (greet people when joining the channel + auto-op master)
    goBot.AddCallback("JOIN", func (e *irc.Event) {
        if e.Nick == botNickname {
			         goBot.Privmsg(e.Arguments[0], "Hi peeps. I am a NON friendly Go Bot undergoing development.")
		           } else if e.Nick == botMaster {
                 goBot.Notice(e.Nick,"Hello master. Auto-opping you now.")
                 goBot.Mode(e.Arguments[0], "+o", e.Nick)
               } else {
                 goBot.Privmsg(e.Arguments[0], "Hello " + e.Nick + "! Welcome to " + e.Arguments[0] + ".")
               }
             })

      // Recognizes the bot master (based on nickname for now) and performs required action
      goBot.AddCallback("PRIVMSG", func (e *irc.Event) {
        if e.Nick == botMaster {
            var tempEvent = e.Arguments [1];
            privmsg := strings.Split(tempEvent," ")
            // No action specified
            if len(privmsg) <= 0 {
              goBot.Notice(e.Nick,"You must specify a course of action, master (e.g.: join, part, kick, etc)")
            } else if len(privmsg) > 0 {
                     var action = privmsg[0]
                     // Action: kick, no reason specified, using default: Your behavior is not conducive...
                    if action == "kick" && len(privmsg) > 2 && len (privmsg) < 4 {
                    roomName, nickName := privmsg[1], privmsg[2]
                    var reason = "Your behavior is not conducive to the desired environment."
                    if nickName == botMaster {
                      goBot.Notice(e.Nick,"Huh?")
                      goBot.Notice(botMaster,e.Nick + " tried to kick you on " + roomName )
                    } else {
                      goBot.Kick(nickName, roomName, reason)
                    }
                    // Action: kick, incorrect syntax (missing channelName or nickName)
                    } else if action == "kick" && len(privmsg) < 3 {
                      goBot.Notice(e.Nick,"SYNTAX must be: kick #channel nick <reason>")
                    // Action: kick, reason specified
                    } else if action == "kick" && len(privmsg) == 4 {
                      roomName, nickName, reason := privmsg[1], privmsg[2], privmsg[3]
                      if nickName == botMaster {
                        goBot.Notice(e.Nick,"Huh?")
                        goBot.Notice(botMaster,e.Nick + " tried to kick you on " + roomName + " with reason: " + reason)
                      } else {
                        goBot.Kick(nickName, roomName, reason)
                      }
                    // Action: nick change, incorrect syntax
                    } else if action == "nick" && len(privmsg) != 2 {
                      goBot.Notice(e.Nick,"SYNTAX must be: nick new_nickname")
                    // Action: nick change, OK syntax
                    } else if action == "nick" && len(privmsg) == 2 {
                      newNickname := privmsg[1]
                      goBot.Notice(e.Nick,"Changing nickname to: " + newNickname)
                      goBot.Nick(newNickname)
                    // Action: part channel, incorrect syntax
                    } else if action == "part" && len(privmsg) != 2 {
                      goBot.Notice(e.Nick,"SYNTAX must be: part #channel_name")
                    // Action: part channel, OK syntax
                    } else if action == "part" && len(privmsg) == 2 {
                      channelToPart := privmsg[1]
                      goBot.Notice(e.Nick,"Leaving channel: " + channelToPart)
                      goBot.Part(channelToPart)
                    // Action: join channel, incorrect syntax
                    } else if action == "join" && len(privmsg) != 2 {
                      goBot.Notice(e.Nick,"SYNTAX must be: join #channel_name")
                    // Action: join channel, OK syntax
                    } else if action == "join" && len(privmsg) == 2 {
                      channelToJoin := privmsg[1]
                      goBot.Notice(e.Nick,"Joining channel: " + channelToJoin)
                      goBot.Join(channelToJoin)
                    // Action: op, incorrect syntax.
                    } else if action == "op" && len(privmsg) != 3 {
                      goBot.Notice(e.Nick,"SYNTAX must be: op #channel_name nickname")
                    // Action: op, OK syntax
                    } else if action == "op" && len(privmsg) == 3 {
                      channelName, nickName := privmsg[1], privmsg[2]
                      goBot.Mode(channelName, "+o", nickName)
                    // Action: deop, incorrect syntax
                    } else if action == "deop" && len(privmsg) != 3 {
                      goBot.Notice(e.Nick,"SYNTAX must be: deop #channel_name nickname")
                    // Action: deop, OK syntax
                    } else if action == "deop" && len(privmsg) == 3 {
                      channelName, nickName := privmsg[1], privmsg[2]
                      goBot.Mode(channelName, "-o", nickName)
                    // Action: change channel mode, incorrect syntax
                    } else if action == "mode" && len(privmsg) != 3 {
                      goBot.Notice(e.Nick,"SYNTAX must be: mode #channel_name channel_mode_here. e.g.: mode #test +m")
                    // Action: change channel mode, OK syntax
                    } else if action == "mode" && len(privmsg) == 3 {
                      channelName, channelMode := privmsg[1], privmsg[2]
                      goBot.Mode(channelName, channelMode)
                    }
                  } else {
                    goBot.Notice(e.Nick,"Something went wrong.")
                  }
              // Master not recognized
              } else if e.Nick != botMaster && e.Arguments[0] == botNickname {
                goBot.Notice(e.Nick, "I refuse to listen to you! You're not my master.")
                } else {
                }
              })

    // Handle VERSION & PING events
    goBot.AddCallback("CTCP_VERSION", func(e *irc.Event) {
	        goBot.Notice(e.Nick, "Currently running GoBot version: " + botVersion)
	        })
    goBot.AddCallback("CTCP_PING", func(e *irc.Event) {
          goBot.SendRawf("NOTICE %s :\x01%s\x01", e.Nick, e.Message())
          })

    goBot.Loop()

} // THE END
