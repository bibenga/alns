package main

import (
	"cmp"
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	"github.com/bibenga/alns"
)

func main() {
	// https://alns.readthedocs.io/en/latest/examples/travelling_salesman_problem.html
	// go run examples/tsp/main.go && neato -Tpng examples/tsp/tsp.dot -o examples/tsp/tsp.png

	dists := dists(Coords)
	nodes := make([]int, len(Coords))
	for i := range len(Coords) {
		nodes[i] = i
	}

	// rnd := rand.New(rand.NewPCG(1, 2))
	rnd := alns.RuntimeRand

	// var initSol alns.State
	initSol := NewTspState(nodes, map[int]int{}, dists)
	if initSolG, err := greedyRepair(initSol, rnd); err != nil {
		panic(err)
	} else {
		initSol = initSolG.(*TspState)
	}

	fmt.Println("optimal solution: 564")
	fmt.Printf("initial solution: %.4f\n", initSol.Objective())

	destroyOperatorNames := []string{"randomRemoval", "pathRemoval", "worstRemoval"}
	destroyOperators := []alns.Operator{randomRemoval, pathRemoval, worstRemoval}
	repairOperatorNames := []string{"greedyRepair"}
	repairOperators := []alns.Operator{greedyRepair}

	sel, err := alns.NewRouletteWheel(
		[4]float64{3, 2, 1, 0.5},
		0.8,
		len(destroyOperators),
		len(repairOperators),
		nil,
	)
	if err != nil {
		panic(err)
	}
	accept := alns.HillClimbing{}
	stop := alns.MaxRuntime{MaxRuntime: 1 * time.Second}

	a := alns.ALNS{
		Rnd:               rnd,
		CollectObjectives: false,
		DestroyOperators:  destroyOperators,
		RepairOperators:   repairOperators,
		Selector:          &sel,
		Acceptor:          &accept,
		Stop:              &stop,
		InitialSolution:   initSol,
	}

	// print progress
	started := time.Now()
	prevLoggedPercent := 0
	lastBestOutcome := initSol.Objective()
	a.Listener = func(outcome alns.Outcome, cand alns.State) error {
		if outcome == alns.Best {
			lastBestOutcome = cand.Objective()
		}
		elapsed := time.Since(started)
		percent := int(min(elapsed.Seconds()/stop.MaxRuntime.Seconds(), 1) * 10)
		if percent > prevLoggedPercent {
			prevLoggedPercent = percent
			fmt.Printf("\rprogress: %3d%%; lastBest: %.4f", percent*10, lastBestOutcome)
		}
		return nil
	}
	result, err := a.Iterate()
	if err != nil {
		panic(err)
	}

	fmt.Println("") // after progress we should make a new line because we use "\r"

	// print result
	// result := &a.Result
	statistics := &result.Statistics
	best := result.BestState.(*TspState)

	fmt.Printf("best solution: %.4f\n", best.Objective())

	fmt.Printf("statistics: IterationCount=%d; TotalRuntime=%s\n",
		statistics.IterationCount,
		statistics.TotalRuntime,
	)
	fmt.Println("  destroy operators")
	for i, name := range destroyOperatorNames {
		fmt.Printf("    %d: %14s; %s\n", i, name, statistics.DestroyOperatorCounts[i])
	}
	fmt.Println("  repair operators")
	for i, name := range repairOperatorNames {
		fmt.Printf("    %d: %14s; %s\n", i, name, statistics.RepairOperatorCounts[i])
	}
	if len(statistics.Objectives) > 0 {
		fmt.Println("objectives")
		for i, objective := range statistics.Objectives {
			elapsed := statistics.Runtimes[i]
			fmt.Printf("%4d: %12s - %.4f\n", i, elapsed, objective)
		}
	} else {
		fmt.Println("the objectives were not collected")
	}

	writeDotFile("examples/tsp/tsp.dot", Coords, best.edges)
}

