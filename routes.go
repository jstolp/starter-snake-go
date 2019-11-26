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
var HUNGRY_TRESHOLD int = 40; // defines when snake goes looking for food.

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

	log.Print("Should Be 2... " + strconv.Itoa(len(getOpenAjdacentNodes(Coord{0,0}))))


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
	enemySnakes = len(decoded.Board.Snakes) - 1
	enemyHeadPosList := getEnemyHeadPos(decoded)
	validMoves := len(getPossibleMoves(decoded))

	if (validMoves == 0) {
		// DEAD!
		log.Print("I'm dead this round... 0 valid moves")
	}

	// above hungry? above  chaise tail! (Set this low maybe? i used to do 90... but FOOD == DANGER)
	if (decoded.You.Health > HUNGRY_TRESHOLD || !iAmTheBiggestSnakeAlive(decoded) )  && ( len(decoded.You.Body) >= 4 || len(decoded.Board.Food) == 0 ) {

			if ( iAmTheBiggestSnakeAlive(decoded) && enemySnakes > 0 && !isThereSafeFood(decoded)) {
				log.Print("Biggest Snake Alive... and there are enemies! And No Easy Food")
				if ( nil != AstarBoard(decoded, enemyHeadPosList[0]) ) {
					log.Print("In for the kill...")
					moveCoord = AstarBoard(decoded, enemyHeadPosList[0])
				} else {
					log.Print("Wanted to kill... but...")
					moveCoord = AstarBoard(decoded, tailPos)
				}

			} else {
				// FOOD!!!
					if (isThereSafeFood(decoded)) {
						moveCoord = AstarBoard(decoded, SafeFoodHead(decoded))
					} else {
					log.Print("no SAFE ... food")
					// NO FOOD...
					if (len(decoded.You.Body) <= 4) {
						// too small? WALK! LONG PATH
							moveCoord = LongestPath(decoded, tailPos)
					} else {
							moveCoord = AstarBoard(decoded, tailPos)
					}
					// still chaseTail
				}
			}

	} else {
		// Head or Tail Food...? What's best?
		// BUG: INDEX OUT OF RANGE... Maybe food is not safe
		//if(SafeFoodHead(decoded))
		if (isThereSafeFood(decoded)) {
			moveCoord = AstarBoard(decoded, SafeFoodHead(decoded))
		}
		//moveCoord = AstarBoard(decoded, SafeFoodTail(decoded))
		if moveCoord == nil {
			//moveCoord = algorithm.Astar(decoded.Board.Height, decoded.Board.Width, decoded.You, decoded.Board.Snakes, algorithm.ChaseTail(decoded.You.Body))
			if (len(decoded.You.Body) <= 4) {
				// too small? WALK! LONG PATH
					moveCoord = LongestPath(decoded, tailPos)
			} else {
					moveCoord = AstarBoard(decoded, tailPos)
			}
		}
	}

	if AstarBoard(decoded, tailPos) == nil {
		log.Print("TAIL NOT REACHABLE...") // Should do a fillMove (or should've done it the previous move...)
	}

	if (nil != moveCoord) {
		if (!isSafe(moveCoord[1], decoded)) {
				log.Print("NEXT MOVE in MoveCoord Defintion... WAS NOT SAFE! That killed me...")
				nextMove = getRandomValidMove(decoded)
		} else {
			nextMove = goToDir(headPos, moveCoord[1])
		}

	} else {
		// Path to tail, or food or otherwise not found! fallback
		// MoveCoord was 0
		coordList := getPossibleMoves(decoded)
		if (len(coordList) >= 1 && isSafe(coordList[0], decoded)) {
			nextMove = goToDir(headPos, coordList[0])
		}

	}


	if (isNextMoveFatal(me, prevMove, nextMove)) {
		log.Println("Next move (" + nextMove + ") was fatal... new move is: ")
		nextMove = newPossibleMoves(decoded)
		if ("no" == nextMove) {
			// no safe nextMove... Let's gamble!
			nextMove = getRandomValidMove(decoded)
		}
		log.Print(nextMove)
	}

//	if (len(getSafeCoordList(decoded)) >= 2) {
		// Let's check if there are both "Safe" next round...
		// i.e. is a tail (but moving, since health > 99), so not blockedByTail...

		//  both are valid in 1 move...
		// but with one is best? (highest floodFill... is better? let's test)
		// foreach validMoves as Move
		// check number of nextMoves if reaches tail after 1 spot, then 99999 (for infinite)
		// best floodFill everrrr because of moving to my own tail)
		// move with highest score wins...
		// if equal... then the current valid move stays unchanged

