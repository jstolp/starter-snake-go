package main

import (
	"log"
	"encoding/json"
	"net/http"
	. "github.com/jstolp/pofadder-go/api"
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
		log.Print("no food on board... chasing tail...")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, tailPos)
	} else if (health > HUNGRY_TRESHOLD) &&  len(decoded.You.Body) >= 4 {
		// there is food but i'm not hungry, still chase tail only if i'm big enough.
		log.Print("Not hungry... chasing tail")
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, closestCorner(boardHeight, boardWidth, headPos))
	} else if (health > HUNGRY_TRESHOLD) {
		log.Print("not big enough... moving to corner")
		targetCorner := closestCorner(boardHeight, boardWidth, headPos)
		moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, targetCorner)
		if moveCoord == nil {
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, ChaseTail(me.Body))
		}
	} else {
		log.Print("Going to FAR food...")
		
		if(health < 50) { // almost starving... no detours anymoe go to nearest Food.
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, NearestFood(foodList,headPos))
		} else {
			log.Print("Going to CORNER")
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, closestCorner(boardHeight, boardWidth, headPos))
		}
		
		if moveCoord == nil {
			log.Print("TAIL NearestFood!! FALBACK")
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, NearestFood(foodList, tailPos))
		}
		
		if moveCoord == nil {
			log.Print("CANT REACH FOOD, chaseTail Fallback")
			moveCoord = Astar(boardHeight, boardWidth, me, enemySnakes, ChaseTail(me.Body))
		}
	}

	grid := mapToGrid(moveCoord[len(moveCoord) - 1], decoded, boardHeight)
	PrintGrid(grid)
	
	nextMove = Heading(headPos, moveCoord[1])

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
			Coord{boardWidth-2,1},
			Coord{1, boardHeight - 2},
			Coord{boardHeight - 2, boardWidth -2},
		}

		for _, coord := range corners {
			
			if( dist(headPos, coord) >= distToCorner) {
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

/* Dist to function in steps (int) */
func dist(a Coord, b Coord) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
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

func mapToGrid(target Coord, decoded SnakeRequest, grid_size int) ([][]string) {

  
  grid := make([][]string, grid_size)
  me := decoded.You
  foodList := decoded.Board.Food

  
  for i := 0; i < len(grid); i++ {
      grid[i] = make([]string, grid_size)
  }

  for i := 0; i < grid_size; i++ {
     grid[0][i] = "."

     grid[i][0] = "."

     grid[grid_size-1][i] = "."

     grid[i][grid_size-1] = "."
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
//if (len(foodList) > 0) {
  // there is food on the board.
  for _, coord := range foodList {
     c := coord.X
     r := coord.Y

     if grid[r][c] != "#" {
        grid[r][c] = "!"
      }
  }
//}

// set target to grid
  grid[target.Y][target.X] = "$"
  


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

 return grid;
}

func PrintGrid(grid [][]string) {
    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[0]); j++ {
            if grid[i][j] == "" {
                fmt.Printf(".")
            } else {
                fmt.Print(grid[i][j])
            }
        }
        fmt.Print("\n")
    }
    fmt.Print("\n")
}
