package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Coord struct {
	x,y int
}

func parseCubes(input []string) [][][]Coord{
	cubes := [][][]Coord{}
	r1 := [][]Coord{}
	r2 := [][]Coord{}
	for y, line := range input {
		c1 := []Coord{}
		c2 := []Coord{}
		for x, chRaw := range line {
			ch := string(chRaw)
			if ch != " " {
				cor := Coord{x, y}
				if len(c1) == 50 {
					c2 = append(c2, cor)
				} else {
					c1 = append(c1, cor)
				}
			}
		}
		if len(c2) > 0 {
			r2 = append(r2, c2)
		}
		r1 = append(r1, c1)
		if len(r1) == 50 {
			if len(r2) > 0 {
				// cube := r2
				// fmt.Printf("New r2! rows: %d, columns: %d, firstVal: %d, secondVal: %d, lastVal:%d\n", len(cube), len(cube[0]), cube[0][0], cube[0][1], cube[49][49])
				cubes = append(cubes, r2)
				r2 = nil
			}
			cubes = append(cubes, r1)
			// cube := r1
			// fmt.Printf("New r1! rows: %d, columns: %d, firstVal: %d, secondVal: %d, lastVal:%d\n", len(cube), len(cube[0]), cube[0][0], cube[0][1], cube[49][49])
			r1 = nil
		}
	}
	return cubes
}

func parseBoard(input []string) (map[Coord]string, int){
	board := make(map[Coord]string)
	maxX := 0
	for i, line := range input {
		for j, chRaw := range line {
			ch := string(chRaw)
			if ch != " " {
				c := Coord{j,i}
				if j > maxX {
					maxX = j
				}
				board[c] = string(ch)
			}
		}
	}
	return board, maxX
}

func add(c1 Coord, c2 Coord) Coord {
	return Coord{c1.x+c2.x, c1.y+c2.y}
}

func wrapLine(cur Coord, dir Coord, board map[Coord]string, maxX int, maxY int) (Coord, bool) {
	// Need to wrap around
	runner := cur
	newCur := cur
	stop := false
	for k := 0; k < maxX + maxY; k++ {
		next := add(runner, dir)
		if next.y < 0 {
			next.y = maxY
		}
		if (next.x < 0) {
			next.x = maxX
		}
		next.y = next.y % (maxY+1)
		next.x = next.x % (maxX+1)
		if _, ok2 := board[next]; ok2 {
			if board[next] == "#" {
				stop = true
				break
			}
			if board[next] == "." {
				newCur = next
				break
			}
		}
		runner = next
	}
	return newCur, stop
}

type coordFlip func(Coord) Coord

type NewCube struct {
	cube int
	dir Coord
	coords coordFlip
}

func flip(x Coord,y Coord) (Coord, Coord){
	return x,y 
}

