//event handler when 'post' button is pressed
function postButtonPressed () {
  var post_title = document.getElementById("text-title").value;
  var post_body = document.getElementById("text-body").value;
  var tripcode_name = document.getElementById("text-tripcode-un").value;
  var tripcode_pass = document.getElementById("text-tripcode-pw").value;
  var current_board = document.getElementById("current-board").innerHTML;
  //create the json object and send to server
  var obj = {Title: post_title, Body: post_body, TripcodeUsername: tripcode_name, TripcodePassword: tripcode_pass, Board: current_board};
  sendCreatePost(obj);
}

//for sending a create post ajax request to server
function sendCreatePost(obj) {
  var jsn = JSON.stringify(obj);
  //console.log(jsn);
  var xmlhttp = new XMLHttpRequest();

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
      if (xmlhttp.status == 200) {
        //console.log(xmlhttp.responseText);
        createPostRequestReplied(xmlhttp.responseText);
      } else {
        createPostStatusMessage(1, "error: can't communicate with server");
      }
    }
  };

  xmlhttp.open("POST", "/ajax/createpost/", true);
  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(jsn);
}

//when the create post request get its reply, show message
function createPostRequestReplied(r) {
  var res = JSON.parse(r);
  createPostStatusMessage(res.Code, res.Message)
}

function createPostStatusMessage(code, msg) {
  //get the status message div
  var status_div = document.getElementById("status-message");
  //change color base on code
  if (code === 0) {
    status_div.style.color = "#25f54b";
  } else if (code === 1) {
    status_div.style.color = "#fc2323";
  } else {
    status_div.style.color = "#a1ecff";
  }
  status_div.innerHTML = msg;
}

//these are for sending a reply

//when the 'add reply' div is clicked
function toggleReplySection(event) {
  event.target.parentElement.getElementsByClassName("add-reply-element")[0].style.display = "block";
}

function replyButtonPressed(event) {
  var parent_element = event.target.parentElement.parentElement;
  var reply_body = parent_element.getElementsByClassName("reply-box")[0].value;
  var tripcode_name = parent_element.getElementsByClassName("reply-tripcode-un")[0].value;
  var tripcode_pass = parent_element.getElementsByClassName("reply-tripcode-pw")[0].value;
  var parent_thread = parent_element.getAttribute("pid");
  //create the json object and send to server
  var obj = {Body: reply_body, TripcodeUsername: tripcode_name, TripcodePassword: tripcode_pass, Parent: parent_thread};
  sendCreateReply(obj);
}

//for sending a create reply ajax request to server
function sendCreateReply(obj) {
  var jsn = JSON.stringify(obj);
  //console.log(jsn);
  var xmlhttp = new XMLHttpRequest();

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
      if (xmlhttp.status == 200) {
        //console.log(xmlhttp.responseText);
        createReplyRequestReplied(xmlhttp.responseText);
      } else {
        createPostStatusMessage(1, "error: can't communicate with server");
      }
    }
  };

  xmlhttp.open("POST", "/ajax/createreply/", true);
  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(jsn);
}

//when the create reply request get its reply, show message
function createReplyRequestReplied(r) {
  var res = JSON.parse(r);
  createReplyStatusMessage(res.PostId, res.Code, res.Message)
}

function createReplyStatusMessage(pid, code, message) {
  //for each post element
  var post_elements = document.getElementsByClassName("post-entry");
  for (var i=0; i<post_elements.length; i++) {
    //if it's the right post
    if (post_elements[i].getAttribute("pid") === pid) {
      var status_div = post_elements[i].getElementsByClassName("reply-status-message")[0];

      //change color base on code
      if (code === 0) {
        status_div.style.color = "#25f54b";
      } else if (code === 1) {
        status_div.style.color = "#fc2323";
      } else {
        status_div.style.color = "#a1ecff";
      }
      status_div.innerHTML = message;
    }
  }
}

