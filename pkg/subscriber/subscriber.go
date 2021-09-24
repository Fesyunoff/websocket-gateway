package subscriber

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/fesyunoff/websocket-gateway/pkg/config"
	t "github.com/fesyunoff/websocket-gateway/pkg/types"
	"github.com/gorilla/websocket"
)

type Subscriber interface {
	Subscribe(w http.ResponseWriter, r *http.Request)
	Get() map[string]*websocket.Conn
}

var _ Subscriber = (*Gateway)(nil)

type Gateway struct {
	Clients  map[string]*websocket.Conn
	mux      sync.Mutex
	upgrader websocket.Upgrader
}

//TODO: add subscribe to selected topics
func (g *Gateway) Subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var resp string
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		conn_id := r.Header.Get("Sec-Websocket-Key")
		req := &t.ClientRequest{}
		err = json.Unmarshal(message, req)
		if err != nil {
			log.Println("unmarshal:", err)
			continue
		}
		if req.Action == "subscribe" {
			resp = "subscribed"
			g.add(conn_id, c)

		}
		if req.Action == "unsubscribe" {
			resp = "unsubscribed"
			g.delite(conn_id)
		}

		log.Printf("recv: %s, %s", conn_id, string(resp))
		err = c.WriteMessage(websocket.TextMessage, []byte(resp))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (g *Gateway) Get() map[string]*websocket.Conn {
	g.mux.Lock()
	defer g.mux.Unlock()
	return g.Clients
}

func (g *Gateway) add(id string, c *websocket.Conn) {
	g.mux.Lock()
	g.Clients[id] = c
	g.mux.Unlock()
}

func (g *Gateway) delite(id string) {
	g.mux.Lock()
	delete(g.Clients, id)
	g.mux.Unlock()
}

func NewGatewayConfig(host string,
	port int) *config.Gateway {
	return &config.Gateway{
		Host: host,
		Port: port,
	}
}
