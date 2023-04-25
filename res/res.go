package res

import (
	_ "embed"
)

//go:embed template.html
var HtmlTemplate string

//go:embed edit.html
var Edit string

//go:embed main.css
var CSS []byte
