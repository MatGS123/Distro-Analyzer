// Package score aplica reglas determinísticas para calcular el puntaje final.
//
// Este paquete recibe señales estructuradas y devuelve un resultado numérico.
// Las reglas son explícitas, testeables y no dependen de decisiones externas.
package score

import (
	"log"
	"math"
	"strings"

	"distroanalyzer/profile"
)

// Engine calcula el puntaje final aplicando reglas determinísticas.
type Engine struct {
	distros []Distro
}

type ScoreOutput struct {
	Result         *profile.Result
	BestDistroID   string
	BestDistroName string
}



// NewEngine crea un motor de scoring con la base de distros.
func NewEngine(distros []Distro) *Engine {
	return &Engine{
		distros: distros,
	}
}

// Score calcula el resultado final para un perfil basado en sus señales.
func (e *Engine) Score(signals *profile.Signals) *ScoreOutput {
	// Calcular dimensiones del usuario
	dimensions := e.calculateDimensions(signals)

	// DEBUG: Ver dimensiones calculadas
	log.Printf("DEBUG - User dimensions: Rolling=%d, DIY=%d, Perf=%d, Dev=%d",
		   dimensions.RollingScore, dimensions.DIYScore,
	    dimensions.PerformanceScore, dimensions.DevScore)


	// Encontrar mejor match
	bestMatch := e.findBestMatch(dimensions, signals)

	log.Printf("DEBUG - Best match: %s (score: %.2f)",
		   bestMatch.distro.Name, bestMatch.matchScore)

	// Calcular score final (0-100)
	finalScore := e.calculateFinalScore(bestMatch, dimensions, signals)

	// Determinar categoría
	category := e.determineCategory(finalScore)

	// Generar explicación
	explanation := e.buildExplanation(bestMatch, dimensions, finalScore)

	return &ScoreOutput{
		Result: &profile.Result{
			Score:       finalScore,
			Category:    category,
			Explanation: explanation,
			Confidence:  bestMatch.matchScore,
		},
		BestDistroID:   bestMatch.distro.ID,
		BestDistroName: bestMatch.distro.Name,
}
}

// calculateDimensions extrae dimensiones del perfil a partir de señales.
func (e *Engine) calculateDimensions(signals *profile.Signals) UserDimensions {
	dims := UserDimensions{}

	// 1. Ciclo de vida (Rolling vs LTS)
	dims.RollingScore = calculateRollingPreference(signals)

	// 2. Personalización (DIY vs Easy)
	dims.DIYScore = calculateDIYPreference(signals)

	// 3. Performance/Gaming
	dims.PerformanceScore = calculatePerformanceNeed(signals)

	// 4. Developer focus
	dims.DevScore = calculateDevFocus(signals)

	return dims
}

// calculateRollingPreference detecta preferencia por rolling/bleeding edge.
func calculateRollingPreference(signals *profile.Signals) int {
	score := 5 // neutral

	// Tecnologías bleeding edge
	bleedingTech := []string{"rust", "mojo", "zig", "deno", "bun"}
	for _, tech := range signals.TechStack {
		for _, bleeding := range bleedingTech {
			if tech == bleeding {
				score += 2
			}
		}
	}

	// Experience senior = más tolerancia al cambio
	if signals.ExperienceLevel == profile.ExpSenior {
		score += 1
	}

	// Keywords de estabilidad
	stableKeywords := []string{"production", "enterprise", "stable", "lts"}
	for _, kw := range signals.Keywords {
		for _, stable := range stableKeywords {
			if kw == stable {
				score -= 2
			}
		}
	}

	return clamp(score, 0, 10)
}

