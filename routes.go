package main

import (
	"log"
	"encoding/json"
	"net/http"
	. "github.com/jstolp/pofadder-go/api"
	. "github.com/jstolp/pofadder-go/astar"
	"fmt"
	"math"
	"strconv"
)

type CoordList []Coord // remove if api is imported correctly

var edgeSnakeLimit int = 0;
var turn int = 0;
var move string = "down";
var nextMove string = "";
var prevMove string = "";
var headPos Coord;
var tailPos Coord;
var health int = 100;
var numOfStartingSnakes int = 1;
var numSnakesLeft int = 1;
var enemySnakes int = 0;
var foodPointList []Coord;
var endCicle bool = false;
var boardHeight int = 0;
var boardWidth int = 0;

/* SNAKE SETUP */
var HUNGRY_TRESHOLD  int = 90; // defines when snake goes looking for food.


/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func Start(res http.ResponseWriter, req *http.Request) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}

	headPos = decoded.You.Body[0]
	boardHeight, boardWidth = decoded.Board.Height, decoded.Board.Width // SE corner X, Y
	
	log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

	fmt.Print("Start Pos: " + strconv.Itoa(headPos.X) + "," + strconv.Itoa(headPos.Y))
	
	
	if(numOfStartingSnakes == 1) {
		log.Print("\n\n It's Gonna be a SOLO GAME \n")
	}
	/*
	 e19c41 - orange test 2
   00ff55 - green
   ff4f00 - orange test 1 -nee te rood
	*/
	respond(res, StartResponse{
		Color: "#e19c41",
		HeadType: "sand-worm",
		TailType: "curled",
	})
}


func Move(res http.ResponseWriter, req *http.Request) {
	nextMove = prevMove
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}

	var moveCoord []Coord
	turn = decoded.Turn
	me := decoded.You
	health = me.Health
	headPos = decoded.You.Body[0]
	tailPos = getTailPos(me)
	enemySnakes := decoded.Board.Snakes
	foodList := decoded.Board.Food
	
	if (len(decoded.Board.Food) == 0) && len(decoded.You.Body) >= 4 {
		// no food? chase tail. but only if i'm big enough
		log.Print("no food on board... chasing tail...")
		moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
	} else if (health > HUNGRY_TRESHOLD) &&  len(decoded.You.Body) >= 4 {
		// there is food but i'm not hungry, still chase tail only if i'm big enough.
		log.Print("chasing tail")
		moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
	} else if (health > HUNGRY_TRESHOLD) {
		log.Print("not big enough... moving to corner")
		targetCorner := closestCorner(boardHeight, boardWidth, headPos)
		moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, targetCorner)
		if moveCoord == nil {
			moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, astar.ChaseTail(me.Body))
		}
	} else {
		moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, astar.NearestFood(foodList, headPos))
		if moveCoord == nil {
			moveCoord = astar.Astar(boardHeight, boardWidth, me, enemySnakes, astar.ChaseTail(me.Body))
		}
	}

	nextMove = astar.Heading(headPos, moveCoord[1])

	closestCorner(boardHeight, boardWidth, headPos)
	
	fmt.Print("T " + strconv.Itoa(turn) + " Health:" + strconv.Itoa(health) + " Move: " + nextMove + "\n")
	respond(res, MoveResponse{
		Move: nextMove,
	})
}

// see if i can attach these methods to the struct Snake or something..
// func (target Snake) Head() Coord { return target.Body[0] }

func getHeadPos(target Snake) Coord {
	body := target.Body
  return body[0]
}
func getTailPos(target Snake) Coord {
	body := target.Body
  return body[len(body) - 1]
}

func isNextMoveFatal(me Snake, currentDir string, targetDir string) bool {
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

		// if dist to my own tail is 1, and i'm going in the same direction...
		// i'll die...
		if (dist(headPos, tailPos) == 1 && targetDir == goToDir(headPos, tailPos)) {
			log.Print("CRASHING INTO MY OWN TAIL IN ... 3 . 2.. .1.. no... next MOVE ahhaah \n\n")
			log.Print()
			return true
		}

		//log.Print("The move " + targetDir + " is safe...\n")
		return false
}

func isMoveOOB(headPos Coord, direction string) bool {
	switch direction {
		case "down":
			if (headPos.Y + 1 < boardHeight) {
					return false
			}
		case "up":
			if (headPos.Y + 1 > 1) {
				return false
			}
		case "left":
			if (headPos.X + 1 > 1) {
				return false
			}
		case "right":
			if (headPos.X + 1 < boardWidth) {
				return false
			}
	}
	return true
}
/* Returns the closestCorner */
func closestCorner(boardHeight int, boardWidth int, headPos Coord) Coord {
		distToCorner := -1
		targetCoord := Coord{0,0}
		corners := CoordList {
			Coord{1,1},
			Coord{boardWidth-1,1},
			Coord{1, boardHeight - 1},
			Coord{boardHeight - 1, boardWidth -1},
		}

		for _, coord := range corners {
			
			if( dist(headPos, coord) > distToCorner) {
				distToCorner = dist(headPos, coord)
				targetCoord = coord
			}
		} // end for "valid" coords...


		//log.Print("closestCord is: " + strconv.Itoa(targetCoord.X) + "," + strconv.Itoa(targetCoord.Y))
		return targetCoord	
}

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

/* Dist to function in steps (int) */
func dist(a Coord, b Coord) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

// closestFoodPoint
func minDistFood(headPos Coord, food []Coord) Coord {
	min := food[0]
	for _, f := range food {
		if dist(min, headPos) < dist(f, headPos) {
			min = f
		}
	}
	return min
}

// just a testing function to dump a object../
func dd(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err == nil {
		log.Printf(string(data))
	}
}

func Ping(res http.ResponseWriter, req *http.Request) {
	log.Print("PONG to a server ping... \n")
	return
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func End(res http.ResponseWriter, req *http.Request) {
	log.Print("The game has ended.... \n\n")
	return
}
