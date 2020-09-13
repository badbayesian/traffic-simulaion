package main

import (
	"flag"
	"math/rand"
	"proj3/graph"
	"time"
)

func fanIn(experiments, size int, experimentType string, input *graph.Input,
	p *int, done <-chan bool) <-chan *graph.Graph {
	tasks := make(chan *graph.Graph)
	go func() {
		defer close(tasks)
		for i := 0; i < experiments; i++ {
			select {
			case <-done:
				return
			case tasks <- graph.GenerateExperiment(size, i,
				experimentType, input):
			}
		}
	}()
	return tasks
}

// merge channels
func merge(tasks chan *graph.Graph, c <-chan *graph.Graph,
	doneB chan bool, done <-chan bool) {
	for i := range c {
		select {
		case <-done:
			return
		case tasks <- i:
		}
	}
	doneB <- true
}

func simulate(done <-chan bool, tasks <-chan *graph.Graph) <-chan *graph.Graph {
	taskOut := make(chan *graph.Graph)

	go func() {
		defer close(taskOut)
		for task := range tasks {
			select {
			case <-done:
				return
			case taskOut <- task.Simulate():
			}
		}
	}()
	return taskOut
}

//
func fanOut(done <-chan bool, p int,
	channels ...<-chan *graph.Graph) <-chan *graph.Graph {
	doneB := make(chan bool, p)
	tasks := make(chan *graph.Graph)

	for _, c := range channels {
		go merge(tasks, c, doneB, done)
	}

	go func() {
		defer close(tasks)
		for i := 0; i < len(channels); i++ {
			<-doneB
		}
	}()
	return tasks

}

func main() {
	var p = flag.Int("p", 4, "Number of threads")
	var n = flag.Int("n", 100, "Number of experiments")
	var r = flag.Int("r", 100, "Number of rounds in experiment")
	var c = flag.Int("c", 1000, "Number of cars per experiment")
	var s = flag.Float64("s", 1, "Penalty of traffic")
	var size = flag.Int("size", 10, "Size of graph")
	var t = flag.String("t", "city",
		"Type of Experiment (complete, random, city)")
	var loc = flag.String("loc", "tmp", "Save location for json")
	flag.Parse()

	splitC := *c / *p

	input := graph.NewInput(*r, splitC, *p, *s, "0", "1")

	rand.Seed(time.Now().UnixNano())
	done := make(chan bool)
	defer close(done)

	tasks := fanIn(*n, *size, *t, input, p, done)

	workers := make([]<-chan *graph.Graph, *p)
	for i := 0; i < *p; i++ {
		workers[i] = simulate(done, tasks)
	}

	for simulation := range fanOut(done, *p, workers...) {
		simulation.SaveJSON(*loc)
		//simulation.PrintEdges(true)
	}

}
