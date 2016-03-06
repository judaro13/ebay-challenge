package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // go get github.com/mattn/go-sqlite3
	flag "github.com/ogier/pflag"   //go get github.com/ogier/pflag
)

var (
	xmlContent = ``
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rebuildOpt := flag.Bool("rebuild", false, "Recreate the category DB tree")
	renderOpt := flag.Int("render", 0, "Render a category with the given ID. User as --render=12345 ")
	flag.Parse()

	if *rebuildOpt {
		rebuild()
	} else if *renderOpt > 0 {
		render(*renderOpt)
	} else {
		flag.PrintDefaults()
	}
}
