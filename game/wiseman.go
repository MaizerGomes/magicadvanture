package game

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const WiseManSettingKey = "wiseman_config"

const (
	WiseManProviderLegacy     = "legacy"
	WiseManProviderGemini     = "gemini"
	WiseManProviderCloudflare = "cloudflare"
)

const (
	DefaultGeminiModel     = "gemini-2.5-flash"
	DefaultCloudflareModel = "@cf/meta/llama-3.1-8b-instruct"
)

type WiseManConfig struct {
	Enabled   bool   `json:"enabled"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	APIKey    string `json:"api_key"`
	AccountID string `json:"account_id,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
}

type EducationalChallenge struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectIndex  int      `json:"correct_index"`
	CorrectAnswer string   `json:"correct_answer"`
}

type WiseManService struct {
	Store *Store
}

func NewWiseManService(store *Store) *WiseManService {
	return &WiseManService{Store: store}
}

func (w *WiseManService) GenerateEducationalChallenge(gs *GameState) (EducationalChallenge, error) {
	if w.Configured() {
		prompt := "Generate a multiple-choice educational question for a 9-year-old child. " +
			"Subjects can be math, history, science, or geography. " +
			"Format: Return ONLY a raw JSON object with 'question' (string), 'options' (array of 4 strings), 'correct_index' (number 0-3), and 'correct_answer' (string matching the option). " +
			"Do not include any other text or markdown blocks."

		if NormalizeLanguage(gs.Language) == LanguagePortuguese {
			prompt = "Gere uma pergunta educacional de múltipla escolha para uma criança de 9 anos. " +
				"Os assuntos podem ser matemática, história, ciência ou geografia. " +
				"Formato: Retorne APENAS um objeto JSON bruto com 'question' (string), 'options' (array de 4 strings), 'correct_index' (número 0-3) e 'correct_answer' (string correspondente à opção). " +
				"Não inclua nenhum outro texto ou blocos de markdown."
		}

		response, err := w.Ask(gs, prompt)
		if err == nil {
			clean := strings.TrimSpace(response)
			clean = strings.TrimPrefix(clean, "```json")
			clean = strings.TrimPrefix(clean, "```")
			clean = strings.TrimSuffix(clean, "```")
			clean = strings.TrimSpace(clean)

			var challenge EducationalChallenge
			if err := json.Unmarshal([]byte(clean), &challenge); err == nil && len(challenge.Options) == 4 {
				return challenge, nil
			}
		}
	}

	// Fallback to hardcoded list if AI fails or is not configured
	lang := gs.Language
	challenges := []EducationalChallenge{
		{
			Question:      "What is 12 + 15?",
			Options:       []string{"25", "27", "30", "22"},
			CorrectIndex:  1,
			CorrectAnswer: "27",
		},
		{
			Question:      "How many planets are in our solar system?",
			Options:       []string{"7", "8", "9", "10"},
			CorrectIndex:  1,
			CorrectAnswer: "8",
		},
		{
			Question:      "What is 8 x 7?",
			Options:       []string{"42", "54", "56", "64"},
			CorrectIndex:  2,
			CorrectAnswer: "56",
		},
		{
			Question:      "Who was the first person to walk on the moon?",
			Options:       []string{"Buzz Aldrin", "Yuri Gagarin", "Neil Armstrong", "John Glenn"},
			CorrectIndex:  2,
			CorrectAnswer: "Neil Armstrong",
		},
		{
			Question:      "What is 100 - 37?",
			Options:       []string{"63", "73", "53", "67"},
			CorrectIndex:  0,
			CorrectAnswer: "63",
		},
		{
			Question:      "Which is the largest ocean on Earth?",
			Options:       []string{"Atlantic", "Indian", "Arctic", "Pacific"},
			CorrectIndex:  3,
			CorrectAnswer: "Pacific",
		},
	}

	if NormalizeLanguage(lang) == LanguagePortuguese {
		challenges = []EducationalChallenge{
			{
				Question:      "Quanto é 12 + 15?",
				Options:       []string{"25", "27", "30", "22"},
				CorrectIndex:  1,
				CorrectAnswer: "27",
			},
			{
				Question:      "Quantos planetas existem no nosso sistema solar?",
				Options:       []string{"7", "8", "9", "10"},
				CorrectIndex:  1,
				CorrectAnswer: "8",
			},
			{
				Question:      "Quanto é 8 x 7?",
				Options:       []string{"42", "54", "56", "64"},
				CorrectIndex:  2,
				CorrectAnswer: "56",
			},
			{
				Question:      "Quem foi a primeira pessoa a caminhar na lua?",
				Options:       []string{"Buzz Aldrin", "Yuri Gagarin", "Neil Armstrong", "John Glenn"},
				CorrectIndex:  2,
				CorrectAnswer: "Neil Armstrong",
			},
			{
				Question:      "Quanto é 100 - 37?",
				Options:       []string{"63", "73", "53", "67"},
				CorrectIndex:  0,
				CorrectAnswer: "63",
			},
			{
				Question:      "Qual é o maior oceano da Terra?",
				Options:       []string{"Atlântico", "Índico", "Ártico", "Pacífico"},
				CorrectIndex:  3,
				CorrectAnswer: "Pacífico",
			},
		}
	}

	return challenges[time.Now().UnixNano()%int64(len(challenges))], nil
}

