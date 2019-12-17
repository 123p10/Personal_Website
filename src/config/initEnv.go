package config

import (
	"os"
)

func LoadEnv() {
	os.Setenv("ARTICLES_PATH", "/home/owen/go/src/github.com/Personal_Website/static/articles/")
	os.Setenv("TEMPLATES_PATH", "/home/owen/go/src/github.com/Personal_Website/static/templates/")
	os.Setenv("PAGES_PATH", "/home/owen/go/src/github.com/Personal_Website/static/pages/")
	os.Setenv("STATIC", "/home/owen/go/src/github.com/Personal_Website/static/")

	os.Setenv("ARTICLES_SUFFIX", ".html")
}
