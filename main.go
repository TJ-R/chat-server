package main

import (
	"time"
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"net/http"
)

type Message struct {
	Msg string `json:"message"`
	Error bool `json:"error"`
}

func main() {
	log.Println("Chat server started...")
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connecting To Websocket")
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:5173"},
		})
		if err != nil {
			log.Println(err)
			return
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*45)
		defer cancel()

		for {
			select {
			case <- ctx.Done():
				log.Println("Websocket connection timed out:", ctx.Err())
				return

			default:
				var msg Message	
				err = wsjson.Read(ctx, c, &msg)
				log.Println(msg)
				if err != nil {
					log.Println("Websocket read error:", err)
					return
				}

				log.Printf("Recieved message %v\n", msg.Msg)

				wsjson.Write(ctx, c, "Fuck You")
				

				// Cancel the existing context
				cancel()
				ctx, cancel = context.WithTimeout(r.Context(), time.Second*45)
				defer cancel()
			}
		}
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