//function for when the 'show replies' is clicked
function show_replies_clicked(event) {
  var parent_element = event.target.parentElement;
  var post_id = parent_element.getAttribute("pid");
  parent_element.getElementsByClassName("post-replies")[0].innerHTML = "Loading...";
  send_show_replies({PostId: post_id});
}

function send_show_replies(obj) {
  var jsn = JSON.stringify(obj);
  //console.log(jsn);
  var xmlhttp = new XMLHttpRequest();

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
      if (xmlhttp.status == 200) {
        //console.log(xmlhttp.responseText);
        showRepliesRequestReplied(xmlhttp.responseText);
      } else {
        //createPostStatusMessage(1, "error: can't communicate with server");
        console.log("show reply error: " + xmlhttp.status);
      }
    }
  };

  xmlhttp.open("POST", "/ajax/showreplies/", true);
  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(jsn);
}

function showRepliesRequestReplied(data) {
  var res = JSON.parse(data);
  var post_elems = document.getElementsByClassName("post-entry");
  var target = false;
  //run through all the posts to match the id
  for (var i=0; i<post_elems.length; i++) {
    var cur = post_elems[i];
    if (cur.getAttribute("pid") === res.PostId) {
      target = cur;
    }
  }
  //if post don't exist, return
  if (target === false) {
    return;
  }
  //now add the replies or show error message
  var replies_div = target.getElementsByClassName("post-replies")[0];
  //clear any loading message
  replies_div.innerHTML = "";
  if (res.Status === 1) {
    replies_div.innerHTML = res.Message;
    return;
  } else { //if all good
    //get the template element
    var template = document.getElementById("element-templates").getElementsByClassName("post-reply")[0];
    //all the replies are in this array
    var replies_array = res.Replies;
    console.log(replies_array);
    for (var j=0; j<replies_array.length; j++) {
      var clon = template.cloneNode(true);
      //insert text into elements
      clon.getElementsByClassName("reply-author")[0].innerHTML = replies_array[j].Author;
      clon.getElementsByClassName("reply-body")[0].innerHTML = replies_array[j].Body.replaceAll("\n", "<br>");
      clon.getElementsByClassName("reply-time")[0].innerHTML = replies_array[j].TimePosted;
      //append to the div
      replies_div.appendChild(clon);
    }
    //change the button to 'reload' instead of 'view'
    target.getElementsByClassName("show-replies")[0].innerHTML = "> reload all replies"
  }
}

function jump_butn_pressed() {
  //the destination
  var dest = document.getElementById("jump-page").value;
  var curBoard = document.getElementById("current-board").innerHTML;
  //redirect to the page of the board
  window.location.href = "/nocache/board_template/" + curBoard + "/" + dest + "/";
}

function boardjs_setup() {
  document.getElementById("butn-post").addEventListener("click", postButtonPressed);
  document.getElementById("jump-butn").addEventListener("click", jump_butn_pressed);

  var add_reply_toggles = document.getElementsByClassName("add-reply");
  for (var i=0; i<add_reply_toggles.length; i++) {
    add_reply_toggles[i].addEventListener("click", toggleReplySection);
  }

  var add_reply_buttons = document.getElementsByClassName("butn-reply");
  for (var i=0; i<add_reply_buttons.length; i++) {
    add_reply_buttons[i].addEventListener("click", replyButtonPressed);
  }

  var show_replies_links = document.getElementsByClassName("show-replies");
  for (var i=0; i<show_replies_links.length; i++) {
    show_replies_links[i].addEventListener("click", show_replies_clicked);
  }

  //replace all the linebreak with proper html <br> s
  var postbodies = document.getElementsByClassName("post-body");
  for (var i=0; i<postbodies.length; i++) {
    postbodies[i].innerHTML = postbodies[i].innerHTML.replaceAll("\n", "<br>");
  }
}

document.addEventListener("DOMContentLoaded", boardjs_setup);
