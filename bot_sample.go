// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"os"
	"os/signal"
	"regexp"
	"strings"
  "net/http"
	"github.com/mattermost/mattermost-server/model"
	"io/ioutil"
	"encoding/json"
	"fmt"
	//"encoding/base64"
	"log"
	//"reflect"
)

const (
	SAMPLE_NAME = "Mr. Clanky"

	USER_EMAIL    = "mschult4@nd.edu"
	USER_PASSWORD = "password"
	USER_NAME     = "mr_clanky"
	USER_FIRST    = "Madalyn"
	USER_LAST     = "Schulte"

	TEAM_NAME        = "NDLUG"
	CHANNEL_LOG_NAME = "bots"
)

var client *model.Client4
var webSocketClient *model.WebSocketClient

var botUser *model.User
var botTeam *model.Team
var debuggingChannel *model.Channel

type Sports struct {
IdTeam string `json:"idTeam"`
IdSoccerXML string `json:"idSoccerXML"`
IntLoved string `json:"intLoved"`
StrTeam string `json:"strTeam"`
StrTeamShort string `json:"strTeamShort"`
StrAlternate string `json:"strAlternate"`
IntFormedYear string `json:"intFormedYear"`
StrSport string `json:"strSport"`
StrLeague string `json:"strLeague"`
IdLeague string `json:"idLeague"`
StrDivision string `json:"strDivision"`
StrManager string `json:"strManager"`
StrStadium string `json:"strStadium"`
StrKeywords string `json:"strKeywords"`
StrRSS string `json:"strRSS"`
StrStadiumThumb string `json:"strStadiumThumb"`
StrStadiumDescription string `json:"strStadiumDescription"`
StrStadiumLocation string `json:"strStadiumLocation"`
IntStadiumCapacity string `json:"intStadiumCapacity"`
StrWebsite string `json:"strWebsite"`
StrFacebook string `json:"strFacebook"`
StrTwitter string `json:"strTwitter"`
StrInstagram string `json:"strInstagram"`
StrDescriptionEN string `json:"strDescriptionEN"`
StrDescriptionDE string `json:"strDescriptionDE"`
StrDescriptionFR string `json:"strDescriptionFR"`
StrDescriptionCN string `json:"strDescriptionCN"`
StrDescriptionIT string `json:"strDescriptionIT"`
StrDescriptionJP string `json:"strDescriptionJP"`
StrDescriptionRU string `json:"strDescriptionRU"`
StrDescriptionES string `json:"strDescriptionES"`
StrDescriptionPT string `json:"strDescriptionPT"`
StrDescriptionSE string `json:"strDescriptionSE"`
StrDescriptionNL string `json:"strDescriptionNL"`
StrDescriptionHU string `json:"strDescriptionHU"`
StrDescriptionNO string `json:"strDescriptionNO"`
StrDescriptionIL string `json:"strDescriptionIL"`
StrDescriptionPL string `json:"strDescriptionPL"`
StrGender string `json:"strGender"`
StrCountry string `json:"strCountry"`
StrTeamBadge string `json:"strTeamBadge"`
StrTeamJersey string `json:"strTeamJersey"`
StrTeamLogo string `json:"strTeamLogo"`
StrTeamFanart1 string `json:"strTeamFanart1"`
StrTeamFanart2 string `json:"strTeamFanart2"`
StrTeamFanart3 string `json:"strTeamFanart3"`
StrTeamFanart4 string `json:"strTeamFanart4"`
StrTeamBanner string `json:"strTeamBanner"`
StrYoutube string `json:"strYoutube"`
StrLocked string `json:"strLocked"`
}

type SportsAPIResponse struct {
  Teams []Sports `json:"teams"`
}

type dataMsg struct {
	teams json.RawMessage
}

type dstStruct struct {
	idLeague string
}

type Person struct {
    Name string
    Parents map[string]string
}

//type levelone struct {
//	name string
//	sublevel map[string]
//}


// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	println(SAMPLE_NAME)

	SetupGracefulShutdown()

	client = model.NewAPIv4Client("https://chat.ndlug.org")

	// Lets test to see if the mattermost server is up and running
	MakeSureServerIsRunning()

	// lets attempt to login to the Mattermost server as the bot user
	// This will set the token required for all future calls
	// You can get this token with client.AuthToken
	LoginAsTheBotUser()

	// If the bot user doesn't have the correct information lets update his profile
	UpdateTheBotUserIfNeeded()

	// Lets find our bot team
	FindBotTeam()

	// This is an important step.  Lets make sure we use the botTeam
	// for all future web service requests that require a team.
	//client.SetTeamId(botTeam.Id)

	// Lets create a bot channel for logging debug messages into
	CreateBotDebuggingChannelIfNeeded()
	SendMsgToDebuggingChannel("_"+SAMPLE_NAME+" says 'Stay in school!'_", "")

	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4("wss://chat.ndlug.org", client.AuthToken)
	if err != nil {
		println("We failed to connect to the web socket")
		PrintError(err)
	}

	webSocketClient.Listen()

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				HandleWebSocketResponse(resp)
			}
		}
	}()

	// You can block forever with
	select {}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser() {
	if user, resp := client.Login(USER_EMAIL, USER_PASSWORD); resp.Error != nil {
		println(USER_EMAIL)
		println(USER_PASSWORD)

		println(user)
		println(resp)
		println(resp.Error)
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botUser = user
	}
}

