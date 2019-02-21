package main

import (
	"log"
	"net/http"
	"github.com/jstolp/pofadder-go/api"
	"github.com/tkanos/gonfig"
	"math/rand"
	"fmt"
)

func Info(res http.ResponseWriter, req *http.Request) {
	configuration := api.Configuration{}
	errConf := gonfig.GetConf("config/config.json", &configuration)
	if errConf != nil {
		log.Printf("Bad configuration in config.json: %v", errConf)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Jay's battleSnake " + configuration.Home_Route))
	res.Write([]byte("Jay's battleSnake int: " + rand.Intn(100)))
	res.Write([]byte("Jay's battleSnake " + configuration.Home_Route))
	fmt.Print(rand.Intn(100))
 	fmt.Println()
}

func Index(res http.ResponseWriter, req *http.Request) {
	configuration := api.Configuration{}
	errConf := gonfig.GetConf("config/config.json", &configuration)
	if errConf != nil {
		log.Printf("Bad configuration in config.json: %v", errConf)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Jay's battleSnake " + configuration.Home_Route))
	fmt.Print(rand.Intn(100))
 	fmt.Println()
}
/* Battlesnake documentation can be found at <a href=\"https://docs.battlesnake.io\">https://docs.battlesnake.io</a>. */

/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func Start(res http.ResponseWriter, req *http.Request) {
	decoded := api.SnakeRequest{}
	err := api.DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}
	dump(decoded)

	respond(res, api.StartResponse{
		Color: "#ffffff",
		HeadType: "fang",
		TailType: "bolt",
	})

}

func Move(res http.ResponseWriter, req *http.Request) {
	decoded := api.SnakeRequest{}
	err := api.DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	dump(decoded)

	respond(res, api.MoveResponse{
		Move: "down",
	})
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}
