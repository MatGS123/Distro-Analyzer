package score

// Distro representa una distribución de Linux con sus características.
type Distro struct {
	ID          string
	Name        string
	Rolling     int  // 0-10: cuán rolling/bleeding edge es
	Easy        int  // 0-10: facilidad de uso
	DIY         int  // 0-10: nivel de personalización
	Performance int  // 0-10: optimización de rendimiento
	DevFocus    int  // 0-10: orientación a desarrollo
	Popularity  int  // HPD (Hits Per Day) de DistroWatch
	Trend       Trend
}

// Trend representa la tendencia de popularidad.
type Trend int

const (
	TrendDown Trend = -1
	TrendStable Trend = 0
	TrendUp   Trend = 1
)
