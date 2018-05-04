package im

import (
	"net/http"
	log "github.com/flywithbug/log4go"
	"fmt"
)

//TODO
//1. apns
// 2.外部接口，输入参数，发送应用内消息或者apns推送给对应的人

func startHttpServer(port  int)  {
	http.HandleFunc("/summary", Summary)
	addr := fmt.Sprintf("localhost:%d", port)
	handler := loggingHandler{http.DefaultServeMux}
	go func() {
		err := http.ListenAndServe(addr, handler)
		if err != nil {
			log.Fatal("http server err:", err)
		}
		panic(err)
	}()

}

type loggingHandler struct {
	handler http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("http request:%s %s %s", r.RemoteAddr, r.Method, r.URL)
	h.handler.ServeHTTP(w, r)
}