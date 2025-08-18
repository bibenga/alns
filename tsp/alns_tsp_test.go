package tsp

import (
	"alns"
	"cmp"
	"log/slog"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

func TestTsp(t *testing.T) {
	a := alns.NewWithPCGRandom(1, 2)

	a.AddDestroyOperator(randomRemoval, "randomRemoval")
	a.AddDestroyOperator(pathRemoval, "pathRemoval")
	a.AddDestroyOperator(worstRemoval, "worstRemoval")

	a.AddRepairOperator(greedyRepair, "greedyRepair")

	a.OnOutcome = func(outcome alns.Outcome, cand alns.State) {
		if outcome == alns.BEST {
			slog.Debug("New best", "Objective", cand.Objective())
		}
	}

	coords := COORDS
	dists := dists(coords)

	nodes := make([]int, len(coords))
	for i := range len(coords) {
		nodes[i] = i
	}

	var initSol alns.State
	initSol = NewTspState(nodes, map[int]int{}, dists)
	initSol = greedyRepair(initSol, a.Rnd)

	slog.Info("initial solution", "Objective", initSol.Objective())

	sel := alns.NewRouletteWheel([]float64{3, 2, 1, 0.5}, 0.8, 3, 1, nil)
	accept := alns.HillClimbing{}
	stop := alns.MaxRuntime{MaxRuntime: 1 * time.Second}
	// stop := MaxIterations{MaxIterations: 10}
	result := a.Iterate(initSol, &sel, &accept, &stop)
	slog.Info("best solution", "Objective", result.BestState.Objective())

	slog.Info("statistics",
		"IterationCount", result.Statistics.IterationCount(),
		"TotalRuntime", result.Statistics.TotalRuntime(),
		"DestroyOperatorCounts", result.Statistics.DestroyOperatorCounts,
		"RepairOperatorCounts", result.Statistics.RepairOperatorCounts,
	)
}

var COORDS = [][2]float64{
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

func euclidean(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
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

func greedyRepair(state alns.State, rnd *rand.Rand) alns.State {
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
			panic(-1)
		}

		var unvisited []int
		for _, other := range current.nodes {
			if other != node && !slices.Contains(visited, other) && !wouldFormSubcycle(node, other, current) {
				unvisited = append(unvisited, other)
			}
		}
		if len(unvisited) == 0 {
			panic("len(unvisited) == 0")
		}

		nearest := slices.MinFunc(unvisited, func(a, b int) int {
			return cmp.Compare(current.dists[node][a], current.dists[node][b])
		})

		current.edges[node] = nearest
		visited = append(visited, nearest)
	}

	return state
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

func randomRemoval(state alns.State, rnd *rand.Rand) alns.State {
	destroyed := state.(*TspState).Clone()

	toRemove := edgesToRemove(destroyed)

	removed := 0
	for removed != toRemove {
		idx := rnd.IntN(len(destroyed.nodes))
		node := destroyed.nodes[idx]
		if _, ok := destroyed.edges[node]; ok {
			removed += 1
			delete(destroyed.edges, node)
		}
	}

	return destroyed
}

func pathRemoval(state alns.State, rnd *rand.Rand) alns.State {
	destroyed := state.(*TspState).Clone()

	nodeIdx := rnd.IntN(len(destroyed.nodes))
	node := destroyed.nodes[nodeIdx]

	toRemove := edgesToRemove(destroyed)

	for range toRemove {
		nextNode := destroyed.edges[node]
		delete(destroyed.edges, node)
		node = nextNode
	}

	return destroyed
}

func worstRemoval(state alns.State, rnd *rand.Rand) alns.State {
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

	return destroyed
}
