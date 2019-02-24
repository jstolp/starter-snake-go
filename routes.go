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

var leftBound int = 1;
var topBound int = 1;
var rightBound int = 0;
var botBound int = 0;

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

/* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */

func shrinkArena() {
	leftBound = leftBound + 1
	topBound = topBound + 1
	rightBound = rightBound - 1
	botBound = botBound - 1
	// edgeSnakeLimit = ((botBound - 1) * (rightBound - 1)) - FALSE ASSUMPTION. it doesn't work if you shrink, because you are bigger

	log.Print("BOARD Size: TOP LEFT  NW Corner x:" + strconv.Itoa(topBound) + " , " + strconv.Itoa(leftBound))
	log.Print("BOARD Size: BOT RIGHT SE Corner x:" + strconv.Itoa(botBound) + "," + strconv.Itoa(rightBound))
	log.Println("Snake Edge Limit: " + strconv.Itoa(edgeSnakeLimit))

	log.Print("Shrunk the Area by 1x1... new SIZES \n")
}

func Start(res http.ResponseWriter, req *http.Request) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	decoded := SnakeRequest{}
	err := DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}

  headPos = decoded.You.Body[0]
  //foodPointList := decoded.Board.Food
  numOfStartingSnakes = len(decoded.Board.Snakes)
  topBound, leftBound = 1, 1; // Set NW bound, X, Y
	botBound, rightBound = decoded.Board.Height, decoded.Board.Width // SE corner X, Y

	edgeSnakeLimit = (((botBound - 1) * 2) + ((rightBound - 1) * 2))

	log.Print("BOARD Size: TOP LEFT  NW Corner x:" + strconv.Itoa(topBound) + " , " + strconv.Itoa(leftBound))
	log.Print("BOARD Size: BOT RIGHT SE Corner x:" + strconv.Itoa(botBound) + "," + strconv.Itoa(rightBound))
	log.Println("Snake Edge Limit: " + strconv.Itoa(edgeSnakeLimit))
	log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

	log.Print("Start Pos: " + strconv.Itoa(headPos.X) + "," + strconv.Itoa(headPos.Y))
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
		HeadType: "tongue",
		TailType: "curled",
	})
}

// Check if MoveIs Out of Bounds...
// What a horror function.... v0.2.0 consider refactor
func isMoveOOB(headPos Coord, direction string) bool {
	switch direction {
		case "down":
			if (headPos.Y + 1 < botBound) {
					return false
			}
		case "up":
			if (headPos.Y + 1 > topBound) {
				return false
			}
		case "left":
			if (headPos.X + 1 > leftBound) {
				return false
			}
		case "right":
			if (headPos.X + 1 < rightBound) {
				return false
			}
	}
	return true
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

// TIME TO DEPRICATE THIS FUNCTION. I NEED A MOVE THAT IS... 1 NOT OOB.
// 1. Not into a wall!
// 2. not into Myself.
// 3. not into my Body.
// 4 (Battle) - not into another snake.
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

	me := decoded.You
	headPos := getHeadPos(me)
	foodPointList = decoded.Board.Food
//	tailPos := getTailPos(me)
	health = me.Health
	myLength := len(me.Body)
	numSnakesLeft = len(decoded.Board.Snakes)
	enemySnakes = numSnakesLeft - 1
	turn = decoded.Turn

  // IF at 0,0 I'm in the TOP-left corner
	if (headPos.X == 0 && headPos.Y == 0) {
		log.Printf("I'm in the TOP-LEFT NW CORNER AT TURN %d", turn)
	}

	if (headPos.X == rightBound - 1 && headPos.Y == botBound - 1) {
		log.Printf("I'm in the BOT-RIGHT SE CORNER AT TURN %d", turn)
	}

/*
	if (me.Body[0].X == 0 && me.Body[0].Y == 0 && myLength == edgeSnakeLimit) {
		log.Print("IM TOP LEFT... \n\n")
		shrinkArena()
	} else {
		dd(me.Body[0])
	}
*/

	if (myLength == edgeSnakeLimit) {
			log.Println("circle JErk")
	}



	if (myLength > edgeSnakeLimit) {
			log.Println("DEATH DEATH DEATH TAILCRASH")
	}


	if (enemySnakes < 1) {
		// SOLO MODE!
		} else {
			// BATTLE  MODE
		//log.Print("TURN " + strconv.Itoa(turn) + " e: "+ strconv.Itoa(enemySnakes)+" h: "+ strconv.Itoa(health) + "\n")
	}


	nextMoveIsOOB := isMoveOOB(headPos, nextMove)
	if (nextMoveIsOOB) {
		// CLOCKWISE: invDir(randomNOOBmove(headPos, move))
		nextMove = randomNOOBmove(headPos, move)
		// COUNTER-CLOCKWISE: randomNOOBmove(headPos, move)
	}

	if (endCicle) {
		endCicle = false;
		log.Print("END CIRCLE COMMAND. SHRuhk the arena... i should VEERE... \n")
	}
	// if im bigger... i can't do the edge snake strategy...
  if (endCicle == false && myLength == edgeSnakeLimit) {
			log.Println("CirleJerk... Infinity SNAKEEE... let's switch the strat.")
			shrinkArena()
			edgeSnakeLimit = (((botBound - 1) * 2) + ((rightBound - 1) * 2))
			nextMove = randomNOOBmove(headPos, move)
			//nextMove = randomNOOBmove(headPos, move)
	}

	if (health < 60) {
		log.Print("im hungry... \n\n")

		if(len(foodPointList) > 0) {
		closestFoodPoint := minDistFood(headPos,foodPointList)
			foodDir := goToDir(headPos,closestFoodPoint)
			dd(foodDir)

			fmt.Print("im gooing to " + foodDir + "seems to be a good idea...")
			if(!isNextMoveFatal(me, prevMove, foodDir)) {
					nextMove = foodDir
			} else {
				fmt.Print("STOP STOP STOP " + foodDir + " is fatal!!!!")
				//	nextMove = randomNOOBmove(headPos, prevMove)
					fmt.Print("OK... ive selected " +  nextMove + "as the next move")
			}
	} else {
			log.Print("IM HUNGRY BUT THERE IS NO FOOD \n\n")
		}
	}


  test := isNextMoveFatal(me, prevMove, nextMove)
	move = nextMove // finalise the move
	fmt.Print(strconv.Itoa(turn) + "Move: " + nextMove + "\n Is fatal: ")
	fmt.Print(test)
	fmt.Println()
	respond(res, MoveResponse{
		Move: nextMove,
	})
	prevMove = nextMove // Re-allocate move command to prev/last move\
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

func GetBodies(snakes SnakesList) []Coord {
  list := make([]Coord, 0)
  for _, s := range snakes {
    list = append(list, s.Body...)
  }
  return list
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
	res.Write([]byte("Jay's battleSnake mk III self aware: " + configuration.Home_Route))
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
