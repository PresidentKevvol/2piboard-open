package main

import (
  "fmt"
  "os"
  "time"
  "context"
  "math/rand"
  "github.com/jackc/pgx/v4/pgxpool"
  "golang.org/x/crypto/bcrypt"
)

type CreatePost struct {
  Title string              `json:"Title"`
  Body string               `json:"Body"`
  TripcodeUsername string   `json:"TripcodeUsername"`
  TripcodePassword string   `json:"TripcodePassword"`
  Board string              `json:"Board"`
}

type CreatePostRequestResponse struct {
  Code int64
  Message string
}

type CreateReply struct {
  Body string               `json:"Body"`
  TripcodeUsername string   `json:"TripcodeUsername"`
  TripcodePassword string   `json:"TripcodePassword"`
  Parent string             `json:"Parent"`
}

type CreateReplyRequestResponse struct {
  PostId string
  Code int64
  Message string
}

func addPost(p CreatePost, ip string) CreatePostRequestResponse {
  //fmt.Println("Post create received: ", p)
  //the result object to return
  var res CreatePostRequestResponse

  if p.Title == "" {
    res = CreatePostRequestResponse {Code: 1, Message: "error: Can't create a post without a title."}
  } else if p.TripcodeUsername != "" && p.TripcodePassword == "" {
    res = CreatePostRequestResponse {Code: 1, Message: "error: Can't use a tripcode without a password"}
  } else if p.TripcodeUsername == "" && p.TripcodePassword != "" {
    res = CreatePostRequestResponse {Code: 1, Message: "error: Can't use a tripcode without a username"}
  } else {
    //if all good, insert post into database
    //add code to insert post to database
    err := addPostToDB(p, ip)
    if err != nil {//if error exists
      if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
        res = CreatePostRequestResponse {Code: 1, Message: "error: this tripcode is taken, or your username does not match password"}
      } else {
        res = CreatePostRequestResponse {Code: 1, Message: "error: database problems"}
      }
    } else {//if no error
      res = CreatePostRequestResponse {Code: 0, Message: "Post created"}
    }
  }
  return res
}

func addReply(r CreateReply, ip string) CreateReplyRequestResponse {
  //fmt.Println("Reply create received: ", r)
  //the result object to return
  var res CreateReplyRequestResponse

  if r.TripcodeUsername != "" && r.TripcodePassword == "" {
    res = CreateReplyRequestResponse {PostId: r.Parent, Code: 1, Message: "error: Can't use a tripcode without a password"}
  } else if r.TripcodeUsername == "" && r.TripcodePassword != "" {
    res = CreateReplyRequestResponse {PostId: r.Parent, Code: 1, Message: "error: Can't use a tripcode without a username"}
  } else {
    //if all good, insert post into database
    //add code to insert post to database
    err := addReplyToDB(r, ip)
    if err != nil {//if error exists
      if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
        res = CreateReplyRequestResponse {PostId: r.Parent, Code: 1, Message: "error: this tripcode is taken, or your username does not match password"}
      } else {
        res = CreateReplyRequestResponse {PostId: r.Parent, Code: 1, Message: "error: database problems"}
      }
    } else {//if no error
      res = CreateReplyRequestResponse {PostId: r.Parent, Code: 0, Message: "Reply created"}
    }
  }
  return res
}

var DBUsername string
var DBPassword string
var DBURL string
var DBDatabaseName string

var connPool *pgxpool.Pool

var QUERY_INSERT_POST string = "INSERT INTO post (id, title, body, time_posted, author, in_board, poster_ip) VALUES ($1, $2, $3, $4, $5, $6, $7);"
var QUERY_INSERT_REPLY string = "INSERT INTO reply (id, body, time_posted, author, parent, poster_ip) VALUES ($1, $2, $3, $4, $5, $6);"
var QUERY_CHECK_TRIPCODE_EXIST string = "SELECT id, username, hash FROM tripcode WHERE username=$1;"
var QUERY_CREATE_TRIPCODE string = "INSERT INTO tripcode (id, username, hash) VALUES ($1, $2, $3);"

//sets up postgreSQL database access
//called in the main function
func DBAccessSetup() {
  //get the values from the .env file
  DBUsername = os.Getenv("database_username")
  DBPassword = os.Getenv("database_password")
  DBURL = os.Getenv("database_url")
  DBDatabaseName = os.Getenv("database_name")
  full_url := fmt.Sprintf("postgres://%s:%s@%s/%s", DBUsername, DBPassword, DBURL, DBDatabaseName)
  var err error
  connPool, err = pgxpool.Connect(context.Background(), full_url)
	if err != nil {
    //if database cannot be accessed, abort
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

  //test code
  var entre string
  err = connPool.QueryRow(context.Background(), "SELECT (entry) FROM access_check;").Scan(&entre)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	} else {
    //if no error, this is to test if the database can be accessed and accessed properly
    //exit will be removed later
    fmt.Println("Database access successful: ", entre)
  }
}

