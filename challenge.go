package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // go get github.com/mattn/go-sqlite3
	flag "github.com/spf13/pflag"   //go get github.com/spf13/pflag
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
	renderOpt := flag.Int("render", 0, "Render a category with the given integer ID ")
	flag.Parse()

	i := NewIndexer("challenge.db")

	if *rebuildOpt {
		t := time.Now()
		i.Build()
		println("Time taken to build index:", time.Now().Sub(t).String())
	} else if *renderOpt > 0 {
		r := NewRenderer(i)
		r.RenderToFile(*renderOpt)
	} else {
		flag.PrintDefaults()
	}
}
