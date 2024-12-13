package languagemodels

import (
	"context"
	"fmt"
	"time"

	"github.com/aisalamdag23/promptme-cli/internal/domain"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/caching"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/config"
	ratelimit "github.com/aisalamdag23/promptme-cli/internal/infrastructure/rate_limit"
	"github.com/google/generative-ai-go/genai"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
)

// NewGenerator determines which service to create based on the configuration/provider
func NewGenerator(ctx context.Context, cfg *config.Config, log *logrus.Entry, cache caching.Cache) (domain.LanguageModelService, error) {
	log.Infoln("newgenerator.provider:", cfg.LLM.Provider)
	// Check the model provider in the config
	switch cfg.LLM.Provider {
	case domain.LLM_PROVIDER_GEMINI:
		client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.General.APIKey))
		if err != nil {
			log.Errorln("genai.newclient.failed:", err)
			return nil, err
		}

		return &geminiService{
			cfg:     cfg,
			log:     log,
			client:  client,
			cache:   cache,
			limiter: ratelimit.NewRateLimiter(log, rate.Every(time.Minute), cfg.LLM.Gemini.MaxRequestsPerMinute),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported language model provider")
	}
}
