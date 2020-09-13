package graph

import (
	hp "container/heap"
	"math/rand"
	"time"
)

// Shufflable interface to shuffle arrays
type Shufflable interface {
	Len() int
	Swap(i, j int)
}

// Shuffle array
func Shuffle(s Shufflable) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := s.Len(); n > 0; n-- {
		s.Swap(r.Intn(n), n-1)
	}
}

// Path holds the optimal distance and path for dijkstra
type Path struct {
	Distance float64
	Vertices []string
}

type minPath []Path

// NewPath creates a new path
func NewPath(root string) *Path {
	path := new(Path)
	path.Distance = 0
	path.Vertices = []string{root}
	return path
}

func (h minPath) Len() int {
	return len(h)
}

func (h minPath) Less(i, j int) bool {
	return h[i].Distance < h[j].Distance
}
func (h minPath) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minPath) Push(p interface{}) {
	*h = append(*h, p.(Path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Heap stored on array
type Heap struct {
	minPath *minPath
}

// NewHeap inits a new heap
func NewHeap(root string) *Heap {
	h := new(Heap)
	h.minPath = &minPath{}
	h.Push(*NewPath(root))
	return h
}

// Push pushes onto heap
func (h *Heap) Push(p Path) {
	hp.Push(h.minPath, p)
}

// Pop pops on heap
func (h *Heap) Pop() Path {
	return hp.Pop(h.minPath).(Path)
}

// ShortestPath uses Dijkstra to get from A to B
func (graph *Graph) ShortestPath(root, destination string) (float64, []string) {
	graph.lock.RLock()
	defer graph.lock.RUnlock()

	h := NewHeap(root)
	visited := make(map[string]bool)
	for len(*h.minPath) > 0 {
		H := h.Pop()
		vertex := H.Vertices[len(H.Vertices)-1]

		// Dijkstra guarantees best path as long as there is no negative cycles
		if !visited[vertex] {
			if vertex == destination {
				return H.Distance, H.Vertices
			}

			// Shuffle adj nodes
			neighborCount := len(graph.adj[vertex])
			neighbors := make([]string, neighborCount)
			i := 0
			for k := range graph.adj[vertex] {
				neighbors[i] = k
				i++
			}
			rand.Shuffle(neighborCount, func(i, j int) {
				neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
			})

			// Dijkstra on shuffled neighbors
			for _, adjVertex := range neighbors {
				dist := graph.adj[vertex][adjVertex].prevWeight
				if !visited[adjVertex] {
					h.Push(Path{Distance: H.Distance + dist,
						Vertices: append([]string{},
							append(H.Vertices, adjVertex)...)})
				}
			}

			visited[vertex] = true
		}
	}

	return 0, nil
}
