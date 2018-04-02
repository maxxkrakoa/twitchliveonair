package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var yellowLED = 8
var redLED = 25

var pinYellow, pinRed rpio.Pin

type streams struct {
	Data []data `json:"data"`
}

type data struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Type   string `json:"type"`
}

func main() {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	primaryUserLogin := os.Args[1]

	fmt.Println("Twitch Live On Air")
	fmt.Println(primaryUserLogin)
	fmt.Println(clientID)

	// Unmap gpio memory when done
	defer rpio.Close()

	// init notifications
	notifyInit()

	for {
		// check Twitch status
		primaryUserLive := isPrimaryUserLive(primaryUserLogin, clientID)
		followingUserLive := isFollowingUserLive()

		// update notifications
		notifyPrimaryUserLive(primaryUserLive)
		notifyFollowingUserLive(followingUserLive)

		// take a break
		time.Sleep(time.Second * 10)
	}
}

// TODO: split out into notifications package with three functions:
func notifyInit() {
	if err := rpio.Open(); err != nil {
		fmt.Println("Not running on Raspberry Pi")
	} else {
		// running on Raspberry Pi
		fmt.Println("Running on Raspberry Pi")
		// setup pins, etc
		pinYellow = rpio.Pin(yellowLED)
		pinYellow.Output()
		pinRed = rpio.Pin(redLED)
		pinRed.Output()
	}
}

func notifyPrimaryUserLive(live bool) {
	fmt.Printf("Primary User Live: %t\n", live)
	if pinRed != 0 {
		if live {
			pinRed.High()
		} else {
			pinRed.Low()
		}
	}
}

func notifyFollowingUserLive(live bool) {
	fmt.Printf("Following User Live: %t\n", live)
	if pinYellow != 0 {
		if live {
			pinYellow.High()
		} else {
			pinYellow.Low()
		}
	}
}

func isPrimaryUserLive(userLogin string, clientID string) bool {
	userStreams := getUserStreams(userLogin, clientID)

	liveStream := false
	if len(userStreams.Data) > 0 {
		fmt.Println("Found streams")
		for _, dataItem := range userStreams.Data {
			liveStream = liveStream || (dataItem.Type == "live")
		}
	} else {
		fmt.Println("No streams")
	}

	return liveStream
}

// TODO: actually make this work
func isFollowingUserLive() bool {
	return false

	if rand.Intn(2) == 1 {
		return true
	} else {
		return false
	}
}

func getUserStreams(user string, clientID string) streams {
	url := "https://api.twitch.tv/helix/streams?user_login=" + user
	fmt.Println(url)

	twitchClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Client-ID", clientID)

	res, getErr := twitchClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	streams1 := streams{}
	marshErr := json.Unmarshal(body, &streams1)
	if marshErr != nil {
		fmt.Println(marshErr)
		return streams1
	}

	return streams1
}
