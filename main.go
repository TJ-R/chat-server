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
				messageType, message, err := c.Reader(ctx)
				if err != nil {
					log.Println("Websocket read error:", err)
					return
				}

				log.Printf("Receved message: %s (Type: %d)", message, messageType)

				// Cancel the existing context
				cancel()
				ctx, cancel = context.WithTimeout(r.Context(), time.Second*45)
				defer cancel()
			}
		}

		for {
			var v interface{}
			err = wsjson.Read(ctx, c, &v)
			if err != nil {
				log.Println(err)
				return
			}

			wsjson.Write(ctx, c, "Fuck You")
			log.Println(v)

		}

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
