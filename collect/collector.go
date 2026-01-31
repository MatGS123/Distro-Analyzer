// Package collect define la interfaz y tipos para recolectar datos de perfiles.
package collect

import "distroanalyzer/profile"

// Collector representa cualquier fuente capaz de recolectar datos crudos.
type Collector interface {
	// Collect obtiene datos a partir de un identificador (username, URL, etc).
	Collect(input string) (*profile.RawData, error)
}
