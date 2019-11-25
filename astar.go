package main

import (
	. "github.com/jstolp/pofadder-go/api"
	"math"
	"fmt"
	"log"
)

/*
     ripped from alecj1240/alec-snake, astar
	G - the amount of steps its taken to get to that Node
	H - the heuristic - estimation from this Node to the target
	F - the sum of G and H
	ParentCoords - the coordinates of the previous step (Node)
*/

// Node holds: coordinates, G, H, F, parent coords
type Node struct {
	Coord        Coord
	G            int
	H            int
	F            int
	ParentCoords Coord
}

func GetAdjacentCoords(Location Coord) []Coord {
	var adjacentCoords = make([]Coord, 0)
	adjacentCoords = append(adjacentCoords, Coord{X: Location.X + 1, Y: Location.Y})
	adjacentCoords = append(adjacentCoords, Coord{X: Location.X - 1, Y: Location.Y})
	adjacentCoords = append(adjacentCoords, Coord{X: Location.X, Y: Location.Y + 1})
	adjacentCoords = append(adjacentCoords, Coord{X: Location.X, Y: Location.Y - 1})

	return adjacentCoords
}

func removeFromOpenList(removalNode Node, openList []Node) []Node {
	var newOpenList = make([]Node, 0)
	for i := 0; i < len(openList); i++ {
		if removalNode.Coord != openList[i].Coord {
			newOpenList = append(newOpenList, openList[i])
		}
	}
	return newOpenList
}

func appendList(appendingNode Node, Snakes []Snake, List []Node, BoardHeight int, BoardWidth int) []Node {
	if NodeBlocked(appendingNode.Coord, Snakes) == false && OnBoard(appendingNode.Coord, BoardHeight, BoardWidth) {
		List = append(List, appendingNode)
	}

	return List
}

// reverseCoords reverses the path of coordinates so it's in chronological order
func reverseCoords(path []Coord) []Coord {
	for a := 0; a < len(path)/2; a++ {
		b := len(path) - a - 1
		path[a], path[b] = path[b], path[a]
	}
	return path
}

// countDirection gives the counts for the moves going into that direction
func CountDirectionFloodFill(game SnakeRequest, moveCoord Coord) int {
		closedList := make(map[Coord]bool)
		openList := make([]Node, 0)
		pathTracker := make(map[Coord]Coord)

		me := game.You
		myTailPos := me.Body[len(me.Body) - 1]
		firstMove  := Node{Coord: moveCoord, G: 0, H: 0, F: 0}
		openList = append(openList, firstMove)
		
		floodFillNodeCount := 1

		for len(openList) > 0 {
			// get the biggest distance!
			var closeNode = openList[0]
			for _, openItem := range openList {
				if openItem.F > closeNode.F {
					closeNode = openItem
				}
			}

			// put it on the closed list
			closedList[closeNode.Coord] = true
			openList = removeFromOpenList(closeNode, openList)
			// loop through leastNodes's adjacent tiles -- call them T
			closeNeighbours := GetAdjacentCoords(closeNode.Coord)

			for _, neighbour := range closeNeighbours {

				// 1. If T on the closed list, ignore it
				if closedList[neighbour] {
					continue
				}



				if neighbour == myTailPos {
					// we don't have a target... since we want the MAX!
					// if we hit out tail (and we are <100 health, return 99999 for infinite *opt for -1* )
					// if neighbour is my tail, i'm always safe.
					// maybe also implement enemyTails for escape.
					floodFillNodeCount++
					closedList[neighbour] = true
					//return 99999 // infinite! since it's 1 off my tail.
					return floodFillNodeCount + 999
					//maxDist = 999
				}


				for _, item := range openList {
					if neighbour == item.Coord {
						if NodeBlockedExceptTail(neighbour, game.Board.Snakes) == false && OnBoard(neighbour, game.Board.Height, game.Board.Width) {
							// floodFillNodeCount++ because it's a tile.
							//floodFillNodeCount++
							if (closeNode.G+1)+Dist(neighbour, moveCoord) > item.F {
								// count dist from MoveCoord to The new item...
								item.F = (closeNode.G + 1) + Dist(neighbour, moveCoord)
								item.G = closeNode.G + 1
								item.H = Dist(neighbour, moveCoord)
								item.ParentCoords = neighbour
								pathTracker[item.Coord] = closeNode.Coord
								floodFillNodeCount++
							}
							// if(isThisATailOfEnemy) return 9999
							// but we just have to count it!

						}
					}
				}


				// 2. If T is not on the open list add it
				var openNode = Node{
					Coord:        neighbour,
					G:            closeNode.G + 1,
					H:            Dist(neighbour, moveCoord),
					F:            (closeNode.G + 1) + (Dist(neighbour, moveCoord)),
					ParentCoords: closeNode.Coord,
				}

				pathTracker[neighbour] = closeNode.Coord
				openList = appendList(openNode, game.Board.Snakes, openList, game.Board.Height, game.Board.Width)
			}
		} // end for openList, we are done for this direction.

		return floodFillNodeCount // return longestPath
}