//		if (countLongestMoveDir(decoded) != "invalid") {
//				log.Print("2 or more SAFE Moves: Going to tail!!! " + countLongestMoveDir(decoded))
//				nextMove = countLongestMoveDir(decoded)
//		}
//	}

	fmt.Print("T " + strconv.Itoa(decoded.Turn) + " H:" + strconv.Itoa(health) + " E:" + strconv.Itoa(enemySnakes) + " Move: " + nextMove + "\n")

	respond(res, MoveResponse{
		Move: nextMove,
	})
}

/*
Only go in for the kill (La Roux... ;-)
If i'm the biggest alive, else eat food
*/
func iAmTheBiggestSnakeAlive(game SnakeRequest) bool {
		biggest := true
		 for i := 0; i < len(game.Board.Snakes); i++ {
				 if (game.Board.Snakes[i].ID != game.You.ID && len(game.Board.Snakes[i].Body) >= len(game.You.Body) ) {
					 // if a snake is equal or bigger... i'm not the biggest!
					 biggest = false
				 }
		 }
		 return biggest
}

func getEnemyHeadPos(game SnakeRequest) []Coord {
		coordList := make([]Coord, 0)
		snakeList := game.Board.Snakes
		 for i := 0; i < len(snakeList); i++ {
				 if ( snakeList[i].ID != game.You.ID && len(snakeList[i].Body) < len(game.You.Body) ) {
					 // only if snake is 2 smaller!!! 2 < (3+1)
					 // if the snake is not you and smaller.. it's food
					 coordList = append(coordList, snakeList[i].Body[0])
				 }
		 }
		 return coordList
}

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
	var safeFoodDist = Dist(foodArray[0], You)

	for i := 0; i < len(foodArray); i++ {
		if Dist(foodArray[i], You) < safeFoodDist {
			if (countEscapeRoutesFromCoord(foodArray[i], req) > 1) {
				safeFood = foodArray[i]
				safeFoodDist = Dist(foodArray[i], You)
			} else {
				log.Print("TAIL food was not safe... skipping");
			}
		}
	}

	return safeFood
}
//safeClosestFood

func isThereSafeFood(game SnakeRequest) bool {
	foodArray := game.Board.Food
	if (0 == len(foodArray)) {
		// no food. no safe food.
		return false
	}

	if ( 1 == len(foodArray) && iAmTheBiggestSnakeAlive(game) == false && len(game.Board.Snakes) > 2 && game.You.Health > 20 ) {
		// More Than 2 is a party... let's avoid the last spick food
		// if more then two snakes.. unless i'm under 20 health...
		// Then any Food is OK! (i'm not Undefined Behavior).
		if ItsMyFood(foodArray[0], game) {
			return true
		}
	}



	for i := 0; i < len(foodArray); i++ {
			if (ItsMyFood(foodArray[i], game) && isSafe(foodArray[i], game) && countEscapeRoutesFromCoord(foodArray[i], game) > 1) {
				// there is safe food...
				return true
			}
	}

	return false
}

func SafeFoodHead(game SnakeRequest) Coord {
	You := game.You.Body[0]
	foodArray := game.Board.Food

	var safeFoodDist = Dist(foodArray[0], You)
	var safeFood = foodArray[0]

	for i := 0; i < len(foodArray); i++ {
		if(ItsMyFood(foodArray[i], game)) {
			safeFood = foodArray[i] // this is the closest food to ME!
			safeFoodDist = Dist(foodArray[i], You)
		}

		if Dist(foodArray[i], You) < safeFoodDist {

			// only return safeFood && 2 escape routes...
			if (isSafe(foodArray[i], game) && countEscapeRoutesFromCoord(foodArray[i], game) > 1) {
				//var safeFood = foodArray[0] // do i want the closest food?
				safeFood = foodArray[i] // this is the closest food
				safeFoodDist = Dist(foodArray[i], You)
			}
		}
	}

	return safeFood
}

// returns true, if you're head is the closest to the food of all the snakes.
func ItsMyFood(foodCoord Coord, game SnakeRequest) bool {
		snakeList := game.Board.Snakes
		myDist := dist(game.You.Body[0], foodCoord)
		 for i := 0; i < len(snakeList); i++ {
				 if ( snakeList[i].ID != game.You.ID ) {
					 if ( dist(foodCoord, snakeList[i].Body[0]) < myDist || len(game.You.Body) <= len(snakeList[i].Body) && dist(foodCoord, snakeList[i].Body[0]) <= myDist )  {
						 // if a enemy is a step closer... it's not my food...
						  // neither claims from bigger snakes (same distance... let's not fight for that food...)
						 return false
					 }
				 }
		 }
		 // i am the closest snake to that food Coordinate
		 return true
}

