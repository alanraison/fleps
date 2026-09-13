package model

func CalculatePlayerRoundScore(predictions GamePredictions, results Result) int {
	score := 0
	for _, pred := range predictions {
		if result, ok := results, true; ok {
			if pred.HomeGoals == result.HomeGoals && pred.AwayGoals == result.AwayGoals {
				score += 3
			} else if (pred.HomeGoals-result.HomeGoals)*(pred.AwayGoals-result.AwayGoals) > 0 {
				score += 1
			}
		}
	}
	return score
}