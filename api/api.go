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

type RootInfoResponse struct {
  ApiVersion string `json:"apiversion"`
  Author string `json:"author"`
  Color string `json:"color"`
  HeadType string `json:"head"`
  TailType string `json:"tail"`
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
