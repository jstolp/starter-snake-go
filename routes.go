package main

import (
	"log"
	"encoding/json"
	"net/http"
	. "github.com/jstolp/pofadder-go/api"
	"fmt"
	"math"
	"math/rand"
	"time"
	"strconv"
	"strings"
)

type CoordList []Coord // remove if api is imported correctly (WHY WINDOWS?!?

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
  numOfStartingSnakes = len(decoded.Board.Snakes)
	log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

	fmt.Print("Start Pos: " + strconv.Itoa(headPos.X) + "," + strconv.Itoa(headPos.Y))


	if (numOfStartingSnakes == 1) {
		log.Print("\n\n It's Gonna be a SOLO GAME \n")
		HUNGRY_TRESHOLD = 90
	}

	respond(res, StartResponse{
		Color: "#000000",
		HeadType: "evil",
		TailType: "sharp",
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
	//validMoves := len(getPossibleMoves(decoded))
	//numberOfSnakes := len(decoded.Board.Snakes)
	//foodList := decoded.Board.Food

	//log.Print("ESCAPE: around my head: " + strconv.Itoa(countEscapeRoutesFromCoord(headPos, decoded)))



	//enemyHead :=
	if (health > 20) {

		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, getEnemyHeadPos(decoded))
		if moveCoord == nil {
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
		} else {
			if (numOfStartingSnakes == 1) {
				// SOLO game... moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
				moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
			} else {
				log.Print("ATTACK!!! HEAD HUNTING!")
				nextMove = Heading(headPos, moveCoord[1])

			}

		}
	} else if (len(decoded.Board.Food) == 0) && len(decoded.You.Body) >= 4 {
		// NO FOOD... Bigger than 4 BodyParts,  No food on the board
		log.Print("no food on board... chasing tail...")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
	} else if (health < 20) {
		// THERE IS FOOD, under HUNGRY_TRESHOLD
		log.Print("Hunting for food! I'm below HUNGRY_TRESHOLD")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, SafeFoodHead(decoded))
		if (dist(headPos, tailPos) == 1) {
			//log.Print("Grabbing Food Close to TAIL!")
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, SafeFoodTail(decoded))
		}

		if moveCoord == nil {
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
		}
	}



	if (moveCoord == nil || len(moveCoord) < 1) {
		nextMove = getRandomValidMove(decoded)
		log.Print("Used random valid Move: " + nextMove)
	} else {
		nextMove = Heading(headPos, moveCoord[1])
	}


	fmt.Print("T " + strconv.Itoa(turn) + " Health:" + strconv.Itoa(health) + " Move: " + nextMove + "\n")

	respond(res, MoveResponse{
		Move: nextMove,
	})
}