func isNextMoveFatal(me Snake, currentDir string, targetDir string) bool {
		// doing a 180 is never safe, so check for that...
		flipDir := invDir(currentDir)
		if (flipDir == targetDir) {
			// log.Print("The move is " + targetDir + "but in going " + currentDir + "That would be fatal...\n")
			return true
		}
		// check if a move is NOT_OUT_OF_BOUNDS (hit a wall) WALL SNAKE
		if (isMoveOOB(headPos, targetDir)) {
		//	log.Print("Next Move is Fatal because of a BOUNDARY " + targetDir + "\n")
			return true
		}

		return false
}

func newPossibleMoves(game SnakeRequest) string {
	allCoords := SurroundingCoordinates(game.You.Body[0])
	onlyOption := "no"
	for _,coord := range allCoords {
			if (isSafe(coord, game)) {
				// if it's safe, let's GO
				return goToDir(game.You.Body[0], coord)
			} else if (isNodeOnBoard(coord) && !NodeBlockedExceptTail(coord,game.Board.Snakes) ) {
				onlyOption = goToDir(game.You.Body[0], coord)
				log.Print("My only option is... " + onlyOption)
			}
		}
		return onlyOption
}
/**
returns list of coords that are possible
*/
func getSafeCoordList(game SnakeRequest) []Coord {
	var validCoords = make([]Coord, 0)
	headPos := game.You.Body[0]
	//tailPos := getTailPos(game.You)
	enemySnakes := game.Board.Snakes
	allCoords := getOpenAjdacentNodes(headPos)
	for _,coord := range allCoords {
		dir := Heading(headPos, coord)
		if (false == isMoveOOB(headPos, dir)) {
			if (false == NodeBlockedExceptTail(coord, enemySnakes)) {
				if (isSafe(coord, game)) {
					validCoords = append(validCoords,coord)
				}
			}
		}
  }
		return validCoords
}

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

func countLongestMoveDir(game SnakeRequest) string {
	headPos := game.You.Body[0]
	allCoords := getOpenAjdacentNodes(headPos) // get all open Adjects from my thing
	longestDir := "invalid"
	longestDist := -1
	for _,coord := range allCoords {

		dir := Heading(headPos, coord)
		dist := CountDirectionFloodFill(game, coord)

			if (isSafe(coord, game)) {
				// count root from there to the coord!
				if (dist > longestDist) {
					// this distance is bigger! let's print it here
					longestDist = dist
					longestDir = dir
					//log.Print("New LONG:" + longestDir + " " + strconv.Itoa(longestDist) + "T")
				} else {
					//log.Print("" + dir + ":" + strconv.Itoa(dist) + "T")
				}

				//log.Print("Got the Most Optimal Route in random move...")
				//return Heading(headPos, coord)
				// false == NodeBlockedExceptTail(coord, enemySnakes) ?
				// if we have a safe Move.. return that one, else... each one is as bad af them..
			}
		} // for adjacentNodes

		return longestDir
}



func getRandomValidMove(game SnakeRequest) string {
	headPos := game.You.Body[0]
	//tailPos := getTailPos(game.You)
	enemySnakes := game.Board.Snakes
	allCoords := getOpenAjdacentNodes(headPos)
	nextDir := "invalid"

	for _,coord := range allCoords {
			if (isSafe(coord, game) && false == NodeBlockedExceptTail(coord, enemySnakes)) {
				log.Print("Got the Most Optimal Route in random move...")
				return Heading(headPos, coord)
				// if we have a safe Move.. return that one, else... each one is as bad af them..
			}

			if (false == NodeBlocked(coord, enemySnakes)) {
				dir := Heading(headPos, coord)
				//if(game.You.Health > 99 && dist(headPos, tailPos) == 1 && dir == goToDir(headPos, tailPos)) {
				//	log.Print("skipping " + dir + " as it would crash into tail")
					// This is Incorrect i guess...
				//	continue
				//}

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
	//log.Print("new possbile: " + newPossibleMoves(game))
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
		// Let's do suicide instead ... do not check if it's on the board...
		// if (!isNodeOnBoard(point)) {
		//	return false
		// }

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
					 //log.Print("turn ")
					 if (game.Turn > 1 && len(snakeList[i].Body)-1 == j && snakeList[i].Health < 100) {
						 // Tail is safe...
						 //log.Print("tail is safe to step... except in trn 1...")
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
// manhatten dist... snakes can't slither diagonally...
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
