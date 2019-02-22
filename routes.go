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

var turn int = 0
var move string = "down"
var nextMove string = ""
var prevMove string = ""
var headPos Coord
var tailPos Coord


func dd(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err == nil {
		log.Printf(string(data))
	}
}

func getTailPos(target Snake) Coord {
	body := target.Body
  return body[len(body) - 1]
}

func getHeadPos(target Snake) Coord {
	body := target.Body
  return body[0]
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
}

func isMoveOOB(headPos Coord, direction string) bool {
	switch direction {
		case "down":
			if (headPos.Y + 1 < botBound) {
					return false
			}
		case "up":
			if (headPos.Y > 0) {
				return false
			}
		case "left":
			if(headPos.X + 1 > 1) {
				return false
			}
		case "right":
			if(headPos.X + 1 < rightBound) {
				return false
			}
	}
return true
}

func randomNOOBmove(headPos Coord, currentDir string) string {

  //randomInt = rand.Intn(100)
	switch currentDir {
		case "down":
				return "right"
		case "up":
				return "left"
		case "left":
				return "down"
		case "right":
				return "up"
	}
	return "down"
}

func Move(res http.ResponseWriter, req *http.Request) {
	nextMove = prevMove
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	turn = decoded.Turn
	log.Print("TURN " + strconv.Itoa(turn) + "\n")

	headPos := getHeadPos(decoded.You)
	nextMoveOOB := isMoveOOB(headPos, nextMove)
	if (nextMoveOOB) {
		nextMove = randomNOOBmove(headPos, move)
	}

	move = nextMove // finalise the move
	fmt.Print("Move: " + move)
	fmt.Println()
	respond(res, MoveResponse{
		Move: move,
	})
	prevMove = move // Re-allocate move command to prev/last move\
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
