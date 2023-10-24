package layout

import _ "embed"

// supplementary
var (
	//go:embed opabinia.yml
	Config []byte
)

// markdown pages
var (
	//go:embed index.md
	Index []byte

	//go:embed pages/other_page.md
	OtherPage []byte
)

// assets

var (
	//go:embed assets/favicon/apple-touch-icon.png
	AppleTouchIcon []byte

	//go:embed assets/favicon/favicon-16x16.png
	Favicon16x16 []byte

	//go:embed assets/favicon/favicon-32x32.png
	Favicon32x32 []byte
)

//go:embed templates/default.tmpl
var Default []byte
