package embedded

import (
	_ "embed"
)

//go:embed docs/rbedit_synopsis.md
var DocsRbeditSynopsisMarkdown string

//go:embed docs/rbedit_example.md
var DocsRbeditExampleMarkdown string
