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
	} else {
		err := connect("tcp", fmt.Sprintf("%s:%d", *ip, *port))
		if err != nil {
			fmt.Printf("Failed to connect to %s:%d\n", *ip, *port)
			os.Exit(2)
		}
	}
}

func connect(protocol string, host string) error {
	fmt.Printf("Connecting to %s..\n", host)

	conn, err := net.Dial(protocol, host)
	if err != nil {
		return err
	}

	defer conn.Close()

	clientLoop(conn)
	return nil
}

func listen(protocol string, host string) error {
	ln, err := net.Listen(protocol, host)
	if err != nil {
		return err
	}

	defer ln.Close()

	fmt.Printf("Listening on %s..\n", host)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("Connection accepted from %s!\n", conn.RemoteAddr().String())
		fmt.Fprintf(conn, "You have successfully connected to %s!\n", conn.LocalAddr().String())
		go serverLoop(conn)
	}
}

func serverLoop(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		command := strings.Split(strings.Trim(msg, "\n"), " ")
		if strings.ToLower(command[0]) == "exit" || strings.ToLower(command[0]) == "quit" {
			fmt.Println("Closing the connection!")
			conn.Close()
			return nil
		}
		// fmt.Println("DEBUG: ", command)
		if command[0] == "" {
			continue
		}
		var cmd *exec.Cmd
		if len(command) > 1 {
			cmd = exec.Command(command[0], command[1:]...)
		} else {
			cmd = exec.Command(command[0])
		}

		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(conn, "%s: command not found\n", command[0])
		} else {
			fmt.Fprintf(conn, out.String())
		}
	}
}

func clientLoop(conn net.Conn) error {
	scanner := bufio.NewScanner(conn)
	textScanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	for {
		if textScanner.Scan() {
			line := textScanner.Text()
			if strings.Trim(line, "\n") == "exit" {
				fmt.Fprintf(conn, "exit\n")
				return nil
			}
			fmt.Fprintf(conn, "%s\n", line)
		}
	}
}