func (w *WiseManService) GenerateMillionaireChallenge(gs *GameState) (EducationalChallenge, error) {
	if w.Configured() {
		prompt := "Generate a CHALLENGING multiple-choice general knowledge question for a 'Who Wants to Be a Millionaire' game. " +
			"Subjects can be history, science, geography, or culture. " +
			"Format: Return ONLY a raw JSON object with 'question' (string), 'options' (array of 4 strings), 'correct_index' (number 0-3), and 'correct_answer' (string matching the option). " +
			"Do not include any other text or markdown blocks."

		if NormalizeLanguage(gs.Language) == LanguagePortuguese {
			prompt = "Gere uma pergunta de conhecimentos gerais DESAFIADORA para um jogo 'Quem Quer Ser um Milionário'. " +
				"Os assuntos podem ser história, ciência, geografia ou cultura. " +
				"Formato: Retorne APENAS um objeto JSON bruto com 'question' (string), 'options' (array de 4 strings), 'correct_index' (número 0-3) e 'correct_answer' (string correspondente à opção). " +
				"Não inclua nenhum outro texto ou blocos de markdown."
		}

		response, err := w.Ask(gs, prompt)
		if err == nil {
			clean := strings.TrimSpace(response)
			clean = strings.TrimPrefix(clean, "```json")
			clean = strings.TrimPrefix(clean, "```")
			clean = strings.TrimSuffix(clean, "```")
			clean = strings.TrimSpace(clean)

			var challenge EducationalChallenge
			if err := json.Unmarshal([]byte(clean), &challenge); err == nil && len(challenge.Options) == 4 {
				return challenge, nil
			}
		}
	}

	// Fallback
	return EducationalChallenge{
		Question:      "Which element has the atomic number 1?",
		Options:       []string{"Helium", "Hydrogen", "Oxygen", "Lithium"},
		CorrectIndex:  1,
		CorrectAnswer: "Hydrogen",
	}, nil
}

func (w *WiseManService) AskWiseManWithSearch(gs *GameState, question string) (string, error) {
	if !w.Configured() {
		return "", errors.New("wise man ai is not configured")
	}

	prompt := fmt.Sprintf("Search for the answer to this question and provide a clear, concise explanation: %s", question)
	if NormalizeLanguage(gs.Language) == LanguagePortuguese {
		prompt = fmt.Sprintf("Pesquise a resposta para esta pergunta e forneça uma explicação clara e concisa: %s", question)
	}

	return w.Ask(gs, prompt)
}

func (w *WiseManService) LoadConfig() (*WiseManConfig, error) {
	if w == nil || w.Store == nil {
		return nil, nil
	}

	raw, err := w.Store.GetSetting(WiseManSettingKey)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw) != "" {
		var cfg WiseManConfig
		if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
			return nil, err
		}
		cfg.normalize()
		if cfg.Enabled {
			return &cfg, nil
		}
	}

	return envWiseManConfig(), nil
}

