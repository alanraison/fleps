package model

func signum(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func CalculatePlayerRoundScore(predictions GamePredictions, results Result) int {
	score := 0
	for _, pred := range predictions {
		if result, ok := results, true; ok {
			if pred.HomeGoals == result.HomeGoals && pred.AwayGoals == result.AwayGoals {
				score += 3
			} else if (signum(result.HomeGoals-result.AwayGoals) == signum(pred.HomeGoals-pred.AwayGoals)) {
				score += 1
			}
		}
	}
	return score
}