# ALNS in Go 

This is a **partial port/adaptation** of the Python library [N-Wouda/ALNS](https://github.com/N-Wouda/ALNS) to the **Go** programming language.  

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
initSol := NewMyProblemState(...)

destroyOperators := []alns.Operator{randomRemoval, pathRemoval, worstRemoval}
repairOperators := []alns.Operator{greedyRepair}

selector, err := alns.NewRouletteWheel(
    [4]float64{3, 2, 1, 0.5},
    0.8,
    len(destroyOperators),
    len(repairOperators),
    nil,
)
if err != nil {
    ...
}
acceptor := alns.HillClimbing{}
stop := alns.MaxRuntime{MaxRuntime: 1 * time.Second}

a := alns.ALNS{
    Rnd:               rnd,
    DestroyOperators:  destroyOperators,
    RepairOperators:   repairOperators,
    Selector:          &selector,
    Acceptor:          &acceptor,
    Stop:              &stop,
    InitialSolution:   initSol,
}

if result, err := a.Iterate(); err != nil {
    ...
} else {
    // do something with result
}
```