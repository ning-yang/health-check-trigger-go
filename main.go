package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var node1 = flag.String("node1", "104.198.66.26:1819", "remote ip address + port")
var node2 = flag.String("node2", "104.197.41.219:1819", "remote ip address + port")

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	c := make(chan int, 2)
	var node1HealthState string
	var node2HealthState string
	for { // break the loop if text == "q"
		fmt.Print("Enter node1 health state (1 or 0):")
		scanner.Scan()
		node1HealthState = scanner.Text()
		fmt.Print("Enter node2 health state (1 or 0):")
		scanner.Scan()
		node2HealthState = scanner.Text()

		go sendReplyToRemoteNode(node1, node1HealthState, c)
		go sendReplyToRemoteNode(node2, node2HealthState, c)
		<-c
		<-c
	}
}

func sendReplyToRemoteNode(addr *string, data string, c chan int) {
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		": sending data to",
		*addr,
		"-",
		data)

	conn.Write([]byte(data))

	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		": get reply from",
		*addr,
		"-",
		reply)

	err = conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	c <- 1
}
