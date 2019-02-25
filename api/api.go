package api

import (
  "encoding/json"
  "net/http"
)
/* via github.com/tkanos/gonfig */
type Configuration struct {
    Port       int
    Home_Route string
}

/* snakeAPi jay (with head and tail type) */
type Coord struct {
  X int `json:"x"`
  Y int `json:"y"`
}

// via Joram/COORDS.go
func (c Coord) Adjacent() []Coord {
	return []Coord{
		{c.X + 0, c.Y + 1},
		{c.X + 0, c.Y - 1},
		{c.X + 1, c.Y + 0},
		{c.X - 1, c.Y + 0},
	}
}

// need probably need such scaffolding. this for BFS/ A* / Pathing via Dijkstra
func (c Coord) SurroundingCoords() []Coord {
	return []Coord{
		{c.X + 0, c.Y + 1},
		{c.X + 0, c.Y - 1},
		{c.X + 1, c.Y + 0},
		{c.X - 1, c.Y + 0},

		{c.X + 1, c.Y + 1},
		{c.X - 1, c.Y - 1},
		{c.X + 1, c.Y - 1},
		{c.X - 1, c.Y + 1},

	}

}

// handy handy function to check NW,SE Corner for example
func (c Coord) Equal(other Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

type Snake struct {
  ID     string  `json:"id"`
  Name   string  `json:"name"`
  Health int     `json:"health"`
  Body   []Coord `json:"body"`
}

type Board struct {
  Height int     `json:"height"`
  Width  int     `json:"width"`
  Food   []Coord `json:"food"`
  Snakes []Snake `json:"snakes"`
}

type Game struct {
  ID string `json:"id"`
}

type SnakeRequest struct {
  Game  Game  `json:"game"`
  Turn  int   `json:"turn"`
  Board Board `json:"board"`
  You   Snake `json:"you"`
}

type StartResponse struct {
  Color string `json:"color,omitempty"`
  HeadType string `json:"headType,omitempty"`
  TailType string `json:"tailType,omitempty"`
}

type MoveResponse struct {
  Move string `json:"move"`
}

type CoordList []Coord

func (list *CoordList) UnmarshalJSON(data []byte) error {
  var obj struct {
    Data []Coord `json:"data"`
  }
  if err := json.Unmarshal(data, &obj); err != nil {
    return err
  }
  *list = obj.Data
  return nil
}

type SnakesList []Snake

func (list *SnakesList) UnmarshalJSON(data []byte) error {
  var obj struct {
    Data []Snake `json:"data"`
  }
  if err := json.Unmarshal(data, &obj); err != nil {
    return err
  }
  *list = obj.Data
  return nil
}

func DecodeSnakeRequest(req *http.Request, decoded *SnakeRequest) error {
  err := json.NewDecoder(req.Body).Decode(&decoded)
  return err
}
