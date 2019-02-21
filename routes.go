package main

import (
	"log"
	"net/http"
	"github.com/jstolp/pofadder-go/api"
	"math/rand"
	"fmt"
	"strconv"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}


func Index(res http.ResponseWriter, req *http.Request) {
	/* Battlesnake documentation can be found at <a href=\"https://docs.battlesnake.io\">https://docs.battlesnake.io</a>. */
	/* configuration := api.Configuration{}
	errConf := gonfig.GetConf("config/config.json", &configuration)
	if errConf != nil {
		log.Printf("Bad configuration in config.json: %v", errConf)
	}*/
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Jay's battleSnake mk 1 self aware"))
}

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
		Color: "#fefefe",
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
	dump(decoded.You)
	dump(decoded.You.body[0])
	respond(res, api.MoveResponse{
		Move: "down",
	})
}
/*
var prevDir := "na"
var currentDir := "na"
var currentPos := api.Coord{}
*/

/* Dist to function in steps (int) */
func dist(a Coord, b Coord) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

/* move from coord to coord -> returns MOVE */
func GoToDir(curr Coord, next Coord) string {
	dir := ""
	if curr.X < next.X {
		dir = "right"
	} else if curr.X > next.X {
		dir = "left"
	} else if curr.Y < next.Y {
		dir = "down"
	} else if curr.Y > next.Y {
		dir = "up"
	}
	return dir
}
/* Inverses direction */
func invDir(currentDir String) string {
		if(currentDir == "down") {
			return "up"
		}
		if(currentDir == "up") {
			return "down"
		}
		if(currentDir == "left") {
			return "right"
		}
		if(currentDir == "right") {
			return "left"
		}
}