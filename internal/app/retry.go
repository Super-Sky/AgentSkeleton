package app

type RetryDecision string

const (
	RetryDecisionAccept     RetryDecision = "accept"
	RetryDecisionRetry      RetryDecision = "retry"
	RetryDecisionUnresolved RetryDecision = "unresolved"
)

type RetryPolicy struct {
	MaxAutomaticRetries int `yaml:"max_automatic_retries" json:"max_automatic_retries"`
}

type RetryResult struct {
	Decision       RetryDecision `yaml:"decision" json:"decision"`
	Attempt        int           `yaml:"attempt" json:"attempt"`
	Remaining      int           `yaml:"remaining" json:"remaining"`
	ValidationErrs []string      `yaml:"validation_errors" json:"validation_errors"`
}

func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAutomaticRetries: 2,
	}
}

func EvaluateResponse(policy RetryPolicy, attempt int, response ResponseEnvelope) RetryResult {
	result := RetryResult{
		Attempt:   attempt,
		Remaining: max(policy.MaxAutomaticRetries-attempt, 0),
	}

	if err := response.Validate(); err == nil && response.Status == "ok" {
		result.Decision = RetryDecisionAccept
		return result
	} else {
		if err != nil {
			result.ValidationErrs = []string{err.Error()}
		} else if response.Status != "ok" {
			result.ValidationErrs = []string{"response status is not acceptable for downstream state updates"}
		}
	}

	if attempt < policy.MaxAutomaticRetries {
		result.Decision = RetryDecisionRetry
		return result
	}

	result.Decision = RetryDecisionUnresolved
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
