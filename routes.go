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

	log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

	fmt.Print("Start Pos: " + strconv.Itoa(headPos.X) + "," + strconv.Itoa(headPos.Y))

	//log.Print("Should Be 2... " + strconv.Itoa( countOpenAjdacents(Coord{0,0}))
 //log.Print("would be 2: " + strconv.Itoa(countOpenAjdacents(countOpenAjdacents(topLeft))))
 //log.Print("would be 4: " + strconv.Itoa(countOpenAjdacents(countOpenAjdacents(Coord{1,1}))))
 //log.Print("would be 3: " + strconv.Itoa(countOpenAjdacents(countOpenAjdacents(Coord{4,0}))))


	if (numOfStartingSnakes == 1) {
		log.Print("\n\n It's Gonna be a SOLO GAME \n")
		HUNGRY_TRESHOLD = 90
	}
	/*
	 e19c41 - orange test 2
   00ff55 - green
   ff4f00 - orange test 1 -nee te rood
	*/
	respond(res, StartResponse{
		Color: "#ff00aa",
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
	validMoves := len(getPossibleMoves(decoded))
	//numberOfSnakes := len(decoded.Board.Snakes)
	//foodList := decoded.Board.Food

	//log.Print("ESCAPE: around my head: " + strconv.Itoa(countEscapeRoutesFromCoord(headPos, decoded)))
/*
	if (health > 98) {
		log.Print("goin to closestCorner")
		targetCorner := closestCorner(boardHeight, boardWidth, headPos)
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, targetCorner)
		if moveCoord == nil {
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
		}
	} else if (len(decoded.Board.Food) == 0) && len(decoded.You.Body) >= 4 {
		// NO FOOD... Bigger than 4 BodyParts,  No food on the board
		log.Print("no food on board... chasing tail...")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
	} else if (health < HUNGRY_TRESHOLD) {
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
	} else {

		targetCorner := closestCorner(boardHeight, boardWidth, headPos)
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, targetCorner)
		if moveCoord == nil {
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
		}
		// Long En
	}



	if (moveCoord == nil || len(moveCoord) < 1) {
		nextMove = getRandomValidMove(decoded)
		log.Print("Used random valid Move: " + nextMove)
	} else {
		nextMove = Heading(headPos, moveCoord[1])
	}

	if(isMoveOOB(headPos, nextMove)) {
		log.Print("NEXT move is OOB detection + next:Move" + nextMove)
	}

	if (health > 99 && dist(headPos, tailPos) == 1 && nextMove == goToDir(headPos, tailPos)) {
		// select a different move, as i'm heading For my Own Tail...
		nextMove = getRandomValidMove(decoded)
		//dd(decoded)
	}

	if (nextMove == "invalid") {
		log.Print("Turn "+ strconv.Itoa(turn) + " is my last... Dag mooie wereld!")
	}


	// Check if we still have a path to tail... if not.... let's switch tactics:
	if (nil == Astar(boardHeight, boardWidth, me, enemySnakes, getTailPos(me))) {
		// my Tail is not reachable by shortest path!
		log.Print("Switch Tactic to LONGEST PATH!!!!")

		//dd(decoded)
	}

	*/

	if (validMoves == 0) {
		// easy let's move that way
		log.Print("I'm dead this round...")
	}

	if (validMoves == 1) {
		// easy let's move that way
		log.Print("Only one move possible")
	}

	if(validMoves == 2) {
		log.Print("2 Valid moves... let's decide!")
		allMoves := getPossibleMoves(decoded)
		for _,checkMove := range allMoves {
			dir := goToDir(headPos, checkMove)
			if (nil == AStarBoardFromTo(decoded, checkMove, tailPos)) {
				// no path from next valid move to tail...
				log.Print(dir + " NO PATH to tail")
			} else {
				log.Print(dir + " has PATH to my tail")
			}
		}
	}

	if (health < 20) {
		//log.Print("HUNTING AFTER LONGEST PATH!!!")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, SafeFoodHead(decoded))
		if (nil == moveCoord) {
			moveCoord = AstarBoard(decoded, tailPos)
		}
		nextMove = Heading(headPos, moveCoord[1])
		if (isNodeOnBoard(moveCoord[1]) && isFree(moveCoord[1], decoded)) {
			//log.Print("FoodMove is safe")
		} else {
			log.Print("Food is FATAL... other move")
			nextMove = getRandomValidMove(decoded)
		}
	} else {
		if (nil != AstarBoard(decoded, tailPos)) {
			moveCoord = LongestPath(decoded, tailPos)
		}

		if (nil == moveCoord) {
			log.Print("LONGEST PATH to TAIL NOT FOUND... DEAD?")
			nextMove = getRandomValidMove(decoded)
			log.Print("LONGEST PATH fatal... get random")
		} else {
			nextMove = Heading(headPos, moveCoord[1])
			if (isNodeOnBoard(moveCoord[1]) && isFree(moveCoord[1], decoded)) {
				//log.Print("longest path is SAFE")
			} else {
				nextMove = getRandomValidMove(decoded)
				log.Print("LONGEST PATH fatal... got random move")
			}
		}
   }

	//mapToGrid(decoded)
	//minifyPrint(decoded)


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
			if (countEscapeRoutesFromCoord(foodArray[i], req) > 1){
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
	for _,coord := range allCoords {
		if (false == NodeBlocked(coord, enemySnakes)) {
			dir := Heading(headPos, coord)
			if(game.You.Health > 99 && dist(headPos, tailPos) == 1 && dir == goToDir(headPos, tailPos)) {
				log.Print("skipping " + dir + " as it would crash into tail")
				continue
			}

			if (false == isMoveOOB(headPos, dir)) {
				return dir
			}
		}
	}

	log.Print("INVALID MOVE IN: getRandomValidMove")
 	return "invalid" // invalid move
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


 func isFree(point Coord, req SnakeRequest) bool {
 snakeList := req.Board.Snakes
 	for i := 0; i < len(snakeList); i++ {
 		for j := 0; j < len(snakeList[i].Body); j++ {
 			if snakeList[i].Body[j].X == point.X && snakeList[i].Body[j].Y == point.Y {
				if (len(snakeList[i].Body)-1 == j && snakeList[i].Health < 99) {
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