func LongestPath(game SnakeRequest, target Coord) []Coord {
		closedList := make(map[Coord]bool)
		openList := make([]Node, 0)
		pathTracker := make(map[Coord]Coord)

		me := game.You
		headPos := Node{Coord: me.Body[0], G: 0, H: 0, F: 0}
		openList = append(openList, headPos)

		for len(openList) > 0 {
			// BIGGEST F!!!
			var closeNode = openList[0]
			for _, openItem := range openList {
				if openItem.F > closeNode.F {
					closeNode = openItem
				}
			}

			// put it on the closed list
			closedList[closeNode.Coord] = true
			openList = removeFromOpenList(closeNode, openList)
			// loop through leastNodes's adjacent tiles -- call them T
			closeNeighbours := GetAdjacentCoords(closeNode.Coord)

			for _, neighbour := range closeNeighbours {

				// 1. If T on the closed list, ignore it
				if closedList[neighbour] {
					continue
				}

				if neighbour == target {

					closedList[neighbour] = true

					path := make([]Coord, 0)
					path = append(path, target)
					path = append(path, neighbour)
					current := closeNode.Coord
					path = append(path, current)

					_, pathway := pathTracker[current]

					for ; pathway; _, pathway = pathTracker[current] {
						current = pathTracker[current]
						path = append(path, current)
					}

					return reverseCoords(path)
				}

				// 2. If T is not on the open list add it
				for _, item := range openList {
					if neighbour == item.Coord {
						if NodeBlocked(neighbour, game.Board.Snakes) == false && OnBoard(neighbour, game.Board.Height, game.Board.Width) {
							if (closeNode.G+1)+Dist(neighbour, target) > item.F {
								item.F = (closeNode.G + 1) + Dist(neighbour, target)
								item.G = closeNode.G + 1
								item.H = Dist(neighbour, target)
								item.ParentCoords = neighbour

								pathTracker[item.Coord] = closeNode.Coord
							}
						}
					}
				}

				var openNode = Node{
					Coord:        neighbour,
					G:            closeNode.G + 1,
					H:            Dist(neighbour, target),
					F:            (closeNode.G + 1) + (Dist(neighbour, target)),
					ParentCoords: closeNode.Coord,
				}

				pathTracker[neighbour] = closeNode.Coord
				openList = appendList(openNode, game.Board.Snakes, openList, game.Board.Height, game.Board.Width)

			}

		}
		return nil
}

func AStarBoardFromTo(game SnakeRequest, source Coord, target Coord) []Coord {
	closedList := make(map[Coord]bool)
	openList := make([]Node, 0)
	pathTracker := make(map[Coord]Coord)

	headPos := Node{Coord: source, G: 0, H: 0, F: 0}
	openList = append(openList, headPos)

	for len(openList) > 0 {
		// find the Node the least F on the open list
		var closeNode = openList[0]
		for _, openItem := range openList {
			if openItem.F < closeNode.F {
				closeNode = openItem
			}
		}

		// put it on the closed list
		closedList[closeNode.Coord] = true
		openList = removeFromOpenList(closeNode, openList)
		// loop through leastNodes's adjacent tiles -- call them T
		closeNeighbours := GetAdjacentCoords(closeNode.Coord)

		for _, neighbour := range closeNeighbours {

			// 1. If T on the closed list, ignore it
			if closedList[neighbour] {
				continue
			}

			if neighbour == target {

				closedList[neighbour] = true

				path := make([]Coord, 0)
				path = append(path, target)
				path = append(path, neighbour)
				current := closeNode.Coord
				path = append(path, current)

				_, pathway := pathTracker[current]

				for ; pathway; _, pathway = pathTracker[current] {
					current = pathTracker[current]
					path = append(path, current)
				}

				return reverseCoords(path)
			}

			// 2. If T is not on the open list add it
			for _, item := range openList {
				if neighbour == item.Coord {
					if NodeBlockedExceptTail(neighbour, game.Board.Snakes) == false && OnBoard(neighbour, game.Board.Height, game.Board.Width) {
						if (closeNode.G+1)+Dist(neighbour, target) < item.F {
							item.F = (closeNode.G + 1) + Dist(neighbour, target)
							item.G = closeNode.G + 1
							item.H = Dist(neighbour, target)
							item.ParentCoords = neighbour

							pathTracker[item.Coord] = closeNode.Coord
						}
					}
				}
			}

			var openNode = Node{
				Coord:        neighbour,
				G:            closeNode.G + 1,
				H:            Dist(neighbour, target), // Heuristic!
				F:            (closeNode.G + 1) + (Dist(neighbour, target)),
				ParentCoords: closeNode.Coord,
			}

			pathTracker[neighbour] = closeNode.Coord
			openList = appendList(openNode, game.Board.Snakes, openList, game.Board.Height, game.Board.Width)

		}

	}
	return nil
}

