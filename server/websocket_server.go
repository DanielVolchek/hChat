package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
This file is a very basic first version of a websocket server
It is a one to one clone of a video I found on https://tutorialedge.net/golang/go-websocket-tutorial/
Down the line I will adapt it to fit our needs
*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  MSG_LEN,
	WriteBufferSize: MSG_LEN,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client succesfully connected...")
	processWSConn(ws)

}

func setupHTTPRoutes() {
	log.Println("Setting up routes")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	log.Println("Listening on port: " + HTTP_PORT)
	log.Fatal(http.ListenAndServeTLS(":"+HTTP_PORT,"/etc/letsencrypt/live/hchat.org/cert.pem","/etc/letsencrypt/live/hchat.org/privkey.pem", nil))
}

func processWSConn(ws_conn *websocket.Conn) error {
	// 	// func first needs to process message
	log.Println("WS: client connected: " + (*ws_conn).RemoteAddr().String())
	msgbuf, err := ReadWSMessage(ws_conn)
	if err != nil {
		return err
	}
	room_num, username, _, err := ProcessInput(msgbuf)
	if err != nil {
		(*ws_conn).WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return err
	}
	// Validate username and room existence
	p_user, err := NewUser(nil, ws_conn, username, room_num)
	if err != nil {
		(*ws_conn).WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return err
	}
	fmt.Println("Adding user to room")
	AddUserToRoom(p_user, room_num)
	//  seperate user into appropriate room
	return nil

}
