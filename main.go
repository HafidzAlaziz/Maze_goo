package main

import (
	"fmt"
	"log"
)

func main() {
	// Menampilkan Menu Kesulitan
	ClearTerminal()
	fmt.Println("========================================")
	fmt.Println("       MAZE SOLVER AI - PROYEK UAS      ")
	fmt.Println("========================================")
	fmt.Println("Pilih Tingkat Kesulitan:")
	fmt.Println("1. Mudah (11x11)")
	fmt.Println("2. Sedang (21x21)")
	fmt.Println("3. Sulit (51x21 - Landscape)")
	fmt.Print("\nMasukkan pilihan (1/2/3): ")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("Input tidak valid. Menggunakan tingkat Sedang.")
		choice = 2
	}

	// 1. Inisialisasi Dimensi Labirin Berdasarkan Pilihan
	width, height := 21, 21
	switch choice {
	case 1:
		width, height = 11, 11
	case 2:
		width, height = 21, 21
	case 3:
		width, height = 51, 21
	default:
		fmt.Println("Pilihan tidak ada. Menggunakan tingkat Sedang.")
	}

	// Persiapan Animasi
	HideCursor()
	defer ShowCursor()
	ClearTerminal()

	fmt.Printf("Menghasilkan Labirin Acak (%dx%d)...\n", width, height)
	fmt.Println("Tingkat Kesulitan:", map[int]string{1: "Mudah", 2: "Sedang", 3: "Sulit"}[choice])

	// 2. Maze Generation: Membuat labirin menggunakan algoritma Recursive Backtracking
	// Algoritma ini menjamin semua area terhubung dan memiliki satu jalur unik (sebelum kita modifikasi)
	maze := NewMaze(width, height)
	maze.Generate()

	// 3. Pathfinding: Mencari jalur terpendek menggunakan algoritma A* (A-Star) dengan Animasi
	// Algoritma ini menggunakan Heuristic Manhattan Distance untuk efisiensi pencarian
	fmt.Println("Mencari jalur tercepat dengan Algoritma A*...")

	// Callback untuk visualisasi live di terminal
	animationCallback := func(closedSet map[string]bool, current *Node) {
		LivePrint(maze, closedSet, current)
	}

	path := SolveAStar(maze.Grid, 0, 0, maze.Width-1, maze.Height-1, animationCallback)

	// Tampilkan kembali kursor dan geser ke bawah hasil animasi
	ShowCursor()
	fmt.Println("\n\nSelesai mencari jalur!")

	if path == nil {
		log.Fatal("Gagal menemukan jalur dari Start ke Finish.")
	}

	// 4. Output Visualisasi: Terminal (ASCII)
	// Berguna untuk pengecekan cepat di lingkungan command line
	fmt.Println("\nVisualisasi Akhir di Terminal:")
	PrintTerminal(maze, path)

	// 5. Output Visualisasi: File PNG
	// Mengekspor labirin ke gambar dengan jalur berwarna merah
	outputFile := "maze_result.png"
	cellSize := 20 // Ukuran piksel per sel labirin
	err = ExportPNG(outputFile, maze, path, cellSize)
	if err != nil {
		log.Fatalf("Gagal mengekspor gambar: %v", err)
	}

	fmt.Printf("\nBerhasil! Hasil labirin telah disimpan ke: %s\n", outputFile)
	fmt.Println("Jalur tercepat ditandai dengan warna MERAH pada file gambar.")
}
