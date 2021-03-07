package main
/*
functions for rendering the posts in a board page
*/

import (
  "fmt"
  //"os"
  "time"
  "context"
  "strconv"
)

//represents a post
type PostItem struct {
  PostId int64
  Title string
  Body string
  Author string
  TimePosted string
  Replies []PostReplyItem
  NumTotalReplies int32
}

//represents a reply to a post
type PostReplyItem struct {
  Body string
  Author string
  TimePosted string
}

//represents the contents of a page of a board
type BoardPageObject struct {
  Url string
  Title string
  Posts []PostItem
  CurrentPage int
}

type ShowRepliesRequest struct {
  PostId string           `json:"PostId"`
}

type ShowRepliesRequestResponse struct {
  PostId string
  Replies []PostReplyItem
  Status int
  Message string
}

func getPostsInBoardTest(url string) []PostItem {
  //test code
  reply1 := PostReplyItem {Body: "Butter chicken", Author: "Anonymous", TimePosted: time.Now().Format("2006-01-02 15:04:05")}
  reply2 := PostReplyItem {Body: "REEEEEEE my waifu", Author: "Sharkquille", TimePosted: time.Now().Format("2006-01-02 15:04:05")}
  replies := []PostReplyItem {reply1, reply2,}
  poste := PostItem {PostId: 1122, Title: "Po Po Po Po?", Body: "when I will grow up?", Author: "papius", TimePosted: time.Now().Format("2006-01-02 15:04:05"), Replies:replies, NumTotalReplies: 2}
  return []PostItem {poste,}
}

func getBoardContentTest(url string, page int) BoardPageObject {
  return BoardPageObject {Url: url, Title: "Video Gaymes", Posts: getPostsInBoardTest(""), CurrentPage: page}
}

var QUERY_GET_BOARDS string = "SELECT url, board.title FROM board;"
var QUERY_GET_BOARD_DESCRIPTION string = "SELECT title FROM board WHERE url = $1;"
var QUERY_GET_POSTS_IN_BOARD string = "SELECT post.id, title, body, time_posted, tripcode.username FROM post, tripcode WHERE in_board = $1 and post.author = tripcode.id ORDER BY time_posted DESC LIMIT 10 OFFSET $2;"
var QUERY_GET_REPLY_TO_POST string = "SELECT body, time_posted, tripcode.username FROM reply, tripcode WHERE parent = $1 and reply.author = tripcode.id ORDER BY time_posted DESC;"
var PAGE_LENGTH int = 10

//get all posts in a specific page from a board
func getPostsInBoard(url string, page_num int) ([]PostItem, error) {
  //page number determine offset
  offset := (page_num - 1) * PAGE_LENGTH
  //get the rows from the database
  rows, err := connPool.Query(context.Background(), QUERY_GET_POSTS_IN_BOARD, url, offset)
  if err != nil {//if the db can't be accessed at all, return the err
    fmt.Printf("querying error: %s", err)
    return make([]PostItem, 0), err
  }
  //the result array to send back to client
  res := make([]PostItem, 0)
  for i:=0; i<PAGE_LENGTH; i+=1 {
    if rows.Next() == false {
      break
    }
    var PostId int64
    var Title string
    var Body string
    var Author string
    var TimePosted time.Time
    //scan to items
    rows.Scan(&PostId, &Title, &Body, &TimePosted, &Author)
    post := PostItem {PostId: PostId, Title: Title, Body: Body, Author: Author, TimePosted: TimePosted.Format("2006-01-02 15:04:05"), Replies:make([]PostReplyItem, 0), NumTotalReplies: 0}
    res = append(res, post)
  }
  //close the connection
  rows.Close()

  return res, nil
}

func getBoardContent(url string, page int) BoardPageObject {
  var boardTitle string
  err := connPool.QueryRow(context.Background(), QUERY_GET_BOARD_DESCRIPTION, url).Scan(&boardTitle)
  if err != nil {
    boardTitle = ""
  }
  posts, queryErr := getPostsInBoard(url, page)
  if queryErr != nil {
    return BoardPageObject {Url: url, Title: boardTitle, Posts: nil, CurrentPage: page}
  }
  return BoardPageObject {Url: url, Title: boardTitle, Posts: posts, CurrentPage: page}
}

//get all reply to a post by the post's id
func getReplyToPost(postId int64) ([]PostReplyItem, error) {
  //get the rows from the database
  rows, err := connPool.Query(context.Background(), QUERY_GET_REPLY_TO_POST, postId)
  if err != nil {//if the db can't be accessed at all, return the err
    fmt.Printf("querying error: %s", err)
    return make([]PostReplyItem, 0), err
  }
  //the result array to send back to client
  res := make([]PostReplyItem, 0)
  for i:=0; ; i+=1 {
    if rows.Next() == false {
      break
    }
    var Body string
    var Author string
    var TimePosted time.Time
    //scan to items
    rows.Scan(&Body, &TimePosted, &Author)
    post := PostReplyItem {Body: Body, Author: Author, TimePosted: TimePosted.Format("2006-01-02 15:04:05")}
    res = append(res, post)
  }
  //close the connection
  rows.Close()

  return res, nil
}

func showReplies(s ShowRepliesRequest) ShowRepliesRequestResponse {
  //get the id
  id, err := strconv.ParseInt(s.PostId, 10, 64)
  // return in case error
  errReturn := ShowRepliesRequestResponse {PostId: s.PostId, Replies: make([]PostReplyItem, 0), Status: 1, Message: "Can't get replies"}
  if err != nil {
    return errReturn
  }
  //now fetch the replies
  aray, getErr := getReplyToPost(id)
  if getErr != nil {
    return errReturn
  }
  //if all good, we return
  return ShowRepliesRequestResponse {PostId: s.PostId, Replies: aray, Status: 0, Message: "success"}
}

type Board struct {
  Url string
  Description string
  //PostCount int
}

//function for getting all available boards for rendering the front page
func getAvailableBoards() []Board {
  //get all the boards
  rows, err := connPool.Query(context.Background(), QUERY_GET_BOARDS)
  if err != nil {
    return nil
  }
  //the result array to send back to client
  res := make([]Board, 0)
  for i:=0; ; i+=1 {
    if rows.Next() == false {
      break
    }
    var url string
    var desc string
    //var count int
    //scan to items
    rows.Scan(&url, &desc)
    b := Board {Url: url, Description: desc}
    res = append(res, b)
  }
  //close the connection
  rows.Close()
  //record time and return the array
  lastCached = time.Now()
  return res
}

var boardsListCached []Board = nil
var lastCached time.Time = time.Now()

type FrontPageContext struct {
  Boards []Board
}

func getFrontPageContext() FrontPageContext {
  //everytime the front page is accessed, check cached array
  if boardsListCached == nil {
    //if it is null, try reloading it
    boardsListCached = getAvailableBoards()
  } else if lastCached.Add(time.Second * 60).Before(time.Now()) {
    //when time up, reload it
    boardsListCached = getAvailableBoards()
  }
  //can add other code here if needed

  return FrontPageContext {Boards: boardsListCached}
}
