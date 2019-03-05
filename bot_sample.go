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
	//"fmt"
	//"encoding/base64"
	"log"
	"math"
	"time"
	//"reflect"
	"strconv"
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

type Event struct {
	IdEvent string `json:"idEvent"`
	IdSoccerXML string `json:"idSoccerXML"`
	StrEvent string `json:"strEvent"`
	StrFilename string `json:"strFilename"`
	StrSport string `json:"strSport"`
	IdLeague string `json:"idLeague"`
	StrLeague string `json:"strLeague"`
	StrSeason string `json:"strSeason"`
	StrDescriptionEN string `json:"strDescriptionEN"`
	StrHomeTeam string `json:"strHomeTeam"`
	StrAwayTeam string `json:"strAwayTeam"`
	IntHomeScore string `json:"intHomeScore"`
	IntRound string `json:"intRound"`
	IntAwayScore string `json:"intAwayScore"`
	IntSpectators string `json:"intSpectators"`
	StrHomeGoalDetails string `json:"strHomeGoalDetails"`
	StrHomeRedCards string `json:"strHomeRedCards"`
	StrHomeYellowCards string `json:"strHomeYellowCards"`
	StrHomeLineupGoalkeeper string `json:"strHomeLineupGoalkeeper"`
	StrHomeLineupDefense string `json:"strHomeLineupDefense"`
	StrHomeLineupMidfield string `json:"strHomeLineupMidfield"`
	StrHomeLineupForward string `json:"strHomeLineupForward"`
	StrHomeLineupSubstitutes string `json:"strHomeLineupSubstitutes"`
	StrHomeFormation string `json:"strHomeFormation"`
	StrAwayRedCards string `json:"strAwayRedCards"`
	StrAwayYellowCards string `json:"strAwayYellowCards"`
	StrAwayGoalDetails string `json:"strAwayGoalDetails"`
	StrAwayLineupGoalkeeper string `json:"strAwayLineupGoalkeeper"`
	StrAwayLineupDefense string `json:"strAwayLineupDefense"`
	StrAwayLineupMidfield string `json:"strAwayLineupMidfield"`
	StrAwayLineupForward string `json:"strAwayLineupForward"`
	StrAwayLineupSubstitutes string `json:"strAwayLineupSubstitutes"`
	StrAwayFormation string `json:"strAwayFormation"`
	IntHomeShots string `json:"intHomeShots"`
	IntAwayShots string `json:"intAwayShots"`
	DateEvent string `json:"dateEvent"`
	StrDate string `json:"strDate"`
	StrTime string `json:"strTime"`
	StrTVStation string `json:"strTVStation"`
	IdHomeTeam string `json:"idHomeTeam"`
	IdAwayTeam string `json:"idAwayTeam"`
	StrResult string `json:"strResult"`
	StrCircuit string `json:"strCircuit"`
	StrCountry string `json:"strCountry"`
	StrCity string `json:"strCity"`
	StrPoster string `json:"strPoster"`
	StrFanart string `json:"strFanart"`
	StrThumb string `json:"strThumb"`
	StrBanner string `json:"strBanner"`
	StrMap string `json:"strMap"`
	StrLocked string `json:"strLocked"`
}

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

type EventAPIResponse struct {
	Event []Event `json:"event"`
}

type SportsAPIResponse struct {
  Teams []Sports `json:"teams"`
}

type Leagues struct {
IdEvent string `json:"idEvent"`
IdSoccerXML string `json:"idSoccerXML"`
StrEvent string `json:"strEvent"`
StrFilename string `json:"strFilename"`
StrSport string `json:"strSport"`
IdLeague string `json:"idLeague"`
StrLeague string `json:"strLeague"`
StrSeason string `json:"strSeason"`
StrDescriptionEN string `json:"strDescriptionEN"`
StrHomeTeam string `json:"strHomeTeam"`
StrAwayTeam string `json:"strAwayTeam"`
IntHomeScore string `json:"intHomeScore"`
IntRound string `json:"intRound"`
IntAwayScore string `json:"intAwayScore"`
IntSpectators string `json:"intSpectators"`
StrHomeGoalDetails string `json:"strHomeGoalDetails"`
StrHomeRedCards string `json:"strHomeRedCards"`
StrHomeYellowCards string `json:"strHomeYellowCards"`
StrHomeLineupGoalkeeper string `json:"strHomeLineupGoalkeeper"`
StrHomeLineupDefense string `json:"strHomeLineupDefense"`
StrHomeLineupMidfield string `json:"strHomeLineupMidfield"`
StrHomeLineupForward string `json:"strHomeLineupForward"`
StrHomeLineupSubstitutes string `json:"strHomeLineupSubstitutes"`
StrHomeFormation string `json:"strHomeFormation"`
StrAwayRedCards string `json:"strAwayRedCards"`
StrAwayYellowCards string `json:"strAwayYellowCards"`
StrAwayGoalDetails string `json:"strAwayGoalDetails"`
StrAwayLineupGoalkeeper string `json:"strAwayLineupGoalkeeper"`
StrAwayLineupDefense string `json:"strAwayLineupDefense"`
StrAwayLineupMidfield string `json:"strAwayLineupMidfield"`
StrAwayLineupForward string `json:"strAwayLineForward"`
StrAwayLineupSubstitutes string `json:"strAwayLineup"`
StrAwayFormation string `json:"strAwayFormation"`
IntHomeShots string `json:"intHomeShots"`
IntAwayShots string `json:"intAwayShots"`
DateEvent string `json:"dateEvent"`
StrDate string `json:"strDate"`
StrTime string `json:"strTime"`
StrTVStation string `json:"strTVStation"`
IdHomeTeam string `json:"idHomeTeam"`
IdAwayTeam string `json:"idAwayTeam"`
StrResult string `json:"strResult"`
StrCircuit string `json:"strCircuit"`
StrCountry string `json:"strCountry"`
StrCity string `json:"strCity"`
StrPoster string `json:"strPoster"`
StrFanart string `json:"strFanart"`
StrThumb string `json:"strThumb"`
StrBanner string `json:"strBanner"`
StrMap string `json:"strMap"`
StrLocked string `json:"strLocked"`

}

