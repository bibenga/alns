# ALNS in Go (Partial Implementation)

This is a **partial port/adaptation** of the Python library **ALNS** (Adaptive Large Neighbourhood Search) to the **Go** programming language.  
The original implementation can be found here: [N-Wouda/ALNS](https://github.com/N-Wouda/ALNS).

## Overview
- Implements core components of the ALNS metaheuristic: **destroy operators**, **repair operators**, **acceptance criteria**, and the **operator selection mechanism**.
- Can be used to solve complex combinatorial optimization problems such as TSP, VRP, and others, similar to the Python version.


## Import the module in your Go project:
```go
import "github.com/bibenga/alns"
```

## Exmaple
```go
solver := alns.NewDefault()

solver.AddDestroyOperator(randomRemoval, "randomRemoval")
solver.AddDestroyOperator(pathRemoval, "pathRemoval")
solver.AddRepairOperator(greedyRepair, "greedyRepair")

initSolution := NewMyProblemState(...)

sel := alns.NewRouletteWheel([]float64{3, 2, 1, 0.5}, 0.8, 3, 1, nil)
accept := alns.HillClimbing{}
stop := alns.MaxIterations{MaxIterations: 10}
result := solver.Iterate(initSolution, &sel, &accept, &stop)

```