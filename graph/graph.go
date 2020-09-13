package graph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
)

// Node struct
type Node struct {
	name  string
	value float64
}

// NewNode creates a new vertex with weight
func NewNode(name string, value float64) *Node {
	node := new(Node)
	node.name = name
	node.value = value
	return node
}

// Adj defines the weight of an edge
type Adj struct {
	prevWeight float64
	weight     float64
}

// NewAdj creates a new weight
func NewAdj(weight float64) *Adj {
	return &Adj{weight: weight, prevWeight: weight}
}

// Graph struct
type Graph struct {
	lock        *sync.RWMutex
	adj         map[string]map[string]Adj
	vertex      map[string]float64
	visited     map[string]bool
	visitedEdge map[string]map[string]bool
	input       *Input
	id          int
	cityType    string
}

// NewGraph creates a simple graph
func NewGraph(input *Input) *Graph {
	graph := new(Graph)
	graph.lock = &sync.RWMutex{}
	graph.adj = make(map[string]map[string]Adj)
	graph.visited = make(map[string]bool)
	graph.vertex = make(map[string]float64)
	graph.visitedEdge = make(map[string]map[string]bool)
	graph.input = input
	return graph
}

// AddVertex adds vertex
func (graph *Graph) AddVertex(node *Node) {
	graph.lock.Lock()
	defer graph.lock.Unlock()

	graph.vertex[node.name] = node.value

	graph.adj[node.name] = map[string]Adj{}
	graph.visitedEdge[node.name] = map[string]bool{}
}

// AddEdges adds bidirectional edges with equal weights
func (graph *Graph) AddEdges(nodeA, nodeB *Node, weight float64) {
	graph.lock.Lock()
	defer graph.lock.Unlock()

	_, okA := graph.vertex[nodeA.name]
	_, okB := graph.vertex[nodeB.name]
	if !okA && !okB {
		graph.vertex[nodeA.name] = nodeA.value
		graph.vertex[nodeB.name] = nodeB.value

		graph.adj[nodeA.name] = map[string]Adj{nodeB.name: *NewAdj(weight)}
		graph.adj[nodeB.name] = map[string]Adj{nodeA.name: *NewAdj(weight)}

		graph.visitedEdge[nodeA.name] = map[string]bool{nodeB.name: false}
		graph.visitedEdge[nodeB.name] = map[string]bool{nodeA.name: false}
	} else if !okA {
		graph.vertex[nodeA.name] = nodeA.value
		graph.adj[nodeA.name] = map[string]Adj{nodeB.name: *NewAdj(weight)}
		graph.adj[nodeB.name][nodeA.name] = *NewAdj(weight)

		graph.visitedEdge[nodeA.name] = map[string]bool{nodeB.name: false}
		graph.visitedEdge[nodeB.name][nodeA.name] = false
	} else if !okB {
		graph.vertex[nodeB.name] = nodeB.value
		graph.adj[nodeB.name] = map[string]Adj{nodeA.name: *NewAdj(weight)}
		graph.adj[nodeA.name][nodeB.name] = *NewAdj(weight)

		graph.visitedEdge[nodeA.name][nodeB.name] = false
		graph.visitedEdge[nodeB.name] = map[string]bool{nodeA.name: false}
	} else {
		graph.adj[nodeA.name][nodeB.name] = *NewAdj(weight)
		graph.adj[nodeB.name][nodeA.name] = *NewAdj(weight)

		graph.visitedEdge[nodeA.name][nodeB.name] = false
		graph.visitedEdge[nodeB.name][nodeA.name] = false

	}
}

// UpdateEdge atomically updates edge
func (graph *Graph) UpdateEdge(nodeA, nodeB string, weight float64) {
	graph.lock.Lock()
	defer graph.lock.Unlock()

	// :( CONCURRENT MAPS ARENT ALLOWED NEED TO ADD THE LOCKS
	//atomic.AddInt64(graph.adj[nodeA][nodeB].weight, weight)
	pw, w := graph.adj[nodeA][nodeB].prevWeight, graph.adj[nodeA][nodeB].weight
	graph.adj[nodeA][nodeB] = Adj{weight: weight + w, prevWeight: pw}
	graph.visitedEdge[nodeA][nodeB] = true
}

// PrintEdges prints the edges with option to print the weights as well
func (graph *Graph) PrintEdges(withWeight bool) {
	graph.lock.RLock()
	defer graph.lock.RUnlock()
	var output string

	for vertex, edges := range graph.adj {
		output += vertex
		for edge, weight := range edges {
			output += " ->" + edge
			if withWeight {
				output += ":" + fmt.Sprintf("%d", int(weight.prevWeight))
			}
		}
		output += "\n"
	}
	fmt.Println(output)
}

// PrintVertices prints the vertices of the graph
func (graph *Graph) PrintVertices() {
	graph.lock.RLock()
	defer graph.lock.RUnlock()

	for vertex, weight := range graph.vertex {
		fmt.Printf("%v: %v\n", vertex, weight)
	}
}

// SaveJSON saves the adj matrix as a json
func (graph *Graph) SaveJSON(loc string) {
	graph.lock.RLock()
	defer graph.lock.RUnlock()

	fullLoc := fmt.Sprintf("../data/%v_%v_%v_cars_%v_threads_%v_size_%v.json", loc,
		graph.id, graph.cityType, graph.input.carCount, graph.input.threads, len(graph.vertex))
	outWriter, _ := os.Create(fullLoc)
	defer outWriter.Close()

	output := make(map[string]map[string]float64)
	for vertexA, edges := range graph.adj {
		output[vertexA] = make(map[string]float64)
		for vertexB, weight := range edges {
			output[vertexA][vertexB] = weight.prevWeight
		}
	}

	jsonString, _ := json.Marshal(output)
	ioutil.WriteFile(fullLoc, jsonString, 0644)
}

// Equal checks if two graphs are equal
func (graph *Graph) Equal(graphOther *Graph) bool {
	graph.lock.RLock()
	graphOther.lock.RLock()
	defer graph.lock.RUnlock()
	defer graphOther.lock.RUnlock()

	return reflect.DeepEqual(graph.adj, graphOther.adj)
}