func wrapCube(outOfBounds Coord, cubes [][][]Coord, dir Coord, cube int, board map[Coord]string, max int) (Coord, int, Coord, bool) {
	// Hard coded mapping of the next cube, coordinate and direction when wrapping a cube
	// { [cube] : { [incomingDirection]: nextCube } }
	cubeDirections := map[int]map[Coord]NewCube{}
	up := Coord{0,-1}
	right := Coord{1,0}
	down := Coord{0,1}
	left := Coord{-1,0}
	// up, right, down, left
	cubeDirections[0] = map[Coord]NewCube{}
	cubeDirections[0][up] = NewCube{cube: 5, dir: up, coords: func (c Coord) Coord{return Coord{c.x, max} }}
	cubeDirections[0][right] = NewCube{cube: 3, dir: left, coords: func (c Coord) Coord{return Coord{max, max-c.y} }}
	cubeDirections[0][down] = NewCube{cube: 2, dir: left, coords: func (c Coord) Coord{return Coord{max, c.x} }}
	cubeDirections[0][left] = NewCube{cube: 1, dir: left, coords: func (c Coord) Coord{return Coord{max, c.y} }}
	cubeDirections[1] = map[Coord]NewCube{}
	cubeDirections[1][up] = NewCube{cube: 5, dir: right, coords: func (c Coord) Coord{return Coord{0, c.x} }}
	cubeDirections[1][right] = NewCube{cube: 0, dir: right, coords: func (c Coord) Coord{return Coord{0, c.y} }}
	cubeDirections[1][down] = NewCube{cube: 2, dir: down, coords: func (c Coord) Coord{return Coord{c.x, 0} }}
	cubeDirections[1][left] = NewCube{cube: 4, dir: right, coords: func (c Coord) Coord{return Coord{0, max-c.y} }}
	cubeDirections[2] = map[Coord]NewCube{}
	cubeDirections[2][up] = NewCube{cube: 1, dir: up, coords: func (c Coord) Coord{return Coord{c.x, max} }}
	cubeDirections[2][right] = NewCube{cube: 0, dir: up, coords: func (c Coord) Coord{return Coord{c.y, max} }}
	cubeDirections[2][down] = NewCube{cube: 3, dir: down, coords: func (c Coord) Coord{return Coord{c.x, 0} }}
	cubeDirections[2][left] = NewCube{cube: 4, dir: down, coords: func (c Coord) Coord{return Coord{c.y, 0} }}
	cubeDirections[3] = map[Coord]NewCube{}
	cubeDirections[3][up] = NewCube{cube: 2, dir: up, coords: func (c Coord) Coord{return Coord{c.x, max} }}
	cubeDirections[3][right] = NewCube{cube: 0, dir: left, coords: func (c Coord) Coord{return Coord{max, max-c.y} }}
	cubeDirections[3][down] = NewCube{cube: 5, dir: left, coords: func (c Coord) Coord{return Coord{max, c.x} }}
	cubeDirections[3][left] = NewCube{cube: 4, dir: left, coords: func (c Coord) Coord{return Coord{max, c.y} }}
	cubeDirections[4] = map[Coord]NewCube{}
	cubeDirections[4][up] = NewCube{cube: 2, dir: right, coords: func (c Coord) Coord{return Coord{0, c.x} }}
	cubeDirections[4][right] = NewCube{cube: 3, dir: right, coords: func (c Coord) Coord{return Coord{0, c.y} }}
	cubeDirections[4][down] = NewCube{cube: 5, dir: down, coords: func (c Coord) Coord{return Coord{c.x, 0} }}
	cubeDirections[4][left] = NewCube{cube: 1, dir: right, coords: func (c Coord) Coord{return Coord{0, max-c.y} }}
	cubeDirections[5] = map[Coord]NewCube{}
	cubeDirections[5][up] = NewCube{cube: 4, dir: up, coords: func (c Coord) Coord{return Coord{c.x, max} }}
	cubeDirections[5][right] = NewCube{cube: 3, dir: up, coords: func (c Coord) Coord{return Coord{c.y, max} }}
	cubeDirections[5][down] = NewCube{cube: 0, dir: down, coords: func (c Coord) Coord{return Coord{c.x, 0} }}
	cubeDirections[5][left] = NewCube{cube: 1, dir: down, coords: func (c Coord) Coord{return Coord{c.y, 0} }}

	nc := cubeDirections[cube][dir]
	newCur := nc.coords(outOfBounds)
	newCube := nc.cube
	newDir := nc.dir
	boardCoord := cubes[newCube][newCur.y][newCur.x]
	stop := board[boardCoord] == "#"
	// fmt.Printf("Stop: %s\n", stop)
	// if !stop {
	// 	fmt.Printf("prev: %d, dir: %d, cube: %d, newCur: %d, newCube: %d, newDir: %d\n", outOfBounds, dir, cube, newCur, newCube, newDir)
	// }

	return newCur, newCube, newDir, stop
}

func isInCube(coord Coord, max int) bool {
	if coord.x < 0 || coord.y < 0 || coord.x > max || coord.y > max {
		// fmt.Println("not in cube!")
		return false
	}
	return true
}

type seenKey struct {
	cube int
	dir Coord
}

