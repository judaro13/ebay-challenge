package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
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
	BestOfferEnabled string      `xml:"BestOfferEnabled"`
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

func findParent(parentID int, categories []*Category) *Category {
	for _, category := range categories {
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
		}
		category = category.Parent
	}
	return ancestors
}

func (c *Category) index(categories []*Category) {
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
