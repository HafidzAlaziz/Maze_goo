package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ResetCursor memindahkan kursor ke posisi (0,0) di terminal.
func ResetCursor() {
	fmt.Print("\033[H")
}

// HideCursor menyembunyikan kursor terminal agar animasi lebih bersih.
func HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor menampilkan kembali kursor terminal.
func ShowCursor() {
	fmt.Print("\033[?25h")
}

// ClearTerminal awal untuk membersihkan layar sekali saja.
func ClearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// LivePrint menampilkan kondisi labirin secara in-place menggunakan ANSI escape codes.
func LivePrint(m *Maze, closedSet map[string]bool, current *Node) {
	ResetCursor()
	var sb strings.Builder

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			char := ' '
			if m.Grid[y][x] == 1 {
				char = '#' // Dinding
			} else {
				// Cek apakah ini posisi AI saat ini
				if current != nil && current.X == x && current.Y == y {
					char = '@'
				} else if closedSet[fmt.Sprintf("%d,%d", x, y)] {
					char = '.' // Dieksplorasi
				}
			}
			sb.WriteRune(char)
			sb.WriteRune(' ')
		}
		sb.WriteRune('\n')
	}

	fmt.Print(sb.String())
	time.Sleep(30 * time.Millisecond) // Jeda sedikit lebih cepat
}

// PrintTerminal prints the maze with the path to the console.
func PrintTerminal(m *Maze, path []*Node) {
	// Create a copy of the grid to mark the path
	displayGrid := make([][]rune, m.Height)
	for i := range m.Grid {
		displayGrid[i] = make([]rune, m.Width)
		for j := range m.Grid[i] {
			if m.Grid[i][j] == 1 {
				displayGrid[i][j] = '#' // Wall
			} else {
				displayGrid[i][j] = ' ' // Path
			}
		}
	}

	// Mark the path with '*'
	for _, n := range path {
		displayGrid[n.Y][n.X] = '*'
	}

	// Print the grid
	for i := range displayGrid {
		for j := range displayGrid[i] {
			fmt.Printf("%c ", displayGrid[i][j])
		}
		fmt.Println()
	}
}

// ExportPNG creates a PNG image of the maze with the path.
func ExportPNG(filename string, m *Maze, path []*Node, cellSize int) error {
	imgWidth := m.Width * cellSize
	imgHeight := m.Height * cellSize

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Colors
	wallColor := color.RGBA{0, 0, 0, 255}       // Black
	pathColor := color.RGBA{255, 255, 255, 255} // White
	solveColor := color.RGBA{255, 0, 0, 255}    // Red

	// Draw Maze
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			c := pathColor
			if m.Grid[y][x] == 1 {
				c = wallColor
			}
			fillCell(img, x, y, cellSize, c)
		}
	}

	// Draw Path
	for _, n := range path {
		fillCell(img, n.X, n.Y, cellSize, solveColor)
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	return nil
}

func fillCell(img *image.RGBA, x, y, size int, c color.Color) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			img.Set(x*size+i, y*size+j, c)
		}
	}
}
