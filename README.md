# This repo contains my final project for parallel programming class in go.  
You can find the full report of the project along with design decisions on [report/report.pdf](https://github.com/badbayesian/traffic-simulaion/blob/master/report/report.pdf)

## Traffic simulation

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
