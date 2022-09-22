package ws

import (
	"github.com/lyzco/multiPersonChat/model"
	"sync"
	"time"
)

type ClientManager struct {
	Users      map[string]*Client
	UserLock   sync.RWMutex
	Register   chan *Client
	Unregister chan *Client
	//Clients    map[*Client]bool
	//ClientLock sync.RWMutex
	Rooms    map[int]*Room
	RoomLock sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		//Clients:    map[*Client]bool{},
		Users:      map[string]*Client{},
		Register:   make(chan *Client, 10),
		Unregister: make(chan *Client, 10),
		Rooms:      map[int]*Room{},
	}
}

func (cm *ClientManager) Start() {
	for {
		select {
		case c := <-Manager.Register:

			//Manager.AddClient(c)
			Manager.AddUser(c)

			room := Manager.GetRoomById(c.RoomId)
			if room == nil {
				//room = NewRoom(c.RoomId)
				//Manager.AddRoom(c.RoomId, room)
			}

			room.AddUser(c.Username, c)

			room.Send(model.MsgTypeBot, "welcome "+c.Username, "chatBot", c)

			room.SendRoomInfo(10, c)

		case c := <-Manager.Unregister:

			room := Manager.GetRoomById(c.RoomId)
			if room != nil {
				room.DeleteUser(c.Username)

				if len(room.Users) == 0 {
					Manager.DeleteRoom(room)
				} else {
					room.SendRoomInfo(20, c)
					room.Send(model.MsgTypeBot, c.Username+" exit", "chatBot", c)
				}
			}

			c.Socket.Close()
			Manager.DeleteUser(c)
		}
	}
}

func (cm *ClientManager) AddUser(c *Client) {
	cm.UserLock.Lock()
	defer cm.UserLock.Unlock()
	cm.Users[c.Username] = c
}

func (cm *ClientManager) AddRoom(rooId int, r *Room) {
	cm.RoomLock.Lock()
	defer cm.RoomLock.Unlock()
	cm.Rooms[rooId] = r
}

func (cm *ClientManager) GetRoomById(rooId int) *Room {
	if room, ok := cm.Rooms[rooId]; !ok {
		return nil
	} else {
		return room
	}
}

func (cm *ClientManager) CloseSocket(c *Client) {
	for {
		if len(c.SendChain) == 0 {
			Manager.Unregister <- c
		}
		time.Sleep(1 * time.Second)
	}
}

func (cm *ClientManager) DeleteUser(c *Client) {
	cm.UserLock.Lock()
	defer cm.UserLock.Unlock()
	if _, ok := cm.Users[c.Username]; ok {
		delete(cm.Users, c.Username)
	}
}

func (cm *ClientManager) DeleteRoom(room *Room) {
	cm.RoomLock.Lock()
	defer cm.RoomLock.Unlock()
	if _, ok := cm.Rooms[room.Id]; ok {
		delete(cm.Rooms, room.Id)
	}
}
