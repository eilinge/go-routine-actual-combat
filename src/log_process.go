package main

/*
将复杂的任务拆分, 通过goroutine去并发执行
通过channel做数据通信
struct
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
	"flag"

	"github.com/influxdata/influxdb/client/v2"
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

// Writer :write interface, 
type Writer interface {
	Write(wc chan *Message)
	// Write(wc chan []byte)
}

// ReadFromFile :read file path
type ReadFromFile struct {
	path string
}

// WriteToInfluxDB :store data to influxdb
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
		// fmt.Println(string(v))
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
	// for v := range wc {
	// 	fmt.Println(v)
	// }
	infSli := strings.Split(w.influxDBDsn, "@")
	addr := infSli[0]
	username := infSli[1]
	password := infSli[2]
	database := infSli[3]
	precision := infSli[4]
 
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	
	fmt.Println("wc length:", len(wc))

	for v := range wc {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  database,
			Precision: precision,
		})
		if err != nil {
			log.Fatal(err)
			continue
		}
 
		// Create a point and add to batch
		//Tags:Path,Method,Scheme,Status
		tags := map[string]string{
			"Path": v.Path,
			"Method": v.Method,
			"Scheme": v.Scheme,
			"Status": v.Status,
			}
 
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}
 
		fmt.Println("taps:",tags)
		fmt.Println("fields:",fields)
 
		pt, err := client.NewPoint("log_access", tags, fields, v.Timelocal)
		if err != nil {
			log.Fatal(err)
			continue
		}
		bp.AddPoint(pt)
 
		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
			continue
		}
 
		// Close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
			continue
		}
 
		log.Println("write success")
	}
}

func main() {
	var path, influxDBDsn string
	flag.StringVar(&path, "path", "./access.log", "write file")

	flag.StringVar(&influxDBDsn, "influxDBDsn", "http://192.168.212.133:8086@root@tester@mydb@s", "influxdb data source")

	flag.Parse()
	p := &ReadFromFile{path: path}
	w := &WriteToInfluxDB{influxDBDsn: influxDBDsn}
	lp := &LogProcess{
		rc: make(chan []byte, 300),
		wc: make(chan *Message, 300),
		// wc: make(chan []byte),
		reader: p,
		writer: w,
	}

	go lp.reader.Read(lp.rc)
	for v := 0; v <=2; v++ {
		go lp.Process()
	}
	for v := 0; v <=4; v++ {
		go lp.writer.Write(lp.wc)
	}

	time.Sleep(20 * time.Second)
}