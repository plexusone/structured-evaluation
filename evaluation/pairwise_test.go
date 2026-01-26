package evaluation

import "testing"

func TestNewPairwiseComparison(t *testing.T) {
	p := NewPairwiseComparison("input", "output A", "output B")

	if p.Input != "input" {
		t.Errorf("Expected input 'input', got %s", p.Input)
	}
	if p.OutputA != "output A" {
		t.Errorf("Expected output A 'output A', got %s", p.OutputA)
	}
	if p.OutputB != "output B" {
		t.Errorf("Expected output B 'output B', got %s", p.OutputB)
	}
}

func TestPairwiseComparison_SetWinner(t *testing.T) {
	p := NewPairwiseComparison("input", "A", "B")
	p.SetWinner(WinnerA, "A is better because...", 0.9)

	if p.Winner != WinnerA {
		t.Errorf("Expected winner A, got %s", p.Winner)
	}
	if p.Confidence != 0.9 {
		t.Errorf("Expected confidence 0.9, got %f", p.Confidence)
	}
	if p.Reasoning != "A is better because..." {
		t.Errorf("Unexpected reasoning: %s", p.Reasoning)
	}
}

func TestPairwiseComparison_AddCategoryScore(t *testing.T) {
	p := NewPairwiseComparison("input", "A", "B")
	p.AddCategoryScore("accuracy", WinnerA, "A is more accurate", 0.3)
	p.AddCategoryScore("clarity", WinnerB, "B is clearer", 0.5)

	if len(p.CategoryScores) != 2 {
		t.Errorf("Expected 2 category scores, got %d", len(p.CategoryScores))
	}

	if p.CategoryScores[0].Winner != WinnerA {
		t.Errorf("Expected first category winner A, got %s", p.CategoryScores[0].Winner)
	}
}

func TestPairwiseComparison_SwappedComparison(t *testing.T) {
	p := NewPairwiseComparison("input", "A", "B")
	p.ID = "original"

	swapped := p.SwappedComparison()

	if swapped.OutputA != "B" {
		t.Errorf("Expected swapped output A to be 'B', got %s", swapped.OutputA)
	}
	if swapped.OutputB != "A" {
		t.Errorf("Expected swapped output B to be 'A', got %s", swapped.OutputB)
	}
	if swapped.Metadata["swapped_from"] != "original" {
		t.Errorf("Expected swapped_from metadata to be 'original'")
	}
}

func TestComputePairwiseResult(t *testing.T) {
	comparisons := []PairwiseComparison{
		{Winner: WinnerA},
		{Winner: WinnerA},
		{Winner: WinnerB},
		{Winner: WinnerTie},
	}

	result := ComputePairwiseResult(comparisons)

	if result.WinRateA != 0.5 {
		t.Errorf("Expected win rate A 0.5, got %f", result.WinRateA)
	}
	if result.WinRateB != 0.25 {
		t.Errorf("Expected win rate B 0.25, got %f", result.WinRateB)
	}
	if result.TieRate != 0.25 {
		t.Errorf("Expected tie rate 0.25, got %f", result.TieRate)
	}
	if result.OverallWinner != WinnerA {
		t.Errorf("Expected overall winner A, got %s", result.OverallWinner)
	}
}

func TestComputePairwiseResult_Tie(t *testing.T) {
	comparisons := []PairwiseComparison{
		{Winner: WinnerA},
		{Winner: WinnerB},
	}

	result := ComputePairwiseResult(comparisons)

	if result.OverallWinner != WinnerTie {
		t.Errorf("Expected overall winner tie, got %s", result.OverallWinner)
	}
}

func TestComputePairwiseResult_Empty(t *testing.T) {
	result := ComputePairwiseResult([]PairwiseComparison{})

	if result.WinRateA != 0 {
		t.Errorf("Expected 0 win rate for empty input, got %f", result.WinRateA)
	}
}