var Coords = [][2]float64{
	{0, 13},
	{0, 26},
	{0, 27},
	{0, 39},
	{2, 0},
	{5, 13},
	{5, 19},
	{5, 25},
	{5, 31},
	{5, 37},
	{5, 43},
	{5, 8},
	{8, 0},
	{9, 10},
	{10, 10},
	{11, 10},
	{12, 10},
	{12, 5},
	{15, 13},
	{15, 19},
	{15, 25},
	{15, 31},
	{15, 37},
	{15, 43},
	{15, 8},
	{18, 11},
	{18, 13},
	{18, 15},
	{18, 17},
	{18, 19},
	{18, 21},
	{18, 23},
	{18, 25},
	{18, 27},
	{18, 29},
	{18, 31},
	{18, 33},
	{18, 35},
	{18, 37},
	{18, 39},
	{18, 41},
	{18, 42},
	{18, 44},
	{18, 45},
	{25, 11},
	{25, 15},
	{25, 22},
	{25, 23},
	{25, 24},
	{25, 26},
	{25, 28},
	{25, 29},
	{25, 9},
	{28, 16},
	{28, 20},
	{28, 28},
	{28, 30},
	{28, 34},
	{28, 40},
	{28, 43},
	{28, 47},
	{32, 26},
	{32, 31},
	{33, 15},
	{33, 26},
	{33, 29},
	{33, 31},
	{34, 15},
	{34, 26},
	{34, 29},
	{34, 31},
	{34, 38},
	{34, 41},
	{34, 5},
	{35, 17},
	{35, 31},
	{38, 16},
	{38, 20},
	{38, 30},
	{38, 34},
	{40, 22},
	{41, 23},
	{41, 32},
	{41, 34},
	{41, 35},
	{41, 36},
	{48, 22},
	{48, 27},
	{48, 6},
	{51, 45},
	{51, 47},
	{56, 25},
	{57, 12},
	{57, 25},
	{57, 44},
	{61, 45},
	{61, 47},
	{63, 6},
	{64, 22},
	{71, 11},
	{71, 13},
	{71, 16},
	{71, 45},
	{71, 47},
	{74, 12},
	{74, 16},
	{74, 20},
	{74, 24},
	{74, 29},
	{74, 35},
	{74, 39},
	{74, 6},
	{77, 21},
	{78, 10},
	{78, 32},
	{78, 35},
	{78, 39},
	{79, 10},
	{79, 33},
	{79, 37},
	{80, 10},
	{80, 41},
	{80, 5},
	{81, 17},
	{84, 20},
	{84, 24},
	{84, 29},
	{84, 34},
	{84, 38},
	{84, 6},
	{107, 27},
}

func dists(coords [][2]float64) [][]float64 {
	dist := make([][]float64, len(coords))
	for row, coord1 := range coords {
		dist[row] = make([]float64, len(coords))
		for col, coord2 := range coords {
			dist[row][col] = euclidean(coord1[0], coord1[1], coord2[0], coord2[1])
		}
	}
	return dist
}