type LeaguesAPIResponse struct {
  Events []Leagues `json:"events"`
}

type LeagueSearch struct {
IdLeague string `json:"idLeague"`
StrLeague string `json:"strLeague"`
StrSport string `json:"strSport"`
StrLeagueAlternate string `json:"strLeagueAlternate"`
}

type LeagueSearchAPIResponse struct {
  Leagues []LeagueSearch `json:"leagues"`
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
			go LeagueScores(post, "4380")
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [nN][fF][lL](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the nfl scores", post.Id)
			go LeagueScores(post, "4391")
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [nN][bB][aA](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the nba scores", post.Id)
			go LeagueScores(post, "4387")
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [mM][lL][bB](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the mlb scores", post.Id)
			go LeagueScores(post, "4424")
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [mM][lL][sS](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the mls scores", post.Id)
			go LeagueScores(post, "4346")
			return
		}

		if matched, _ := regexp.MatchString(`!scores? [eE][pP][lL](?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("these will be the epl scores", post.Id)
			go LeagueScores(post, "4328")
			return
		}


		if matched, _ := regexp.MatchString(`!scores? aoe2(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("who plays that game anymore, nerd?", post.Id)
			return
		}

		// if you see any word matching 'score' then respond
		/*if matched, _ := regexp.MatchString(`(?:^|\W)score(?:$|\W)`, post.Message); matched {
			SendMsgToDebuggingChannel("Here's a score", post.Id)
			return
		}*/

		if matched, _ := regexp.MatchString(`!scores? team(?:$|\W)`, post.Message); matched {
			client := &http.Client{}
			split_str := strings.Split(post.Message, " ")
			request_str := ""
			for i := 2; i < len(split_str); i++ {
			   if i != 1{
					 request_str += "_"
				 }
				 request_str += split_str[i]
			}

			req_str := "https://www.thesportsdb.com/api/v1/json/1/searchevents.php?e=" + request_str
			req, _ := http.NewRequest("GET", req_str, nil)

			res, err := client.Do(req)
			if err != nil {
				log.Fatal("res error: ", err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal("body error: ", err)
			}
			var s = new(EventAPIResponse)
			err = json.Unmarshal([]byte(body), &s)
			if err != nil {
				log.Fatal("s error: ", err)
			}

			current_time := time.Now().Local()
			first_past := 0
			curr_date := current_time.Format("2006-01-02")
			str_date := s.Event[0].DateEvent
			for ; curr_date <= str_date; first_past++{
				first_past = first_past + 1
				str_date = s.Event[first_past].DateEvent
			}

			if len(s.Event) >= 1{
				message := ""
				//SendMsgToDebuggingChannel("3 Most Recent Team Events:", post.Id)
				for i := first_past; i < int(math.Min(float64(len(s.Event)), 3)) + first_past; i++{

					away_score, _ := strconv.Atoi(s.Event[i].IntAwayScore)
					home_score, _ := strconv.Atoi(s.Event[i].IntHomeScore)
					home_team := s.Event[i].StrHomeTeam
					away_team := s.Event[i].StrAwayTeam
					if home_score >= away_score && s.Event[i].IntHomeScore != "" {
						home_team = "**"+home_team+"**"
					}
					if away_score >= home_score && s.Event[i].IntAwayScore != "" {
						away_team = "**"+away_team+"**"
					}

					score_string := ""
					if (s.Event[i].IntHomeScore == "" && s.Event[i].IntAwayScore == "") {
						score_string += "_(not reported)_"
					} else {
						score_string += s.Event[i].IntHomeScore + " - " + s.Event[i].IntAwayScore
					}

					message += s.Event[i].DateEvent + " | " + home_team + " vs. " + away_team + " | score: " + score_string + "\n"



					/*if (s.Event[i].IntHomeScore > s.Event[i].IntAwayScore){
						SendMsgToDebuggingChannel(s.Event[i].DateEvent + " | **" + s.Event[i].StrHomeTeam + "** vs. " + s.Event[i].StrAwayTeam + " | score: " + s.Event[i].IntHomeScore + " - " + s.Event[i].IntAwayScore, post.Id)
          }else if (s.Event[i].IntHomeScore < s.Event[i].IntAwayScore){
						SendMsgToDebuggingChannel(s.Event[i].DateEvent + " | " + s.Event[i].StrHomeTeam + " vs. **" + s.Event[i].StrAwayTeam + "** | score: " + s.Event[i].IntHomeScore + " - " + s.Event[i].IntAwayScore, post.Id)
					}else{
						SendMsgToDebuggingChannel(s.Event[i].DateEvent + " | **" + s.Event[i].StrHomeTeam + " vs. **" + s.Event[i].StrAwayTeam + "** | score: " + s.Event[i].IntHomeScore + " - " + s.Event[i].IntAwayScore, post.Id)
					}*/
				}

				SendMsgToDebuggingChannel(message, post.Id)
				return
			}
      //SendMsgToDebuggingChannel(s.Event[0].DateEvent, post.Id)
		}

		if matched, _ := regexp.MatchString(`!team(?:$|\W)`, post.Message); matched {
			client := &http.Client{}
      split_str := strings.Split(post.Message, " ")
			request_str := ""
			for i := 1; i < len(split_str); i++ {
			   if i != 1{
					 request_str += "_"
				 }
				 request_str += split_str[i]
			}
			req_str := "https://www.thesportsdb.com/api/v1/json/1/searchteams.php?t=" + request_str
			req, _ := http.NewRequest("GET", req_str, nil)

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
			if len(s.Teams) > 1{
				SendMsgToDebuggingChannel("Multiple teams were returned. To disambiguate, type in one of these queries:", post.Id)
				for i := 0; i < len(s.Teams); i++{
					SendMsgToDebuggingChannel("!team " + s.Teams[i].StrTeam, post.Id)
				}
				return
			}
      if len(s.Teams) > 0{
				SendMsgToDebuggingChannel(s.Teams[0].StrDescriptionEN, post.Id)
			}else{
				SendMsgToDebuggingChannel("I'm unable to understand your query as written. The proper format is ```!team cityname teamname```", post.Id)
			}
			return
		}


		if matched, _ := regexp.MatchString(`!leagues? (?:$|\W)`, post.Message); matched {
			client := &http.Client{}
      split_str := strings.Split(post.Message, " ")
			request_str := ""
			for i := 1; i < len(split_str); i++ {
			   if i != 1{
					 request_str += "_"
				 }
				 request_str += split_str[i]
			}
			req_str := "https://www.thesportsdb.com/api/v1/json/1/searchteams.php?t=" + request_str
			req, _ := http.NewRequest("GET", req_str, nil)

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
			if len(s.Teams) > 1{
				SendMsgToDebuggingChannel("Multiple teams were returned. To disambiguate, type in one of these queries:", post.Id)
				for i := 0; i < len(s.Teams); i++{
					SendMsgToDebuggingChannel("!team " + s.Teams[i].StrTeam, post.Id)
				}
				return
			}
      if len(s.Teams) > 0{
				SendMsgToDebuggingChannel(s.Teams[0].StrDescriptionEN, post.Id)
			}else{
				SendMsgToDebuggingChannel("I'm unable to understand your query as written. The proper format is ```!team cityname teamname```", post.Id)
			}
			return
		}
	}

	//SendMsgToDebuggingChannel("I did not understand you!", post.Id)
}

func LeagueScores(post *model.Post, league_id string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.thesportsdb.com/api/v1/json/1/eventspastleague.php?id="+league_id, nil)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("res error: ", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("body error: ", err)
	}

    var s = new(LeaguesAPIResponse)
	err = json.Unmarshal([]byte(body), &s)
	if err != nil {
		log.Fatal("s error: ", err)
	}
    //fmt.Println(s.Teams[0].StrDescriptionEN)


	message := ""
	for i := 0; i < 10; i++ {
		away_score, _ := strconv.Atoi(s.Events[i].IntAwayScore)
		home_score, _ := strconv.Atoi(s.Events[i].IntHomeScore)
		home_team := s.Events[i].StrHomeTeam
		away_team := s.Events[i].StrAwayTeam
		if home_score >= away_score && s.Events[i].IntHomeScore != "" {
			home_team = "**"+home_team+"**"
		}
		if away_score >= home_score && s.Events[i].IntAwayScore != "" {
			away_team = "**"+away_team+"**"
		}

		score_string := ""
		if (s.Events[i].IntHomeScore == "" && s.Events[i].IntAwayScore == "") {
			score_string += "_(not reported)_"
		} else {
			score_string += s.Events[i].IntHomeScore + " - " + s.Events[i].IntAwayScore
		}

		message += s.Events[i].DateEvent + " | " + home_team + " vs. " + away_team + " | score: " + score_string + "\n"

	}
	SendMsgToDebuggingChannel(message, post.Id)
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
