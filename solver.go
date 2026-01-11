package main

import (
	"container/heap"
	"fmt"
	"math"
)

// Node represents a coordinate in the grid for pathfinding.
type Node struct {
	X, Y    int
	G, H, F float64
	Parent  *Node
	Index   int // Required for container/heap
}

// PriorityQueue implements heap.Interface and holds Nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].F < pq[j].F }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := x.(*Node)
	n.Index = len(*pq)
	*pq = append(*pq, n)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Heuristic calculates Manhattan Distance.
func Heuristic(a, b *Node) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

// SolveAStar finds the shortest path from start to goal.
// callback is an optional function for real-time visualization.
func SolveAStar(grid [][]int, startX, startY, endX, endY int, callback func(closedSet map[string]bool, current *Node)) []*Node {
	startNode := &Node{X: startX, Y: startY}
	endNode := &Node{X: endX, Y: endY}

	openSet := &PriorityQueue{}
	heap.Init(openSet)
	heap.Push(openSet, startNode)

	closedSet := make(map[string]bool)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node)

		if callback != nil {
			callback(closedSet, current)
		}

		if current.X == endNode.X && current.Y == endNode.Y {
			// Path found, reconstruct it
			path := []*Node{}
			for current != nil {
				path = append(path, current)
				current = current.Parent
			}
			return path
		}

		closedSet[posKey(current.X, current.Y)] = true

		// Neighbors: Up, Down, Left, Right
		neighbors := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, d := range neighbors {
			nx, ny := current.X+d[0], current.Y+d[1]

			// Check boundaries and if it's a wall (1)
			if ny < 0 || ny >= len(grid) || nx < 0 || nx >= len(grid[0]) || grid[ny][nx] == 1 {
				continue
			}

			if closedSet[posKey(nx, ny)] {
				continue
			}

			gScore := current.G + 1
			neighbor := &Node{X: nx, Y: ny, Parent: current, G: gScore}
			neighbor.H = Heuristic(neighbor, endNode)
			neighbor.F = neighbor.G + neighbor.H

			// Simplified: just push to open set.
			// In a more robust implementation, we'd check if it's already in openSet with a lower G.
			heap.Push(openSet, neighbor)
		}
	}

	return nil // No path found
}

func posKey(x, y int) string {
	// Simple key for the closed set map
	return fmt.Sprintf("%d,%d", x, y)
}
