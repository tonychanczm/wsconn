package main

import (
	"log"

	"github.com/gorilla/websocket"

	"wsconn"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:9029/echo", nil)
	if err != nil {
		log.Fatalln(err)
	}
	conn := wsconn.NewWSConn(c)
	_, err = conn.Write([]byte("111a7s8dasudhiasduhasdjnasdsdj111111"))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("111a7s8dasudhiasduhasdjnasdsdj111111"))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("111a7s8dasudhiasduhasdjnasdsdj111111"))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("111a7s8dasudhiasduhasdjnasdsdj111111"))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte("111a7s8dasudhiasduhasdjnasdsdj111111"))
	if err != nil {
		log.Fatalln(err)
	}


	err = conn.Close()
	log.Println(err)
	log.Println("done")

}
