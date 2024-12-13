package domain

import "context"

type (
	// LanguageModelService handles the call to language model
	LanguageModelService interface {
		GenerateResponse(ctx context.Context, userPrompt UserPrompt) (string, error)
	}
)
