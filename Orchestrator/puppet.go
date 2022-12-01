package main

import (
    // "bytes"
    "fmt"
    // "io"
    // "io/ioutil"
    // "os"
    "time"
    "log"


	"golang.org/x/crypto/ssh"
	// kh "golang.org/x/crypto/ssh/knownhosts"
)

func runPuppet(experiment string, config *ssh.ClientConfig, duration time.Duration, param OverlayParams) {

    results := make(chan string, 10)
    timeout := time.After(duration)

    cmd := "cd "+PATH+"/Orchestrator && "+GOPATH+"go run . -type="+experiment+" -machine=puppet -d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout+"\n"
    hostname := PUPPET
    
    go func(hostname string) {
        results <- executeCmd(cmd, hostname, config)
    }(hostname)
    
    select {
        case res := <-results:
            fmt.Print(res)
        case <-timeout:
            log.Println(hostname, ": Timed out!")
            return
    }
}