func UpdateTheBotUserIfNeeded() {
	if botUser.FirstName != USER_FIRST || botUser.LastName != USER_LAST || botUser.Username != USER_NAME {
		botUser.FirstName = USER_FIRST
		botUser.LastName = USER_LAST
		botUser.Username = USER_NAME

		if user, resp := client.UpdateUser(botUser); resp.Error != nil {
			println("We failed to update the Sample Bot user")
			PrintError(resp.Error)
			os.Exit(1)
		} else {
			botUser = user
			println("Looks like this might be the first run so we've updated the bots account settings")
		}
	}
}

func FindBotTeam() {
	if team, resp := client.GetTeamByName(TEAM_NAME, ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + TEAM_NAME + "'")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botTeam = team
	}
}

func CreateBotDebuggingChannelIfNeeded() {
	if rchannel, resp := client.GetChannelByName(CHANNEL_LOG_NAME, botTeam.Id, ""); resp.Error != nil {
		println("We failed to get the channels")
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = CHANNEL_LOG_NAME
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = botTeam.Id
	if rchannel, resp := client.CreateChannel(channel); resp.Error != nil {
		println("We failed to create the channel " + CHANNEL_LOG_NAME)
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + CHANNEL_LOG_NAME)
	}
}

func SendMsgToDebuggingChannel(msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = debuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		PrintError(resp.Error)
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	HandleMsgFromDebuggingChannel(event)
}

func HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ingore it
	if event.Broadcast.ChannelId != debuggingChannel.Id {
		return
	}

	// Lets only reponded to messaged posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	println("responding to debugging channel msg")

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == botUser.Id {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("It's aliiiiiive", post.Id)
			return
		}

		// if you see any word matching 'up' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm up up up!", post.Id)
			return
		}

		// if you see any word matching 'running' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Yes I'm running (away!)", post.Id)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Hello to you too!", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [nN][hH][lL](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the nhl scores", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [nN][fF][lL](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the nfl scores", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [nN][bB][aA](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the nba scores", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [mM][lL][bB](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the mlb scores", post.Id)
			return
		}


		if matched, _ := regexp.MatchString(`!scores? [mM][lL][sS](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the mls scores", post.Id)
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [eE][pP][lL](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the epl scores", post.Id)
			return
		}


		if matched, _ := regexp.MatchString(`!scores? aoe2(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("who plays that game anymore, nerd?", post.Id)
			return
		}

		// if you see any word matching 'score' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)score(?:$|\W)`, post.Message); matched {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", "https://www.thesportsdb.com/api/v1/json/1/searchteams.php?t=Arsenal", nil)

			res, err := client.Do(req)
			if err != nil {
				log.Fatal("res error: ", err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal("body error: ", err)
			}

      var s = new(SportsAPIResponse)
			err = json.Unmarshal([]byte(body), &s)
			if err != nil {
				log.Fatal("s error: ", err)
			}


      fmt.Println(s.Teams[0].StrDescriptionEN)


			//m := make(map[string]interface{})
			//n := make(map[string]map[string]interface{})
			//m := n["teams"]

			//var m dataMsg
			//err = json.Unmarshal(body, &m)
			//if err != nil {
			//	log.Fatal("m error: ", err)
			//}

			//strs := m["teams"].([]interface{})
			//str1 := strs[0]
			//fmt.Println(str1["idLeague"].(map[string]string))
			//n, _ := m["teams"]
      //var n dataMsg
			//err = json.Unmarshal(*m["teams"], &n)
      //var dst interface{}
			//dst = new(dstStruct)
			//err = json.Unmarshal(m.teams, &dst)
			//if err != nil {
			//	log.Fatal("dst error: ", err)
			//}

      //var str string
			//err = json.Unmarshal(*n["idLeague"], &str)

			//fmt.Println(m["teams"])



			SendMsgToDebuggingChannel("Here's a score", post.Id)
			return
		}

	}

	//SendMsgToDebuggingChannel("I did not understand you!", post.Id)
}

func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}

			SendMsgToDebuggingChannel("_"+SAMPLE_NAME+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}
