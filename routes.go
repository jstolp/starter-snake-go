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

var rightBound int = 0
var botBound int = 0

var turn int = 0
var move string = "down"
var nextMove string = ""
var prevMove string = ""
var headPos Coord
var tailPos Coord
var health int = 100;
var numOfStartingSnakes int = 1;
var numSnakesLeft int = 1;
var enemySnakes int = 0;
var foodPointList []Coord;

/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func Start(res http.ResponseWriter, req *http.Request) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}

  //foodPointList := decoded.Board.Food
  numOfStartingSnakes = len(decoded.Board.Snakes)
	rightBound = decoded.Board.Width
	botBound = decoded.Board.Height

	log.Print("BOARD Size: " + strconv.Itoa(botBound) + " by " + strconv.Itoa(rightBound))
	log.Println("")
	log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

	if(numOfStartingSnakes == 1) {
		log.Print("It's Gonna be a SOLO GAME \n")
	}
	/*
	 e19c41 - orange test 2
   00ff55 - green
   ff4f00 - orange test 1 -nee te rood
	*/
	respond(res, StartResponse{
		Color: "#e19c41",
		HeadType: "tongue",
		TailType: "curled",
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

func minDistFood(headPos Coord, food []Coord) Coord {
	min := food[0]
	for _, f := range food {
		if dist(min, headPos) < dist(f, headPos) {
			min = f
		}
	}
	return min
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
	foodPointList = decoded.Board.Food
	health = decoded.You.Health
	myLength := len(decoded.You.Body)
	//foodPointList := decoded.Board.Food
	numSnakesLeft = len(decoded.Board.Snakes)
	enemySnakes = numSnakesLeft - 1
	turn = decoded.Turn
	if (enemySnakes < 1) {
		log.Print("SOLO " + strconv.Itoa(turn) + "MY LENGTH: " + strconv.Itoa(myLength) +" h: "+ strconv.Itoa(health) + "\n")
	} else {
		log.Print("TURN " + strconv.Itoa(turn) + " e: "+ strconv.Itoa(enemySnakes)+" h: "+ strconv.Itoa(health) + "\n")
	}

	headPos := getHeadPos(decoded.You)
	nextMoveOOB := isMoveOOB(headPos, nextMove)
	if (nextMoveOOB) {
		nextMove = randomNOOBmove(headPos, move)
	}

/*
	if (health < 30) {
		closestFoodPoint := minDistFood(headPos,foodPointList)
		dd(closestFoodPoint)
		log.Print("Im going to die of starvation in " + strconv.Itoa(health) + " turns \n\n")
		foodDir := goToDir(headPos,closestFoodPoint)

		if(!isMoveOOB(headPos, foodDir)) {
				nextMove = foodDir
		} else {
				nextMove = randomNOOBmove(headPos, move)
		}
		if (isNextMoveFatal(10, headPos, prevMove, nextMove)) {
			// last ditch effort to correct...
			nextMove = invDir(nextMove)
		}
	} // end HEALTH LOW
*/

	move = nextMove // finalise the move
	fmt.Print(strconv.Itoa(turn) + "Move: " + nextMove)
	fmt.Println()
	respond(res, MoveResponse{
		Move: nextMove,
	})
	prevMove = nextMove // Re-allocate move command to prev/last move\
}


func isNextMoveFatal(health int, headPos Coord, currentDir string, targetDir string) bool {
    // doing a 180 is never safe, so check for that...

		flipDir := invDir(currentDir)
		if(flipDir == targetDir) {
			log.Print("The move is " + targetDir + "but in going " + currentDir + "That would be fatal...\n")
			return true
		}
		// check if a move is NOT_OUT_OF_BOUNDS (hit a wall) WALL SNAKE
		if (isMoveOOB(headPos, targetDir)) {
			log.Print("Next Move is Fatal because of a BOUNDARY " + targetDir + "\n")
			return true
		}

		if (health == 1) {
			log.Print("Dag Mooie Wereld... Hongersnood is geen grapje... \n\n")
			return true
		}

		log.Print("The move " + targetDir + " is safe...\n")
		return false
}


// see if i can attach these methods to the struct Snake or something..
func getHeadPos(target Snake) Coord {
	body := target.Body
  return body[0]
}
func getTailPos(target Snake) Coord {
	body := target.Body
  return body[len(body) - 1]
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
func goToDir(curr Coord, next Coord) string {
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

// just a testing function to dump a object../
func dd(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err == nil {
		log.Printf(string(data))
	}
}

// Extra route
func Index(res http.ResponseWriter, req *http.Request) {
	/* Battlesnake documentation can be found at <a href=\"https://docs.battlesnake.io\">https://docs.battlesnake.io</a>. */
	configuration := Configuration{}
	errConf := gonfig.GetConf("config/config.json", &configuration)
	if errConf != nil {
		log.Printf("Bad configuration in config.json: %v", errConf)
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Jay's battleSnake mk 2 self aware: " + configuration.Home_Route))
}

func Ping(res http.ResponseWriter, req *http.Request) {
	log.Print("PONG to a server ping... \n")
	return
}

func End(res http.ResponseWriter, req *http.Request) {
	log.Print("The game has ended.... \n\n")
	return
}
