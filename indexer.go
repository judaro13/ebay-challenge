package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// SQL constants
const (
	createCategoriesTable = `create table categories (
        id integer not null primary key, 
        bestOfferEnabled text,
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

//Indexer : Structure definer
type Indexer struct {
	dbName string
}

//NewIndexer : Generate a new object
func NewIndexer(dbName string) *Indexer {
	indexer := &Indexer{
		dbName: dbName,
	}
	return indexer
}

//Build : Index all categories
func (indexer *Indexer) Build() {
	indexer.createdDB()
	downloader := NewDownloader()
	categories := downloader.GetCategories()

	for _, category := range categories {
		category.index(categories)
		fmt.Printf(".")
	}

	db, err := sql.Open("sqlite3", indexer.dbName)
	checkErr(err)
	defer db.Close()

	for _, category := range categories {
		indexer.insertCategory(db, category)
	}
}

// createdDB delete and generate a new DB
func (indexer *Indexer) createdDB() {
	os.Remove(indexer.dbName)
	db, err := sql.Open("sqlite3", indexer.dbName)
	defer db.Close()
	checkErr(err)

	_, err = db.Exec(createCategoriesTable)
	checkErr(err)
}

// insertCategory: Insert a given category to the DB
func (indexer *Indexer) insertCategory(db *sql.DB, category *Category) {

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

// getCategory: Retrieve a category and index all related categories
func (indexer *Indexer) getCategory(id int) *Category {
	db, err := sql.Open("sqlite3", indexer.dbName)
	checkErr(err)

	defer db.Close()
	rows, err := db.Query(queryCategory, id, fmt.Sprintf(`%%/%d/%%`, id))
	if err != nil {
		log.Fatal(err)
		println("Try running first --rebuild")
	}

	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		category := &Category{}
		err = rows.Scan(&category.ID, &category.BestOfferEnabled, &category.Level, &category.ParentID, &category.Name, &category.LeafCategory, &category.LSD)
		checkErr(err)
		if category.ID >= id {
			categories = append(categories, category)
		}
	}

	for _, category := range categories {
		category.index(categories)
		fmt.Printf(".")
	}

	fmt.Println("")
	if len(categories) > 0 {
		return categories[0]
	}
	fmt.Printf("category %d was not found\n", id)
	return nil
}
