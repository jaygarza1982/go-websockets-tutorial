package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//We will use this as a set
var socketsSet = map[*websocket.Conn]bool{}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	conn.SetCloseHandler(func(code int, text string) error {
		textString := fmt.Sprintf("Closing client code %d because %s", code, text)

		delete(socketsSet, conn)

		return errors.New(textString)
	})

	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))

		//Write to other sockets
		for socket := range socketsSet {
			//Don't write to self
			if socket != conn {
				socket.WriteMessage(1, p)
				fmt.Printf("Writing to %d sockets\n", len(socketsSet))
			}
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)

	// ws.

	// Add our current websocket to the sockets array
	clientId := fmt.Sprintf("%d", rand.Int())
	socketsSet[ws] = true
	// sockets = append(sockets, ws)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Client %s connected\n", clientId)

	//Listen on this socket forever until it closes
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
