package main

import (
  "log"
  "encoding/json"
  "net/http"
//  . "github.com/jstolp/pofadder-go/api" // Heroku wants this
. "./api" // local wants this
  "github.com/tkanos/gonfig"
// . "github.com/maximelamure/algorithms/datastructure" // BFS
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

  edgeSnakeLimit = (((botBound - 1) * 2) + ((rightBound - 1) * 2))

  log.Print("BOARD Size: TOP LEFT  NW Corner x:" + strconv.Itoa(topBound) + " , " + strconv.Itoa(leftBound))
  log.Print("BOARD Size: BOT RIGHT SE Corner x:" + strconv.Itoa(botBound) + "," + strconv.Itoa(rightBound))
  log.Println("Snake Edge Limit: " + strconv.Itoa(edgeSnakeLimit))
  log.Print("Enemy Snakes: " + strconv.Itoa(numOfStartingSnakes - 1) + "\n\n")

  log.Print("Start Pos: " + strconv.Itoa(headPos.X) + "," + strconv.Itoa(headPos.Y))
  if(numOfStartingSnakes == 1) {
    log.Print("\n\n It's Gonna be a SOLO GAME \n")
  }
  /* heads: "beluga" "bendr" "dead" "evil" "fang" "pixel" "regular" "safe" "sand-worm" "shades" "silly" "smile" "tongue"
  tails: "block-bum" "bolt" "curled" "fat-rattle" "freckled" "hook" "pixel" "regular" "round-bum" "sharp" "skinny" "small-rattle" */
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
//  tailPos := getTailPos(me)

  health = me.Health
  myLength := len(me.Body)
  numSnakesLeft = len(decoded.Board.Snakes)
  enemySnakes = numSnakesLeft - 1
  turn = decoded.Turn
  if (health > 99) {
    fmt.Print("I JUST ATE FOOOOOOD \n\n")
    // i just ate. reset foodPoint
      noTargetFood = true
  }
  // IF at 0,0 I'm in the TOP-left corner
  if (headPos.X == 0 && headPos.Y == 0) {
    //log.Printf("I'm in the TOP-LEFT NW CORNER AT TURN %d", turn)
  }

  if (headPos.X == rightBound - 1 && headPos.Y == botBound - 1) {
    //log.Printf("I'm in the BOT-RIGHT SE CORNER AT TURN %d", turn)
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
    //
  }
log.Print("TURN " + strconv.Itoa(turn) + " e: "+ strconv.Itoa(enemySnakes)+" h: "+ strconv.Itoa(health) + "\n")

  nextMoveIsOOB := isMoveOOB(headPos, nextMove)
  if (nextMoveIsOOB) {
    // CLOCKWISE: invDir(randomNOOBmove(headPos, move))
    nextMove = randomNOOBmove(headPos, move)
    // COUNTER-CLOCKWISE: randomNOOBmove(headPos, move)
  }

  if(len(foodPointList) > 0) {

    if (health < 40) {

      if (noTargetFood) {
        closestFoodPoint := minDistFood(headPos,foodPointList)
        selectedFood = closestFoodPoint
        fmt.Print("I've selected food...")
        dd(selectedFood)
      }
      //log.Print("im hungry and there is food... \n\n")
        foodDir := goToDir(headPos, selectedFood)

        //fmt.Print("im gooing to " + foodDir + "seems to be a good idea...")
        if(!isNextMoveFatal(me, prevMove, foodDir)) {
            nextMove = foodDir
        } else {
          fmt.Print("STOP STOP STOP " + foodDir + " is fatal!!!! \n")
          //  nextMove = randomNOOBmove(headPos, prevMove)
          fmt.Print("OK... ive selected " +  nextMove + " as the next move \n")
        }
      } // end if health < 60
    }  else {
    log.Print("THERE IS NO FOOD \n\n")
    noTargetFood = true;
  }


  //test := isNextMoveFatal(me, prevMove, nextMove)
  move = nextMove // finalise the move
  fmt.Print(strconv.Itoa(turn) + ": Move " + move)
  //fmt.Print(test)
  //fmt.Println()
  respond(res, MoveResponse{
    Move: move,
  })
  prevMove = move // Re-allocate move command to prev/last move\
} // END MOVE COMMAND

