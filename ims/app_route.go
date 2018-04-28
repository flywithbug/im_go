package ims

import "sync"



type AppRoute struct {
	mutex sync.Mutex
	apps map[int64]*Route
}

func NewAppRoute()*AppRoute  {
	r := new(AppRoute)
	r.apps = make(map[int64]*Route)
	return r
}

//若已有则返回已有，无则new一个返回
func (app *AppRoute)FindOrAddRoute(appId int64)(*Route)  {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if route ,ok := app.apps[appId];ok {
		return route
	}
	route := NewRoute(appId)
	app.apps[appId] = route
	return route
}
func (app *AppRoute) FindRoute(appid int64) *Route{
	app.mutex.Lock()
	defer app.mutex.Unlock()
	return app.apps[appid]
}

func (app *AppRoute) AddRoute(route *Route) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	app.apps[route.appId] = route
}


//ClientSet
type ClientSet map[*Client]struct{}

func NewClientSet() ClientSet {
	return make(map[*Client]struct{})
}

func (set ClientSet) Add(c *Client) {
	set[c] = struct{}{}
}

func (set ClientSet) IsMember(c *Client) bool {
	if _, ok := set[c]; ok {
		return true
	}
	return false
}

func (set ClientSet) Remove(c *Client) {
	if _, ok := set[c]; !ok {
		return
	}
	delete(set, c)
}

func (set ClientSet) Count() int {
	return len(set)
}

func (set ClientSet) Clone() ClientSet {
	n := make(map[*Client]struct{})
	for k, v := range set {
		n[k] = v
	}
	return n
}



