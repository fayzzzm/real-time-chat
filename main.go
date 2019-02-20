package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var index int

func main() {
	netListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("There is something error -> %s", err.Error())
	}

	users := make(map[net.Conn]string)    /*Users in our connection*/
	newUsers := make(chan net.Conn)       /*New uses for adding to our connection */
	connectedUsers := make(chan net.Conn) /*Connected users */
	disconnection := make(chan net.Conn)  /*Disconnected users*/
	messages := make(chan string)         /*Messages from users*/

	defer netListener.Close()

	go handleConnection(newUsers, netListener)

	for {
		select {
		// newUser Handlers
		case user := <-newUsers:
			{
				go func(account net.Conn) {
					reader := bufio.NewReader(account)
					io.WriteString(account, "Enter your username\n")
					userName, _ := reader.ReadString('\n')
					userName = strings.TrimSpace(userName)
					users[account] = userName
					connectedUsers <- account
				}(user)
			}
		// Created user message handler
		case user := <-connectedUsers:
			{
				go handleMessage(user, disconnection, messages, users[user])
			}
			// Message writer
		case message := <-messages:
			{
				// Akai Muhammadali users[key] key is Conn
				// Baroi hami naviwtestam
				for user := range users {
					io.WriteString(user, message)
				}
			}
			// Delete connection
		case exit := <-disconnection:
			{
				log.Printf("The Client %v, has gone", users[exit])
				delete(users, exit)
			}
		}

	}
}

func handleConnection(newUsers chan net.Conn, netListener net.Listener) {
	// Forever accept connections
	for {
		newUser, err := netListener.Accept()
		if err != nil {
			log.Printf("There is some error -> %s", err.Error())
		}
		// Through into NewUser
		newUsers <- newUser
	}
}

func handleMessage(connectUsers net.Conn, disconnect chan net.Conn, messages chan string, userName string) {
	read := bufio.NewReader(connectUsers)
	for {
		// Reading string intil denominator
		ms, err := read.ReadString('\n')
		if err != nil {
			log.Printf("Your text is too hard for reading, broo")
			break
		}
		fmt.Printf("Client %s: %s", userName, ms)
		messages <- fmt.Sprintf("Client %s: %s", userName, ms)
	}
	// if user discconnects
	disconnect <- connectUsers
}