//creates(registers) a new tripcode in the database
func createTripcode(username, password string) (int64, error) {
  //hash the password using bcrypt with cost 14
  hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  if err != nil {
    fmt.Println("hash failure")
    return 0, err
  }

  //convert the byte array to string
  hash := string(hashBytes)
  //generate a random (+ve 63 bit) number as id
  id := rand.Int63()

  //insert the id, username, and hash into DB
  _, execErr := connPool.Exec(context.Background(), QUERY_CREATE_TRIPCODE, id, username, hash)
  if execErr != nil {
    fmt.Println("tripcode insert failure")
    return 0, execErr
  }
  //return nil for no error
  return id, nil
}

func getTripcodeId(username, password string) (int64, error) {
  var id int64
  var passHash string
  //try to get from the db a tripcode with the same username
  rows, err := connPool.Query(context.Background(), QUERY_CHECK_TRIPCODE_EXIST, username)
  if err != nil {//if the db can't be accessed at all, return the err
    fmt.Printf("querying error: %s", err)
    return 0, err
  }

  if rows.Next() == false {//if there is 'no more' rows (not even the 1st row)
    //return -1 denoting there is no id for this tripcode, and nil means no error
    return -1, nil
  }

  //now scan the row to get the hash and numeric id
  scanErr := rows.Scan(&id, nil, &passHash)
  //close the query so resources can be freed up
  rows.Close()
  if scanErr != nil {//if there is an error in scaning, return with the scan error
    fmt.Printf("scanning error: %s", scanErr)
    return 0, scanErr
  }

  //after that, check if the hash match
  hashErr := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
  if hashErr != nil {//if hash does not match, it will throw an error
    return 0, hashErr
  }

  //if everything is good, we can now return the id for this tripcode
  return id, nil
}

//function to add the post to the DB
func addPostToDB(p CreatePost, ip string) error {
  //if no tripcode supplied (anonymous post)
  if p.TripcodeUsername == "" {
    //insert the post into db anonymously
    id := rand.Int63() //the post's id
    time := time.Now() //the time of the post's creation
    _, err := connPool.Exec(context.Background(), QUERY_INSERT_POST, id, p.Title, p.Body, time, -1, p.Board, ip)
    if err != nil  {
      //if error, return with error
      return err
    }
  } else {//if tripcode supplied
    //grab the id corresponding to the tripcode from the DB
    trip_id, err := getTripcodeId(p.TripcodeUsername, p.TripcodePassword)
    if err != nil  {
      //if error, return with error
      return err
    } else if trip_id == -1 {//or, if the tripcode have not existed yet
      //create one in the database
      trip_id, err = createTripcode(p.TripcodeUsername, p.TripcodePassword)
      if err != nil  {
        //if error, return with error
        return err
      }
    }
    //now we can insert the post into db
    id := rand.Int63() //the post's id
    time := time.Now() //the time of the post's creation
    _, err = connPool.Exec(context.Background(), QUERY_INSERT_POST, id, p.Title, p.Body, time, trip_id, p.Board, ip)
    if err != nil  {
      //if error, return with error
      return err
    }
  }
  //return nil for no error
  return nil
}

//function to add a reply to a post to the DB
func addReplyToDB(p CreateReply, ip string) error {
  //if no tripcode supplied (anonymous post)
  if p.TripcodeUsername == "" {
    //insert the post into db anonymously
    id := rand.Int63() //the post's id
    time := time.Now() //the time of the post's creation
    _, err := connPool.Exec(context.Background(), QUERY_INSERT_REPLY, id, p.Body, time, -1, p.Parent, ip)
    if err != nil  {
      //if error, return with error
      return err
    }
  } else {//if tripcode supplied
    //grab the id corresponding to the tripcode from the DB
    trip_id, err := getTripcodeId(p.TripcodeUsername, p.TripcodePassword)
    if err != nil  {
      //if error, return with error
      return err
    } else if trip_id == -1 {//or, if the tripcode have not existed yet
      //create one in the database
      trip_id, err = createTripcode(p.TripcodeUsername, p.TripcodePassword)
      if err != nil  {
        //if error, return with error
        return err
      }
    }
    //now we can insert the post into db
    id := rand.Int63() //the post's id
    time := time.Now() //the time of the post's creation
    _, err = connPool.Exec(context.Background(), QUERY_INSERT_REPLY, id, p.Body, time, trip_id, p.Parent, ip)
    if err != nil  {
      //if error, return with error
      return err
    }
  }
  //return nil for no error
  return nil
}
