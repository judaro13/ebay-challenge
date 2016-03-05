package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // go get github.com/mattn/go-sqlite3
	flag "github.com/ogier/pflag"   //go get github.com/ogier/pflag
)

var (
	xmlContent = `<?xml version="1.0" encoding="UTF-8"?>
<GetCategoriesResponse xmlns="urn:ebay:apis:eBLBaseComponents">
  <Timestamp>2016-03-05T19:19:08.199Z</Timestamp>
  <Ack>Success</Ack>
  <Version>949</Version>
  <Build>E949_CORE_API_17785418_R1</Build>
  <CategoryArray>
    <Category>
      <BestOfferEnabled>true</BestOfferEnabled>
      <CategoryID>10542</CategoryID>
      <CategoryLevel>1</CategoryLevel>
      <CategoryName>Real Estate</CategoryName>
      <CategoryParentID>10542</CategoryParentID>
      <LSD>true</LSD>
    </Category>
    <Category>
      <BestOfferEnabled>true</BestOfferEnabled>
      <CategoryID>15825</CategoryID>
      <CategoryLevel>2</CategoryLevel>
      <CategoryName>Commercial</CategoryName>
      <CategoryParentID>10542</CategoryParentID>
      <LeafCategory>true</LeafCategory>
      <LSD>true</LSD>
    </Category>
    <Category>
      <BestOfferEnabled>true</BestOfferEnabled>
      <CategoryID>1607</CategoryID>
      <CategoryLevel>2</CategoryLevel>
      <CategoryName>Other Real Estate</CategoryName>
      <CategoryParentID>10542</CategoryParentID>
      <LeafCategory>true</LeafCategory>
      <LSD>true</LSD>
    </Category>
  </CategoryArray>
  <CategoryCount>7</CategoryCount>
  <UpdateTime>2015-09-01T22:57:09.000Z</UpdateTime>
  <CategoryVersion>113</CategoryVersion>
  <ReservePriceAllowed>true</ReservePriceAllowed>
  <MinimumReservePrice>0.0</MinimumReservePrice>
</GetCategoriesResponse>
`
)

const (
	insertCategorySQL = `insert into category(id, bestOfferEnabled, level, parentId, name, LeafCategory, lsd ) values(?, ?, ?, ?, ?, ?, ?)`
)

// GetCategoriesResponse contains the categories from the ebay API.
type GetCategoriesResponse struct {
	XMLName       xml.Name    `xml:"GetCategoriesResponse"`
	Timestamp     string      `xml:"Timestamp"`
	CategoryArray []*Category `xml:"CategoryArray>Category"`
}

// Category contains the categories
type Category struct {
	XMLName          xml.Name `xml:"Category"`
	BestOfferEnabled bool     `xml:"BestOfferEnabled"`
	ID               int      `xml:"CategoryID"`
	Level            int      `xml:"CategoryLevel"`
	Name             string   `xml:"CategoryName"`
	ParentID         int      `xml:"CategoryParentID"`
	LeafCategory     bool     `xml:"LeafCategory"`
	LSD              bool     `xml:"LSD"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rebuildOpt := flag.Bool("rebuild", false, "Recreate the category DB tree")
	renderOpt := flag.Int("render", 0, "Render a category with the given ID")
	flag.Parse()

	if *rebuildOpt {
		rebuild()
	} else if *renderOpt > 0 {

	} else {
		flag.PrintDefaults()
	}
}

func rebuild() {
	rebuildDB()
	dict := &GetCategoriesResponse{}
	xml.Unmarshal([]byte(xmlContent), dict)

	db, err := sql.Open("sqlite3", "challenge.db")
	defer db.Close()
	checkErr(err)

	for _, category := range dict.CategoryArray {
		insertCategory(db, category)
		fmt.Printf("%#v\n", category)
	}

}

func rebuildDB() {
	os.Remove("challenge.db")
	db, err := sql.Open("sqlite3", "challenge.db")
	defer db.Close()
	checkErr(err)

	sqlStmt := `
	create table category (
        id integer not null primary key, 
        bestOfferEnabled boolean,
        level integer,
        parentId integer,
        name text,
        LeafCategory boolean,
        lsd boolean, 
        FOREIGN KEY(parentId) REFERENCES category(idartistid)
    );
	`
	_, err = db.Exec(sqlStmt)
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
		category.LSD)
	checkErr(err)
	tx.Commit()

}