func getEnemyHeadPos(game SnakeRequest) []Coord {
		coordList := make(Coord[], 0)
		snakeList := game.Board.Snakes
		 for i := 0; i < len(snakeList); i++ {
				 if ( snakeList[i].ID != game.You.ID && len(snakeList[i].Body) <= len(game.You.Body) ) {
					 coordList = append(coordList, snakeList[i].Body[0])
				 }
		 }
		 return coordList
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

func SafeFoodTail(req SnakeRequest) Coord {
	You := getTailPos(req.You)
	foodArray := req.Board.Food
	var safeFood = foodArray[0]
	var safeFoodF = Dist(foodArray[0], You)

	for i := 0; i < len(foodArray); i++ {
		if Dist(foodArray[i], You) < safeFoodF {
			if (countEscapeRoutesFromCoord(foodArray[i], req) > 1) {
				safeFood = foodArray[i]
				safeFoodF = Dist(foodArray[i], You)
			} else {
				log.Print("TAIL food was not safe... skipping");
			}
		}
	}

	return safeFood
}
//safeClosestFood
func SafeFoodHead(req SnakeRequest) Coord {
	You := req.You.Body[0]
	foodArray := req.Board.Food
	var safeFood = foodArray[0]
	var safeFoodF = Dist(foodArray[0], You)

	for i := 0; i < len(foodArray); i++ {
		if Dist(foodArray[i], You) < safeFoodF {

			if (!isSafe(foodArray[i], req) || countEscapeRoutesFromCoord(foodArray[i], req) > 1){
				safeFood = foodArray[i]
				safeFoodF = Dist(foodArray[i], You)
			} else {
				log.Print("HEAD food was not safe... skipping");
			}

		}
	}

	return safeFood
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

func newPossibleMoves(game SnakeRequest) string {
	allCoords := SurroundingCoordinates(game.You.Body[0])
	for _,coord := range allCoords {
			if (isSafe(coord, game)) {
				return goToDir(game.You.Body[0], coord)
			}
		}
		return "no"
}
/**
returns list of coords that are possible
*/
func getPossibleMoves(game SnakeRequest) []Coord {
	var validCoords = make([]Coord, 0)

	headPos := game.You.Body[0]
	//tailPos := getTailPos(game.You)
	enemySnakes := game.Board.Snakes
	allCoords := getOpenAjdacentNodes(headPos)
	for _,coord := range allCoords {
		dir := Heading(headPos, coord)
		if (false == isMoveOOB(headPos, dir)) {
			if (false == NodeBlockedExceptTail(coord, enemySnakes)) {
				validCoords = append(validCoords,coord)
			}
		}
  }
		return validCoords
}

func getRandomValidMove(game SnakeRequest) string {
	headPos := game.You.Body[0]
	//tailPos := getTailPos(game.You)
	enemySnakes := game.Board.Snakes
	allCoords := getOpenAjdacentNodes(headPos)
	nextDir := "invalid"


	for _,coord := range allCoords {
			if (isSafe(coord, game) && false == NodeBlocked(coord, enemySnakes)) {
				// is isSafe correct U
				log.Print("Got the Most Optimal Route in random move...")
				return Heading(headPos, coord)
				// if we have a safe Move.. return that one, else... each one is as bad af them..
			}

			if (false == NodeBlocked(coord, enemySnakes)) {
				dir := Heading(headPos, coord)
				if(game.You.Health > 99 && dist(headPos, tailPos) == 1 && dir == goToDir(headPos, tailPos)) {
					log.Print("skipping " + dir + " as it would crash into tail")
					continue
				}

				if (false == isMoveOOB(headPos, dir)) {
					log.Println("false is move OOB")
					// if nothing better aries, this is the one...
					nextDir = dir
				}
		}
	}

	if (nextDir != "invalid") {
		// if the loop didn't find a safe coordinate... use the next best thing...
		return nextDir
	}

	log.Print("INVALID MOVE IN: getRandomValidMove")
	escapeRoutes := countEscapeRoutesFromCoord(headPos, game)
	//dd(GetAdjacentCoords(headPos))
	log.Print("escape routes: " + strconv.Itoa(escapeRoutes))
	log.Print("new possbile: " + newPossibleMoves(game))
	return newPossibleMoves(game)
	//return "invalid"
 	//return getAnyMove(game) // invalid move
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
			Coord{boardWidth-2,1},
			Coord{1, boardHeight - 2},

			Coord{boardHeight - 2, boardWidth -1},
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

func isNodeOnBoard(target Coord) bool {
	if target.X < 0 || target.Y < 0 { // TOP LEFT CORDER NE
		return false
	}
	if target.X > boardWidth - 1 || target.Y > boardHeight - 1 { // OOB Protection
		return false
	}
	return true
}

// Shuffle... For use in find random direction

func shuffle(src []string) []string {
  final := make([]string, len(src))
  rand.Seed(time.Now().UTC().UnixNano())
  perm := rand.Perm(len(src))

  for i, v := range perm {
 final[v] = src[i]
  }
  return final
 }


	func isSafe(point Coord, game SnakeRequest) bool {
		//if (!isNodeOnBoard(point)) {
	//		return false
	//	}

		snakeList := game.Board.Snakes
		 for i := 0; i < len(snakeList); i++ {
			 for j := 0; j < len(snakeList[i].Body); j++ {
				 if ( snakeList[i].ID != game.You.ID && len(snakeList[i].Body) >= len(game.You.Body)) {
					 // this snake is longer than you, avoid it's next steps if possible!
					 for _, snakeCoord := range GetAdjacentCoords(snakeList[i].Body[0]) {
						 if (point.X == snakeCoord.X && point.Y == snakeCoord.Y) {
							 // if the next move is a possibleOne for a enemy... it's not safe
							 return false
						 }
					 } // end for "adjacentCoords" for the enemy snake coords...

				 }

				 if snakeList[i].Body[j].X == point.X && snakeList[i].Body[j].Y == point.Y {
					 if (len(snakeList[i].Body)-1 == j && snakeList[i].Health < 99) {
						 log.Print("tail is safe to step...")
						 return true // this is the tail... YES
					 } else {
						 // snake just ate... tail is fatal!
						 return false
					 }

					 return false
				 }
			 }
		 }
		 return true // free node
	}


 func isFree(point Coord, req SnakeRequest) bool {
 snakeList := req.Board.Snakes
 	for i := 0; i < len(snakeList); i++ {
 		for j := 0; j < len(snakeList[i].Body); j++ {
 			if snakeList[i].Body[j].X == point.X && snakeList[i].Body[j].Y == point.Y {
				if (len(snakeList[i].Body)-1 == j && snakeList[i].Health < 99) {
					log.Print("tail is safe to step...")
 					return true // this is the tail... YES
 				} else {
					// snake just ate... tail is fatal!
					return false
				}

 				return false
 			}
 		}
 	}
 	return true // free node
 }


func countEscapeRoutesFromCoord(search Coord, req SnakeRequest) int {
	i := 0

	if (isNodeOnBoard(Coord{X: search.X + 1, Y: search.Y})) {
		if (isFree(Coord{X: search.X + 1, Y: search.Y}, req)) {
			i++
		}
	}

	if (isNodeOnBoard(Coord{X: search.X - 1, Y: search.Y})) {
			if (isFree(Coord{X: search.X - 1, Y: search.Y}, req)) {
				i++
		  }
	}

	if (isNodeOnBoard(Coord{X: search.X, Y: search.Y + 1})) {
		if (isFree(Coord{X: search.X, Y: search.Y + 1}, req)) {
			i++
		}
	}

	if (isNodeOnBoard(Coord{X: search.X, Y: search.Y - 1})) {
		if (isFree(Coord{X: search.X, Y: search.Y - 1}, req)) {
			i++
		}
	}
	return i
}

func getOpenAjdacentNodes(search Coord) []Coord {
	var validCoords = make([]Coord, 0)

	if (isNodeOnBoard(Coord{X: search.X + 1, Y: search.Y})) {
		validCoords = append(validCoords, Coord{X: search.X + 1, Y: search.Y})
	}

	if (isNodeOnBoard(Coord{X: search.X - 1, Y: search.Y})) {
		validCoords = append(validCoords, Coord{X: search.X - 1, Y: search.Y})
	}

	if (isNodeOnBoard(Coord{X: search.X, Y: search.Y + 1})) {
		validCoords = append(validCoords, Coord{X: search.X, Y: search.Y + 1})
	}
	if (isNodeOnBoard(Coord{X: search.X, Y: search.Y - 1})) {
		validCoords = append(validCoords, Coord{X: search.X, Y: search.Y - 1})
	}

	return validCoords
}


func SurroundingCoordinates(search Coord) []Coord {
	return []Coord{
		{search.X + 1, search.Y + 0},
		{search.X - 1, search.Y + 0},
		{search.X + 0, search.Y + 1},
		{search.X + 0, search.Y - 1},
	}
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

/* Dist to function in steps (int) */
func dist(a Coord, b Coord) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

// just a testing function to  a object../
func dd(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err == nil {
		log.Printf(string(data))
	}
}

func minifyPrint(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "")
	if err == nil {
		fmt.Println(strings.Replace(string(data), " ", "", -1))
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

func mapToGrid(decoded SnakeRequest) ([][]string) {

 // use decoded.Board.BoardHeight x decode.Board.BoardWidth for grid!
  grid := make([][]string, decoded.Board.Height)
  me := decoded.You
  foodList := decoded.Board.Food


  for i := 0; i < len(grid); i++ {
      grid[i] = make([]string, decoded.Board.Width)
  }

  for i := 0; i < decoded.Board.Height; i++ {
     grid[0][i] = "."

     grid[i][0] = "."

     grid[decoded.Board.Height-1][i] = "."

     grid[i][decoded.Board.Width-1] = "."
 }

otherSnakes := decoded.Board.Snakes

for _, snake := range otherSnakes {
  for i, coord := range snake.Body {
    c := coord.X
    r := coord.Y

    if grid[r][c] != "#" {
      if(i == 0) {
        grid[r][c] = "h"
      } else if(i == len(snake.Body) - 1) {
        grid[r][c] = "t"
      } else {
        grid[r][c] = "+"
      }
     }
  }
}

/**
 * H -> Head
 * h -> enemy head
 * T -> Tail
 * T -> enemy tail
 * ! -> Food
 * # -> Wall
 * * snakeBody
 * + enemySnake Body
 * $ - Target
 */
  for _, coord := range foodList {
     c := coord.X
     r := coord.Y

     if grid[r][c] != "#" {
        grid[r][c] = "!"
      }
  }

 myBody := me.Body;
 for _, coord := range myBody {
    c := coord.X
    r := coord.Y

    if (grid[r][c] != "#") {
       grid[r][c] = "*"
     }
 }

 headPos := getHeadPos(me)

 c := headPos.X
 r := headPos.Y

 if grid[r][c] != "#" {
    grid[r][c] = "H"
  }

  tailPos := getTailPos(me)

  c = tailPos.X
  r = tailPos.Y

  if grid[r][c] != "#" {
     grid[r][c] = "T"
   }
	 PrintGrid(grid)
 return grid;
}

func PrintGrid(grid [][]string) {
	fmt.Print("|")
	for i := 0; i < len(grid); i++ {
			fmt.Print(strconv.Itoa(i))
		}
		fmt.Print("|\n")
    for i := 0; i < len(grid); i++ {
				fmt.Print("|")
        for j := 0; j < len(grid[0]); j++ {

            if grid[i][j] == "" {
                fmt.Printf(".")
            } else {
                fmt.Print(grid[i][j])
            }

        }
        fmt.Print("|"+ strconv.Itoa(i)+"\n")
    }
    fmt.Print("\n")
}
