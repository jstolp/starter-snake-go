# Pofadder

De Gewone (Domme) PofAdder

## roadmap

https://github.com/maximelamure/algorithms/blob/master/README.md

Hamilton Path.

Zap (uber.org) logging

answer this with source code: https://stackoverflow.com/questions/32999136/perfect-snake-ai

https://en.m.wikipedia.org/wiki/A*_search_algorithm

A* search is what i need....

https://github.com/beefsack/go-astar


in order to do that... i need a priority queue in GO

https://youtu.be/dl4gycknzYY graphs in go

https://golang.org/src/container/heap/example_pq_test.go

https://www.redblobgames.com/pathfinding/a-star/introduction.html


## Changelog

### v0.3.0

- i need new logic to select a point to my Tail (else the infinite loop will never be done)
- there is still some bugs with the logic in selecting a point... maybe the ABS/DIST function is wrong.
- Got some Equal/Neighbor functions for Coords borrowed from joram/jsnek
- and i need some function to check a PATH (of coords) extract my body, other snake bodies and is not OOB
- and then take the LONGEST path to that node (hopefully food) or a POI (snake to kill?) (terrain to win?)

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
