package main

import (
  "fmt"
  "os"
  "strings"
  "strconv"
  "log"
  "math/rand"
  "time"
  "net/http"
  "path/filepath"
  "html/template"
  "encoding/json"
  "github.com/joho/godotenv"
	//"io/ioutil"
  //"github.com/gorilla/websocket"
)

//template for the pages
var ex, exerr = os.Executable()

func getWorkDir() {
  workdir := filepath.Dir(ex)
  //for when the app is dockerized, prevent the strings for paths become like "//views/index.html"
  if workdir == "/" {
    workdir = ""
  }
  return workdir
}
var workdir = getWorkDir()

var chat_template = template.Must(template.ParseFiles(workdir + "/views/index.html", workdir + "/views/board_template.html", workdir + "/views/footer.html", workdir + "/views/header_all.html"))

//the handler for the index page
func handleIndex(w http.ResponseWriter, r *http.Request) {
  //get the stuff needed
  contex := getFrontPageContext()
  //render template
  err := chat_template.ExecuteTemplate(w, "index.html", contex)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
//handler for a board's page with rendering the template but no cache
func handleBoard(w http.ResponseWriter, r *http.Request) {
  //get the link extention from the url
  stention := r.URL.Path[len("/board/"):]
  var board string
  var page int
  //split to get the page number
  splits := strings.Split(stention, "/")
  if len(splits) == 1 {
    board = splits[0]
    page = 1
  } else {
    board = splits[0]
    paj, err := strconv.Atoi(splits[1])
    if err != nil {
      page = 0
    } else {
      page = paj
    }
  }

  //get the content of the board with the specified page number
  content := getBoardContent(board, page)
  if content.Posts == nil {
    http.Error(w, "can't get board", http.StatusInternalServerError)
  }
  //execute the template
  err := chat_template.ExecuteTemplate(w, "board_template.html", content)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
//the handler for the footer
func handleFooter(w http.ResponseWriter, r *http.Request) {
  err := chat_template.ExecuteTemplate(w, "footer.html", nil)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
//the handler for the header
func handleHeaderAll(w http.ResponseWriter, r *http.Request) {
  err := chat_template.ExecuteTemplate(w, "header_all.html", nil)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

//function to fetch the client's ip from the http header
func ReadUserIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
}

//the handler function for the ajax request to make a post
func handleAjaxCreatePost(w http.ResponseWriter, r *http.Request) {
  var p CreatePost
  //fmt.Println("got Request: ", r.Body)
  //parse the json in the body text into a Post object
  err := json.NewDecoder(r.Body).Decode(&p)
  if err != nil {
    //fmt.Println("Decoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //get the ip from the request header
  ip := ReadUserIP(r)
  //the reply object
  var res CreatePostRequestResponse
  //call the addPost method to add the post to db
  res = addPost(p, ip)

  // create json response from serializing a CreatePostReply struct
  a, err := json.Marshal(res)
  if err != nil {
    //fmt.Println("Encoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(a)
}
//the handler function for the ajax request to make a reply to a post
func handleAjaxCreateReply(w http.ResponseWriter, r *http.Request) {
  var rp CreateReply
  //fmt.Println("got Request: ", r.Body)
  //parse the json in the body text into a Post object
  err := json.NewDecoder(r.Body).Decode(&rp)
  if err != nil {
    //fmt.Println("Decoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //get the ip from the request header
  ip := ReadUserIP(r)
  //the reply object
  var res CreateReplyRequestResponse
  //call the addPost method to add the post to db
  res = addReply(rp, ip)

  // create json response from serializing a CreatePostReply struct
  a, err := json.Marshal(res)
  if err != nil {
    //fmt.Println("Encoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(a)
}

//for handling a show replies request
func handleAjaxShowReplies(w http.ResponseWriter, r *http.Request) {
  var rp ShowRepliesRequest
  //fmt.Println("got Request: ", r.Body)
  //parse the json in the body text into a Post object
  err := json.NewDecoder(r.Body).Decode(&rp)
  if err != nil {
    //fmt.Println("Decoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  //get the ip from the request header
  //ip := ReadUserIP(r)
  //the reply object
  var res ShowRepliesRequestResponse
  //call the addPost method to add the post to db
  res = showReplies(rp)

  // create json response from serializing a CreatePostReply struct
  a, err := json.Marshal(res)
  if err != nil {
    //fmt.Println("Encoding error: ", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Write(a)
}

//handler for index but no cache
func handleIndexNocache(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("views/index.html")
  err := t.Execute(w, nil)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
//handler for a board's page but no cache
func handleBoardNocache(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("views/board.html")
  err := t.Execute(w, nil)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
//handler for a board's page with rendering the template but no cache
func handleBoardTemplateNocache(w http.ResponseWriter, r *http.Request) {
  //get the link extention from the url
  stention := r.URL.Path[len("/nocache/board_template/"):]
  var board string
  var page int
  //split to get the page number
  splits := strings.Split(stention, "/")
  if len(splits) == 1 {
    board = splits[0]
    page = 1
  } else {
    board = splits[0]
    paj, err := strconv.Atoi(splits[1])
    if err != nil {
      page = 0
    } else {
      page = paj
    }
  }

  t, _ := template.ParseFiles("views/board_template.html")
  //get the content of the board with the specified page number
  content := getBoardContent(board, page)
  if content.Posts == nil {
    http.Error(w, "can't get board", http.StatusInternalServerError)
  }
  //execute the template
  err := t.Execute(w, content)
  if err != nil { //if there is an error
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {
  //load the .env file
  err := godotenv.Load(workdir + "/.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  //setup the database interface
  DBAccessSetup()
  //set the seed for the standard prng
  rand.Seed(time.Now().UnixNano())

  //print the current working directory
  fmt.Println("directory of this executable: " + workdir)

  //define handler functions for pages
  http.HandleFunc("/", handleIndex)
  http.HandleFunc("/board/", handleBoard)
  http.HandleFunc("/footer", handleFooter)
  http.HandleFunc("/header_all", handleHeaderAll)

  //the handler for ajax requests
  http.HandleFunc("/ajax/createpost/", handleAjaxCreatePost)
  http.HandleFunc("/ajax/createreply/", handleAjaxCreateReply)
  http.HandleFunc("/ajax/showreplies/", handleAjaxShowReplies)

  /*
  http.HandleFunc("/nocache/", handleIndexNocache)
  http.HandleFunc("/nocache/board", handleBoardNocache)
  http.HandleFunc("/nocache/board_template/", handleBoardTemplateNocache)
  */

  //for handling static css and js files
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(workdir + "/static/css"))))
  http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(workdir + "/static/js"))))
  http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(workdir + "/static/fonts"))))
  http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(workdir + "/static/images"))))
  //setup the hostname env variable
  hostname := os.Getenv("host_name")
  if hostname == "" {
    hostname = ":8880"
  }
  //see if there is a ssl cert to be used
  ssl_cert := os.Getenv("ssl_cert")
  ssl_key := os.Getenv("ssl_key")
  if ssl_cert != "" && ssl_key != "" {
    //start the server program with ssl
    log.Fatal(http.ListenAndServeTLS(hostname, ssl_cert, ssl_key, nil))
  } else {
    //start the server program
    log.Fatal(http.ListenAndServe(hostname, nil))
  }
}
