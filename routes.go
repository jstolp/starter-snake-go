package main

import (
	"log"
	"net/http"
	. "github.com/jstolp/pofadder-go/api"
	"github.com/tkanos/gonfig"
	"fmt"
	"math"
	/* "strconv" */
)

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}


func Index(res http.ResponseWriter, req *http.Request) {
	/* Battlesnake documentation can be found at <a href=\"https://docs.battlesnake.io\">https://docs.battlesnake.io</a>. */
	configuration := Configuration{}
	errConf := gonfig.GetConf("config/config.json", &configuration)
	if errConf != nil {
		log.Printf("Bad configuration in config.json: %v", errConf)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Jay's battleSnake mk 1 self aware: " + configuration.Home_Route))
}

/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func Start(res http.ResponseWriter, req *http.Request) {
	log.Print("START: start")
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}
	dump(decoded)

	respond(res, StartResponse{
		Color: "#fefefe",
		HeadType: "fang",
		TailType: "bolt",
	})
	log.Print("START: end")
}

func Move(res http.ResponseWriter, req *http.Request) {
	log.Print("MOVE: start")
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	dump(decoded.You)
	dump(decoded.You.Body[0])
	fmt.Print("Going down...")
	respond(res, MoveResponse{
		Move: "down",
	})
	log.Print("MOVE: end")
}
/*
var prevDir := "na"
var currentDir := "na"
var currentPos := Coord{}
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

/*  sl[len(sl)-1] READ LAST SLICE
sl = sl[:len(sl)-1] RM last SLICE
https://github.com/golang/go/wiki/SliceTricks
*/

/* Inverses direction */
func invDir(currentDir string) string {
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
		return "invalid"
}
