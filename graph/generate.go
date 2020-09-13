package graph

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

// Input are the inputs for each experiment
type Input struct {
	rounds, carCount, threads int
	trafficSlowdown           float64
	root, destination         string
}

// NewInput defines a new Input struct
func NewInput(rounds, carsCount, threads int,
	trafficSlowdown float64, root, destination string) *Input {
	input := new(Input)
	input.rounds = rounds
	input.carCount = carsCount
	input.threads = threads
	input.trafficSlowdown = trafficSlowdown
	input.root = root
	input.destination = destination
	return input
}

// Complete creates a complete graph with equal (1) weights on the verteces and equal (0) on edges
func Complete(n int, input *Input) *Graph {
	if n < 0 {
		panic("")
	}

	var vertexA, vertexB string
	graph := NewGraph(input)
	for i := 0; i < n; i++ {
		vertexA = strconv.Itoa(i)
		graph.AddVertex(NewNode(vertexA, 1))
		graph.adj[vertexA] = make(map[string]Adj)
		graph.visitedEdge[vertexA] = make(map[string]bool)
		for j := 0; j < i; j++ {
			vertexA = strconv.Itoa(i)
			vertexB = strconv.Itoa(j)
			graph.adj[vertexA][vertexB] = *NewAdj(0)
			graph.adj[vertexB][vertexA] = *NewAdj(0)
			graph.visitedEdge[vertexA][vertexB] = false
			graph.visitedEdge[vertexB][vertexA] = false
		}
	}
	return graph
}

// Complement return a new graph which is the complement of the input graph.
func (graph *Graph) Complement(input *Input) *Graph {
	graph.lock.RLock()
	defer graph.lock.RUnlock()

	graphC := NewGraph(input)
	for vertexA, edges := range graph.adj {
		for vertexB, weight := range edges {
			if graphC.adj[vertexB] == nil {
				graphC.adj[vertexB] = make(map[string]Adj)
			}
			graphC.adj[vertexB][vertexA] = weight
		}
	}
	return graphC
}

// NewRandom returns a graph with edges randomly placed throughout the graph
func NewRandom(maxVertices, edgeCount int,
	weightRange []float64, input *Input) *Graph {

	graph := NewGraph(input)
	min, max := weightRange[0], weightRange[1]
	var vertexA, vertexB string
	for i := 0; i < edgeCount; i++ {
		vertexA = strconv.Itoa(rand.Intn(maxVertices))
		vertexB = strconv.Itoa(rand.Intn(maxVertices))
		graph.AddEdges(NewNode(vertexA, 0), NewNode(vertexB, 0),
			rand.Float64()*max+min)
	}
	return graph
}

// NewCity creates a city like graph (e.g. each vertex is connected with betweem 2 - 6 other vertices)
func NewCity(intersections int, input *Input) *Graph {
	graph := NewGraph(input)
	var vertex string
	var crossroads int
	for i := 0; i < intersections; i++ {
		graph.AddVertex(NewNode(strconv.Itoa(i), 0))
	}
	for i := 0; i < intersections; i++ {
		crossroads = rand.Intn(4) + 2
		for j := 0; j < crossroads; j++ {
			vertex = strconv.Itoa(rand.Intn(intersections))
			graph.AddEdges(
				NewNode(strconv.Itoa(i), 0),
				NewNode(vertex, 0),
				1.0)
		}
	}
	return graph
}

// Drive simulates one round of cars driving and then updating the time it took to drive on each weighted edge
func (graph *Graph) Drive() ([][]float64, [][][]string) {
	carsCount, threads := graph.input.carCount, graph.input.threads
	trafficSlowdown, root := graph.input.trafficSlowdown, graph.input.root
	destination := graph.input.destination

	var wg sync.WaitGroup
	distances := make([][]float64, threads)
	paths := make([][][]string, threads)

	// Cars driving around
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		distances[i] = make([]float64, carsCount)
		paths[i] = make([][]string, carsCount)
		go func(i int) {
			for j := 0; j < carsCount; j++ {
				distances[i][j], paths[i][j] =
					graph.ShortestPath(root, destination)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Updating weight of edges of graph
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func(i int) {
			for j := 0; j < carsCount; j++ {
				pathLength := len(paths[i][j]) - 1
				for k := 0; k < pathLength; k++ {
					nodeA, nodeB := paths[i][j][k], paths[i][j][k+1]
					graph.UpdateEdge(nodeA, nodeB, trafficSlowdown)
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Reset any unvisited edges to 1
	for vertexA, edges := range graph.visitedEdge {
		for vertexB, visited := range edges {
			if visited {
				graph.saveWeight(vertexA, vertexB)
			} else {
				graph.resetWeight(vertexA, vertexB)
			}
		}
	}
	return distances, paths
}

func (graph *Graph) resetWeight(vertexA, vertexB string) {
	graph.adj[vertexA][vertexB] = Adj{weight: 1, prevWeight: 1}
}

func (graph *Graph) saveWeight(vertexA, vertexB string) {
	weight := graph.adj[vertexA][vertexB].weight
	graph.adj[vertexA][vertexB] = Adj{weight: 1, prevWeight: weight}
}

// Simulate is the driver to run the simulation some number of times
func (graph *Graph) Simulate() *Graph {
	for i := 0; i < graph.input.rounds; i++ {
		graph.Drive()
	}
	return graph
}

// GenerateExperiment chooses which experiment to run
func GenerateExperiment(size, i int, experimentType string, input *Input) *Graph {
	if experimentType == "city" {
		g := NewCity(size, input)
		g.id = i
		g.cityType = experimentType
		return g
	} else if experimentType == "complete" {
		g := Complete(size, input)
		g.id = i
		g.cityType = experimentType
		return g
	} else if experimentType == "random" {
		g := NewRandom(size, rand.Intn(size*size), []float64{0, 1}, input)
		g.id = i
		g.cityType = experimentType
		return g
	} else {
		valid := "Use city, complete, or random."
		panic(fmt.Sprintf("%v is not a valid experiment. %v",
			experimentType, valid))
	}
}
