# Pofadder

Go-Pofadder! (Team: DutchViper)

## Badges


[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/jstolp/pofadder-go/?ref=repository-badge)


## Bugs

- When food spawns in a corner i.e. (0,0),(10,0),(0,10),(10,10) something bugs out and i'll commit suicide... (OOPS)

## roadmap

- Flood fill, in order to detect stupid moves in the future... https://en.m.wikipedia.org/wiki/Flood_fill
(consider each move of an enemySnake a blocked tile, just as the walls, enemyHead+Body but not it's tail... unless it just ate (100 Health)!)
Flood_fill the tiles both directions (usually you are on a tipping point, two directions still possible)
Goto the one with the highest count (There is more room to wiggle here...)
Allthough there are situations where a bigger count is not really correct. You could have one step indefinely (because you are chasing your tail)
but the other side has 3 steps available (only never grows) so that's something to watch out 

example;
https://codereview.stackexchange.com/questions/123581/golang-flood-fill

### Tournament day 2


## Changelog

### 2021-02-28

Updated to V1 of the BattleSnakeAPI

### [Upcoming Release] v0.9.0 Tournament day one

2019-12-10 ranking

position on leaderboard: 16/63 with 45.15 pts

### v0.8.7 (2019-11-25) FloodFill

hahaha what a joke this "floodFill"

It just detects a path to my tail in certail steps now....

It does so when above 2 or 3 valid moves.

if SafeMoves >= 2 { ChoosePathToTailIfFoundOnSafeNode }

not quite ready yet...

Score of today:  ( #12 / 52  @ 45.45 pt) ~ 45 point median

### v0.8.5 (2019-11-19) [Beta] - Back to Basics...

- Changed path to Longest Path To tail when 4 length or less... (MORE MOVEMENT in EARLY PHASE)

Check if a bigger snake is making a claim for it's food... else it's not safe.

updated safeFood Logic...( #25 / 57 @ 43.48 points)
FOOD == DANGER (for sure with a multi-snake party... let's avoid that...)

Back to Basics... Boring minimal stuff...

37/60 - 41 points... let's improve!


### v0.8.1 (2019-11-16)
Tournament Day...
### v0.8.0 (2019-11-12)

- Added "HeadHunterZZZ" capability (Shortest Path to FOOD ENEMY... let's go!)
Agressive snake.

### v0.6.0 (2019-11-12)
- Tagged before Changes
47. jstolp / GoPofAdderGames 37.58

### v0.5.5 (2019-11-09)
39. 37.16 (11-11-2019)
32. jstolp / GoPofAdder Pofadder 38.79 ranking before LONGEST PATH


### v0.5.x (2019-11-04 21:20)

v0.5.1 (adjusted number of snakes to hungry Threshold) testing:

--> BattleSnakeRank before commit: 2019-11-06: #26. jstolp / GoPofAdder - 38.87 points


v0.5
Done:
 - AStar search introduced
 - Crash into Own tail if health > 99 fixed

Todo:
- update Astar to use board coords
- update printGrid to cope with non-square board (Width x Height) instead of GridSize
- add escapeRoute (isFoodSafe) when determining closest food

### v0.2.x

v0.2
 - Self Awwareness... doesn't crash into tail.

Mark II - Self AWARE

### v0.1.x - De Gewone (Domme) PofAdder

v0.1.5
 2019-02-22 21:00
 Cleanup, NOOB Random moves.
v0.1.0

Mark I - Wall aware
  Finally one that didn't crash a wall...

v1 - Updated, so that the snake can avoid the wallzzzz... MARK-1

v1 - govendor, config added - 2019-02-21

v1 - start modding - 2019-02-21  

v0 - fork - 2019-02-20



----

## Good Reads, Helpfull resources


https://dave.cheney.net/practical-go/presentations/qcon-china.html


### misc

https://github.com/maximelamure/algorithms/blob/master/README.md

Hamilton Path.

Zap (uber.org) logging

answer this with source code: https://stackoverflow.com/questions/32999136/perfect-snake-ai

Hamilton path???
