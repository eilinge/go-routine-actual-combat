package main

/*
将复杂的任务拆分, 通过goroutine去并发执行
通过channel做数据通信
*/
import (
	"time"
	"fmt"
	"strings"
	"os"
	"bufio"
	"io"
	"log"
	"regexp"
	"strconv"
	"net/url"
)

// LogProcess ...
type LogProcess struct {
	rc chan []byte
	wc chan *Message
	reader Reader
	writer Writer
}

type Message struct {
	Timelocal time.Time
	BytesSent int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime float64
}

// Reader :reader interface
type Reader interface {
	Read(rc chan []byte)
}

// Writer :write interface
type Writer interface {
	Write(wc chan *Message)
	// Write(wc chan []byte)
}

// ReadFromFile :read file
type ReadFromFile struct {
	path string
}

// WriteToInfluxDB :store influxdb
type WriteToInfluxDB struct {
	influxDBDsn string
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// rc <- "message test"
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("Open file error:%s", err.Error()))
	}

	// Read file contents line by line from end of file
	f.Seek(0, 2)
	// get a new *Reader, can call more methods
	rd := bufio.NewReader(f)

	for {
		// Line break: not "/n"
		line, err := rd.ReadBytes('\n')
		// end of file
		if err == io.EOF {
			// If the pointer reaches the end, continue to wait
			time.Sleep(5 * time.Millisecond)
			continue
		}else if  err != nil {
		panic(fmt.Sprintf("Open file error:%s", err.Error()))
		}

		// line = '172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854'
		rc <- line
		
	}
}

// Process :method log process
/*
1.Read log data for each row from the Read Channel
2.Regular extraction of required monitoring data(date protocol method...)
*/
func (l *LogProcess) Process() {
	/*
	172.0.0.12 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854

	([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^""]+)
	\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	*/

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^""]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	// get value from chan: l.rc continuous
	// r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^""]+)\"\s+.*`)
	loc , _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		fmt.Println(string(v))
		ret := r.FindStringSubmatch(string(v))
		// fmt.Println("ret", ret[0])
		// byte conver to string
		// l.wc <- strings.ToUpper(string(v))
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}

		message.Timelocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("reqSli fail:", err.Error(), reqSli)
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("reqSli fail:", err.Error(), reqSli[1])
			continue
		}
		message.Path = u.Path

		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		message.UpstreamTime = upstreamTime
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.RequestTime = requestTime

		l.wc <- message
	}
}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
	for v := range wc {
		fmt.Println(v)
	}
}

func main() {
	p := &ReadFromFile{"./access.log"}
	w := &WriteToInfluxDB{influxDBDsn: "usernam@passwd"}
	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan *Message),
		// wc: make(chan []byte),
		reader: p,
		writer: w,
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)

	time.Sleep(20 * time.Second)
}