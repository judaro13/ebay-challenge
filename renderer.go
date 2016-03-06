package main

import (
	"fmt"
	"html/template"
	"os"
)

func render(id int) {
	category := getCategory(id)

	if category == nil {
		fmt.Printf("Category %d not found\n", id)
		return
	}
	f, err := os.Create(fmt.Sprintf(`%d.html`, id))
	defer f.Close()
	checkErr(err)
	t := template.Must(template.ParseFiles("category.html"))
	err = t.Execute(f, category)
	if err != nil {
		fmt.Println("executing template:", err)
	}
}
