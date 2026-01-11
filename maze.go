package main

import (
	"math/rand"
	"time"
)

// Maze structure represents the grid where 1 is Wall and 0 is Road.
type Maze struct {
	Width  int
	Height int
	Grid   [][]int
}

// NewMaze initializes a maze with all walls.
func NewMaze(width, height int) *Maze {
	// Ensure dimensions are odd for proper recursive backtracking
	if width%2 == 0 {
		width++
	}
	if height%2 == 0 {
		height++
	}

	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
		for j := range grid[i] {
			grid[i][j] = 1 // All walls initially
		}
	}

	return &Maze{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

// Generate uses Recursive Backtracking algorithm to carve paths.
func (m *Maze) Generate() {
	rand.Seed(time.Now().UnixNano())

	// Start carving from (1, 1)
	m.carve(1, 1)

	// Set Start (0,0) and Finish (Width-1, Height-1)
	// We ensure they are 0 (Path)
	m.Grid[0][0] = 0
	// Ensure a path from (0,0) to (1,1) if needed, 
	// but (1,1) is already path. Let's make (0,1) or (1,0) path too.
	m.Grid[0][1] = 0 // Entrance
	
	m.Grid[m.Height-1][m.Width-1] = 0
	m.Grid[m.Height-1][m.Width-2] = 0 // Exit
}

func (m *Maze) carve(x, y int) {
	m.Grid[y][x] = 0 // Mark as path

	// Directions: Left, Right, Up, Down (moving 2 units)
	dirs := [][2]int{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}
	
	// Shuffle directions
	rand.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]

		if nx > 0 && nx < m.Width-1 && ny > 0 && ny < m.Height-1 && m.Grid[ny][nx] == 1 {
			// Remove wall between current cell and next cell
			m.Grid[y+d[1]/2][x+d[0]/2] = 0
			m.carve(nx, ny)
		}
	}
}