// calculateDIYPreference detecta si prefiere customización o simplicidad.
func calculateDIYPreference(signals *profile.Signals) int {
	score := 5 // neutral

	// Keywords clave
	diyKeywords := []string{
		"dotfiles", "rice", "customization", "tiling", "window manager",
		"kernel", "arch", "gentoo", "nixos", "low-level", "assembly",
		"hyprland", "sway", "i3", "awesome", "dwm", "qtile", "bspwm", // tiling WMs top
		"wayland", "x11", "compositor", "eww", "polybar", "waybar", "rofi", "wofi", // bars y launchers
		"minimal", "minimalism", "void", "artix", "crux", "alpine", "kiss", // ultra-minimal
		"lfs", "linux from scratch", "custom kernel", "musl", "glibc hardening",
		"immutable", "atomic", "silverblue", "kinoite", "bazzite", "ublue", // atomic desktops
		"home-manager", "flakes", "nix", "guix", // declarative config
		"ricing", "unixporn", "gruvbox", "catppuccin", "tokyonight", // themes populares
	}

	easyKeywords := []string{"beginner", "simple", "easy", "user-friendly"}

	for _, kw := range signals.Keywords {
		for _, diy := range diyKeywords {
			if kw == diy {
				score += 3
			}
		}
		for _, easy := range easyKeywords {
			if kw == easy {
				score -= 2
			}
		}
	}

	// Tech stack de scripting
	scriptingLangs := []string{"bash", "lua", "python"}
	scriptCount := 0
	for _, tech := range signals.TechStack {
		for _, script := range scriptingLangs {
			if tech == script {
				scriptCount++
			}
		}
	}
	if scriptCount >= 2 {
		score += 2
	}

	return clamp(score, 0, 10)
}

// calculatePerformanceNeed detecta necesidad de alto rendimiento/gaming.
func calculatePerformanceNeed(signals *profile.Signals) int {
	score := 3 // bajo por defecto

	perfKeywords := []string{"gaming", "performance", "gpu", "vulkan", "shader", "godot", "unreal"}

	perfTech := []string{
		"c", "c++", "rust", "vulkan", "opengl", "gpu",
		"cuda", "rocm", "opencl", "metal", // compute/GPGPU
		"directx", "dx12", "webgpu", // si menciona ports o gamedev
		"assembly", "asm", "x86", "arm", "riscv", // low-level
		"hpc", "mpi", "openmp", "simd", "avx", "avx512",
		"zig", "c++20", "c++23", "cpp", // modern perf langs
		"ispc", "halide", // domain-specific perf langs
		"game dev", "godot", "unreal", "unity", // engines que piden perf
	}

	for _, kw := range signals.Keywords {
		for _, perf := range perfKeywords {
			if kw == perf {
				score += 2
			}
		}
	}

	for _, tech := range signals.TechStack {
		for _, perf := range perfTech {
			if tech == perf {
				score += 1
			}
		}
	}

	return clamp(score, 0, 10)
}

// calculateDevFocus detecta orientación a desarrollo/DevOps.
func calculateDevFocus(signals *profile.Signals) int {
	score := 5

	devTech := []string{
		"c", "c++", "go", "rust", "python", "ruby", "javascript", "typescript",
		"java", "kotlin", "swift", "php", "perl", "shell", "bash", "lua",
		"docker", "kubernetes", "terraform", "ansible", "vagrant", "chef", "puppet",
		"jenkins", "gitlab", "github actions", "circleci",
		"aws", "gcp", "azure", "cloud",
		"git", "make", "cmake", "gradle", "maven", "npm", "yarn", "pip",
	}

	devKeywords := []string{
		"devops", "backend", "infrastructure", "sre", "platform",
		"rails", "web", "api", "microservices", "containers", "orchestration",
		"automation", "ci/cd", "deployment", "ansible", "kubernetes", "k8s",
	}

	// Bonus especial por keywords de alto nivel
	criticalKeywords := []string{"kernel", "ansible", "kubernetes", "k8s", "docker", "devops"}

	for _, kw := range signals.Keywords {
		kwLower := strings.ToLower(kw)
		for _, critical := range criticalKeywords {
			if strings.Contains(kwLower, critical) {
				score += 2  // Bonus grande
			}
		}
	}

	for _, tech := range signals.TechStack {
		techLower := strings.ToLower(tech)
		for _, dev := range devTech {
			if techLower == dev || strings.Contains(techLower, dev) {
				score += 1
			}
		}
	}

	for _, kw := range signals.Keywords {
		kwLower := strings.ToLower(kw)
		for _, dev := range devKeywords {
			if kwLower == dev || strings.Contains(kwLower, dev) {
				score += 1
			}
		}
	}

	return clamp(score, 0, 10)
}

