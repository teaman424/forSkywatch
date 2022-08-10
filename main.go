package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var isReady = false
var resultAverage = make(chan float64, 1)
var resultMedian = make(chan int, 1)
var resultMode = make(chan int, 1)

func main() {
	// 建立 udp 服务器
	listen, err := net.Listen("tcp", "localhost:1888")
	if err != nil {
		fmt.Println("listen failed error:%v\n", err)
		return
	}
	defer listen.Close() // close service
	fmt.Println("udp server is ready")
	//go thread socket client
	go client_socket()
	///go thread pipe client
	go client_pipe()
	///go thread shared memory client
	go client_shared_memory()
	time.Sleep(1 * time.Second)

	isInputString := false
	inputData := ""
	fmt.Print("You can type intergers and then click [ENTER] \n")
	for isInputString == false {

		//get user keyboard input
		reader := bufio.NewReader(os.Stdin)
		inputData, _ = reader.ReadString('\n')
		inputData = strings.Replace(inputData, "\n", "", -1)
		//check input type using regexp.
		isMatch, err := regexp.MatchString("[0-9]+(\\s)?", inputData)
		if err != nil {
			panic("wow!!")
		}
		isInputString = isMatch
		if !isInputString {
			fmt.Print("You can only enter integers and separate them by using space ! try again \n")
		}
		//fmt.Println("match:", isMatch)
	}

	fmt.Print("User input : " + inputData)
	res := strings.Split(inputData, " ")
	fmt.Println("\nResult 1: ", res)
	values := make([]int, 0, len(res))
	for _, raw := range res {
		v, err := strconv.Atoi(raw)
		if err != nil {
			log.Print(err)
			continue
		}
		values = append(values, v)
	}

	//get want result
	resultAverage <- average(values)
	resultMedian <- median(values)
	resultMode <- mode(values)
	fmt.Println("\n result median : ", resultMedian)
	fmt.Println("\n result average : ", resultAverage)
	fmt.Println("\n result mode : ", resultMode)
	isReady = true
	//send result to client
	for {
		// waitting clent connect
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		fmt.Println(conn.RemoteAddr().String(), " tcp connect success")
		//using goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {

		_, err := conn.Read(buffer[:])

		if err != nil {
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		//recv := string(buffer[:n])
		//fmt.Println(conn.RemoteAddr().String(), "\n receive data string:", recv)
		averageVlaue := <-resultAverage
		//fmt.Println("resultAverage : ", averageVlaue)
		buffer = []byte(fmt.Sprintf("%v", averageVlaue))
		//fmt.Println("resultAverage : ", fmt.Sprintf("%v", averageVlaue))
		//fmt.Println("buffer : ", buffer)
		//binary.LittleEndian.PutUint64(buffer[:], math.Float64bits(resultAverage))
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Printf("write from conn failed, err:%v\n", err)
			break
		}
		isReady = false

	}

}

func client_socket() {

	// 建立client
	listen, err := net.Dial("tcp", "127.0.0.1:1888")
	if err != nil {
		fmt.Printf("listen udp server error:%v\n", err)
	}
	defer listen.Close()

	fmt.Println("socket client is ready")
	// 接收data
	data := make([]byte, 4096)
	for {

		if isReady {
			sendData := []byte("Get ANS")
			_, err = listen.Write(sendData) // 发送数据
			if err != nil {
				fmt.Println("send data error，err:", err)
				return
			}
			n, err := listen.Read(data[:]) // 接收数据
			if err != nil {
				fmt.Println("get data error:", err)
				return
			}
			//fmt.Println("client  Mean value byte : ", (data[:n]))
			fmt.Println("socket client  Mean value  : ", string(data[:n]))
		}
		time.Sleep(1 * time.Second)
	}
}

func client_pipe() {

	// 建立client
	fmt.Println("pipe client is ready")

	for {

		if isReady {
			medianValue := <-resultMedian
			fmt.Println("pipe client  Median value  : ", medianValue)
		}
		time.Sleep(1 * time.Second)
	}
}

func client_shared_memory() {

	// 建立client
	fmt.Println("shared memory client is ready")

	for {

		if isReady {
			modeValue := <-resultMode
			fmt.Println("shared memory client  Mode value  : ", modeValue)
		}
		time.Sleep(1 * time.Second)
	}
}

//取中位數
func median(data []int) int {
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)

	sort.Ints(dataCopy)

	var median int
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

//get average value
func average(data []int) float64 {
	var averageValue float64
	var total int

	for _, item := range data {
		total = total + item
	}
	averageValue = float64(float64(total) / (float64(len(data))))
	return averageValue
}

//get mode value
func mode(testArray []int) (mode int) {
	//create map ,record every number of count total
	countMap := make(map[int]int)
	for _, value := range testArray {
		countMap[value] += 1
	}

	max_key := 0
	max_value := 0
	for key, value := range countMap {
		if max_value < value {
			max_key = key
			max_value = value
		}
	}
	return max_key
}