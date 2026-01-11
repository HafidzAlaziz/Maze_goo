package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// GenerateRequest represents the expected payload for /generate
type GenerateRequest struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Difficulty string `json:"difficulty"`
}

// Coordinate represents a position in the maze
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// SolveRequest represents the expected payload for /solve
type SolveRequest struct {
	Grid  [][]int    `json:"grid"`
	Start Coordinate `json:"start"`
	End   Coordinate `json:"end"`
}

func main() {
	r := gin.Default()

	// CORS Setup to allow React (Vite usually uses 5173, standard React uses 3000)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Endpoint: POST /generate
	r.POST("/generate", func(c *gin.Context) {
		var req GenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
			return
		}

		fmt.Printf("Generating Maze: Width=%d, Height=%d, Difficulty=%s\n", req.Width, req.Height, req.Difficulty)

		// Ensure dimensions are odd
		if req.Width%2 == 0 {
			req.Width++
		}
		if req.Height%2 == 0 {
			req.Height++
		}

		maze := NewMaze(req.Width, req.Height)
		maze.Generate(req.Difficulty)

		c.JSON(http.StatusOK, gin.H{
			"width":  maze.Width,
			"height": maze.Height,
			"grid":   maze.Grid,
		})
	})

	r.POST("/solve", func(c *gin.Context) {
		var req SolveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grid or coordinates"})
			return
		}

		pathNodes := SolveAStar(req.Grid, req.Start.X, req.Start.Y, req.End.X, req.End.Y, nil)
		if pathNodes == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No path found"})
			return
		}

		path := make([]Coordinate, len(pathNodes))
		for i, n := range pathNodes {
			path[i] = Coordinate{X: n.X, Y: n.Y}
		}

		c.JSON(http.StatusOK, gin.H{"path": path})
	})

	// Jalankan di port 8080
	r.Run(":8080")
}
