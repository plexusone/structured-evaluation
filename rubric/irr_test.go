package rubric

import (
	"testing"
)

func TestComputeIRR(t *testing.T) {
	t.Run("perfect agreement", func(t *testing.T) {
		pairs := []RatingPair{
			{Rater1: 5, Rater2: 5, Category: "cat1"},
			{Rater1: 4, Rater2: 4, Category: "cat2"},
			{Rater1: 3, Rater2: 3, Category: "cat3"},
		}

		metrics := ComputeIRR(pairs)

		if metrics.ExactAgreement != 1.0 {
			t.Errorf("ExactAgreement = %v, want 1.0", metrics.ExactAgreement)
		}

		if metrics.AdjacentAgreement != 1.0 {
			t.Errorf("AdjacentAgreement = %v, want 1.0", metrics.AdjacentAgreement)
		}

		if metrics.MeanAbsoluteDifference != 0.0 {
			t.Errorf("MeanAbsoluteDifference = %v, want 0.0", metrics.MeanAbsoluteDifference)
		}

		if metrics.PearsonCorrelation != 1.0 {
			t.Errorf("PearsonCorrelation = %v, want 1.0", metrics.PearsonCorrelation)
		}
	})

	t.Run("adjacent agreement", func(t *testing.T) {
		pairs := []RatingPair{
			{Rater1: 5, Rater2: 4, Category: "cat1"},
			{Rater1: 4, Rater2: 5, Category: "cat2"},
			{Rater1: 3, Rater2: 4, Category: "cat3"},
		}

		metrics := ComputeIRR(pairs)

		if metrics.ExactAgreement != 0.0 {
			t.Errorf("ExactAgreement = %v, want 0.0", metrics.ExactAgreement)
		}

		if metrics.AdjacentAgreement != 1.0 {
			t.Errorf("AdjacentAgreement = %v, want 1.0", metrics.AdjacentAgreement)
		}
	})

	t.Run("no agreement", func(t *testing.T) {
		pairs := []RatingPair{
			{Rater1: 5, Rater2: 1, Category: "cat1"},
			{Rater1: 1, Rater2: 5, Category: "cat2"},
		}

		metrics := ComputeIRR(pairs)

		if metrics.ExactAgreement != 0.0 {
			t.Errorf("ExactAgreement = %v, want 0.0", metrics.ExactAgreement)
		}

		if metrics.AdjacentAgreement != 0.0 {
			t.Errorf("AdjacentAgreement = %v, want 0.0", metrics.AdjacentAgreement)
		}

		if metrics.MeanAbsoluteDifference != 4.0 {
			t.Errorf("MeanAbsoluteDifference = %v, want 4.0", metrics.MeanAbsoluteDifference)
		}

		// Perfect negative correlation
		if metrics.PearsonCorrelation != -1.0 {
			t.Errorf("PearsonCorrelation = %v, want -1.0", metrics.PearsonCorrelation)
		}
	})

	t.Run("empty pairs", func(t *testing.T) {
		metrics := ComputeIRR([]RatingPair{})

		if metrics.SampleSize != 0 {
			t.Errorf("SampleSize = %v, want 0", metrics.SampleSize)
		}
	})
}

