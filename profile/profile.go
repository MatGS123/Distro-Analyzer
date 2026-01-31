// Paquete profile define el modelo de dominio central del sistema.

package profile

import "time"

// ExperienceLevel representa una estimación gruesa del nivel de experiencia.

type ExperienceLevel string

// FitCategory representa la clasificación final del perfil.

type FitCategory string

// Sentiment representa la polaridad emocional general inferida
// a partir del contenido analizado del perfil.
type Sentiment string

// Valores posibles para el nivel de experiencia.
const (
	ExpJunior ExperienceLevel = "junior"
	ExpMid    ExperienceLevel = "mid"
	ExpSenior ExperienceLevel = "senior"
)

// Valores posibles para la categoría final de encaje del perfil.
const (
	FitStrong    FitCategory = "strong_fit"
	FitPotential FitCategory = "potential"
	FitNone      FitCategory = "not_fit"
)

// Valores posibles para el sentimiento general detectado.
const (
	SentimentPos Sentiment = "positive"
	SentimentNeu Sentiment = "neutral"
	SentimentNeg Sentiment = "negative"
)

// Profile representa el perfil analizado completo de una persona o cuenta.

type Profile struct {
	// Identificador o nombre del perfil (usuario, handle, etc.).
	Username string

	// Se utiliza para trazabilidad, no para lógica de negocio.
	Source string

	// Datos crudos recolectados antes o durante el análisis.
	RawData RawData

	// Señales estructuradas extraídas a partir de los datos crudos.
	Signals Signals

	// Resultado final obtenido a partir de reglas determinísticas.
	Result Result

	// Recomendacion
	Recommendation Recommendation

	// Fecha y hora de creación del perfil.
	CreatedAt time.Time
}

// RawData contiene los datos de entrada recolectados del perfil.

type RawData struct {
	Bio          string
	Repositories []string
	Website      string
	Location     string
	Email        string

	// ReadmeText puede ser un campo pesado.
	// Se define como puntero para poder liberarlo (nil)
	// una vez que ya no sea necesario.
	ReadmeText *string
}

// Signals representa señales estructuradas y normalizadas
// extraídas a partir de RawData.

type Signals struct {
	Topics          []string
	Sentiment       Sentiment
	ExperienceLevel ExperienceLevel
	Keywords        []string
	TechStack       []string
}

//Recomendacion de la distro principal

type Recommendation struct {
	DistroID   string
	DistroName string
}

// Result representa el resultado final del análisis del perfil.

type Result struct {
	// Puntaje numérico interno utilizado para comparaciones y umbrales.
	Score int

	// Categoría cualitativa derivada del puntaje.
	Category FitCategory

	// Explicación legible para humanos del resultado obtenido.
	Explanation string

	// Nivel de confianza del análisis, esperado en el rango [0.0, 1.0].
	Confidence float64
}

// ClearLargeData elimina datos crudos voluminosos que ya no son necesarios.

func (p *Profile) ClearLargeData() {
	p.RawData.ReadmeText = nil
}
