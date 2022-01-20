package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	FirstName string
	LastName  string
	age       int
}

//universal variable for db and err
var db *sql.DB
var err error

func setupPostgres() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost/todo?sslmode=disable")
	
	if err != 
}


func ReadItem(c *gin.Context) {
	rows, err := db.Query(`SELECT * FROM profile`)

	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error due to DB"})
	}

	//make a Profile slice that could contain up to 999 elements
	profiles := make([]Profile, 999)

	if rows != nil {
		defer rows.Close()

		//Next prepares the next result row for reading with the Scan method.
		//It returns true on success, or false if there is no next result row or an error happened while preparing it.
		//Err should be consulted to distinguish between the two cases.

		//Every call to Scan, even the first one, must be preceded by a call to Next.
		for rows.Next() {

			//profile yang di kode sebelah sama kek p
			p := Profile{}

			err := rows.Scan(&p.FirstName, &p.LastName, &p.age)

			if err != nil {
				fmt.Println("Error: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error due to DB "})
			}
			profiles = append(profiles, p)
		}

	}
	//print out using JSON
	c.JSON(http.StatusOK, gin.H{"profile": profiles})

}

