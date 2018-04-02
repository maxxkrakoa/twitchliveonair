# Description
A small Go program to monitor the live streams of a Twitch user.

Has logic to turn on GPIOs/LEDs on Raspberry Pi when a live stream is detected.

# Building
```
dep ensure
go build
```

## Cross compiling for Raspberry Pi
`env GOOS=linux GOARCH=arm GOARM=5 go build`

# Running
`export TWITCH_CLIENT_ID="<ClientID>"` before running program

`./twitchliveonair <twitch user id to monitor>`

## Twitch API
https://dev.twitch.tv/docs/authentication

https://dev.twitch.tv/dashboard/apps/

https://dev.twitch.tv/docs/api/reference/#get-streams

For some API calls Client ID is enough.

# Notes

## Basic user info
```
curl \
-H 'Client-ID: <clientID>' \
-X GET 'https://api.twitch.tv/helix/users?login=<login>'
```

## Get stream status
```
curl -H 'Client-ID: <clientID>' \
-X GET 'https://api.twitch.tv/helix/streams?user_login=<login>'
```
