package main

import (
	"time"
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"net/http"
)

func main() {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connecting To Websocket")
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(v)

		c.Close(websocket.StatusNormalClosure, "")
	})

	err := http.ListenAndServe("localhost:8080", fn)
	log.Fatal(err)
}

// ALL OF THE BELOW CODE IS EXAMPLE OF CONNECTING TO A WEBSOCKET AND SENDING A MESSAGE
// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()
//
// 	c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer c.CloseNow()
//
// 	err = wsjson.Write(ctx, c, "hi")
// 	if err != nil {
// 		log.Fatal()
// 	}
//
// 	c.Close(websocket.StatusNormalClosure, "Closed")
// }
