package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"

	"encoding/json"
	"flag"
	"log"
	"net/http"
)

const (
	TypeHandleLine  = 0
	TypeErrNum      = 1
	TpsIntervalTime = 5
)

var TypeMonitorChan = make(chan int, 200)

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

//系统状态监控
type SystemInfo struct {
	HandleLine    int     `json:"handleLine"`   //总处理日志行数
	Tps           float64 `json:"tps"`          //系统吞吐量
	ReadChanLen   int     `json:"readChanLen"`  //read channel 长度
	WriterChanLen int     `json:"writeChanLen"` //write channel 长度
	RunTime       string  `json:"ruanTime"`     //运行总时间
	ErrNum        int     `json:"errNum"`       //错误数
}

/*
$ curl 127.0.0.1:9193/monitor
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   121  100   121    0     0   118k      0 --:--:-- --:--:-- --:--:--  118k{
        "handleLine": 9108,
        "tps": 315.4,
        "readChanLen": 200,
        "writeChanLen": 0,
        "ruanTime": "31.0062396s",

*/
type Monitor struct {
	startTime time.Time
	data      SystemInfo
	tpsSli    []int
	tps       float64
}

func (m *Monitor) start(lp *LogProcess) {

	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.data.ErrNum += 1

			case TypeHandleLine:
				m.data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * TpsIntervalTime)
	go func() {
		for {
			<-ticker.C
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			if len(m.tpsSli) > 2 {
				m.tpsSli = m.tpsSli[1:]
				m.tps = float64(m.tpsSli[1]-m.tpsSli[0]) / TpsIntervalTime
			}
		}
	}()

	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriterChanLen = len(lp.wc)
		m.data.Tps = m.tps

		ret, _ := json.MarshalIndent(m.data, "", "\t")
		io.WriteString(writer, string(ret))
	})

	http.ListenAndServe(":9193", nil)
}

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Writer(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string //读取文件的路径
}

//读取模块
func (r *ReadFromFile) Read(rc chan []byte) {

	//打开文件
	f, err := os.Open(r.path)
	fmt.Println(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file  err :", err.Error()))
	}

	//从文件末尾开始逐行读取文件内容
	f.Seek(0, 2) //2,代表将指正移动到末尾

	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n') //连续读取内容知道需要'\n'结束
		if err == io.EOF {
			time.Sleep(5000 * time.Microsecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes  err :", err.Error()))
		}

		TypeMonitorChan <- TypeHandleLine
		rc <- line[:len(line)-1]
	}

}

type WriteToinfluxDB struct {
	influxDBDsn string //influx data source
}

//写入模块
/**
    1.初始化influxdb client
	2. 从Write Channel中读取监控数据
	3. 构造数据并写入influxdb
*/
func (w *WriteToinfluxDB) Writer(wc chan *Message) {

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

	for v := range wc {
		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  database,
			Precision: precision,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		//Tags:Path,Method,Scheme,Status
		tags := map[string]string{
			"Path":   v.Path,
			"Method": v.Method,
			"Scheme": v.Scheme,
			"Status": v.Status,
		}

		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		fmt.Println("taps:", tags)
		fmt.Println("fields:", fields)

		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// Close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("write success")
	}

}

//解析模块
func (l *LogProcess) Process() {

	/**
	172.0.012 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-"
	"KeepAliveClient" "-" 1.005 1.854

	([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	*/
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v))
		if len(ret) != 14 {
			TypeMonitorChan <- TypeErrNum
			fmt.Println("FindStringSubmatch fail:", string(v))
			fmt.Println(len(ret))
			continue
		}
		message := &Message{}
		//时间: [04/Mar/2018:13:49:52 +0000]
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			fmt.Println("ParseInLocation fail:", err.Error(), ret[4])
		}
		message.TimeLocal = t
		//字符串长度: 2133
		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent
		//"GET /foo?query=t HTTP/1.0"
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			TypeMonitorChan <- TypeErrNum
			fmt.Println("strings.Split fail:", ret[6])
			continue
		}
		message.Method = reqSli[0]
		u, err := url.Parse(reqSli[1])
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			fmt.Println("url parse fail:", err)
			continue
		}
		message.Path = u.Path
		//http
		message.Scheme = ret[5]
		//code: 200
		message.Status = ret[7]
		//1.005
		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		message.UpstreamTime = upstreamTime
		//1.854
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.RequestTime = requestTime
		//fmt.Println(message)
		l.wc <- message
	}
}

/**
分析监控需求:
	某个协议下的某个请求在某个请求方法的 QPS&响应时间&流量
*/
func main() {
	var path, influDsn string
	flag.StringVar(&path, "path", "./access.log", "read file path")
	flag.StringVar(&influDsn, "influxDsn", "http://127.0.0.1:8086@root@tester@mydb@s", "influx data source")
	flag.Parse()
	r := &ReadFromFile{
		path: path,
	}
	w := &WriteToinfluxDB{
		influxDBDsn: influDsn,
	}
	lp := &LogProcess{
		rc:    make(chan []byte, 200),
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	for i := 1; i < 2; i++ {
		go lp.Process()
	}
	for i := 1; i < 4; i++ {
		go lp.write.Writer(lp.wc)
	}
	fmt.Println("begin !!!")
	m := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	m.start(lp)
}
