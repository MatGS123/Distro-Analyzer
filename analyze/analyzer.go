// Package analyze extrae señales estructuradas a partir de datos crudos.
//
// Este paquete NO decide si un perfil es bueno o malo.
// Solo transforma datos crudos en señales observables y comparables.
package analyze

import "distroanalyzer/profile"

// Analyzer extrae señales estructuradas de datos crudos.
type Analyzer interface {
	// Analyze procesa RawData y devuelve Signals estructurados.
	Analyze(data *profile.RawData) (*profile.Signals, error)
}
