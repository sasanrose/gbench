// Package render defines the Renderer interface to be implemented by different
// renderer drivers.
package render

import "github.com/sasanrose/gbench/report"

// Renderer defines renderer interface.
type Renderer interface {
	Render(result *report.Result) error
}
