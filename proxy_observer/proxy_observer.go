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
	callBuffer    map[string][]int64
	nextBufferIdx map[string]int
	coolDownIP    map[string]int64
	showDetails   bool
}

func NewReverseProxyBalancer() ReverseProxyBalancer {
	revLB := ReverseProxyBalancer{}
	revLB.callBuffer = make(map[string][]int64)
	revLB.coolDownIP = make(map[string]int64)
	revLB.nextBufferIdx = make(map[string]int)
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
	_, exist := lb.callBuffer[ip]

	if !exist {
		lb.callBuffer[ip] = make([]int64, BUFFER_WINDOW)
		lb.nextBufferIdx[ip] = 0
	}
	nextIdx := lb.nextBufferIdx[ip]

	lastTime := lb.callBuffer[ip][nextIdx]
	timeNow := time.Now().UnixMilli()

	if lb.callBuffer[ip][nextIdx] == 0 || timeNow-lastTime > 100 {
		lb.callBuffer[ip][nextIdx] = time.Now().UnixMilli()
		lb.nextBufferIdx[ip] = (lb.nextBufferIdx[ip] + 1) % BUFFER_WINDOW
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

func (lb *ReverseProxyBalancer) InCoolDown(ip string) bool {
	_, exist := lb.coolDownIP[ip]
	return exist
}
