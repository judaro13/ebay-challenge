package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // go get github.com/mattn/go-sqlite3
	flag "github.com/ogier/pflag"   //go get github.com/ogier/pflag
)

var (
	xmlContent = ``
)

// SQL constants
const (
	createCategoriesTable = `create table categories (
        id integer not null primary key, 
        bestOfferEnabled boolean,
        level integer,
        parentId integer,
        name text,
        leafCategory boolean,
        lsd boolean,
        eldersPath text, 
        FOREIGN KEY(parentId) REFERENCES category(id)
    );`
	insertCategorySQL = `insert into categories(id, bestOfferEnabled, level, parentId, name, leafCategory, lsd, eldersPath) values(?, ?, ?, ?, ?, ?, ?, ?)`
	queryCategory     = `SELECT 
     id, bestOfferEnabled, level , parentId, name, leafCategory, lsd 
     FROM categories WHERE id=? OR eldersPath LIKE ?
     ORDER BY id ASC`
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

func rebuild() {
	rebuildDB()
	dict := &GetCategoriesResponse{}
	content, err := ioutil.ReadFile("full")

	xml.Unmarshal([]byte(content), dict)
	// xml.Unmarshal([]byte(xmlContent), dict)

	for _, category := range dict.CategoryArray {
		category.index(dict.CategoryArray)
		fmt.Printf(".")
	}

	db, err := sql.Open("sqlite3", "challenge.db")
	defer db.Close()
	checkErr(err)

	for _, category := range dict.CategoryArray {
		insertCategory(db, category)
	}

}

func rebuildDB() {
	os.Remove("challenge.db")
	db, err := sql.Open("sqlite3", "challenge.db")
	defer db.Close()
	checkErr(err)

	_, err = db.Exec(createCategoriesTable)
	checkErr(err)
}

func insertCategory(db *sql.DB, category *Category) {
	tx, err := db.Begin()
	checkErr(err)

	stmt, err := tx.Prepare(insertCategorySQL)
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(category.ID,
		category.BestOfferEnabled,
		category.Level,
		category.ParentID,
		category.Name,
		category.LeafCategory,
		category.LSD,
		category.findAncestors())
	checkErr(err)
	tx.Commit()

	fmt.Printf("*")
}

func getCategory(id int) *Category {
	categories := findAndIndex(id)
	if len(categories) > 0 {
		return categories[0]
	}
	return nil
}

func findAndIndex(id int) []*Category {
	db, err := sql.Open("sqlite3", "challenge.db")
	checkErr(err)

	defer db.Close()
	rows, err := db.Query(queryCategory, id, fmt.Sprintf(`%%/%d/%%`, id))
	checkErr(err)

	defer rows.Close()

	categoryArray := []*Category{}

	for rows.Next() {
		category := &Category{}
		err = rows.Scan(&category.ID, &category.BestOfferEnabled, &category.Level, &category.ParentID, &category.Name, &category.LeafCategory, &category.LSD)
		checkErr(err)
		if category.ID >= id {
			categoryArray = append(categoryArray, category)
		}
	}

	// fmt.Printf("%v", catArray)

	for _, category := range categoryArray {
		category.index(categoryArray)
		fmt.Printf(".")
	}
	fmt.Println("")
	return categoryArray
}
