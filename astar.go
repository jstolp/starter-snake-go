package main

import (
	. "github.com/jstolp/pofadder-go/api"
	"math"
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

func Astar(BoardHeight int, BoardWidth int, me Snake, Snakes []Snake, target Coord) []Coord {
	closedList := make(map[Coord]bool)
	openList := make([]Node, 0)
	pathTracker := make(map[Coord]Coord)

	myHead := Node{Coord: me.Body[0], G: 0, H: 0, F: 0}
	openList = append(openList, myHead)

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


// ChaseTail returns the coordinate of the position behind my tail
func ChaseTail(You []Coord) Coord {
	return You[len(You)-1]
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



// determines if the square is actually on the board
// copy ownfunction
func OnBoard(square Coord, boardHeight int, boardWidth int) bool {
	if square.X >= 0 && square.X < boardWidth && square.Y >= 0 && square.Y < boardHeight {
		return true
	}

	return false
}


// determines if the square is blocked by a snake
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

// Heading determines the direction between two points - must be side by side
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