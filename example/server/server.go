package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"wsconn"
)

var upgrader = websocket.Upgrader{}
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	conn := wsconn.NewWSConn(c)
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	buf := make([]byte, 10)
	for err == nil {
		_, err = conn.Read(buf)
		log.Println(buf)
	}
	log.Println("quit:", err)
}

func main() {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":9029", nil))
}