func solveCube(cubes [][][]Coord, board map[Coord]string, instructions []string, max int) (Coord, Coord){
	curCube := 0
	cur := Coord{0,0} // Coordinate in cube
	dir := Coord{1,0}
	seen := map[seenKey]bool{}
	// fmt.Printf("Instructions: %s\n", instructions)
	for _, instruction := range instructions {
		// 90° clockwise rotation: (x,y) becomes (y,-x)
		// 90° counterclockwise rotation: (x,y) becomes (-y,x)
		// fmt.Printf("instruction: %s\n", instruction)
		if instruction == "R" {
			dir = Coord{-dir.y, dir.x}
		} else if instruction == "L"{
			dir = Coord{dir.y, -dir.x}
		} else {
			n,_ := strconv.Atoi(instruction)
			// fmt.Printf("cur: %d, n: %d\n", cur, n)
			for i := 0; i < n; i++ {
				next := add(cur, dir)
				// fmt.Printf("next: %d\n", next)
				ok := isInCube(next, max-1)
				if !ok {
					newCur, newCube, newDir, stop := wrapCube(cur, cubes, dir, curCube, board, max-1)
					k := seenKey{curCube, dir}
					_, seenOk := seen[k]
					if !seenOk {
						seen[k] = true
						fmt.Printf("prev: %d, dir: %d, cube: %d, newCur: %d, newCube: %d, newDir: %d\n", cur, dir, curCube, newCur, newCube, newDir)
					}
					if stop {
						break
					}
					if newCur != cur {
						cur = newCur
						curCube = newCube
						dir = newDir
						continue;
					}
				} else {
					nextBoardCoord := cubes[curCube][next.y][next.x]
					nextBoard := board[nextBoardCoord]
					// fmt.Printf("NextBoard: %s, boardCoord: %d, curcube: %d, cur: %d\n", nextBoard, nextBoardCoord, curCube, cur)
					if nextBoard == "#" {
						// Need to stop
						// fmt.Println("breaking!")
						break
					}
					if nextBoard == "." {
						cur = next
						// Continue
						continue
					}
				}
			}
		}
	}
	// Map cur to its board position
	fmt.Printf("cur: %d, curCube: %d, dir: %d\n", cur, curCube, dir)
	boardCur := cubes[curCube][cur.y][cur.x]
	return boardCur, dir
}


func solveBoard(board map[Coord]string, start Coord, instructions []string, maxX int, maxY int) (Coord, Coord) {
	cur := start
	dir := Coord{1,0}
	for _, instruction := range instructions {
		// 90° clockwise rotation: (x,y) becomes (y,-x)
		// 90° counterclockwise rotation: (x,y) becomes (-y,x)
		if instruction == "R" {
			dir = Coord{-dir.y, dir.x}
		} else if instruction == "L"{
			dir = Coord{dir.y, -dir.x}
		} else {
			n,_ := strconv.Atoi(instruction)
			for i := 0; i < n; i++ {
				next := add(cur, dir)
				_, ok := board[next];
				if !ok {
					newCur, stop := wrapLine(cur, dir, board, maxX, maxY)
					if stop {
						break
					}
					if newCur != cur {
						cur = newCur;
						continue;
					}
				} else {
					if board[next] == "#" {
						// Need to stop
						break
					}
					if board[next] == "." {
						cur = next
						// Continue
						continue
					}
				}
			}
		}
	}
	return cur, dir
}

func dirVal(dir Coord) int {
	if dir.x == 0 && dir.y == 1 {
		return 1
	}
	if dir.x == -1 && dir.y == 0 {
		return 2
	}
	if dir.x == 0 && dir.y == -1 {
		return 3
	}
	return 0
}

func main() {
	input, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	in := strings.Split(string(input), "\n")
	maxY := len(in)-2
	// Bottom two lines don't have the board
	board, maxX := parseBoard(in[:len(in)-2])
	cubes := parseCubes(in[:len(in)-2])
	// fmt.Printf("lenCube: %d\n", len(cubes))
	// for i, cube := range cubes {
	// 	fmt.Printf("cube: %d, rows: %d, columns: %d, first: %d, firstVal %s, second: %d, secondVal%s, last:%d\n", i, len(cube), len(cubes[i][0]), cubes[i][0][0],board[cubes[i][0][0]], cubes[i][0][1],board[cubes[i][0][1]], cubes[i][49][49])
	// 	// fmt.Println(cube)
	// }
	start := Coord{}
	for i := 0; i < len(in[0]); i++ {
		c := Coord{i,0}
		if _, ok := board[c]; ok {
			if board[c] == "." {
				start = c
				break
			}
		}
	}
	instructions := []string{}
	curNum := ""
	for _, c := range in[len(in)-1] {
		ch := string(c)
		if ch == "R" || ch == "L" {
			instructions = append(instructions, curNum)
			instructions = append(instructions, ch)
			curNum = ""
		} else {
			curNum += ch
		}
	}
	end, dir := solveBoard(board, start, instructions, maxX, maxY)
	partOneAnswer := 1000 * (end.y+1) + 4 * (end.x+1) + dirVal(dir)
	end2, dir2 := solveCube(cubes, board, instructions, 50)
	fmt.Printf("end2: %d, dir2: %d\n", end2, dir2)
	partTwoAnswer := 1000 * (end2.y+1) + 4 * (end2.x+1) + dirVal(dir2)
	fmt.Printf("Part One Answer: %d\n", partOneAnswer)
	fmt.Printf("Part Two Answer: %d\n", partTwoAnswer)
}