/*
//python bfs

frontier = Queue()
frontier.put(start )
visited = {}
visited[start] = True

while not frontier.empty():
   current = frontier.get()
   for next in graph.neighbors(current):
      if next not in visited:
         frontier.put(next)
         visited[next] = True
   */
/*
   type BFSPath struct {
     Source int
     DistTo map[int]int
     EdgeTo map[int]int
     Path   Queue
     G      *Graph
   }

   func NewBFSPath(g *Graph, source int) *BFSPath {
     bfsPath := &BFSPath{
       DistTo: make(map[int]int),
       EdgeTo: make(map[int]int),
       G:      g,
       Path:   NewQueueLinkedList(),
       Source: source,
     }
     bfsPath.bfs(source)
     return bfsPath
   }
*/

func validMove(headPos Coord, direction string) CoordList {
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
			//	fatalMove := Coord{headPos.X - 1, headPos.Y} // fatal if going left.
		case "right":
			validCoords = CoordList {
			    Coord{headPos.X - 1, headPos.Y},
			    Coord{headPos.X, headPos.Y - 1}, // ONE is fatal
			   // Coord{headPos.X + 1, headPos.Y}, // but which
			    Coord{headPos.X, headPos.Y + 1}, // i need my currentDir
			  }
		//	fatalMove := Coord{headPos.X + 1, headPos.Y} // targetPoint is fatal if going right
	}

  for _, coord := range validCoords {
    if coord.X < 0 || coord.Y < 0 { // TOP LEFT CORDER NE
      continue
    }
    if coord.X > rightBound - 1 || coord.Y > botBound - 1 { // OOB Protection
      continue
    }
    moves = append(moves, coord) // then put coord in the possible move list
		// NEVER OOB!!! :D
  }
  return moves
}

// BFS implementation
/*
func bfs(items []Coord) : Coord {
  pq := make(PriorityQueue, len(items))
  visited := []Coord

  i := 0
  for value, priority := range items {
      pq[i] = &Item{
          value:    value,
          priority: priority,
          index:    i,
      }
      i++
  }
  heap.Init(&pq)
//  heap.Push(&pq, item)
//  frontier :=
}
*/

/*
func (b *BFSPath) bfs(v int) {
  queue := NewQueueLinkedList()
  b.DistTo[v] = 0
  queue.Enqueue(v)
  for {
    if queue.IsEmpty() {
      break
    }
    d := queue.Dequeue().(int)
    b.Path.Enqueue(d)
    for r := range b.G.Adj(d) {
      if _, ok := b.DistTo[r]; !ok {
        queue.Enqueue(r)
        b.EdgeTo[r] = d
        b.DistTo[r] = 1 + b.DistTo[d]
      }
    }
  }
}
*/

//func getFatalHeadPos(headPos Coord, string direction) Coord {


	//return fatalMove
//}

// Check if MoveIs Out of Bounds...
// What a horror function.... v0.2.0 consider refactor
func isMoveOOB(headPos Coord, direction string) bool {
	fmt.Print("\n my head is... and im going -> " + direction)
	dd(headPos)


	fmt.Print( "\n\n")
  switch direction {
    case "down":
			fatalMove := Coord{headPos.X, headPos.Y + 1} // down fatal if going up..
			dd(fatalMove)
			fmt.Print("down fatal...")
      if (headPos.Y + 1 < botBound) {
          return false
      }
    case "up":
			fatalMove := Coord{headPos.X, headPos.Y - 1} // down fatal if going up..
			dd(fatalMove)
			fmt.Print("is a fatal move going up...")
      if (headPos.Y + 1 > topBound) {
        return false
      }
    case "left":
				fatalMove := Coord{headPos.X - 1, headPos.Y} // fatal if going left.
				fmt.Print("fatal move if going left \n")
				dd(fatalMove)
      if (headPos.X + 1 > leftBound) {
        return false
      }
    case "right":
		//	fatalMove := Coord{headPos.X + 1, headPos.Y} // targetPoint is fatal if going right
      if (headPos.X + 1 < rightBound) {
        return false
      }
  }

	// GET VALID NON-OOB MOVES:
 	fmt.Print("with my currentHEADPOINT: \n")
	dd(headPos)
	fmt.Print("The following points are IN BOUND!") // and no Body-Crash to foodPOINT
	dd(validMove(headPos, direction))
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
