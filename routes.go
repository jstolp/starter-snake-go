package main

import (
	"log"
	"encoding/json"
	"net/http"
	. "github.com/jstolp/pofadder-go/api"
	"github.com/tkanos/gonfig"
	"fmt"
	"math"
	"strconv"
)

var topBound int = 1
var leftBound int = 1
var rightBound int = 0
var botBound int = 0

func dd(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err == nil {
		log.Printf(string(data))
	}
}

/* returns SnakeHeadPos COORD*/
func getHeadPos(target Snake) Coord {
	body := target.Body

	bodySlice := make([]Coord, len(body))
	//dd(bodySlice)

  return (bodySlice[len(bodySlice)-1])
}

func getTailPos(target Snake) Coord {
	body := target.Body

	bodySlice := make([]Coord, len(body))
	dd(bodySlice)

  return (bodySlice[0])
}

/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func Start(res http.ResponseWriter, req *http.Request) {


	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Print("START: start")
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}

	rightBound = decoded.Board.Width
	botBound = decoded.Board.Height

	log.Print("BOARD: top: " + strconv.Itoa(topBound) + " bot: " + strconv.Itoa(botBound) + "left: " + strconv.Itoa(leftBound) + " right " + strconv.Itoa(rightBound))

	respond(res, StartResponse{
		Color: "#fefefe",
		HeadType: "fang",
		TailType: "bolt",
	})
	log.Print("START: end \n")
}

func Move(res http.ResponseWriter, req *http.Request) {
	log.Print("MOVE: start \n")
	move := "down"
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	headPos := getHeadPos(decoded.You)
	tailPos := getTailPos(decoded.You)
	dd(headPos)
	dd(tailPos)
	fmt.Print("Going ... " + move + "\n")
	respond(res, MoveResponse{
		Move: "down",
	})
	log.Print("MOVE: end \n\n")
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
		dir := ""
		if(currentDir == "down") {
				dir = "up"
		}
		if(currentDir == "up") {
			dir = "down"
		}
		if(currentDir == "left") {
			dir = "right"
		}
		if(currentDir == "right") {
			dir = "left"
		}
		return dir
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

func Ping(res http.ResponseWriter, req *http.Request) {
	log.Print("PONG to a server ping... \n")
	return
}

func End(res http.ResponseWriter, req *http.Request) {
	log.Print("The game has ended.... \n\n")
	return
}
