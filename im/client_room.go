package im


type ClientROOM struct {
	*Connection
	roomId    int64
}