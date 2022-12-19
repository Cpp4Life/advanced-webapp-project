package websocket

type incomingMessage struct {
	roomId string
	data   []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients by room
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan incomingMessage

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan incomingMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room := h.rooms[client.roomId]
			if room == nil {
				// First client in the room, create a new one
				room = make(map[*Client]bool)
				h.rooms[client.roomId] = room
			}
			room[client] = true
		case client := <-h.unregister:
			room := h.rooms[client.roomId]
			if room != nil {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
					if len(room) == 0 {
						// This was last client in the room, delete the room
						delete(h.rooms, client.roomId)
					}
				}
			}
		case incomingMessage := <-h.broadcast:
			room := h.rooms[incomingMessage.roomId]
			if room != nil {
				for client := range room {
					select {
					case client.send <- incomingMessage.data:
					default:
						close(client.send)
						delete(room, client)
					}
				}
				if len(room) == 0 {
					// The room was emptied while broadcasting to the room
					// Delete the room
					delete(h.rooms, incomingMessage.roomId)
				}
			}
		}
	}
}
