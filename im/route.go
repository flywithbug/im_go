package im

import (
	log "github.com/flywithbug/log4go"
	"sync"
)

type Route struct {
	appid        int64
	mutex        sync.Mutex
	clients      map[int32]ClientSet
	roomClients map[int64]ClientSet
}

func NewRoute(appid int64) *Route {
	route := new(Route)
	route.appid = appid
	route.clients = make(map[int32]ClientSet)
	route.roomClients = make(map[int64]ClientSet)
	return route
}

func (route *Route) AddRoomClient(roomId int64, client *Client) {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	set, ok := route.roomClients[roomId]
	if !ok {
		set = NewClientSet()
		route.roomClients[roomId] = set
	}
	set.Add(client)
}

//todo optimise client set clone
func (route *Route) FindRoomClientSet(roomId int64) ClientSet {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.roomClients[roomId]
	if ok {
		return set.Clone()
	} else {
		return nil
	}
}

func (route *Route) RemoveRoomClient(roomId int64, client *Client) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	if set, ok := route.roomClients[roomId]; ok {
		set.Remove(client)
		if set.Count() == 0 {
			delete(route.roomClients, roomId)
		}
		return true
	}
	log.Info("room client non exists")
	return false
}

func (route *Route) AddClient(client *Client) {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	set, ok := route.clients[client.uid]
	if !ok {
		set = NewClientSet()
		route.clients[client.uid] = set
	}
	set.Add(client)
}

func (route *Route) RemoveClient(client *Client) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	//fmt.Printf("RemoveClient :client uid :%d\n",client.uid)
	if set, ok := route.clients[client.uid]; ok {
		set.Remove(client)
		if set.Count() == 0 {
			delete(route.clients, client.uid)
		}
		return true
	}
	log.Info("client non exists")
	return false
}

func (route *Route) FindClientSet(uid int32) ClientSet {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.clients[uid]
	if ok {
		return set.Clone()
	} else {
		return nil
	}
}

func (route *Route) IsOnline(uid int32) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.clients[uid]
	if ok {
		return len(set) > 0
	}
	return false
}
