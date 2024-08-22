package RPB

import (
	"log"
	"net/http"
	"time"
)

var InitialTimeStamps = [10]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type ReverseProxyBalancer struct {
	CallBuffer    map[string][]int64
	NextBufferIdx map[string]int
	coolDownIP    map[string]int64
}

func (lb *ReverseProxyBalancer) ProcessRequest(req *http.Request) bool {
	go lb.ProcessTelemetry(req)

	return lb.InsertCall(req.RemoteAddr)
}

func (lb *ReverseProxyBalancer) InsertCall(ip string) bool {
	_, exist := lb.CallBuffer[ip]

	if !exist {
		lb.CallBuffer[ip] = InitialTimeStamps[0:10]
		lb.NextBufferIdx[ip] = 0
	}
	nextIdx := lb.NextBufferIdx[ip]

	lastTime := lb.CallBuffer[ip][nextIdx]
	timeNow := time.Now().UnixMilli()

	if lb.CallBuffer[ip][nextIdx] == 0 || timeNow-lastTime > 100 {
		lb.CallBuffer[ip][nextIdx] = time.Now().UnixMilli()
		lb.NextBufferIdx[ip]++
		return true
	}

	lb.coolDownIP[ip] = timeNow
	return false
}

func (lb *ReverseProxyBalancer) ProcessTelemetry(req *http.Request) {
	log.Println(req.Header)
	log.Println(req.Method)
	log.Println(req.URL)
	log.Println(req.Host)
	log.Println(req.RemoteAddr)
}
