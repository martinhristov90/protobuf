package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	"github.com/martinhristov90/protobuf/api/v1"
)

func main() {

	//Getting the mode
	var mode string
	var address string

	flag.StringVar(&mode, "mode", "server", "-mode=server | -mode=client")
	flag.StringVar(&address, "address", "", "-address=127.0.0.1 (flag only used in client mode)")

	//Parsing the flags
	flag.Parse()

	//Check if the address is actual
	ip := net.ParseIP(address)

	if ip == nil && mode == "client" {
		fmt.Printf("Cannot parse IP entered in flag address %s\n\n", address)
		flag.PrintDefaults()
	}

	switch mode {

	case "server":
		server()
	case "client":
		client(ip)
	default:
		// printing the usage if mode is not client or server
		flag.PrintDefaults()
	}
}

func server() {

	listenAddr := net.TCPAddr{
		IP: net.IPv4(0,0,0,0),
		Port: 55555,
	}

	fmt.Println("Starting server on localhost port 55555")
	listener, err := net.ListenTCP("tcp",&listenAddr)

	if err != nil {
		fmt.Println("Cannot start TCP listener")
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(("Error accepting connection"))
		}
		connectionHandlerServer(conn)
	}
}

func connectionHandlerServer (conn *net.TCPConn){


	defer func() {
		fmt.Println("Closing connection from :",conn.RemoteAddr())
		conn.Close()
	}()

	//var raw_data []byte
	raw_data := make([]byte,4096)
	//fmt.Println("TYPe %t",raw_data)

	bytesRead, err := conn.Read(raw_data)
	if err != nil {
		fmt.Println("Error reading raw data")
	}

	record := msg_v1.Record{}
	err = record.Unmarshal(raw_data[:bytesRead]) //Unmarshaling only the protobuf messange,not the whole raw_data
	if err != nil {
		fmt.Println("Error unmarshaling data",err)
	}

	fmt.Printf("Your message is: %s, it was sent at : %v, its length is %v\n",record.Msg,record.GetTimeSent().Local(), bytesRead)
}

func client(ip net.IP) {

	remoteAddr := net.TCPAddr{
		IP: ip,
		Port: 55555,
	}

	tcpDialer, err := net.DialTCP("tcp",nil,&remoteAddr)
	if err != nil {
		fmt.Printf("Error dialing remote address %s",remoteAddr.IP)
	}

	defer tcpDialer.Close()

	text := "Hello"

	raw_data := connectionHandlerClient(text)

	fmt.Println("The raw data of your message looks like",raw_data)
	fmt.Println("Sending the raw data over the wire to server")

	// Writing the protobuf raw data to the socket
	tcpDialer.Write(raw_data)

}

func connectionHandlerClient (text string) []byte {

	now := time.Now().Local()

	//Creating a Record type msg
	message := msg_v1.Record {
		Msg: text,
		TimeSent: &now,
	}

	//Marshaling the message
	raw, err := message.Marshal()
	if err != nil{
		log.Fatal("Error during marshaling")
	}

	//Returning raw data to be send over the wire
	return raw
}
