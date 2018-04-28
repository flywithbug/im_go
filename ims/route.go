package ims

import "sync"


type Route struct {
	appId 	int64
	mutex sync.Mutex
	clients  map[int64]ClientSet
	roomClients map[int64]ClientSet //群聊
}

func NewRoute(appId int64)*Route  {
	route := new(Route)
	route.appId = appId
	route.clients = make(map[int64]ClientSet)
	return route
}


func (route *Route) AddClient(client *Client) {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	set, ok := route.clients[client.uId]
	if !ok {
		set = NewClientSet()
		route.clients[client.uId] = set
	}
	set.Add(client)
}

func (route *Route) RemoveClient(client *Client) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	if set, ok := route.clients[client.uId]; ok {
		set.Remove(client)
		if set.Count() == 0 {
			delete(route.clients, client.uId)
		}
		return true
	}
	//log.Info("client non exists")
	return false
}

func (route *Route) FindClientSet(uid int64) ClientSet {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.clients[uid]
	if ok {
		return set.Clone()
	} else {
		return nil
	}
}

func (route *Route) IsOnline(uid int64) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.clients[uid]
	if ok {
		return len(set) > 0
	}
	return false
}



func (route *Route) AddRoomClient(room_id int64, client *Client) {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	set, ok := route.roomClients[room_id];
	if !ok {
		set = NewClientSet()
		route.roomClients[room_id] = set
	}
	set.Add(client)
}

//todo optimise client set clone
func (route *Route) FindRoomClientSet(room_id int64) ClientSet {
	route.mutex.Lock()
	defer route.mutex.Unlock()

	set, ok := route.roomClients[room_id]
	if ok {
		return set.Clone()
	} else {
		return nil
	}
}

func (route *Route) RemoveRoomClient(room_id int64, client *Client) bool {
	route.mutex.Lock()
	defer route.mutex.Unlock()
	if set, ok := route.roomClients[room_id]; ok {
		set.Remove(client)
		if set.Count() == 0 {
			delete(route.roomClients, room_id)
		}
		return true
	}
	//log.Info("room client non exists")
	return false
}
