package main

import (
	"fmt"
	"log"
)

// Global scope room map
var ROOM_MAP map[int]*Room

// Room struct
// Contains array of connections and helper information
type Room struct {
	room_num   int     // Description of room for server
	conn_users []*User // Dynamic array of User pointers

}

// New Room takes in room number, room capacity, room desc, and pointer to empty slice of users (conn_users := make([]User, 0))
func NewRoom(room_num int, conn_users []*User) *Room {
	return &(Room{room_num: room_num, conn_users: conn_users})
}

// TODO
// When adding user to room create new goroutine to read for new incoming messages and send response back if new message is found
func AddUserToRoom(p_user *User, room_num int) (err error){
	if _, ok := ROOM_MAP[room_num]; !ok {
		ROOM_MAP[room_num] = NewRoom(room_num, make([]*User, 0))
		log.Println("created a new room#" + fmt.Sprint(room_num))
	}
	ROOM_MAP[room_num].conn_users = append(ROOM_MAP[room_num].conn_users, p_user)
	log.Println("len", fmt.Sprint(len((*ROOM_MAP[room_num]).conn_users)))
	log.Println("Added new user to room #" + fmt.Sprint(room_num))

	go func(){
		messages, err := readDB(room_num)
		if err != nil{
			return
		}
		for _, msg := range messages{
			combined_msg := fmt.Sprintf("#%d_%s:%s", room_num, msg.username, msg.message)
			RespondWithString(p_user, combined_msg)
		}
		RespondWithString(p_user, "END_OF_DB")
	}()
	go ReadConnOnLoop(p_user)
	return nil
}

// Search through array of users in room pointer
func RemoveUserFromRoom(p_user *User, room_num int) {
	user_index := -1
	p_room := ROOM_MAP[room_num]
	conn_users := (*p_room).conn_users
	for index, element := range conn_users {
		if *element == *p_user {
			user_index = index
			break
		}
	}
	if user_index == -1 {
		return
	}
	(*p_room).conn_users = append(conn_users[:user_index], conn_users[user_index+1:]...)
	log.Println("removed user from #" + fmt.Sprint(room_num))
	fmt.Println("new len is", len((*p_room).conn_users))
	hasUsers := CheckRoomHasConnection(p_room)
	if !hasUsers {
		DelRoom((*p_room).room_num)
	}
}

// DelRoom() function deletes room from globally scoped map of rooms
// Uses built in delete function to do so
func DelRoom(room_num int) {
	delete(ROOM_MAP, room_num)
}

// TODO
// Rooms should be created and deleted if they contain users to prevent unneccessary bloat
// Room number and room description for all rooms should be stored on git and updated periodically, then searched through when creating new rooms

// CheckRoomHasConnection() function checks the passed in room to see if the user slice contains any users.

func CheckRoomHasConnection(p_room *Room) bool {
	return len((*p_room).conn_users) >= 1
}

func DistributeMessageToRoom(p_room *Room, message string) {
	// loops through range of dereferenced room pointers user list
	n := 0
	usersToRemove := make([]*User, 0)
	for _, user := range (*p_room).conn_users {
		fmt.Println("user: ", n, (*user).username)
		n++
		// attempts to send each user message
		err := RespondWithString(user, message)
		// if response fails passes user into remove user function
		if err != nil {
			usersToRemove = append(usersToRemove, user)
		}
	}
	for _, user := range usersToRemove {
		RemoveUserFromRoom(user, (*p_room).room_num)
	}
}
