# This repo contains my final project for parallel programming class in go.  
You can find the full report of the project along with design decisions and benchmarks on [report/report.pdf.](https://github.com/badbayesian/traffic-simulaion/blob/master/report/report.pdf)

## TL;Dr
### Data Structure and Algorithm
The traffic simulation is fundamentally a undirected graph with weighted edges where cars find a shortest path between two locations using Dijkstra. For the graph structure,I represented the edges as nested maps (hash tables) as I needed fast edge lookups and writes and could leverage that the graph structure did not change (more details in report). After each round, each weighted edge changed it weight with respect to the number of cars that used that edge as an analogy to traffic. Each experiment was iterated until a steady state(s) was reached, wherein the cars do not change paths or change paths in a cyclical manner.


### Concurrency
There are 3 main aspects of concurrency in the simulation.
1. Each car within the experiment runs concurrently with locks to synchronize updating the edge weights. An alternative version used Sync.Cond but it was slower.
2. Each experiment runs concurrently as an embarrassingly parallelizable process as each experiment is independent.
3. FanIn/FanOut infrastructure to manage concurrent elements of 1. and 2. 

## Running Traffic simulation

All datasets are generated from running the experiment

usage statement:

The experiment is defined in /main where a user runs go run simulation.go with the following flags.

`p = flag.Int("p", 4, "Number of threads")`

`n = flag.Int("n", 100, "Number of experiments")`

`r = flag.Int("r", 100, "Number of rounds in experiment")`

`c = flag.Int("c", 1000, "Number of cars per experiment")`

`s = flag.Float64("s", 1, "Penalty of traffic")`

`size = flag.Int("size", 10, "Size of graph")`

`t = flag.String("t", "city","Type of Experiment (complete, random, city)")`

`loc = flag.String("loc", "tmp", "Save location for json")`

Note that the results of the experiment (the weighted graphs) will be saved in data while the timing report from generate/process.batch will be saved in report/report.txt
