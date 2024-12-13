package languagemodels

import (
	"context"
	"fmt"

	"github.com/aisalamdag23/promptme-cli/internal/domain"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/caching"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/config"
	ratelimit "github.com/aisalamdag23/promptme-cli/internal/infrastructure/rate_limit"
	"github.com/sirupsen/logrus"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
)

type geminiService struct {
	cfg     *config.Config
	log     *logrus.Entry
	client  *genai.Client
	cache   caching.Cache
	limiter *ratelimit.RateLimiter
}

// GenerateResponse checks cache for saved responses (if any), send message to gemini for uncached responses
func (s *geminiService) GenerateResponse(ctx context.Context, userPrompt domain.UserPrompt) (string, error) {
	// check cache
	strResp, exists := s.cache.Get(userPrompt.Text)
	if exists {
		s.log.Infoln("getcache.exists.value:", strResp)
		return strResp, nil
	}

	// start lm api
	strResp, err := s.generateGeminiResponse(ctx, userPrompt.Text)
	if err != nil {
		s.log.Errorln("generategeminiresponse.failed", err)
		return "", err
	}

	s.log.Infoln("setcache.newresponse.set")
	s.cache.Set(userPrompt.Text, strResp)

	return strResp, nil
}

// initKeywords init topics that the gemini response will only answer
func (s *geminiService) initKeywords() []*genai.Content {
	if s.cfg.LLM.Keywords != "" {
		s.log.Infoln("keyword.restriction.on:", s.cfg.LLM.Keywords)
		return []*genai.Content{
			{
				Parts: []genai.Part{
					genai.Text(fmt.Sprintf(domain.LLM_RESPONSE_CONTEXT, s.cfg.LLM.Keywords)),
				},
				Role: "user",
			},
		}
	}
	s.log.Infoln("keyword.restriction.off")
	return []*genai.Content{}
}

func (s *geminiService) generateGeminiResponse(ctx context.Context, strPrompt string) (string, error) {
	err := s.limiter.Do(ctx)
	if err != nil {
		s.log.Errorln("ratelimiter.do.failed", err)
		return "", err
	}

	model := s.client.GenerativeModel(s.cfg.LLM.Gemini.Model)
	cs := model.StartChat()
	cs.History = s.initKeywords()
	iter := cs.SendMessageStream(context.Background(), genai.Text(strPrompt))
	strResp := ""
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		if resp != nil {
			for _, cand := range resp.Candidates {
				if cand.Content != nil {
					for _, part := range cand.Content.Parts {
						strResp += fmt.Sprintf("%s", part)
					}
				}
			}
		}
	}

	return strResp, nil
}
