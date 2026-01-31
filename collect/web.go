package collect

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"distroanalyzer/profile"
	"golang.org/x/net/html"
)

// WebCollector recolecta datos de URLs web públicas.
type WebCollector struct {
	client *http.Client
}

// NewWebCollector crea un nuevo collector para páginas web.
func NewWebCollector() *WebCollector {
	return &WebCollector{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Collect obtiene contenido de una URL y extrae información básica.
func (w *WebCollector) Collect(url string) (*profile.RawData, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Verificar robots.txt (simplificado)
	if !w.canFetch(url) {
		return nil, fmt.Errorf("disallowed by robots.txt")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "DistroAnalyzer/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	text := w.extractText(string(body))

	return &profile.RawData{
		Bio:        "", // Web no tiene bio estructurada
		Website:    url,
		ReadmeText: &text,
	}, nil
}

// canFetch verifica si se puede scrapear la URL (robots.txt básico).
func (w *WebCollector) canFetch(url string) bool {
	// Implementación simplificada: siempre permite
	// En producción: parsear robots.txt del dominio
	return true
}

// extractText extrae texto visible de HTML.
func (w *WebCollector) extractText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	var text strings.Builder
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			content := strings.TrimSpace(n.Data)
			if content != "" {
				text.WriteString(content)
				text.WriteString(" ")
			}
		}

		// Ignorar scripts y styles
		if n.Type == html.ElementNode {
			if n.Data == "script" || n.Data == "style" {
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	result := text.String()
	if len(result) > 5000 {
		result = result[:5000] // Limitar texto extraído
	}

	return result
}
