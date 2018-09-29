package render

import "github.com/sasanrose/gbench/report"

type Renderer interface {
	Render(result *report.Result) error
}
