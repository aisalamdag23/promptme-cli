package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aisalamdag23/promptme-cli/internal/domain"
	languagemodels "github.com/aisalamdag23/promptme-cli/internal/usecase/language_models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// careerCmd represents the career command
var careerCmd = &cobra.Command{
	Use:   "career",
	Short: "A command-line career coach that provides personalized advice and actionable tips for career growth, leveraging AI to answer prompts and guide users",
	Long: `This CLI app serves as your personal career coach, designed to offer tailored guidance 
for professional development. Whether you're exploring new job opportunities, preparing for 
interviews, or seeking career advice, this tool leverages AI to generate detailed, actionable 
responses based on your input. With support for caching, error handling, and real-time 
interaction, it ensures a smooth and efficient experience for users aiming to achieve their career goals.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wait, err := time.ParseDuration(fmt.Sprintf("%ds", cfg.General.ShutdownWaitSec))
		if err != nil {
			log.Fatal("time.parseduration:", err)
		}

		// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			processPrompt(ctx)
		}()

		// Block until signal received
		<-c
		log.Infoln("Shutting down server...")

		_, cancel = context.WithTimeout(ctx, wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.

		log.Infoln("Shutdown complete")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(careerCmd)
}

func processPrompt(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for { // Continuous loop to accept input
		fmt.Print("Enter your prompt: ")
		// read input text
		strPrompt, _ := reader.ReadString('\n')
		// for time tracking - start request
		resp := domain.Response{
			StartTime: time.Now(),
		}

		strPrompt = strings.TrimSpace(strPrompt)

		log = loggerWithMetadata(strPrompt)

		log.Debugln("prompt.started.at:", resp.StartTime)

		// call to llm
		srv, err := languagemodels.NewGenerator(ctx, cfg, log, inMemCache)
		if err != nil {
			log.Fatal("languagemodels.newgenerator.failed:", err)
		}
		// generate response - from cache or from llm
		resp.Text, err = srv.GenerateResponse(ctx, domain.UserPrompt{Text: strPrompt})
		resp.EndTime = time.Now()
		if err != nil {
			fmt.Println("Something went wrong. Try again")
			continue
		}
		// for time tracking - end request and return response
		resp.ResponseTime = resp.EndTime.Sub(resp.StartTime)

		// logs and debugs
		log.Debugln("response.returned.at:", resp.ResponseTime)
		log.Infoln("response.time:", resp.ResponseTime)

		// display response and time
		fmt.Println("Coach: " + resp.Text)
		fmt.Println("Response time: " + resp.ResponseTime.String())
	}

}

func loggerWithMetadata(strPrompt string) *logrus.Entry {
	reqID, _ := uuid.NewUUID()
	return log.WithFields(logrus.Fields{
		"user_prompt": strPrompt,
		"request_id":  reqID.String(),
	})
}
