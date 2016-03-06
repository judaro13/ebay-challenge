package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

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

// GetCategoriesResponse contains the categories from the ebay API.
type GetCategoriesResponse struct {
	XMLName xml.Name `xml:"GetCategoriesResponse"`

	Timestamp     string      `xml:"Timestamp"`
	Build         string      `xml:"Build"`
	CategoryArray []*Category `xml:"CategoryArray>Category"`
}

// Category contains the categories
type Category struct {
	XMLName          xml.Name    `xml:"Category"`
	BestOfferEnabled bool        `xml:"BestOfferEnabled"`
	ID               int         `xml:"CategoryID"`
	Level            int         `xml:"CategoryLevel"`
	Name             string      `xml:"CategoryName"`
	ParentID         int         `xml:"CategoryParentID"`
	LeafCategory     bool        `xml:"LeafCategory"`
	LSD              bool        `xml:"LSD"`
	Indexed          bool        `xml:"-"`
	Children         []*Category `xml:"-"`
	Parent           *Category   `xml:"-"`
}

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
	dict := getCategory(id)
	f, err := os.Create(fmt.Sprintf(`%d.html`, id))
	defer f.Close()
	checkErr(err)
	t := template.Must(template.ParseFiles("category.html"))
	err = t.Execute(f, dict.CategoryArray[0])
	if err != nil {
		fmt.Println("executing template:", err)
	}
}

func (c *Category) index(categories *GetCategoriesResponse) {
	if c.Indexed {
		return
	}

	parent := findParent(c.ParentID, categories)
	if parent != nil && c.ID != c.ParentID {
		c.Parent = parent
		parent.Children = append(parent.Children, c)
		parent.index(categories)
	}

	c.Indexed = true

}

func findParent(parentID int, categories *GetCategoriesResponse) *Category {
	for _, category := range categories.CategoryArray {
		if category.ID == parentID {
			return category
		}
	}
	return nil
}

func (c *Category) debug() {
	fmt.Printf("%s%s\n", strings.Repeat("    ", c.Level), c.Name)
	for _, child := range c.Children {
		child.debug()
	}
}

func (c *Category) findAncestors() string {
	ancestors := "/"
	category := c

	for category != nil {
		if category.Parent != nil {
			ancestors = "/" + strconv.Itoa(category.ParentID) + ancestors
			category = category.Parent
		}
	}
	return ancestors
}

func rebuild() {
	rebuildDB()
	dict := &GetCategoriesResponse{}
	content, err := ioutil.ReadFile("full")

	xml.Unmarshal([]byte(content), dict)
	// xml.Unmarshal([]byte(xmlContent), dict)

	for _, category := range dict.CategoryArray {
		category.index(dict)
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

func getCategory(id int) *GetCategoriesResponse {
	db, err := sql.Open("sqlite3", "challenge.db")
	checkErr(err)

	defer db.Close()
	rows, err := db.Query(queryCategory, id, fmt.Sprintf(`%%/%d/%%`, id))
	checkErr(err)

	defer rows.Close()

	dict := &GetCategoriesResponse{
		CategoryArray: []*Category{},
	}

	for rows.Next() {
		category := &Category{}
		err = rows.Scan(&category.ID, &category.BestOfferEnabled, &category.Level, &category.ParentID, &category.Name, &category.LeafCategory, &category.LSD)
		checkErr(err)
		if category.ID >= id {
			dict.CategoryArray = append(dict.CategoryArray, category)
		}
	}

	// fmt.Printf("%v", catArray)

	for _, category := range dict.CategoryArray {
		category.index(dict)
		fmt.Printf(".")
	}
	fmt.Println("")
	return dict
}
