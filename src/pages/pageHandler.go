package pages

import (
	"io/ioutil"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) SavePage() error {
	fileName := os.Getenv("ARTICLES_PATH") + p.Title + os.Getenv("ARTICLES_SUFFIX")
	return ioutil.WriteFile(fileName, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	fileName := os.Getenv("ARTICLES_PATH") + title + os.Getenv("ARTICLES_SUFFIX")
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