func (w *WiseManService) SaveConfig(cfg *WiseManConfig) error {
	if w == nil || w.Store == nil {
		return errors.New("wise man store is not initialized")
	}
	if cfg == nil {
		return w.Store.DeleteSetting(WiseManSettingKey)
	}

	cfg.normalize()
	if !cfg.Enabled {
		return w.Store.DeleteSetting(WiseManSettingKey)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return w.Store.SetSetting(WiseManSettingKey, string(data))
}

func (w *WiseManService) ClearConfig() error {
	if w == nil || w.Store == nil {
		return errors.New("wise man store is not initialized")
	}
	return w.Store.DeleteSetting(WiseManSettingKey)
}

func (w *WiseManService) Configured() bool {
	cfg, err := w.LoadConfig()
	return err == nil && cfg != nil && cfg.Enabled && strings.TrimSpace(cfg.Provider) != ""
}

func (w *WiseManService) EnsureConfigured(reader *bufio.Reader) (string, error) {
	cfg, err := w.LoadConfig()
	if err != nil {
		return "", err
	}
	if cfg != nil && cfg.Enabled {
		return "", nil
	}

	fmt.Print("Would you like to wire an AI bot for the wise man now? (y/n): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		return "", nil
	}
	return w.ConfigureInteractive(reader)
}

func (w *WiseManService) ConfigureInteractive(reader *bufio.Reader) (string, error) {
	if reader == nil {
		return "", errors.New("reader is required")
	}

	current, _ := w.LoadConfig()
	if current != nil && current.Enabled {
		fmt.Printf("Current wise man setup: %s (%s)\n", current.providerLabel(), current.modelLabel())
	}

	fmt.Println("Choose a wise man provider:")
	fmt.Println("1) Gemini")
	fmt.Println("2) Cloudflare Workers AI")
	fmt.Print("Selection: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.ToLower(strings.TrimSpace(choice))

	cfg := &WiseManConfig{Enabled: true}
	switch choice {
	case "", "1", "gemini", "g":
		cfg.Provider = WiseManProviderGemini
		cfg.Model = DefaultGeminiModel
		const geminiURL = "https://aistudio.google.com/app/apikey"
		fmt.Println("Opening Google AI Studio in your browser...")
		if err := openURL(geminiURL); err != nil {
			fmt.Printf("Open this URL manually: %s\n", geminiURL)
		}
		fmt.Println("Create or copy your Gemini API key, then paste it below.")
		fmt.Print("Gemini API key: ")
		cfg.APIKey, _ = reader.ReadString('\n')
		cfg.APIKey = strings.TrimSpace(cfg.APIKey)
		fmt.Printf("Model [%s]: ", cfg.Model)
		model, _ := reader.ReadString('\n')
		model = strings.TrimSpace(model)
		if model != "" {
			cfg.Model = model
		}
	case "2", "cloudflare", "cloud", "cf":
		cfg.Provider = WiseManProviderCloudflare
		cfg.Model = DefaultCloudflareModel
		const cloudflareURL = "https://dash.cloudflare.com/profile/api-tokens"
		fmt.Println("Opening Cloudflare API token settings in your browser...")
		if err := openURL(cloudflareURL); err != nil {
			fmt.Printf("Open this URL manually: %s\n", cloudflareURL)
		}
		fmt.Println("Create or copy a Workers AI API token, then paste it below.")
		fmt.Print("Cloudflare API token: ")
		cfg.APIKey, _ = reader.ReadString('\n')
		cfg.APIKey = strings.TrimSpace(cfg.APIKey)
		fmt.Print("Cloudflare Account ID: ")
		cfg.AccountID, _ = reader.ReadString('\n')
		cfg.AccountID = strings.TrimSpace(cfg.AccountID)
		fmt.Printf("Model [%s]: ", cfg.Model)
		model, _ := reader.ReadString('\n')
		model = strings.TrimSpace(model)
		if model != "" {
			cfg.Model = model
		}
	default:
		return "", errors.New("wise man setup cancelled")
	}

	if strings.TrimSpace(cfg.APIKey) == "" {
		return "", errors.New("wise man api key cannot be empty")
	}
	if cfg.Provider == WiseManProviderCloudflare && strings.TrimSpace(cfg.AccountID) == "" {
		return "", errors.New("cloudflare account id cannot be empty")
	}

	if err := w.SaveConfig(cfg); err != nil {
		return "", err
	}
	return fmt.Sprintf("Wise man AI configured with %s.", cfg.providerLabel()), nil
}

func (w *WiseManService) Ask(gs *GameState, question string) (string, error) {
	cfg, err := w.LoadConfig()
	if err != nil {
		return "", err
	}
	if cfg == nil || !cfg.Enabled {
		return "", errors.New("wise man ai is not configured")
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return "", errors.New("question cannot be empty")
	}

	switch cfg.Provider {
	case WiseManProviderGemini:
		return askWiseManGemini(cfg, gs, question)
	case WiseManProviderCloudflare:
		return askWiseManCloudflare(cfg, gs, question)
	case WiseManProviderLegacy:
		return askWiseManLegacy(cfg, gs, question)
	default:
		return "", fmt.Errorf("unsupported wise man provider %q", cfg.Provider)
	}
}

func (w *WiseManService) Farewell(gs *GameState) (string, error) {
	if gs == nil {
		return "", errors.New("game state is nil")
	}
	if !w.Configured() {
		return "", errors.New("wise man ai is not configured")
	}

	question := "The player has died. Give a short, humble farewell speech in 2 or 3 sentences. Praise the player's journey, avoid revealing any secret about your true nature, and keep the tone gentle."
	if NormalizeLanguage(gs.Language) == LanguagePortuguese {
		question = "O jogador morreu. Faça um breve discurso de despedida em 2 ou 3 frases. Elogie a jornada do jogador, não revele nenhum segredo sobre sua verdadeira natureza e mantenha um tom gentil."
	}
	return w.Ask(gs, question)
}

func askWiseManGemini(cfg *WiseManConfig, gs *GameState, question string) (string, error) {
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", strings.TrimSpace(cfg.Model))
	payload := map[string]any{
		"systemInstruction": map[string]any{
			"parts": []map[string]string{{"text": wiseManSystemPrompt(gs)}},
		},
		"contents": []map[string]any{
			{
				"role":  "user",
				"parts": []map[string]string{{"text": question}},
			},
		},
		"generationConfig": map[string]any{
			"temperature":     0.8,
			"maxOutputTokens": 768,
			"topP":            0.95,
			"topK":            40,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payloadBytes))
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("gemini request failed: %s", msg)
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(payloadBytes, &result); err != nil {
		return "", err
	}
	for _, candidate := range result.Candidates {
		var parts []string
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				parts = append(parts, part.Text)
			}
		}
		reply := strings.TrimSpace(strings.Join(parts, "\n"))
		if reply != "" {
			return reply, nil
		}
	}
	return "", errors.New("gemini returned no response")
}

