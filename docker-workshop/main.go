package main

import (
    "database/sql"
    "encoding/json" // package to encode and decode the json into struct and vice versa
    "fmt"
    "log"
    "net/http" // used to access the request and response object of the api
    "strconv"  // package used to covert string into int type
    "github.com/gorilla/mux" // used to get the params from the route
    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "adm"
    password = "pass"
    dbname   = "songs"
)

// User schema of the user table
type Song struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Singer string `json:"singer"`
    Genre  string  `json:"genre"`
}

// response format
type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

func main() {
    r := Router()
    // fs := http.FileServer(http.Dir("build"))
    // http.Handle("/", fs)
    fmt.Println("Starting server on the port 8080...")

    log.Fatal(http.ListenAndServe(":8080", r))
}

// Router is exported and used in main.go
func Router() *mux.Router {

    router := mux.NewRouter()

	router.HandleFunc("/health", Health)    
    router.HandleFunc("/api/songs", GetAllSongs).Methods("GET", "OPTIONS")
    router.HandleFunc("/api/newsong", CreateSong).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/song/{id}", UpdateSong).Methods("PUT", "OPTIONS")
    router.HandleFunc("/api/deletesong/{id}", DeleteSong).Methods("DELETE", "OPTIONS")

    return router
}

// CreateUser create a user in the postgres db
func Health(w http.ResponseWriter, r *http.Request) {
    // set the header to content type x-www-form-urlencoded
    // Allow all origin to handle cors issue
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    createConnection()

}



// create connection with postgres db
func createConnection() *sql.DB {
    psqlconn := fmt.Sprintf("postgresql://postgres:password@postgres_db:5432?sslmode=disable")
    //psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    //"password=%s dbname=%s sslmode=disable",
    //host, port, user, password)

    // Open the connection
    db, err := sql.Open("postgres", psqlconn)

    if err != nil {
        panic(err)

	}

    // check the connection
    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected to DB")
    // return the connection
    return db
}

// CreateUser create a user in the postgres db
func CreateSong(w http.ResponseWriter, r *http.Request) {
    // set the header to content type x-www-form-urlencoded
    // Allow all origin to handle cors issue
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // create an empty user of type models.User
    var song Song

    // decode the json request to user
    err := json.NewDecoder(r.Body).Decode(&song)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    // call insert user function and pass the user
    insertID := insertSong(song)

    // format a response object
    res := response{
        ID:      insertID,
        Message: "Song created successfully",
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

// GetUser will return a single user by its id
func GetSong(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // get the userid from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id type from string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // call the getUser function with user id to retrieve a single user
    song, err := getSong(int64(id))

    if err != nil {
        log.Fatalf("Unable to get user. %v", err)
    }

    // send the response
    json.NewEncoder(w).Encode(song)
}

// GetAllUser will return all the users
func GetAllSongs(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    var songs []Song
    songs, err := getAllSongs()

    if err != nil {
        log.Fatalf("Unable to get all user. %v", err)
    }

    // send all the users as response
    json.NewEncoder(w).Encode(songs)
}

// UpdateUser update user's detail in the postgres db
func UpdateSong(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // get the userid from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id type from string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // create an empty user of type models.User
    var song Song

    // decode the json request to user
    err = json.NewDecoder(r.Body).Decode(&song)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    // call update user to update the user
    updatedRows := updateSong(int64(id), song)

    // format the message string
    msg := fmt.Sprintf("Song updated successfully. Total rows/record affected %v", updatedRows)

    // format the response message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteSong(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // get the userid from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id in string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // call the deleteUser, convert the int to int64
    deletedRows := deleteSong(int64(id))

    // format the message string
    msg := fmt.Sprintf("Song deleted successfully. Total rows/record affected %v", deletedRows)

    // format the reponse message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one user in the DB
func insertSong(song Song) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the insert sql query
    // returning userid will return the id of the inserted user
    sqlStatement := `INSERT INTO "songs" (name, singer, genre) VALUES ($1, $2, $3) RETURNING songid`

    // the inserted id will store in this id
    var id int64

    // execute the sql statement
    // Scan function will save the insert id in the id
    err := db.QueryRow(sqlStatement, song.Name, song.Singer, song.Genre).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    // return the inserted id
    return id
}

// get one user from the DB by its userid
func getSong(id int64) (Song, error) {
    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create a user of models.User type
    var song Song

    // create the select sql query
    sqlStatement := `SELECT * FROM songs WHERE songid=$1`

    // execute the sql statement
    row := db.QueryRow(sqlStatement, id)

    // unmarshal the row object to user
    err := row.Scan(&song.ID, &song.Name, &song.Singer, &song.Genre)

    switch err {
    case sql.ErrNoRows:
        fmt.Println("No rows were returned!")
        return song, nil
    case nil:
        return song, nil
    default:
        log.Fatalf("Unable to scan the row. %v", err)
    }

    return song, err
}

// get one user from the DB by its userid
func getAllSongs() ([]Song, error) {
    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    var songs []Song

    // create the select sql query
    sqlStatement := `SELECT * FROM songs`

    // execute the sql statement
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // close the statement
    defer rows.Close()

    // iterate over the rows
    for rows.Next() {
        var song Song

        // unmarshal the row object to user
        err = rows.Scan(&song.ID, &song.Name, &song.Singer, &song.Genre)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        // append the user in the users slice
        songs = append(songs, song)

    }

    // return empty user on error
    return songs, err
}

// update user in the DB
func updateSong(id int64, song Song) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the update sql query
    sqlStatement := `UPDATE songs SET name=$2, singer=$3, genre=$4 WHERE songid=$1`

    // execute the sql statement
    res, err := db.Exec(sqlStatement, id, song.Name, song.Singer, song.Genre)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

// delete user in the DB
func deleteSong(id int64) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the delete sql query
    sqlStatement := `DELETE FROM songs WHERE songid=$1`

    // execute the sql statement
    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}