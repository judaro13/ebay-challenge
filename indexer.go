package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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

	fmt.Printf("*")
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
