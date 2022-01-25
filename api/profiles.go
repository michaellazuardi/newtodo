package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

//every variable in the struct MUST have a capital letter to be considered as a public variable
type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

//universal variable for db and err
var db *sql.DB

//POSTGRES COMPLETED --> ABLE TO RUN
// THIS FUNCTION CONNECTS MY API TO THE BACKEND
func SetupPostgres() {
	//Open function opens a database specified by its database driver name and a driver-specific data source name,
	//usually consisting of at least a database name and connection information.
	conn, err := sql.Open("postgres", "postgres://postgres:admin@localhost/test?sslmode=disable")

	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	//Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	err = conn.Ping()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	//is it different if i use log.Println()
	log.Println("connected to postgres")
	db = conn
}

//Read function
func ReadItem(c *gin.Context) {
	rows, err := db.Query(`SELECT * FROM profile`)

	//test error query nya
	if err != nil {
		fmt.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error due to DB"})
	}

	//you need to create the make it to 0 so that it wont be a null pointer
	//EVERYTIME you want to make an empty array, you need to use 'make' function
	profiles := make([]Profile, 0) // -> make an empty array

	if rows != nil {
		defer rows.Close()

		//Next prepares the next result row for reading with the Scan method.
		//It returns true on success, or false if there is no next result row or an error happened while preparing it.
		//Err should be consulted to distinguish between the two cases.

		//Every call to Scan, even the first one, must be preceded by a call to Next.
		for rows.Next() {

			p := Profile{}

			//scan function will copy the firstname, lastname, age to the 'p'
			err = rows.Scan(&p.FirstName, &p.LastName, &p.Age)

			if err != nil {
				fmt.Println("Error: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error due to DB "})
				return
			}
			//since now 'p' has got some elements, we append it to the empty array of 'profiles'
			profiles = append(profiles, p)
		}

	}
	//print out using JSON
	c.JSON(http.StatusOK, gin.H{"profile": profiles})

}

//Create function
func CreateItem(c *gin.Context) {
	//CREATE PROFILE POSTED FROM POSTMAN CLIENT
	//Params function returns the value of the URL param --> I STILL DONT UNDERSTAND, DOES THIS MEAN WE RETURN THE VALUE OF FIRST_NAME
	//ONLY OR WHAT??

	var createProfile Profile

	err := c.ShouldBindJSON(&createProfile)

	//	checker := c.ShouldBindJSON(&createProfile)

	//validate
	if createProfile.FirstName == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter profile"})
		return
	} else {

		//test error query nya
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message-1": err.Error()})
			return
		} else {
			//kalo ga ada masalah sama query nya kita bisa pakai empty array yang type Profile buat panggil FirstName, LastName, Age
			_, err := db.Query("INSERT INTO profile VALUES ($1, $2, $3);", createProfile.FirstName, createProfile.LastName, createProfile.Age)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message-2": err.Error()})
				return
			}
		}

		//insert into DB
		log.Print("Successfuly inserted into DB")

		//show success
		c.JSON(http.StatusCreated, gin.H{"profile": createProfile})
	}
}

//Update function
func UpdateItem(c *gin.Context) {
	var oldProfile Profile
	var updateProfile Profile

	err := c.ShouldBindJSON(&updateProfile)

	if updateProfile.FirstName == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter a first name"})
		return
	} else {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message-1": err.Error()})
			return
		} else {
			//Scan function copies the columns from the match row into the values pointed at the destination
			//old profile cuma sebagai tempat nampung data data yang lama karena disini kita pakai Scan function
			err := db.QueryRow("SELECT * FROM profile WHERE first_name=$1", &updateProfile.FirstName).Scan(&oldProfile.FirstName, &oldProfile.LastName, &oldProfile.Age)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message-2": err.Error()})
				return
			} else {
				//'updateProfile' itu yang berisi permintaan dari post request
				_, err := db.Query("UPDATE profile SET last_name=$1 WHERE first_name=$2", updateProfile.LastName, updateProfile.FirstName)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message-3": "error with DB"})
					return
				}
			}
		}
		//update DB
		log.Print("Successfuly update the DB")

		//show success
		c.JSON(http.StatusOK, gin.H{"profile": updateProfile})
	}
}

func DeleteItem(c *gin.Context) {

	var oldProfile Profile
	var deleteProfile Profile
	err := c.ShouldBindJSON(&deleteProfile)

	if deleteProfile.FirstName == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter a first name"})
		return
	} else {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message-1": err.Error()})
			return
		} else {
			err := db.QueryRow("SELECT * FROM profile WHERE first_name=$1", deleteProfile.FirstName).Scan(&oldProfile.FirstName, &oldProfile.LastName, &oldProfile.Age)
			if err != nil {
				c.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
				return
			} else {
				_, err := db.Query("DELETE FROM profile WHERE first_name=$1", &deleteProfile.FirstName)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			}
		}
	}

	//delete from DB
	log.Print("Deleted from DB")

	//show success
	c.JSON(http.StatusOK, gin.H{"Profile deleted": deleteProfile})
}
