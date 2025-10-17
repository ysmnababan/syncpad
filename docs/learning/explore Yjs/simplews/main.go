package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrade HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow all connections by default
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Println("Client connected")

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		switch msgType {
		case websocket.TextMessage:
			log.Println("TextMessage:", string(msg))
		case websocket.BinaryMessage:
			// Print first bytes in hex (helpful for debugging)
			preview := 64
			if len(msg) < preview {
				preview = len(msg)
			}
			log.Printf("BinaryMessage: len=%d firstByte=%d hex-prefix=%s\n",
				len(msg), msg[0], hex.EncodeToString(msg[:preview]))

			// message type convention used by y-websocket: first byte = message type (0 sync, 1 awareness)
			if len(msg) > 0 {
				switch msg[0] {
				case 0:
					log.Println("-> Yjs sync/update frame (binary).")
					// If you want to persist or forward it, keep it as []byte and
					// send it to other peers unchanged.
				case 1:
					log.Println("-> Yjs awareness frame (likely contains JSON payload).")
					// Attempt a heuristic decode (see below)
					// tryExtractAwarenessJSON(msg)
				default:
					log.Println("-> Unknown first byte:", msg[0])
				}
			}
		default:
			log.Println("Other messageType:", msgType)
		}
	}
}

func main() {
	http.HandleFunc("/ws/codemirror-demo-2025-10-17", handleWebSocket)

	fmt.Println("WebSocket server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
