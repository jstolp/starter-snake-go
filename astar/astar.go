package astar

import (
	"github.com/jstolp/pofadder-go/api"
	"math"
)

/*
     ripped from alecj1240/alec-snake, astar
	G - the amount of steps its taken to get to that Node
	H - the heuristic - estimation from this Node to the destination
	F - the sum of G and H
	ParentCoords - the coordinates of the previous step (Node)
*/

// Node holds: coordinates, G, H, F, parent coords
type Node struct {
	Coord        api.Coord
	G            int
	H            int
	F            int
	ParentCoords api.Coord
}

func GetAdjacentCoords(Location api.Coord) []api.Coord {
	var adjacentCoords = make([]api.Coord, 0)
	adjacentCoords = append(adjacentCoords, api.Coord{X: Location.X + 1, Y: Location.Y})
	adjacentCoords = append(adjacentCoords, api.Coord{X: Location.X - 1, Y: Location.Y})
	adjacentCoords = append(adjacentCoords, api.Coord{X: Location.X, Y: Location.Y + 1})
	adjacentCoords = append(adjacentCoords, api.Coord{X: Location.X, Y: Location.Y - 1})

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

func appendList(appendingNode Node, Snakes []api.Snake, List []Node, BoardHeight int, BoardWidth int) []Node {
	if NodeBlocked(appendingNode.Coord, Snakes) == false && OnBoard(appendingNode.Coord, BoardHeight, BoardWidth) {
		List = append(List, appendingNode)
	}

	return List
}

// reverseCoords reverses the path of coordinates so it's in chronological order
func reverseCoords(path []api.Coord) []api.Coord {
	for a := 0; a < len(path)/2; a++ {
		b := len(path) - a - 1
		path[a], path[b] = path[b], path[a]
	}
	return path
}

func Astar(BoardHeight int, BoardWidth int, MySnake api.Snake, Snakes []api.Snake, Destination api.Coord) []api.Coord {
	closedList := make(map[api.Coord]bool)
	openList := make([]Node, 0)
	pathTracker := make(map[api.Coord]api.Coord)

	myHead := Node{Coord: MySnake.Body[0], G: 0, H: 0, F: 0}
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

			if neighbour == Destination {

				closedList[neighbour] = true

				path := make([]api.Coord, 0)
				path = append(path, Destination)
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
						if (closeNode.G+1)+Dist(neighbour, Destination) < item.F {
							item.F = (closeNode.G + 1) + Dist(neighbour, Destination)
							item.G = closeNode.G + 1
							item.H = Dist(neighbour, Destination)
							item.ParentCoords = neighbour

							pathTracker[item.Coord] = closeNode.Coord
						}
					}
				}
			}

			var openNode = Node{
				Coord:        neighbour,
				G:            closeNode.G + 1,
				H:            Dist(neighbour, Destination),
				F:            (closeNode.G + 1) + (Dist(neighbour, Destination)),
				ParentCoords: closeNode.Coord,
			}

			pathTracker[neighbour] = closeNode.Coord
			openList = appendList(openNode, Snakes, openList, BoardHeight, BoardWidth)

		}

	}
	return nil
}


// ChaseTail returns the coordinate of the position behind my tail
func ChaseTail(You []api.Coord) api.Coord {
	return You[len(You)-1]
}

// abs is build in math.abs

// Dist distance
func Dist(pointA api.Coord, pointB api.Coord) int {
	var DistX = Abs(pointB.X - pointA.X)
	var DistY = Abs(pointB.Y - pointA.Y)
	var DistDistance = DistX + DistY
	return DistDistance
}

// determines if the square is actually on the board
// copy ownfunction
func OnBoard(square api.Coord, boardHeight int, boardWidth int) bool {
	if square.X >= 0 && square.X < boardWidth && square.Y >= 0 && square.Y < boardHeight {
		return true
	}

	return false
}

// determines if the square is blocked by a snake
func SquareBlocked(point api.Coord, Snakes []api.Snake) bool {
	for i := 0; i < len(Snakes); i++ {
		for j := 0; j < len(Snakes[i].Body); j++ {
			if Snakes[i].Body[j].X == point.X && Snakes[i].Body[j].Y == point.Y {
				if len(Snakes[i].Body)-1 == j {
					return false
				}

				return true
			}
		}
	}

	return false
}

// Heading determines the direction between two points - must be side by side
func Heading(startingPoint api.Coord, headingPoint api.Coord) string {
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
func NearestFood(FoodCoords []api.Coord, You api.Coord) api.Coord {
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
