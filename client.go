package main

import (
    "net"
    "fmt"
    "log"
    "bufio"
    "os"
    "sync"
)

var wg sync.WaitGroup

const (
    PORT = ":9500"
)

func ReadMessages(conn net.Conn) {
    for {
        wg.Add(1)
        msgReader := bufio.NewReader(conn)
        msg, err := msgReader.ReadString(byte('\n'))
        if err != nil {
            log.Fatal("[Error]: Unable to read network packet data from net.Conn")
        }
        fmt.Println("\n")
        log.Println(msg)
        wg.Done()
    }
}

func main() {
    // Get the server IP from stdin
    fmt.Print("Enter the IP of the server you wish to connect to: ")
    hostReader := bufio.NewReader(os.Stdin)
    host, err := hostReader.ReadString(byte('\n'))
    if err != nil {
        log.Fatal("[Error]: Unable to obtain host from stdin")
    }

    // Connect to the server
    conn, err := net.Dial("tcp", host + PORT)
    if err != nil {
        log.Fatal("[Error]: Unable to connect to server")
    }
    fmt.Println("[Info]: Success! Connected to server")
    fmt.Print("Please supply your username: ")

    // Get the username
    usernameReader := bufio.NewReader(os.Stdin)
    username, err := usernameReader.ReadBytes(byte('\n'))
    if err != nil {
        log.Fatal("[Error]: Unable to read username from stdin")
    }
    conn.Write(username)

    // Get a separate thread for reading the messages
    go ReadMessages(conn)

    // Enter the message loop
    for {
        wg.Add(1)
        fmt.Print(">>> ")
        msgReader := bufio.NewReader(os.Stdin)
        msg, err := msgReader.ReadBytes(byte('\n'))
        if err != nil {
            log.Fatal("[Error]: Unable to read the message to send from stdin")
        }
        conn.Write(msg)
        wg.Done()
    }
}
