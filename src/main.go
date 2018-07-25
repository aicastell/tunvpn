package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"runcmd"
	"tuntap"
)

const (
	MTU     string = "1300"
	BUFSIZE int    = 1500
)

type Request struct {
	Action string
}

func main() {

	iface := flag.String("i", "tun0", "tunnel interface")
	local_ip := flag.String("l", "10.0.0.1/24", "local ip and netmask")
	remote_ip := flag.String("r", "8.8.8.8", "remote ip")
	port := flag.String("p", "1234", "application port (local and remote)")
	flag.Parse()

	fmt.Println("Configuration in use:")
	fmt.Println("\tTunnel interface: ", *iface)
	fmt.Println("\tLocal IP/netmask: ", *local_ip)
	fmt.Println("\tRemote IP: ", *remote_ip)
	fmt.Println("\tApplication port: ", *port)

	tun, err := tuntap.Tun(*iface)
	if err != nil {
		fmt.Println("error: tun:", err)
		return
	}
	defer tun.Close()

	runcmd.Cmd("/sbin/ip", "link", "set", "dev", *iface, "mtu", MTU)
	runcmd.Cmd("/sbin/ip", "addr", "add", *local_ip, "dev", *iface)
	runcmd.Cmd("/sbin/ip", "link", "set", "dev", *iface, "up")

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", *remote_ip, *port))
	if nil != err {
		log.Fatalln("Unable to prepare remote socket:", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", *port))
	if nil != err {
		log.Fatalln("Unable to prepare local socket:", err)
	}

	listenConn, err := net.ListenUDP("udp", localAddr)
	if nil != err {
		log.Fatalln("Unable to listen on UDP socket:", err)
	}
	defer listenConn.Close()

	// Goroutine to process incoming packets
	go func() {
		inBytes := make([]byte, BUFSIZE)
		for {
			// recv from socket
			n, addr, err := listenConn.ReadFromUDP(inBytes)
			if err != nil || n == 0 {
				fmt.Println("Error: ", err)
				continue
			}

			// decode IPv4 header (debug)
			header, _ := ipv4.ParseHeader(inBytes[:n])
			fmt.Printf("Received %d bytes from %v: %+v\n", n, addr, header)
			if err != nil || n == 0 {
				fmt.Println("Error: ", err)
				continue
			}

			// write to TUN interface
			wb, err := tun.Write(inBytes[:n])
			if err != nil || wb == 0 {
				fmt.Println("Error writting to tunnel: ", err)
				continue
			}
		}
	}()

	// Main loop to process outgoing packets
	outBytes := make([]byte, BUFSIZE)
	for {
		// read from TUN interface
		rb, err := tun.Read(outBytes)
		if err != nil || rb == 0 {
			fmt.Println("Error reading from tunnel: ", err)
			break
		}

		// decode IPv4 header
		header, _ := ipv4.ParseHeader(outBytes[:rb])
		fmt.Printf("Sending %d bytes to remote: %+v\n", rb, header)

		// send to socket
		n, err := listenConn.WriteToUDP(outBytes[:rb], remoteAddr)
		if err != nil || n == 0 {
			fmt.Println("Error: ", err)
			continue
		}
	}
}
