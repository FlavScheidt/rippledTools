package main

import (
    // "bytes"
    // "fmt"
    // "io"
    // "io/ioutil"
    // "os"
    "time"
    // "log"


    "golang.org/x/crypto/ssh"
    // kh "golang.org/x/crypto/ssh/knownhosts"
)

func rippledMonitor(config *ssh.ClientConfig, duration time.Duration) {

    status := RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" server_info & \n"
    //run status 5 minutes for 15 minutes (or whatever duration specified)
   for start := time.Now(); time.Since(start) < duration; {
        go func(hostname string) {
            results <- executeCmd(status, hostname, config)
        }(hostname)

        select {
            case res := <-results:
                fmt.Println(hostname, res)
            case <-timeout:
                log.Println(hostname, ": Timed out!")
                return
        }
        time.Sleep(300 * time.Second)
   }

}

func inspectSync() {

}