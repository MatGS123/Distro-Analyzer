package analyze

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"distroanalyzer/profile"
	"github.com/sashabaranov/go-openai"
)

// AIAnalyzer usa Cerebras (vía OpenAI SDK) para extraer señales.
type AIAnalyzer struct {
	client *openai.Client
	model  string
}

// NewAIAnalyzer crea un analizador configurado para la API de Cerebras.
func NewAIAnalyzer(apiKey, model string) (*AIAnalyzer, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("CEREBRAS_API_KEY is required")
	}

	// Configuramos el cliente para que apunte a Cerebras en lugar de OpenAI
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.cerebras.ai/v1"

	client := openai.NewClientWithConfig(config)

	return &AIAnalyzer{
		client: client,
		model:  model,
	}, nil
}

// Analyze envía datos a Cerebras y parsea la respuesta.
func (a *AIAnalyzer) Analyze(data *profile.RawData) (*profile.Signals, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userPrompt := a.buildPrompt(data)

	// DEBUG: Ver qué enviamos
	log.Printf("DEBUG - Prompt sent to Cerebras:\n%s", userPrompt)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: a.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.1,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("cerebras API error: %w", err)
	}

	responseText := resp.Choices[0].Message.Content
	if responseText == "" {
		return nil, fmt.Errorf("empty response from cerebras")
	}

	// DEBUG: Ver qué responde Cerebras
	log.Printf("DEBUG - Cerebras raw response:\n%s", responseText)

	signals, err := a.parseSignals(responseText)
		if err != nil {
			return nil, err
	}

	// Post-procesamiento: extraer hashtags de la bio si Cerebras los ignoró
	if data.Bio != "" {
		bioTechs := extractHashtagTechs(data.Bio)
		// Agregar a TechStack si no están ya
		for _, tech := range bioTechs {
			if !contains(signals.TechStack, tech) {
				signals.TechStack = append(signals.TechStack, tech)
			}
		}
	}
	// DEBUG: Ver señales parseadas
	log.Printf("DEBUG - Parsed signals: TechStack=%v, Topics=%v, ExpLevel=%s",
		signals.TechStack, signals.Topics, signals.ExperienceLevel)

	return signals, nil
}

func extractHashtagTechs(bio string) []string {
	techs := []string{}
	words := strings.Fields(bio)

	techHashtags := map[string]string{
		"#ansible":     "ansible",
		"#k8s":         "kubernetes",
		"#kubernetes":  "kubernetes",
		"#docker":      "docker",
		"#python":      "python",
		"#go":          "go",
		"#rust":        "rust",
		"#javascript":  "javascript",
		"#typescript":  "typescript",
		"#ruby":        "ruby",
		"#php":         "php",
		"#java":        "java",
		"#devops":      "devops",
		"#terraform":   "terraform",
		"#aws":         "aws",
		"#gcp":         "gcp",
		"#azure":       "azure",
		"#linux":       "linux",
		"#nodejs":      "nodejs",
		"#react":       "react",
		"#vue":         "vue",
		"#django":      "django",
		"#rails":       "rails",
		"#drupal":      "drupal",
	}

	for _, word := range words {
		wordLower := strings.ToLower(word)
		if tech, exists := techHashtags[wordLower]; exists {
			techs = append(techs, tech)
		}
	}

	return techs
}

// contains verifica si un slice contiene un string (case-insensitive)
func contains(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}

// buildPrompt construye el prompt de usuario con los datos crudos.
func (a *AIAnalyzer) buildPrompt(data *profile.RawData) string {
	var parts []string
	parts = append(parts, "Analiza el siguiente perfil y extrae señales estructuradas:")
	parts = append(parts, "Presta ESPECIAL atención a hashtags, menciones de tecnologías en bio, y herramientas usadas.")

	if data.Bio != "" {
		parts = append(parts, "Bio: "+data.Bio)
	}
	if len(data.Repositories) > 0 {
		repos := "Repositorios: " + strings.Join(data.Repositories, ", ")
		parts = append(parts, repos)
	}
	if data.Website != "" {
		parts = append(parts, "Website: "+data.Website)
	}
	if data.Location != "" {
		parts = append(parts, "Location: "+data.Location)
	}
	if data.ReadmeText != nil && *data.ReadmeText != "" {
		readme := *data.ReadmeText
		if len(readme) > 2000 {
			readme = readme[:2000]
		}
		parts = append(parts, "README (primeras líneas): "+readme)
	}

	return strings.Join(parts, "\n\n")
}

// parseSignals limpia markdown y procesa el JSON.
func (a *AIAnalyzer) parseSignals(content string) (*profile.Signals, error) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var raw rawSignals
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Limpiar experience_level si viene con pipes
	expLevel := raw.ExperienceLevel
	if strings.Contains(expLevel, "|") {
		// Si viene "mid|senior", tomar "senior"
		parts := strings.Split(expLevel, "|")
		expLevel = parts[len(parts)-1] // tomar el último (más alto)
	}
	expLevel = strings.TrimSpace(expLevel)

	return &profile.Signals{
		Topics:          raw.Topics,
		Sentiment:       profile.Sentiment(raw.Sentiment),
		ExperienceLevel: profile.ExperienceLevel(raw.ExperienceLevel),
		Keywords:        raw.Keywords,
		TechStack:       raw.TechStack,
	}, nil
}

const systemPrompt = `Eres un asistente que extrae señales estructuradas de perfiles técnicos.

Debes devolver SOLO un JSON válido con este formato exacto:
{
"topics": ["tema1", "tema2"],
"sentiment": "positive|neutral|negative",
"experience_level": "junior|mid|senior",
"keywords": ["palabra1", "palabra2"],
"tech_stack": ["tech1", "tech2"]
}

Reglas ESTRICTAS:
- NO inventes información.
- experience_level debe ser EXACTAMENTE UNO de estos valores: "junior", "mid", "senior". NO uses pipes ni múltiples opciones.
- Si el perfil parece muy experto (kernel developer, creator de frameworks, mantainer de proyectos grandes), usa "senior".
- PRIORIDAD ABSOLUTA: Si la bio contiene hashtags con tecnologías (ejemplo: #ansible, #k8s, #docker), DEBES incluirlas en tech_stack.
- Los hashtags son indicadores directos de experiencia y deben ser incluidos siempre.
- tech_stack debe incluir: lenguajes de programación, frameworks, herramientas DevOps, plataformas mencionadas en bio Y repositorios.
- Extrae tecnologías tanto de la bio como de los nombres de repositorios.
- Responde SOLO con el JSON, sin texto adicional.`

type rawSignals struct {
	Topics          []string `json:"topics"`
	Sentiment       string   `json:"sentiment"`
	ExperienceLevel string   `json:"experience_level"`
	Keywords        []string `json:"keywords"`
	TechStack       []string `json:"tech_stack"`
}