func TestComputeIRRFromResults(t *testing.T) {
	t.Run("with numeric scores", func(t *testing.T) {
		results1 := []CategoryResult{
			*NewCategoryResultWithNumeric("cat1", ScorePass, 5.0, ""),
			*NewCategoryResultWithNumeric("cat2", ScorePartial, 3.0, ""),
		}
		results2 := []CategoryResult{
			*NewCategoryResultWithNumeric("cat1", ScorePass, 4.0, ""),
			*NewCategoryResultWithNumeric("cat2", ScorePartial, 3.0, ""),
		}

		metrics := ComputeIRRFromResults(results1, results2)

		if metrics.SampleSize != 2 {
			t.Errorf("SampleSize = %v, want 2", metrics.SampleSize)
		}

		// One exact match (cat2: 3 vs 3)
		if metrics.ExactAgreement != 0.5 {
			t.Errorf("ExactAgreement = %v, want 0.5", metrics.ExactAgreement)
		}

		// Both are adjacent (diff <= 1)
		if metrics.AdjacentAgreement != 1.0 {
			t.Errorf("AdjacentAgreement = %v, want 1.0", metrics.AdjacentAgreement)
		}
	})

	t.Run("with categorical only", func(t *testing.T) {
		results1 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScoreFail, ""),
		}
		results2 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePartial, ""),
		}

		metrics := ComputeIRRFromResults(results1, results2)

		if metrics.SampleSize != 2 {
			t.Errorf("SampleSize = %v, want 2", metrics.SampleSize)
		}

		// One exact match (cat1: pass=5 vs pass=5)
		if metrics.ExactAgreement != 0.5 {
			t.Errorf("ExactAgreement = %v, want 0.5", metrics.ExactAgreement)
		}
	})

	t.Run("mismatched categories", func(t *testing.T) {
		results1 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat3", ScorePass, ""),
		}
		results2 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePass, ""),
		}

		metrics := ComputeIRRFromResults(results1, results2)

		// Only cat1 matches
		if metrics.SampleSize != 1 {
			t.Errorf("SampleSize = %v, want 1", metrics.SampleSize)
		}
	})
}

func TestComputeCategoricalAgreement(t *testing.T) {
	t.Run("perfect agreement", func(t *testing.T) {
		results1 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePartial, ""),
			*NewCategoryResult("cat3", ScoreFail, ""),
		}
		results2 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePartial, ""),
			*NewCategoryResult("cat3", ScoreFail, ""),
		}

		agreement := ComputeCategoricalAgreement(results1, results2)

		if agreement.ExactAgreement != 1.0 {
			t.Errorf("ExactAgreement = %v, want 1.0", agreement.ExactAgreement)
		}

		if agreement.SampleSize != 3 {
			t.Errorf("SampleSize = %v, want 3", agreement.SampleSize)
		}

		// Check confusion matrix
		if agreement.ConfusionMatrix["pass:pass"] != 1 {
			t.Errorf("ConfusionMatrix[pass:pass] = %v, want 1", agreement.ConfusionMatrix["pass:pass"])
		}
	})

	t.Run("mixed agreement", func(t *testing.T) {
		results1 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePass, ""),
		}
		results2 := []CategoryResult{
			*NewCategoryResult("cat1", ScorePass, ""),
			*NewCategoryResult("cat2", ScorePartial, ""),
		}

		agreement := ComputeCategoricalAgreement(results1, results2)

		if agreement.ExactAgreement != 0.5 {
			t.Errorf("ExactAgreement = %v, want 0.5", agreement.ExactAgreement)
		}

		// Check confusion matrix shows disagreement
		if agreement.ConfusionMatrix["pass:partial"] != 1 {
			t.Errorf("ConfusionMatrix[pass:partial] = %v, want 1", agreement.ConfusionMatrix["pass:partial"])
		}
	})
}

func TestGetNumericOrCategorical(t *testing.T) {
	t.Run("uses numeric if available", func(t *testing.T) {
		cr := *NewCategoryResultWithNumeric("test", ScorePartial, 4.2, "")
		score := getNumericOrCategorical(cr)
		if score != 4.2 {
			t.Errorf("score = %v, want 4.2", score)
		}
	})

	t.Run("converts pass to 5", func(t *testing.T) {
		cr := *NewCategoryResult("test", ScorePass, "")
		score := getNumericOrCategorical(cr)
		if score != 5.0 {
			t.Errorf("score = %v, want 5.0", score)
		}
	})

	t.Run("converts partial to 3", func(t *testing.T) {
		cr := *NewCategoryResult("test", ScorePartial, "")
		score := getNumericOrCategorical(cr)
		if score != 3.0 {
			t.Errorf("score = %v, want 3.0", score)
		}
	})

	t.Run("converts fail to 1", func(t *testing.T) {
		cr := *NewCategoryResult("test", ScoreFail, "")
		score := getNumericOrCategorical(cr)
		if score != 1.0 {
			t.Errorf("score = %v, want 1.0", score)
		}
	})
}
