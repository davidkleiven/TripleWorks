package pkg

type Scorer func(a, b string) float64

type TopSelector struct {
	Num int
}

// Assigns each source to K targets
func (t *TopSelector) Select(sources []string, targets []string, scorer Scorer) [][]int {
	AssertGreater(len(targets), t.Num)

	scores := make([]float64, len(targets))
	result := make([][]int, len(sources))
	for i, source := range sources {
		for j, target := range targets {
			scores[j] = scorer(source, target)
		}
		result[i] = IndirectDescendingSort(scores)[:t.Num]
	}
	return result
}
