# Mattermost Bot Sample

## Overview and Project Goals

This sample Bot shows how to use the Mattermost [Go driver](https://github.com/mattermost/platform/blob/master/model/client.go) to interact with a Mattermost server, listen to events and respond to messages. Documentation for the Go driver can be found [here](https://godoc.org/github.com/mattermost/platform/model#Client).

Highlights of APIs used in this sample:
 - Log in to the Mattermost server
 - Create a channel
 - Modify user attributes 
 - Connect and listen to WebSocket events for real-time responses to messages
 - Post messages to channel in response to user queries

In using this bot framework , our final goal (i.e., the goal of Dan, Grace, Madalyn and Matt) was to create a Mattermost bot that allows the user to query a sports score api for information about their favorite teams and leagues. We provide a simple text-based query system which the bot, through a series of callbacks, formulates an appropriate response to the user, or, if the bot is unable to find an appropriate response, they will ask the user for more specific information in their subsequent queries. This project makes use of thesportsdb.com, a crowd-sourced sports information db which provides streamlined API calls to retrieve their data. 

## Taking advantage of Go's unique features
goroutines are a fat

## Bot Deployment
Bot is currently running on a Windows Azure Linux Virtual Machine. To interact with the bot, type some of the designated commands into the "bots" channel of NDLUG. Mr. Clanky will respond with various sports information.

Example commands:
	!scores (nhl | nba | nfl | mls | epl | mlb)
	!scores team [cityname] [teamname]
	!team [cityname] [teamname]
	!help
