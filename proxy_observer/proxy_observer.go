package RPB

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const BUFFER_WINDOW = 10

type ReverseProxyBalancer struct {
	CallBuffer    map[string][]int64
	NextBufferIdx map[string]int
	coolDownIP    map[string]int64
	showDetails   bool
}

func NewReverseProxyBalancer() ReverseProxyBalancer {
	revLB := ReverseProxyBalancer{}
	revLB.CallBuffer = make(map[string][]int64)
	revLB.coolDownIP = make(map[string]int64)
	revLB.NextBufferIdx = make(map[string]int)
	showDetailValue, err := strconv.ParseBool(os.Getenv("SHOW_REQUEST_DETAIL"))
	if err != nil {
		log.Fatal(err)
	}
	revLB.showDetails = showDetailValue
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
	if lb.showDetails {
		log.Println("From " + req.RemoteAddr + " :: " + req.Method + " " + req.URL.Path + " " + req.Proto + " " + strconv.FormatInt(req.ContentLength, 10))
		log.Println(req.Header)
		log.Println()
	} else {
		log.Println("From " + req.RemoteAddr + " :: " + req.Method + " " + req.URL.Path)
	}
}

func (lb *ReverseProxyBalancer) MonitorCoolDownList() {
	for {
		time.Sleep(time.Second)
		curTime := time.Now().UnixMilli()
		for key, val := range lb.coolDownIP {
			if curTime-val > 5*1000 {
				delete(lb.coolDownIP, key)
				log.Println(key + " deleted from the cool down list")
			}
		}
	}
}
