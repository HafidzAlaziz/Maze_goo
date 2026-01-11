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
func (m *Maze) Generate(difficulty string) {
	rand.Seed(time.Now().UnixNano())

	// Start carving from (1, 1)
	m.carve(1, 1)

	// Braiding for Hard mode
	if difficulty == "Hard" {
		m.braid(0.05) // Probabilitas dikurangi agar lebih terstruktur (5%)
	}
}

func (m *Maze) braid(p float64) {
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			if m.Grid[y][x] == 1 {
				pathNeighbors := 0
				if m.Grid[y-1][x] == 0 {
					pathNeighbors++
				}
				if m.Grid[y+1][x] == 0 {
					pathNeighbors++
				}
				if m.Grid[y][x-1] == 0 {
					pathNeighbors++
				}
				if m.Grid[y][x+1] == 0 {
					pathNeighbors++
				}

				if pathNeighbors >= 2 && rand.Float64() < p {
					m.Grid[y][x] = 0
				}
			}
		}
	}
}

func (m *Maze) carve(x, y int) {
	m.Grid[y][x] = 0 // Mark as path

	dirs := [][2]int{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}
	rand.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, d := range dirs {
		nx, ny := x+d[0], y+d[1]

		if nx > 0 && nx < m.Width-1 && ny > 0 && ny < m.Height-1 && m.Grid[ny][nx] == 1 {
			m.Grid[y+d[1]/2][x+d[0]/2] = 0
			m.carve(nx, ny)
		}
	}
}
