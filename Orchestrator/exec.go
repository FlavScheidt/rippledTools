package main

import (
    "bytes"
    "fmt"
    // "io"
    // "io/ioutil"
    // "os"
    "time"
    "log"


	"golang.org/x/crypto/ssh"
	// kh "golang.org/x/crypto/ssh/knownhosts"
 )

type HeartbeatParams struct {
    InitialDelay time.Duration
    Interval     time.Duration
}


type OverlayParams struct {
    d            string
    dlo          string
    dhi          string
    dscore       string
    dlazy        string
    dout         string
    gossipFactor string
}


func executeCmd(cmd, hostname string, config *ssh.ClientConfig, tp string) string {//, client *ssh.Client) string {
    
    var stdoutBuf bytes.Buffer

    if tp == "gossip"{
        clientG, err := ssh.Dial("tcp", hostname+":2022", config)
        if err != nil {
            log.Fatalf("%s unable to connect: %v", hostname, err)
        }

        defer clientG.Close()
        ssG, err := clientG.NewSession()
        if err != nil {
            log.Fatal("unable to create SSH session: ", err)
        }
        defer ssG.Close()

        // Creating the buffer which will hold the remotly executed command's output.
        ssG.Stdout = &stdoutBuf
        ssG.Run(cmd)
    } else {
        client, err := ssh.Dial("tcp", hostname+":22", config)
        if err != nil {
            log.Fatalf("unable to connect: %v", err)
        }

        defer client.Close()
        ss, err := client.NewSession()
        if err != nil {
            log.Fatal("unable to create SSH session: ", err)
        }
        defer ss.Close()

        // Creating the buffer which will hold the remotly executed command's output.
        ss.Stdout = &stdoutBuf
        ss.Run(cmd)

    }
    // Let's print out the result of command.
    // fmt.Println(stdoutBuf.String())

    return hostname + ": " + stdoutBuf.String()
}

func runParallel(cmd string, hosts []string, config *ssh.ClientConfig, duration time.Duration, tp string) {

    results := make(chan string, 10)
    timeout := time.After(duration)

    for _, hostname := range hosts {
        go func(hostname string) {
            results <- executeCmd(cmd, hostname, config, tp)
        }(hostname)
    }

    for i := 0; i < len(hosts); i++ {
        select {
        case res := <-results:
            fmt.Print(res)
        case <-timeout:
            fmt.Println("Timed out!")
            return
        }
    }
}
