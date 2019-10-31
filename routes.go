package main

import (
  "log"
  "encoding/json"
  "net/http"
  //. "github.com/jstolp/pofadder-go/api" // Heroku wants this
   . "./api" // local wants this
  "github.com/tkanos/gonfig"
// . "github.com/maximelamure/algorithms/datastructure" // BFS
  "github.com/nickdavies/go-astar/astar"
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
var me Snake;
var health int = 100;
var numOfStartingSnakes int = 1;
var numSnakesLeft int = 1;
var enemySnakes int = 0;
var foodPointList []Coord;
var endCicle bool = false;

var selectedFood Coord;
var noTargetFood bool = true; // i have no target food.

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


  log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")
  if(numOfStartingSnakes == 1) {
    log.Print("\n\n It's Gonna be a SOLO GAME \n")
  }
  /* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
  tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */
  /*
   e19c41 - orange test 2
   00ff55 - green
  */
  respond(res, StartResponse{
    Color: "#ff00aa",
    HeadType: "tongue",
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
  
  me = decoded.You
  headPos = getHeadPos(me)
  foodPointList = decoded.Board.Food
  tailPos = getTailPos(me)


  rows, cols := decoded.Board.Height, decoded.Board.Width
  // Build a new AStar object from Height x Width
  astarBoard := astar.NewAStar(rows, cols)
  p2p := astar.NewPointToPoint()
  grid := mapToGrid(astarBoard, decoded, rows)

  


  health = me.Health
  //myLength := len(me.Body)
  numSnakesLeft = len(decoded.Board.Snakes)
  enemySnakes = numSnakesLeft - 1
  turn = decoded.Turn

    if(len(foodPointList) > 0) {
		closestFoodPoint := minDistFood(headPos,foodPointList)
		source := make([]astar.Point, 1)
		source[0].Row = headPos.Y
		source[0].Col = headPos.X
		
		target := make([]astar.Point, 1)
	    target[0].Row = closestFoodPoint.Y
		target[0].Col = closestFoodPoint.X
		path := astarBoard.FindPath(p2p, source, target)
	    for path != nil {
			path = path.Parent
			log.Printf("TRGT At (%d, %d)\n", path.Col, path.Row)
			log.Printf("HEAD At (%d, %d)\n", headPos.X, headPos.Y)
			// target coord: 
			targetCoord := Coord{path.Col, path.Row}
			nextMove = goToDir(headPos, targetCoord)
			log.Print("next move is: " + move)
			
		}
		//DrawPath(grid, path, "$")
		PrintGrid(grid);
	}

//log.Print("TURN " + strconv.Itoa(turn) + " e: "+ strconv.Itoa(enemySnakes)+" h: "+ strconv.Itoa(health) + "\n")


    validMoves := validMoveCoordinates(headPos, nextMove, me.Body)

    //
    // if there is a direct path to the food... Go for it!

    if (health < 40) {
        //closestFoodPoint := minDistFood(headPos,foodPointList)
        //selectedFood = closestFoodPoint
                foodDir := goToDir(headPos, selectedFood)
//        fmt.Print("I've selected food...")
//        dd(selectedFood)


      //log.Print("im hungry and there is food... \n\n")

            if (IsPosInCoordList(selectedFood,validMoves)) {
                    // MOVE IS VALID GO DO IT!

          if(!isNextMoveFatal(me, prevMove, foodDir)) {
            fmt.Print("nextMove: " + foodDir + " because its not fatall.... \n")
            nextMove = foodDir
          } else {
                fmt.Print("This FoodDirection is FATAL: " + foodDir)
              nextMove = randomNOOBmove(headPos, nextMove)
          }
            } else {
                if(!isNextMoveFatal(me, prevMove, foodDir)) {
            nextMove = foodDir
        } else {
          fmt.Print("STOP STOP STOP " + foodDir + " is fatal!!!! \n")
          //  nextMove = randomNOOBmove(headPos, prevMove)
                    // GET VALID NON-OOB MOVES:
                    fmt.Print("with my currentHEADPOINT: \n")
                    dd(headPos)
                    fmt.Print("The following points are IN BOUND!") // and no Body-Crash to foodPOINT
                    /* if(isSafeCoordinate(selectedFood, me.Body)) {
                            foodDir := goToDir(headPos, selectedFood)
                            fmt.Print("foodDir" + foodDir +" is a safe Coordinate")
                    }
 */
          fmt.Print("OK... ive selected " +  nextMove + " as the next move \n")
        }
            }
        //fmt.Print("im gooing to " + foodDir + "seems to be a good idea...")

      } // no health
    
  move = nextMove // finalise the move
  fmt.Print(strconv.Itoa(turn) + ": Move " + move)
  log.Print("Move: " + move)
  
  respond(res, MoveResponse{
    Move: move,
  })
  prevMove = move // Re-allocate move command to prev/last move\
} // END MOVE COMMAND


func isSafeCoordinate(targetcoord Coord, myBodyPoints CoordList) bool {
    validCoords := CoordList {
            Coord{targetcoord.X - 1, targetcoord.Y},
            Coord{targetcoord.X, targetcoord.Y - 1}, // ONE is fatal
            Coord{targetcoord.X + 1, targetcoord.Y}, // but which
        Coord{targetcoord.X, targetcoord.Y + 1}, // i need my currentDir
        }

        for _, coord := range validCoords {
            if coord.X < 0 || coord.Y < 0 { // TOP LEFT CORDER NE
                return false
            }
            if coord.X > rightBound - 1 || coord.Y > botBound - 1 { // OOB Protection
                return false
            }

            if IsPosInCoordList(coord, myBodyPoints) { // if MYBody is in CoordList, no problem
                fmt.Print("isSafeCoord: COORD CRASHING INTO BODY (or tail, or head i guess)")
                return false
            }

        } // end for "valid" coords...
        return true
    }

/**
returns grid,
*/
func mapToGrid(ast astar.AStar, decoded SnakeRequest, grid_size int) ([][]string) {

  
  grid := make([][]string, grid_size)
  me := decoded.You
  foodList := decoded.Board.Food

  for i := 0; i < len(grid); i++ {
      grid[i] = make([]string, grid_size)
  }

  for i := 0; i < grid_size; i++ {
     grid[0][i] = "."
     //ast.FillTile(astar.Point{0, i}, -1)

     grid[i][0] = "."
     //ast.FillTile(astar.Point{i, 0}, -1)

     grid[grid_size-1][i] = "."
     //ast.FillTile(astar.Point{grid_size - 1, i}, -1)

     grid[i][grid_size-1] = "."
     //ast.FillTile(astar.Point{i, grid_size - 1}, -1)
 }

otherSnakes := decoded.Board.Snakes

for _, snake := range otherSnakes {
  for i, coord := range snake.Body {
    c := coord.X
    r := coord.Y

    if grid[r][c] != "#" {
      if(i == 0) {
        grid[r][c] = "h"
        ast.FillTile(astar.Point{r, c}, -1) // not traversable
      } else if(i == len(snake.Body) - 1) {
        grid[r][c] = "t"
      } else {
        ast.FillTile(astar.Point{r, c}, -1) // not traversable
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



 myBody := me.Body;
 for _, coord := range myBody {
    c := coord.X
    r := coord.Y

    ast.FillTile(astar.Point{r, c}, -1)
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

  ast.ClearTile(astar.Point{r, c}) // clear tail tile
  if grid[r][c] != "#" {
     grid[r][c] = "T"
   }

 return grid;
}

func validMoveCoordinates(headPos Coord, direction string, myBodyPoints CoordList) CoordList {
  validCoords := make(CoordList, 0)
   moves := make(CoordList, 0)

    switch direction {
        case "down":
            validCoords = CoordList {
                Coord{headPos.X - 1, headPos.Y},
                Coord{headPos.X, headPos.Y - 1}, // ONE is fatal
                Coord{headPos.X + 1, headPos.Y}, // but which
             //   Coord{headPos.X, headPos.Y + 1}, // i need my currentDir
              }
            // fatalMove := Coord{headPos.X, headPos.Y + 1} // down fatal if going up..
        case "up":
            validCoords =  CoordList {
                Coord{headPos.X - 1, headPos.Y},
                //Coord{headPos.X, headPos.Y - 1}, // ONE is fatal
                Coord{headPos.X + 1, headPos.Y}, // but which
                Coord{headPos.X, headPos.Y + 1}, // i need my currentDir
              }
            //fatalMove := Coord{headPos.X, headPos.Y - 1} // down fatal if going up..
        case "left":
            validCoords = CoordList {
                //Coord{headPos.X - 1, headPos.Y},
                Coord{headPos.X, headPos.Y - 1}, // ONE is fatal
                Coord{headPos.X + 1, headPos.Y}, // but which
                Coord{headPos.X, headPos.Y + 1}, // i need my currentDir
              }
            //    fatalMove := Coord{headPos.X - 1, headPos.Y} // fatal if going left.
        case "right":
            validCoords = CoordList {
                Coord{headPos.X - 1, headPos.Y},
                Coord{headPos.X, headPos.Y - 1}, // ONE is fatal
               // Coord{headPos.X + 1, headPos.Y}, // but which
                Coord{headPos.X, headPos.Y + 1}, // i need my currentDir
              }
        //    fatalMove := Coord{headPos.X + 1, headPos.Y} // targetPoint is fatal if going right
    }

  for _, coord := range validCoords {
    if coord.X < 0 || coord.Y < 0 { // TOP LEFT CORDER NE
      continue
    }
    if coord.X > rightBound - 1 || coord.Y > botBound - 1 { // OOB Protection
      continue
    }
        // REMOVE OWN BODY coordinates... starting with tail! CoordList
        //if(coord == tailPos)
    // COORD TAILPOS IS ACTUALLY OK, because it moves one step next turn! :D
        if(IsPosInCoordList(coord, myBodyPoints)) {
            // CRASHING INTO MYSELF...
            //fmt.Print("VALIDMOVES: COORD CRASHING INTO BODY (or tail, or head i guess)")
            //dd(coord)
            continue
        }
    //    if(coord ) // bodyCoords
        //if(coord in me.Body) -- remove

    moves = append(moves, coord) // then put coord in the possible move list
        // NEVER OOB!!! :D
  } // end for "valid" coords...
//    fmt.Print("The following x:"+ strconv.Itoa(len(moves))+" is valid \n ")
//    dd(moves)
  return moves
}

// returns Index Of[arrayValue]PHP x,y in CoordList
func idx(target Coord, coordinates CoordList) int {
    for index, value := range coordinates {
        if value == target {
            return index // return index
        }
    }
    return -1
}

// returns true if pOsition is in List
func IsPosInCoordList(target Coord, coordList CoordList) bool {
    // create Bool.
    return 0 >= idx(target, coordList) // true if index is found (-1 for not found, 0== found at first value)
}


func isMoveOOB(headPos Coord, direction string) bool {
//    fmt.Print("\n my head is... and im going -> " + direction)
//    dd(headPos)

    if (headPos.Y == 0 && direction == "up") { // TOP LEFT CORDER NE
    fmt.Print("wallcrash north !!!", headPos)
        return true // wallCrash NORTH or LEFT (W)
    }

if headPos.X > rightBound - 1 || headPos.Y > botBound - 1 { // OOB Protection
    fmt.Print("WALLCRASH South OR EAST.", headPos)
        return true // wallCrash South or EAST
    }

    fmt.Print( "\n\n")
  switch direction {
    case "down":
//            fatalMove := Coord{headPos.X, headPos.Y + 1} // down fatal if going up..
//            dd(fatalMove)
//            fmt.Print("down fatal...")
            if headPos.X > rightBound - 1 || headPos.Y > botBound - 1 { // OOB Protection
        fmt.Print("im going down, wallcrash south or east...", headPos)
                return true // wallCrash South or EAST
            }
            if (headPos.Y + 1 < botBound) {
          return false
      }

    case "up":
//            fatalMove := Coord{headPos.X, headPos.Y - 1} // down fatal if going up..
//            dd(fatalMove)
//            fmt.Print("is a fatal move going up...")
            //if headPos.X - 1 < 0 || headPos.Y - 1 < 0 { // TOP LEFT CORDER NE
            //    return true // wallCrash NORTH or LEFT (W)
        //    }
      if (headPos.Y + 1 > topBound) {
        return false
      }

    case "left":
      if (headPos.X > leftBound - 1) {
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

func isNextMoveFatal(me Snake, currentDir string, targetDir string) bool {
    // doing a 180 is never safe, so check for that...
    flipDir := invDir(currentDir)
    if(flipDir == targetDir) {
      //log.Print("The move is " + targetDir + "but in going " + currentDir + "That would be fatal...\n")
      return true
    }
    // check if a move is NOT_OUT_OF_BOUNDS (hit a wall) WALL SNAKE
    if (isMoveOOB(headPos, targetDir)) {
      //log.Print("Next Move is Fatal because of a BOUNDARY " + targetDir + "\n")
      return true
    }

    // if dist to my own tail is 1, and i'm going in the same direction...
    // i'll die...
    if (dist(headPos, tailPos) == 1 && targetDir == goToDir(headPos, tailPos)) {
      log.Print("Dist to TAIL is 0... and i want to go directly to my own tail... i'm dead... \n")
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

// creates a ASTAR BOARD
func createBoard() {

}

// nick-astar
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

// nick-astar
func DrawPath(grid [][]string, path *astar.PathPoint, path_char string) {
    for {
        if grid[path.Row][path.Col] == "#" {
            grid[path.Row][path.Col] = "X"
        } else if grid[path.Row][path.Col] == "" {
            grid[path.Row][path.Col] = path_char
        }

        path = path.Parent
        if path == nil {
            break
        }
    }
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
