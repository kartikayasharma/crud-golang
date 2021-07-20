package middleware

import (
	"crud/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world")
	createConnection()
}

func createConnection() *sql.DB {
	connStr := "user=postgres password=kanishk dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	insertedID := insertUser(user)

	res := response{
		ID:      insertedID,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	db := createConnection()
	defer db.Close()

	var res []models.User

	rs, err := db.Query(`select * from users`)
	defer rs.Close()

	if err != nil {
		log.Fatalf("Unable to execute the get all users query. %v", err)
	}
	for rs.Next() {
		var u models.User
		err := rs.Scan(&u.Name, &u.Location, &u.Age, &u.ID)
		if err != nil {
			log.Fatalf("Unable to scan all users query. %v", err)
		}
		res = append(res, u)
	}
	// fmt.Printf("%+v", send)
	fmt.Printf("Sent complete record")

	json.NewEncoder(w).Encode(res)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	id := params["id"]
	var res models.User

	db := createConnection()
	defer db.Close()

	json.NewDecoder(r.Body).Decode(&id)

	err := db.QueryRow(`select * from users where id = $1`, id).Scan(&res.Name, &res.Location, &res.Age, &res.ID)

	if err != nil {
		log.Fatalf("Unable to execute the select by id query. %v", err)
	}
	fmt.Printf("Sent a single record %v", id)

	json.NewEncoder(w).Encode(res)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	db := createConnection()
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	var u models.User

	json.NewDecoder(r.Body).Decode(&u)

	_, err := db.Query(`update users set name = $2, location = $3, age = $4 where id = $1`, id, u.Name, u.Location, u.Age)

	if err != nil {
		log.Fatalf("Unable to execute the update query. %v", err)
	}

	res := response{
		ID:      id,
		Message: "Added Successfully",
	}
	u.ID = id
	fmt.Printf("Updated a single record %+v", u)

	json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	db := createConnection()

	_, err := db.Query(`delete from users where id = $1`, id)

	if err != nil {
		log.Fatalf("Unable to execute the delete query. %v", err)
	}

	res := response{
		ID:      id,
		Message: "Deleted Successfully",
	}
	fmt.Printf("Deleted a single record %v", id)

	json.NewEncoder(w).Encode(res)
}

func insertUser(u models.User) int64 {
	db := createConnection()
	defer db.Close()

	q := `insert into users (name, location, age) values ($1, $2, $3) returning id`

	var id int64
	err := db.QueryRow(q, u.Name, u.Location, u.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the insert query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}