func askWiseManCloudflare(cfg *WiseManConfig, gs *GameState, question string) (string, error) {
	if strings.TrimSpace(cfg.AccountID) == "" {
		return "", errors.New("cloudflare account id is required")
	}
	endpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", strings.TrimSpace(cfg.AccountID), strings.TrimSpace(cfg.Model))
	payload := map[string]any{
		"prompt":     wiseManSystemPrompt(gs) + "\n\nPlayer question: " + question,
		"max_tokens": 768,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payloadBytes))
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("cloudflare request failed: %s", msg)
	}

	var result struct {
		Result struct {
			Response string `json:"response"`
		} `json:"result"`
	}
	if err := json.Unmarshal(payloadBytes, &result); err != nil {
		return "", err
	}
	reply := strings.TrimSpace(result.Result.Response)
	if reply == "" {
		return "", errors.New("cloudflare returned no response")
	}
	return reply, nil
}

func askWiseManLegacy(cfg *WiseManConfig, gs *GameState, question string) (string, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return "", errors.New("legacy wise man endpoint is not configured")
	}

	payload := map[string]any{
		"model": cfg.modelLabel(),
		"messages": []map[string]string{
			{"role": "system", "content": wiseManSystemPrompt(gs)},
			{"role": "user", "content": question},
		},
		"temperature": 0.8,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey := strings.TrimSpace(cfg.APIKey); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(payloadBytes))
		if msg == "" {
			msg = resp.Status
		}
		return "", fmt.Errorf("wise man request failed: %s", msg)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(payloadBytes, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("wise man returned no response")
	}
	reply := strings.TrimSpace(result.Choices[0].Message.Content)
	if reply == "" {
		return "", errors.New("wise man returned an empty response")
	}
	return reply, nil
}