func AstarBoard(game SnakeRequest, target Coord) []Coord {
		closedList := make(map[Coord]bool)
		openList := make([]Node, 0)
		pathTracker := make(map[Coord]Coord)

		me := game.You
		headPos := Node{Coord: me.Body[0], G: 0, H: 0, F: 0}
		openList = append(openList, headPos)

		for len(openList) > 0 {
			// find the Node the least F on the open list
			var closeNode = openList[0]
			for _, openItem := range openList {
				if openItem.F < closeNode.F {
					closeNode = openItem
				}
			}

			// put it on the closed list
			closedList[closeNode.Coord] = true
			openList = removeFromOpenList(closeNode, openList)
			// loop through leastNodes's adjacent tiles -- call them T
			closeNeighbours := GetAdjacentCoords(closeNode.Coord)

			for _, neighbour := range closeNeighbours {

				// 1. If T on the closed list, ignore it
				if closedList[neighbour] {
					continue
				}

				if neighbour == target {

					closedList[neighbour] = true

					path := make([]Coord, 0)
					path = append(path, target)
					path = append(path, neighbour)
					current := closeNode.Coord
					path = append(path, current)

					_, pathway := pathTracker[current]

					for ; pathway; _, pathway = pathTracker[current] {
						current = pathTracker[current]
						path = append(path, current)
					}

					return reverseCoords(path)
				}

				// 2. If T is not on the open list add it
				for _, item := range openList {
					if neighbour == item.Coord {
						if NodeBlocked(neighbour, game.Board.Snakes) == false && OnBoard(neighbour, game.Board.Height, game.Board.Width) {
							if (closeNode.G+1)+Dist(neighbour, target) < item.F {
								item.F = (closeNode.G + 1) + Dist(neighbour, target)
								item.G = closeNode.G + 1
								item.H = Dist(neighbour, target)
								item.ParentCoords = neighbour

								pathTracker[item.Coord] = closeNode.Coord
							}
						}
					}
				}

				var openNode = Node{
					Coord:        neighbour,
					G:            closeNode.G + 1,
					H:            Dist(neighbour, target),
					F:            (closeNode.G + 1) + (Dist(neighbour, target)),
					ParentCoords: closeNode.Coord,
				}

				pathTracker[neighbour] = closeNode.Coord
				openList = appendList(openNode, game.Board.Snakes, openList, game.Board.Height, game.Board.Width)

			}

		}
		return nil
}

func Astar(BoardHeight int, BoardWidth int, me Snake, Snakes []Snake, target Coord) []Coord {
	closedList := make(map[Coord]bool)
	openList := make([]Node, 0)
	pathTracker := make(map[Coord]Coord)

	headPos := Node{Coord: me.Body[0], G: 0, H: 0, F: 0}
	openList = append(openList, headPos)

	for len(openList) > 0 {

		// find the Node the least F on the open list
		var closeNode = openList[0]
		for _, openItem := range openList {
			if openItem.F < closeNode.F {
				closeNode = openItem
			}
		}

		// put it on the closed list
		closedList[closeNode.Coord] = true
		openList = removeFromOpenList(closeNode, openList)
		// loop through leastNodes's adjacent tiles -- call them T
		closeNeighbours := GetAdjacentCoords(closeNode.Coord)

		for _, neighbour := range closeNeighbours {

			// 1. If T on the closed list, ignore it
			if closedList[neighbour] {
				continue
			}

			if neighbour == target {
				// WE FOUND A PATH TO THE TARGET! :D
				closedList[neighbour] = true

				path := make([]Coord, 0)
				path = append(path, target)
				path = append(path, neighbour)
				current := closeNode.Coord
				path = append(path, current)

				_, pathway := pathTracker[current]

				for ; pathway; _, pathway = pathTracker[current] {
					current = pathTracker[current]
					path = append(path, current)
				}

				return reverseCoords(path)
			}

			// 2. If T is not on the open list add it
			for _, item := range openList {
				if neighbour == item.Coord {
					// only if on Board, and not Blocked By Snakes
					if NodeBlocked(neighbour, Snakes) == false && OnBoard(neighbour, BoardHeight, BoardWidth) {
						if (closeNode.G+1)+Dist(neighbour, target) < item.F {
							item.F = (closeNode.G + 1) + Dist(neighbour, target)
							item.G = closeNode.G + 1
							item.H = Dist(neighbour, target)
							item.ParentCoords = neighbour

							pathTracker[item.Coord] = closeNode.Coord
						}
					}
				}
			}

			var openNode = Node{
				Coord:        neighbour,
				G:            closeNode.G + 1,
				H:            Dist(neighbour, target),
				F:            (closeNode.G + 1) + (Dist(neighbour, target)),
				ParentCoords: closeNode.Coord,
			}

			pathTracker[neighbour] = closeNode.Coord
			openList = appendList(openNode, Snakes, openList, BoardHeight, BoardWidth)

		}

	}
	return nil
}


