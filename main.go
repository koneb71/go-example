package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

var dbhost string = "dbeunorth.mariadb.database.azure.com"
var dbport string = "3306"
var dbname string = "test_go"
var dbuser string = "bgradm@dbeunorth"
var dbpassword string = "sD8f_3Hkjhu-Uy!R"

// Configure the database connection (always check errors)
var db, err = sql.Open("mysql", dbuser + ":" + dbpassword + "@(" + dbhost + ":" + dbport +")/" + dbname + "?parseTime=true")


func populateDb() {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME,
        PRIMARY KEY (id)
    );`

	// Executes the SQL query in our database. Check err to ensure there was no error.
	_, err = db.Exec(query)

	var (
		id int
	)

	query2 := `SELECT id FROM users WHERE id = ?`
	err := db.QueryRow(query2, 1).Scan(&id)
	if err != nil { // if there is an error
		log.Print("id not exist", err) // log it
	}

	if id == 0 {
		user := "admin"
		password := "admin"
		createdAt := time.Now()

		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, user, password, createdAt)
		if err != nil { // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		if result != nil { // if there is an error
			log.Print("template parsing error: ", result) // log it
		}
	}
}

func main() {
	populateDb()
	http.HandleFunc("/", HomePage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type PageVariables struct {
	ID        int
	USERNAME  string
	PASSWORD  string
	CREATEDAT string
}

func HomePage(w http.ResponseWriter, r *http.Request){
	var (
		id        int
		username  string
		password  string
		createdAt string
	)

	// Query the database and scan the values into out variables. Don't forget to check for errors.
	query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
	err := db.QueryRow(query, 1).Scan(&id, &username, &password, &createdAt)

	HomePageVars := PageVariables{ //store the date and time in a struct
		ID: id,
		USERNAME: username,
		PASSWORD: password,
		CREATEDAT: createdAt,
	}

	t, err := template.ParseFiles("templates/index.html") //parse the html file homepage.html
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil { // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}