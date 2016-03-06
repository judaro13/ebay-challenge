package main

import (
	"fmt"
	"html/template"
	"os"
)

const (
	templateFile = `category.html`
)

type Renderer struct {
	t *template.Template
}

func NewRenderer() *Renderer {
	renderer := &Renderer{
		t: template.Must(template.ParseFiles(templateFile)),
	}

	return renderer
}

func (r *Renderer) RenderToFile(id int) {
	category := getCategory(id)

	if category == nil {
		return
	}
	f, err := os.Create(fmt.Sprintf(`%d.html`, id))
	defer f.Close()
	checkErr(err)

	err = r.t.Execute(f, category)
	if err != nil {
		fmt.Println("executing template:", err)
	}
}
