# PromptMe CLI

A command-line tool that leverages a language model API to generate text based on user prompts. 

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Improvements](#improvements)

---

## Features
- **Clean and Modular Design**  
  This project uses Clean Architecture principles to ensure a separation of concerns. Current layers are:  
  - **Domain**: Contains business logic services and core entities.
  - **Usecase**: Application-specific operations (GenerateResponse method for output).
  - **Infrastructure**: Caching, Config, Logging, Rate limit setup. These are pluggable components — caching for get and set operations, centralized config management, logging setup for easier log calls.

- **Config Validation**  
This ensures that all required fields in the configuration file are present and correctly formatted, enhancing the reliability and predictability of the project. The `go-playground/validator` package is used to validate the `.config.yml` file. If any required field is missing or malformed, an error will be raised thus preventing misconfiguration.

- **Structured logging**  
This project uses structured logging where log entries are presented as key-value pairs making filtering, searching, and analysis easy, especially in dubugging and integrating with analysis tools or monitoring systems like Datadog. This project uses [Logrus](https://github.com/sirupsen/logrus) which supports log levels, and easy contextual logging (request_ids, etc. are included in the log) for better traceability.

- **Cobra Package**  
This project uses [cobra package](https://github.com/spf13/cobra) because of:  
  - Ease of use
  - Comes with built-in features such as its support for commands, subcommands, and flag parsing, error handling, etc. reducing the need for boilerplate code
  - Adding new commands or features is straightforward (easily extensible), and
  - Cross-platform compatibility.
    
  To add a new subcommand (for example, you want to add `life` for life coach)
  1. Setup `cobra-cli` first — you can use [this guide](https://www.digitalocean.com/community/tutorials/how-to-use-the-cobra-package-in-go).
  2. Once the `cobra-cli` is ready, run the command below. This will create a new go file `life.go` in `cmd/` directory. For testing purposes, don't add any customizations.
      ```
      cobra-cli add life
      ```
    
  3. Run:
      ```
      go install
      ```
  4. Then run:
      ```
      promptme-cli life
      ```
      Successful creation of subcommand `life` should display:
      ```
      life called
      ```
- **Customizable or Swappable Language Model**  
Currently, this project is using Google Gemini but a different LM can be used by:  
  1. Adding a service layer abstraction (like the `geminiService` struct in `internal/usecase/language_models/gemini.go`), and a method/func `GenerateResponse(ctx context.Context, userPrompt domain.UserPrompt) (string, error)` in that service. For example, you want to add OpenAI (gpt). 
      ```
      type openAIService struct {
      	cfg     *config.Config
      	log     *logrus.Entry
      	client  *genai.Client
      	cache   caching.Cache
      	limiter *ratelimit.RateLimiter
      	...
      }
      func (s *openAIService) GenerateResponse(ctx context.Context, userPrompt domain.UserPrompt) (string, error) {
          ...
      }
      ```
  2. Adding a case in the `NewGenerator` func in `internal/usecase/language_models/language_model.go`. In our example, for code standard and beautification, add the const `LLM_PROVIDER_OPENAI = "openai"` in `internal/domain/shared.go`
      ```
        switch cfg.LLM.Provider {
    	...
        case domain.LLM_PROVIDER_OPENAI:
        // setup initial commands, and connection here etc.
    	...
    	}
      ```
  3. Change config default in `.config.yml` 
      ```
      spec:
        ...
        llm:
          provider: openai
          ...
      ```
---

## Installation
To install PromptMe CLI, you can use the following steps:

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/aisalamdag23/promptme-cli.git
   ```
2. Navigate to the project directory:
   ```bash
   cd promptme-cli
   ```
3. Run command, to generate a `.config.yml` file:
   ```bash
   make setup
   ```
4. Update necessary config values `api_key`, `keywords` (if restricted to a specific topic). For example, in this project I only want anything related to careers, so my keywords are: `career,job,occupation,profession,work,calling,employment,position,vocation` (I got these keywords from searching career synonym)
   ```bash
   header:
      specformat: "1.0"
    
    spec:
      general:
        api_key: SOME_APIKEY_MAKE_SURE_TO_EDIT_THIS
        graceful_shutdown_wait_time_sec: 3
        log_level: debug
      llm:
        provider: gemini
        keywords: career,job,occupation,profession,work,calling,employment,position,vocation
        gemini:
          model: gemini-1.5-flash
          max_requests_per_minute: 5
   ```
5. When you cloned this repository, an executable binary is already in here, you can run it:
   ```bash
   ./promptme-cli career
   ```
   Or you can run:
   ```bash
   go install
   ```
   then:
    ```bash
   promptme-cli career
   ```
   You should see a message saying `Enter your prompt: `, this will allow you to input strings / text

---
## Improvements  
Given more time, here are some improvements that I recommend for this project:  
1. **Unit and Integration Tests**: Implement automated tests to ensure that all components (e.g., language model services, caching, configuration parsing) work as expected. Add unit tests for core business logic and integration tests for the entire flow.
2. **Persistent Caching**: Instead of using in-memory caching, implement persistent caching (like Redis) to store previous responses for faster retrieval when the same prompt is requested.
3. **Better topic restriction setup**: Currently, the project is adding the restriction on the history, setting the topic everytime a prompt is entered/requested `cs.History = s.initKeywords()`. Instead of this approach, implement a more efficient and persistent topic restriction mechanism that initializes the keywords once during setup and applies the restriction dynamically without resetting the history each time a prompt is entered. This ensures better performance and avoids redundant operations while maintaining consistent topic filtering.
4. **Better Caching Logic**: Currently, the project saves only a single value for each matched key. However, a key phrase could have multiple possible values or answers. To improve this, enhance the caching logic to support storing multiple values for a single key. When a key is matched, it can either:
    - **Return all values**: Display all possible answers
    - **Selectively return a value**: Randomize or prioritize responses based on a scoring system, frequency of use, or other configurable logic.
