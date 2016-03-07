package main

import (
	"fmt"
	"html/template"
	"os"
)

const (
	// Path with the html template
	templateFile = `category.html`
)

//Renderer defined struc
type Renderer struct {
	t       *template.Template
	indexer *Indexer
}

//NewRenderer  create a new object Render
func NewRenderer(indexer *Indexer) *Renderer {
	renderer := &Renderer{
		t:       template.Must(template.ParseFiles(templateFile)),
		indexer: indexer,
	}

	return renderer
}

//RenderToFile create a new html with the given category
func (render *Renderer) RenderToFile(id int) {
	category := render.indexer.getCategory(id)

	if category == nil {
		return
	}
	f, err := os.Create(fmt.Sprintf(`%d.html`, id))
	defer f.Close()
	checkErr(err)

	err = render.t.Execute(f, category)
	if err != nil {
		fmt.Println("executing template:", err)
	}
}
