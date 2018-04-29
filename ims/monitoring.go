package ims

import "net/http"
import "encoding/json"
import "os"
import "runtime"
import "runtime/pprof"
import log "github.com/golang/glog"



type ServerSummary struct {
	nConnections      	int64
	nClients          	int64
	inMessageCount  	int64
	outMessageCount 	int64
}


func NewServerSummary() *ServerSummary {
	s := new(ServerSummary)
	return s
}


func Summary(rw http.ResponseWriter, req *http.Request) {
	res, err := json.Marshal(ServerInfo())
	if err != nil {
		log.Info("json marshal:", err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	_, err = rw.Write(res)
	if err != nil {
		log.Info("write err:", err)
	}
	return
}

func ServerInfo()map[string]interface{}  {
	obj := make(map[string]interface{})
	obj["goroutine_count"] = runtime.NumGoroutine()
	obj["connection_count"] = serverSummary.nConnections
	obj["client_count"] = serverSummary.nClients
	obj["in_message_count"] = serverSummary.inMessageCount
	obj["out_message_count"] = serverSummary.outMessageCount
	return obj
}

func Stack(rw http.ResponseWriter, req *http.Request) {
	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	rw.WriteHeader(200)
}

func WriteHttpError(status int, err string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	meta := make(map[string]interface{})
	meta["code"] = status
	meta["message"] = err
	obj["meta"] = meta
	b, _ := json.Marshal(obj)
	w.WriteHeader(status)
	w.Write(b)
}

func WriteHttpObj(data map[string]interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	obj := make(map[string]interface{})
	obj["data"] = data
	b, _ := json.Marshal(obj)
	w.Write(b)
}