func euclidean(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

type TspState struct {
	nodes []int
	edges map[int]int
	dists [][]float64
}

var _ alns.State = &TspState{}

func NewTspState(nodes []int, edges map[int]int, dists [][]float64) *TspState {
	return &TspState{
		nodes: nodes,
		edges: edges,
		dists: dists,
	}
}

func (s *TspState) Clone() *TspState {
	return &TspState{
		nodes: s.nodes,
		edges: maps.Clone(s.edges),
		dists: s.dists,
	}
}

func (s *TspState) Objective() float64 {
	v := 0.0
	for node := range s.edges {
		v += s.dists[node][s.edges[node]]
	}
	return v
}

func greedyRepair(state alns.State, rnd *rand.Rand) (alns.State, error) {
	current := state.(*TspState)

	visited := slices.Collect(maps.Values(current.edges))

	shuffledIndices := rnd.Perm(len(current.nodes))
	nodes := make([]int, len(shuffledIndices))
	for i, ni := range shuffledIndices {
		nodes[i] = current.nodes[ni]
	}

	for len(current.edges) != len(current.nodes) {
		var node = -1
		for _, other := range nodes {
			if _, ok := current.edges[other]; !ok {
				node = other
				break
			}
		}
		if node == -1 {
			panic(fmt.Errorf("node not found"))
		}

		var unvisited []int
		for _, other := range current.nodes {
			if other != node && !slices.Contains(visited, other) && !wouldFormSubcycle(node, other, current) {
				unvisited = append(unvisited, other)
			}
		}
		if len(unvisited) == 0 {
			panic(fmt.Errorf("unvisited is empty"))
		}

		nearest := slices.MinFunc(unvisited, func(a, b int) int {
			return cmp.Compare(current.dists[node][a], current.dists[node][b])
		})

		current.edges[node] = nearest
		visited = append(visited, nearest)
	}

	return state, nil
}

func wouldFormSubcycle(fromNode, toNode int, state *TspState) bool {
	for step := 1; step < len(state.nodes); step++ {
		if _, ok := state.edges[toNode]; !ok {
			return false
		}
		toNode = state.edges[toNode]
		if fromNode == toNode && step != len(state.nodes)-1 {
			return true
		}
	}
	return false
}

const DegreeOfDestruction = 0.1

func edgesToRemove(state *TspState) int {
	return int(float64(len(state.edges)) * DegreeOfDestruction)
}

func randomRemoval(state alns.State, rnd *rand.Rand) (alns.State, error) {
	destroyed := state.(*TspState).Clone()

	toRemove := edgesToRemove(destroyed)

	removed := 0
	for removed != toRemove {
		idx := rnd.IntN(len(destroyed.nodes))
		node := destroyed.nodes[idx]
		if _, ok := destroyed.edges[node]; ok {
			removed++
			delete(destroyed.edges, node)
		}
	}

	return destroyed, nil
}

func pathRemoval(state alns.State, rnd *rand.Rand) (alns.State, error) {
	destroyed := state.(*TspState).Clone()

	nodeIdx := rnd.IntN(len(destroyed.nodes))
	node := destroyed.nodes[nodeIdx]

	toRemove := edgesToRemove(destroyed)

	for range toRemove {
		nextNode := destroyed.edges[node]
		delete(destroyed.edges, node)
		node = nextNode
	}

	return destroyed, nil
}

func worstRemoval(state alns.State, rnd *rand.Rand) (alns.State, error) {
	destroyed := state.(*TspState).Clone()

	worstEdges := slices.Clone(destroyed.nodes)
	slices.SortFunc(worstEdges, func(a, b int) int {
		return cmp.Compare(
			destroyed.dists[a][destroyed.edges[a]],
			destroyed.dists[b][destroyed.edges[b]],
		)
	})

	toRemove := edgesToRemove(destroyed)
	for idx := range toRemove {
		delete(destroyed.edges, worstEdges[len(worstEdges)-(idx+1)])
	}

	return destroyed, nil
}

func writeDotFile(filename string, nodes [][2]float64, edges map[int]int) {
	// scale up
	const k = 3
	nodes = slices.Clone(nodes)
	for n := range nodes {
		nodes[n] = [2]float64{nodes[n][0] * k, nodes[n][1] * k}
	}

	// search min and max
	minX, minY := nodes[0][0], nodes[0][1]
	maxX, maxY := minX, minY
	for _, n := range nodes {
		if n[0] < minX {
			minX = n[0]
		}
		if n[1] < minY {
			minY = n[1]
		}
		if n[0] > maxX {
			maxX = n[0]
		}
		if n[1] > maxY {
			maxY = n[1]
		}
	}
	width := maxX - minX
	if width == 0 {
		panic(fmt.Errorf("zero width"))
	}
	height := maxY - minY
	if height == 0 {
		panic(fmt.Errorf("zero height"))
	}

	// create and write dot file
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "digraph G {")
	fmt.Fprintf(f, "  graph [size=\"%f,%f!\", dpi=20.0];\n", width, height)

	fontsize := 48
	nodeSize := 2
	for i, n := range nodes {
		x := (n[0] - minX)
		y := (n[1] - minY)
		fmt.Fprintf(f,
			"  %d [label=\"%d\", fontsize=%d, pos=\"%f,%f!\", shape=circle, width=%d, height=%d, fixedsize=true];\n",
			i, i, fontsize, x, y, nodeSize, nodeSize)
	}

	arrowsize := 4
	penwidth := 3
	for from, to := range edges {
		fmt.Fprintf(f, "  %d -> %d [arrowsize=%d, penwidth=%d];\n", from, to, arrowsize, penwidth)
	}

	fmt.Fprintln(f, "}")
}
