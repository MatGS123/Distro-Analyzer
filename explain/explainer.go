// Package explain genera explicaciones legibles en lenguaje natural.
//
// Recibe un resultado de scoring y lo convierte en texto comprensible,
// sin realizar análisis ni tomar decisiones.
package explain

import (
	"fmt"
	"strings"

	"distroanalyzer/profile"
)

// Explainer genera explicaciones en lenguaje natural.
type Explainer interface {
	// Explain convierte un Result en texto legible.
	Explain(result *profile.Result, signals *profile.Signals) string
}

// SimpleExplainer genera explicaciones básicas sin IA.
type SimpleExplainer struct{}

// NewSimpleExplainer crea un explainer determinístico.
func NewSimpleExplainer() *SimpleExplainer {
	return &SimpleExplainer{}
}

// Explain genera una explicación estructurada del resultado.
func (e *SimpleExplainer) Explain(result *profile.Result, signals *profile.Signals) string {
	var parts []string

	// Introducción basada en categoría
	intro := e.buildIntro(result)
	parts = append(parts, intro)

	// Razones técnicas
	reasons := e.buildReasons(signals)
	if reasons != "" {
		parts = append(parts, reasons)
	}

	// Nivel de confianza
	confidence := e.buildConfidence(result)
	parts = append(parts, confidence)

	return strings.Join(parts, " ")
}

func (e *SimpleExplainer) buildIntro(result *profile.Result) string {
	switch result.Category {
		case profile.FitStrong:
			return fmt.Sprintf("Excelente coincidencia (score: %d/100).", result.Score)
		case profile.FitPotential:
			return fmt.Sprintf("Coincidencia moderada (score: %d/100).", result.Score)
		case profile.FitNone:
			return fmt.Sprintf("Coincidencia baja (score: %d/100).", result.Score)
		default:
			return fmt.Sprintf("Score: %d/100.", result.Score)
	}
}

func (e *SimpleExplainer) buildReasons(signals *profile.Signals) string {
	var reasons []string

	// Tech stack
	if len(signals.TechStack) > 0 {
		techs := strings.Join(signals.TechStack[:min(3, len(signals.TechStack))], ", ")
		reasons = append(reasons, fmt.Sprintf("Tu stack incluye: %s", techs))
	}

	// Experience level
	if signals.ExperienceLevel != "" {
		var exp string
		switch signals.ExperienceLevel {
			case profile.ExpJunior:
				exp = "junior"
			case profile.ExpMid:
				exp = "intermedio"
			case profile.ExpSenior:
				exp = "senior"
		}
		reasons = append(reasons, fmt.Sprintf("nivel %s", exp))
	}

	// Topics
	if len(signals.Topics) > 0 {
		topics := strings.Join(signals.Topics[:min(2, len(signals.Topics))], ", ")
		reasons = append(reasons, fmt.Sprintf("interés en %s", topics))
	}

	if len(reasons) == 0 {
		return ""
	}

	return strings.Join(reasons, ", ") + "."
}

func (e *SimpleExplainer) buildConfidence(result *profile.Result) string {
	confidencePercent := int(result.Confidence * 100)

	if confidencePercent >= 80 {
		return fmt.Sprintf("Alta confianza (%d%%).", confidencePercent)
	}
	if confidencePercent >= 60 {
		return fmt.Sprintf("Confianza moderada (%d%%).", confidencePercent)
	}
	return fmt.Sprintf("Confianza baja (%d%%).", confidencePercent)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