// findBestMatch encuentra la distro con mejor fit.
func (e *Engine) findBestMatch(dims UserDimensions, signals *profile.Signals) MatchResult {
	var best MatchResult
	bestScore := 0.0

	// máximo teórico de distancia euclidiana en este espacio:
	// cada dimensión 0..10, 4 dimensiones => maxDist = sqrt(4 * 10^2) = 20
	const maxDist = 20.0

	// parámetros (tuneables)
	const alpha = 0.90 // peso para la similitud geométrica
	const beta = 0.10  // peso para la popularidad

	// precomputar maxPopularity
	maxPopularity := 3790.0

	for _, distro := range e.distros {
		// Skip distros muy oscuras para usuarios con perfil claro
		if signals.ExperienceLevel == profile.ExpSenior && distro.Popularity < 500 {
			continue
		}

		// distancia euclidiana simple entre dimensiones
		distance := math.Sqrt(
			math.Pow(float64(dims.RollingScore-distro.Rolling), 2) +
			math.Pow(float64(dims.DIYScore-distro.DIY), 2) +
			math.Pow(float64(dims.PerformanceScore-distro.Performance), 2) +
			math.Pow(float64(dims.DevScore-distro.DevFocus), 2),
		)

		// Normalizar distancia y convertir a similitud [0..1]
		normDist := distance / maxDist
		if normDist < 0 {
			normDist = 0
		}
		if normDist > 1 {
			normDist = 1
		}
		similarity := 1.0 - normDist

		// Normalizar popularidad usando log
		popNorm := 0.0
		if distro.Popularity > 0 {
			popNorm = math.Log(float64(distro.Popularity)+1.0) / math.Log(maxPopularity+1.0)
			if popNorm < 0 {
				popNorm = 0
			}
			if popNorm > 1 {
				popNorm = 1
			}
		}

		// pequeña corrección por tendencia
		trendMultiplier := 1.0
		switch distro.Trend {
			case TrendUp:
				trendMultiplier = 1.08 // +8%
			case TrendDown:
				trendMultiplier = 0.97 // -3%
		}

		// combinar: suma ponderada
		finalScore := (alpha*similarity + beta*popNorm) * trendMultiplier

		if finalScore > bestScore {
			bestScore = finalScore
			best = MatchResult{
				distro:     distro,
				matchScore: finalScore,
			}
		}
	}

	// Después del loop principal, antes de return best:

	// Penalizar distros genéricas para usuarios senior avanzados
	if signals.ExperienceLevel == profile.ExpSenior {
		// Si DevScore es alto pero la distro tiene DevFocus bajo, penalizar
		if dims.DevScore >= 8 && best.distro.DevFocus <= 7 {
			best.matchScore *= 0.85  // -15% penalty
			log.Printf("DEBUG - Penalizing generic distro %s for senior dev profile", best.distro.Name)
		}

		// Si DIY alto pero distro es muy "easy", penalizar
		if dims.DIYScore >= 8 && best.distro.Easy >= 9 {
			best.matchScore *= 0.90  // -10% penalty
			log.Printf("DEBUG - Penalizing too-easy distro %s for DIY user", best.distro.Name)
		}
	}

	return best
}


