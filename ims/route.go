package ims

import "sync"




type Route struct {
	appId 	int64
	mutex sync.Mutex
	clients  map[int64]ClientSet
	//room_clients map[int64]ClientSet
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
	if set, ok := route.clients[client.uid]; ok {
		set.Remove(client)
		if set.Count() == 0 {
			delete(route.clients, client.uid)
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
