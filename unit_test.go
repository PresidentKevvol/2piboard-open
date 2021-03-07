package main

import (
  "fmt"
  "log"
  "math/rand"
  "time"
  "testing"
  "github.com/joho/godotenv"
)

func DBTestsSetup() {
  //load the .env file
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  //setup the database interface
  DBAccessSetup()
  //set the seed for the standard prng
  rand.Seed(time.Now().UnixNano())
  fmt.Printf("DB setup finished\n")
}

//test for creating a tripcode
func TestTripcodeInsert(t *testing.T) {
  DBTestsSetup()

  id, err := createTripcode("Jerry", "pasta_sauce111")
  if err != nil {
    //t.Logf("TestTripcodeInsert error: %s", err)
    t.Errorf("TestTripcodeInsert error: %s", err)
  } else {
    t.Logf("got id for new tripcode: %d\n", id)
  }
}

//test for matching a tripcode
func TestTripcodeMatch(t *testing.T) {
  DBTestsSetup()

  id, err := getTripcodeId("Jerry", "pasta_sauce111")
  if err != nil {
    t.Errorf("TestTripcodeMatch error: %s", err)
  } else {
    t.Logf("got id for tripcode: %d\n", id)
  }
}

func TestCreatePostTrip(t *testing.T) {
  DBTestsSetup()

  p := CreatePost {Title: "Moe?", Body: "", TripcodeUsername:"Jerry", TripcodePassword:"pasta_sauce111", Board:"vg"}
  ip := "192.168.0.72"

  err := addPostToDB(p, ip)
  if err != nil {
    t.Errorf("TestCreatePost error: %s", err)
  } else {
    t.Logf("test post create successful\n")
  }
}

func TestCreatePostAnon(t *testing.T) {
  DBTestsSetup()

  p := CreatePost {Title: "test", Body: "this is a test anonymous post", TripcodeUsername:"", TripcodePassword:"", Board:"vg"}
  ip := "192.168.0.76"

  err := addPostToDB(p, ip)
  if err != nil {
    t.Errorf("TestCreatePostAnon error: %s", err)
  } else {
    t.Logf("test anon post create successful\n")
  }
}

func TestCreatePostTripWrongPW(t *testing.T) {
  DBTestsSetup()

  p := CreatePost {Title: "Moe?", Body: "", TripcodeUsername:"Jerry", TripcodePassword:"pasta_sauce112", Board:"vg"}
  ip := "192.168.0.72"

  err := addPostToDB(p, ip)
  if err != nil {
    t.Errorf("TestCreatePost error: %s", err)
  } else {
    t.Logf("test post create successful\n")
  }
}

func TestAddPostTripWrongPW(t *testing.T) {
  DBTestsSetup()

  p := CreatePost {Title: "Wrong password post test", Body: "", TripcodeUsername:"Jerry", TripcodePassword:"pasta_sauce112", Board:"vg"}
  ip := "192.168.0.72"

  ret := addPost(p, ip)
  t.Logf("addPost(p, ip) with wrong tripcode returns: %v", ret)
  //t.Errorf("TestAddPostTripWrongPW error: no value returned")
}

func TestAddReplyTrip(t *testing.T) {
  DBTestsSetup()

  p := CreateReply {Body: "junja fdjhfdguihfsd", TripcodeUsername:"Dena", TripcodePassword:"whendy_frie5", Parent:"290413355379475811"}
  ip := "192.168.0.90"

  ret := addReply(p, ip)
  t.Logf("addReply(p, ip) with tripcode returns: %v", ret)
}

/*
//the test main function to be called
func TestMain(m *testing.M) {

}
*/
