package RPB

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BUFFER_WINDOW = 10

type ReverseProxyBalancer struct {
	CallBuffer    map[string][]int64
	NextBufferIdx map[string]int
	coolDownIP    map[string]int64
}

func NewReverseProxyBalancer() ReverseProxyBalancer {
	revLB := ReverseProxyBalancer{}
	revLB.CallBuffer = make(map[string][]int64)
	revLB.coolDownIP = make(map[string]int64)
	revLB.NextBufferIdx = make(map[string]int)
	return revLB
}

func (lb *ReverseProxyBalancer) ProcessRequest(req *http.Request) bool {
	go lb.ProcessTelemetry(req)

	return lb.InsertCall(strings.Split(req.RemoteAddr, ":")[0])
}

func (lb *ReverseProxyBalancer) InsertCall(ip string) bool {
	_, exist := lb.CallBuffer[ip]

	if !exist {
		lb.CallBuffer[ip] = make([]int64, BUFFER_WINDOW)
		lb.NextBufferIdx[ip] = 0
	}
	nextIdx := lb.NextBufferIdx[ip]

	lastTime := lb.CallBuffer[ip][nextIdx]
	timeNow := time.Now().UnixMilli()

	if lb.CallBuffer[ip][nextIdx] == 0 || timeNow-lastTime > 100 {
		lb.CallBuffer[ip][nextIdx] = time.Now().UnixMilli()
		lb.NextBufferIdx[ip] = (lb.NextBufferIdx[ip] + 1) % BUFFER_WINDOW
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