// calculateFinalScore convierte el match score a escala 0-100.
func (e *Engine) calculateFinalScore(match MatchResult, dims UserDimensions, signals *profile.Signals) int {
	// baseScore en 0..100
	baseScore := int(math.Round(match.matchScore * 100.0))

	adjustment := 0

	// 1. Ajuste por nivel de experiencia vs facilidad de uso
	distro := match.distro

	if signals.ExperienceLevel == profile.ExpJunior {
		// Usuarios junior: bonus por distros fáciles
		if distro.Easy >= 8 {
			adjustment += 5
		}
		// Penalización por distros DIY extremas
		if distro.DIY >= 9 {
			adjustment -= 10
		}
	} else if signals.ExperienceLevel == profile.ExpSenior {
		// Usuarios senior: bonus por distros con alto DevFocus
		if distro.DevFocus >= 9 {
			adjustment += 5
		}
		// Bonus menor por DIY (aprecian el control)
		if distro.DIY >= 7 {
			adjustment += 3
		}
		// Ligera penalización por distros demasiado simples
		if distro.Easy >= 10 && distro.DIY <= 2 {
			adjustment -= 3
		}
	}

	// 2. Ajuste por coherencia de dimensiones
	// Si el usuario tiene alto DevFocus pero la distro tiene bajo, penalizar
	if dims.DevScore >= 8 && distro.DevFocus <= 5 {
		adjustment -= 5
	}

	// Si el usuario necesita performance pero la distro es débil, penalizar
	if dims.PerformanceScore >= 8 && distro.Performance <= 5 {
		adjustment -= 5
	}

	// 3. Bonus por match perfecto en múltiples dimensiones
	perfectMatches := 0
	if abs(dims.RollingScore - distro.Rolling) <= 1 {
		perfectMatches++
	}
	if abs(dims.DIYScore - distro.DIY) <= 1 {
		perfectMatches++
	}
	if abs(dims.PerformanceScore - distro.Performance) <= 1 {
		perfectMatches++
	}
	if abs(dims.DevScore - distro.DevFocus) <= 1 {
		perfectMatches++
	}

	// Bonus progresivo por matches múltiples
	if perfectMatches >= 3 {
		adjustment += 8
	} else if perfectMatches == 2 {
		adjustment += 4
	}

	// 4. Ajuste por popularidad extrema (evitar distros muy oscuras o moribundas)
	if distro.Popularity < 150 {
		adjustment -= 3 // Distros muy nicho
	}

	if distro.Trend == TrendDown && distro.Popularity < 300 {
		adjustment -= 5 // Distro en declive y poco popular = riesgoso
	}

	finalScore := baseScore + adjustment

	return clamp(finalScore, 0, 100)
}

// Helper: valor absoluto
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
// determineCategory asigna categoría cualitativa.
func (e *Engine) determineCategory(score int) profile.FitCategory {
	if score >= 75 {
		return profile.FitStrong
	}
	if score >= 50 {
		return profile.FitPotential
	}
	return profile.FitNone
}

// buildExplanation genera texto legible del resultado.
func (e *Engine) buildExplanation(match MatchResult, dims UserDimensions, score int) string {
	distro := match.distro

	explanation := "Basado en tu perfil, recomendamos " + distro.Name + ". "

	// Razones principales
	if dims.RollingScore >= 7 && distro.Rolling >= 7 {
		explanation += "Prefieres tecnologías actuales y " + distro.Name + " es rolling release. "
	}
	if dims.DIYScore >= 7 && distro.DIY >= 7 {
		explanation += "Tu perfil indica que te gusta personalizar, " + distro.Name + " ofrece control total. "
	}
	if dims.PerformanceScore >= 7 && distro.Performance >= 8 {
		explanation += "Necesitas alto rendimiento y " + distro.Name + " está optimizada para ello. "
	}

	// Popularidad y tendencia
	if distro.Trend == TrendUp {
		explanation += "Además, está en crecimiento activo. "
	}

	return explanation
}

// clamp limita un valor entre min y max.
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// UserDimensions representa las dimensiones calculadas del usuario.
type UserDimensions struct {
	RollingScore     int // 0-10: LTS(0) → Rolling(10)
	DIYScore         int // 0-10: Easy(0) → DIY(10)
	PerformanceScore int // 0-10: necesidad de rendimiento
	DevScore         int // 0-10: orientación desarrollo
}

// MatchResult representa el resultado del matching.
type MatchResult struct {
	distro     Distro
	matchScore float64 // 0.0-1.0
}