func wiseManSystemPrompt(gs *GameState) string {
	room := "unknown"
	if gs != nil {
		room = gs.CurrentRoomID
	}
	return fmt.Sprintf(
		"You are a wise man in a fantasy terminal RPG. The player is in room %q. They may ask about locations, items, quests, or advice.\n"+
			"Reply in 2 to 5 short sentences as needed. Keep answers complete and coherent, never cut a thought in the middle.\n"+
			"Keep one line practical: if they ask where to find something, give one concrete in-game clue.\n"+
			"If they ask for a precious item or a hidden object, point them toward the harbor, lighthouse, tide caves, ruins, moonwell, observatory, or sage hut using a hint, not a spoiler dump.\n"+
			"Keep the tone warm, wise, and slightly humorous.",
		room,
	)
}

func openURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return errors.New("url cannot be empty")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", rawURL)
	default:
		cmd = exec.Command("xdg-open", rawURL)
	}
	return cmd.Start()
}

func envWiseManConfig() *WiseManConfig {
	provider := strings.TrimSpace(os.Getenv("WISEMAN_AI_PROVIDER"))
	endpoint := strings.TrimSpace(os.Getenv("WISEMAN_AI_URL"))
	apiKey := strings.TrimSpace(os.Getenv("WISEMAN_AI_KEY"))
	model := strings.TrimSpace(os.Getenv("WISEMAN_AI_MODEL"))
	accountID := strings.TrimSpace(os.Getenv("WISEMAN_AI_ACCOUNT_ID"))

	if provider == "" && endpoint == "" && apiKey == "" && model == "" && accountID == "" {
		return nil
	}

	cfg := &WiseManConfig{
		Enabled:   true,
		Provider:  provider,
		APIKey:    apiKey,
		Model:     model,
		AccountID: accountID,
		Endpoint:  endpoint,
	}
	if cfg.Provider == "" {
		if strings.TrimSpace(cfg.AccountID) != "" || strings.Contains(strings.ToLower(cfg.Endpoint), "cloudflare") {
			cfg.Provider = WiseManProviderCloudflare
		} else {
			cfg.Provider = WiseManProviderLegacy
		}
	}
	cfg.normalize()
	return cfg
}

func (cfg *WiseManConfig) normalize() {
	if cfg == nil {
		return
	}
	cfg.Provider = strings.ToLower(strings.TrimSpace(cfg.Provider))
	cfg.Model = strings.TrimSpace(cfg.Model)
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.AccountID = strings.TrimSpace(cfg.AccountID)
	cfg.Endpoint = strings.TrimSpace(cfg.Endpoint)

	if cfg.Provider == "" {
		cfg.Provider = WiseManProviderLegacy
	}
	switch cfg.Provider {
	case WiseManProviderGemini:
		if cfg.Model == "" {
			cfg.Model = DefaultGeminiModel
		}
	case WiseManProviderCloudflare:
		if cfg.Model == "" {
			cfg.Model = DefaultCloudflareModel
		}
	case WiseManProviderLegacy:
		if cfg.Model == "" {
			cfg.Model = "gpt-4o-mini"
		}
	}
	cfg.Enabled = cfg.APIKey != ""
	if cfg.Provider == WiseManProviderCloudflare && cfg.AccountID == "" {
		cfg.Enabled = false
	}
	if cfg.Provider == WiseManProviderLegacy && cfg.Endpoint == "" {
		cfg.Enabled = false
	}
}

func (cfg *WiseManConfig) providerLabel() string {
	if cfg == nil {
		return "unknown"
	}
	switch cfg.Provider {
	case WiseManProviderGemini:
		return "Gemini"
	case WiseManProviderCloudflare:
		return "Cloudflare Workers AI"
	case WiseManProviderLegacy:
		return "legacy endpoint"
	default:
		return cfg.Provider
	}
}

func (cfg *WiseManConfig) modelLabel() string {
	if cfg == nil || strings.TrimSpace(cfg.Model) == "" {
		return "unknown model"
	}
	return cfg.Model
}
