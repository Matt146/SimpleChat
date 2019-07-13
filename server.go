package main

import (
    "net"
    "log"
    "bufio"
)

const (
    PORT = ":9500"
)

type ChatServer struct {
    Users map[net.Conn]string // Map of users from usernames to IP's
}

func (cs *ChatServer) Broadcast(username string, text string) {
    for conn, _ := range cs.Users {
        msg := "--> "
        msg += username + ": " + text
        conn.Write([]byte(msg))
    }
}

func (cs *ChatServer) ConnHandler(conn net.Conn) {
    // Get the username and register to the server
    usernameReader := bufio.NewReader(conn)
    usernameStr, err := usernameReader.ReadString(byte('\n'))
    if err != nil {
        log.Println("[Error]: Unable to read connection to get username")
        return
    }
    usernameStr = usernameStr[:len(usernameStr)-1]
    cs.Users[conn] = usernameStr
    log.Println("[Join]: IP " + conn.RemoteAddr().String() + " joined as '" + usernameStr + "'")

    // Once registered, enter the message loop
    for {
        // Get the message
        msgReader := bufio.NewReader(conn)
        msg, err := msgReader.ReadString(byte('\n'))
        if err != nil {
            log.Println("[Error]: Unabel to read message from '" + usernameStr + "'")
            return
        }
        msgStr := string(msg)

        // Broadcast it
        cs.Broadcast(usernameStr, msgStr)

        // Log the message
        log.Println("[Message]: From '" + usernameStr + "'" + " " + msgStr)
    }
}

func main() {
    cs := &ChatServer{Users: make(map[net.Conn]string)}
    l, err := net.Listen("tcp", PORT)
    if err != nil {
        log.Fatal("[Error]: Unable to listen for connections")
    }
    log.Println("[Info]: Listening for connections!")
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal("[Error]: Unable to accept connection to")
        }
        go cs.ConnHandler(conn)
    }
}