// abs is build in math.abs
/* Dist to function in steps (int) */
func Dist(a Coord, b Coord) int {
	return int(math.Abs(float64(b.X-a.X)) + math.Abs(float64(b.Y-a.Y)))
}

/*
func GetBodies(snakes SnakesList) []Coord {
  list := make([]Coord, 0)
  for _, s := range snakes {
    list = append(list, s.Body...)
  }
  return list
}
*/

// closestFoodPoint
func minDistFood(headPos Coord, food []Coord) Coord {
	min := food[0]
	for _, f := range food {
		if Dist(min, headPos) < Dist(f, headPos) {
			min = f
		}
	}
	return min
}



// determines if the vertex is actually on the board
// copy ownfunction
func OnBoard(vertex Coord, boardHeight int, boardWidth int) bool {
	if vertex.X >= 0 && vertex.X < boardWidth && vertex.Y >= 0 && vertex.Y < boardHeight {
		return true
	}

	return false
}

func NodeDangerous(game SnakeRequest, point Coord) bool {

	myLength := len(game.You.Body)

	for i := 0; i < len(game.Board.Snakes); i++ {

		if (game.Board.Snakes[i].ID != game.You.ID) { // not myself
			// i'm shorter dan this snake!
			if(myLength <= len(game.Board.Snakes[i].Body)) {
					theirMoves := GetAdjacentCoords(game.Board.Snakes[i].Body[0])
					for k := 0; k < len(theirMoves); k++ {
						if (theirMoves[k].X == point.X && theirMoves[k].Y == point.Y) {
							log.Print("This can be a headcrash")
							return true
						}
					}
			}
			//		closeNeighbours := GetAdjacentCoords(game.Board.Snakes[i].Body[0])
		//					for _, neighbour := range closeNeighbours { if point (return True) }
		} else {
			// it's me
		}

		for j := 0; j < len(game.Board.Snakes[i].Body); j++ {
			if game.Board.Snakes[i].Body[j].X == point.X && game.Board.Snakes[i].Body[j].Y == point.Y {
				if len(game.Board.Snakes[i].Body)-1 == j {

					if (game.Board.Snakes[i].Health > 99) {
						return true
					} else {
						//fmt.Print("Tail is safe... will move next turn")
						return false
					}

				}
				// it's in a snakes head.
				return true
			}
		}
	}

	return false
}

func NodeBlockedExceptTail(point Coord, Snakes []Snake) bool {
	for i := 0; i < len(Snakes); i++ {
		for j := 0; j < len(Snakes[i].Body); j++ {
			if Snakes[i].Body[j].X == point.X && Snakes[i].Body[j].Y == point.Y {
				if len(Snakes[i].Body)-1 == j {
					if (Snakes[i].Health > 99) {
						return true
					} else {
						fmt.Print("Tail is safe... will move next turn")
						return false
					}

				}

				return true
			}
		}
	}

	return false
}

// determines if the vertex is blocked
func NodeBlocked(point Coord, Snakes []Snake) bool {
	for i := 0; i < len(Snakes); i++ {
		for j := 0; j < len(Snakes[i].Body); j++ {
			if Snakes[i].Body[j].X == point.X && Snakes[i].Body[j].Y == point.Y {
				if len(Snakes[i].Body)-1 == j {
					return true
				}

				return true
			}
		}
	}

	return false
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

func Heading(startingPoint Coord, headingPoint Coord) string {
	if headingPoint.X > startingPoint.X {
		return "right"
	}
	if headingPoint.X < startingPoint.X {
		return "left"
	}
	if headingPoint.Y > startingPoint.Y {
		return "down"
	}
	if headingPoint.Y < startingPoint.Y {
		return "up"
	}

	return "up"
}

// NearestFood finds the closest food to the head of my snake
func NearestFood(FoodCoords []Coord, You Coord) Coord {
	var nearestFood = FoodCoords[0]
	var nearestFoodF = Dist(FoodCoords[0], You)

	for i := 0; i < len(FoodCoords); i++ {
		if Dist(FoodCoords[i], You) < nearestFoodF {
			nearestFood = FoodCoords[i]
			nearestFoodF = Dist(FoodCoords[i], You)
		}
	}

	return nearestFood
}
