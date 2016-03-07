# Introduction

The application was created in GO language. It includes 5 go files with the solution for the problem exposed.

The main file is `challenge.go` where it takes the params and executes according to them.

In the `downloader.go` file you can find the code to query the ebay API. The data is fetched in XML format and parsed directly
into the Category struct.

The `category.go` file have the definition for the category. All fields from the XML are mapped to the `Category` fields, except for
the `Children` and `Parent` fields, which are populated in code. These fields are used to implement a tree structure.

The `indexer.go` have the function to handle the DB structure and index the categories. It creates
a single table called Categories and it has the next structure:

```sql
CREATE TABLE categories (
        id integer not null primary key, 
        bestOfferEnabled text,
        level integer,
        parentId integer,
        name text,
        leafCategory boolean,
        lsd boolean,
        eldersPath text, 
        FOREIGN KEY(parentId) REFERENCES category(id)
    );
```
    
The `eldersPath` field stores a string with the path of the ancestors. This is generated with the tree schema
of categories and it is used to retrieve all the subcategories related with the given category.

The `renderer.go` file handles the renderization of a given category with all its subcategories, using the `category.html` template file. 
All the data is queried from the database and parsed to the local tree structure for indexing and sorting it.


## Installation

- Be sure you have installed go language in your computer
- Install the go packages.  

`github.com/mattn/go-sqlite3` it can be installed as follows

```shell
    $ go get github.com/mattn/go-sqlite3
```   

`github.com/spf13/pflag` it can be installed as follows
```shell
    $ go get github.com/spf13/pflag
 ```  

## Usage

To compile type the following command:

```shell
$ go build
```

For rebuilding the db:
```shell
$ ./ebay-challenge --rebuild 
```

For rendering a category
```shell
$ ./ebay-challenge --render 12345
```
