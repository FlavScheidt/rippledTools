package main

import (
    // "bytes"
    // "fmt"
    // "io"
    // "io/ioutil"
    "os"
    "time"
    "log"


	"golang.org/x/crypto/ssh"
	// kh "golang.org/x/crypto/ssh/knownhosts"
)


func runParallel(hostname string, config *ssh.ClientConfig, duration time.Duration, param OverlayParams) {
    
    //Start rippled with nohup and disconnect
    log.Println(hostname+" Starting rippled")
    start := "nohup " + RIPPLED_PATH+"rippled --silent --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" & \n"
    go remoteShell(start, hostname, config)

    time.Sleep(60 * time.Second)

    //connect to puppet
    log.Println("Connecting to ", PUPPET)
    // startGossipSub("cd "+GOSSIPSUB_PATH+" && nohup go run . --machine=puppet")
    // go remoteShell()
    go runPuppet(experiment, config, timeout, param)


    // stop  := RIPPLED_PATH+"rippled stop"
    // status := RIPPLED_PATH+"rippled server_info"
    // gossipsub := "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment + "&\n"//+"-d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout
    // gossipsub := "cd "+GOSSIPSUB_PATH+" && nohup go run . -type="+experiment + "&\n"//+"-d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout

    // time.Sleep(60 * time.Second)


    // log.Println(hostname+" Starting gossipsub")
    // go remoteShell(gossipsub, hostname, config)
    // defer client.Close()

}

// func executeCmd(cmd string, ss *ssh.Session) string { //, tp string) {

//     var stdoutBuf bytes.Buffer
//     ss.Stdout = &stdoutBuf
//     err := ss.Run(cmd)
//     if err != nil {
//         log.Fatal(err)
//     }

//     return stdoutBuf.String()
// }

// func remoteCmdDaemon(cmd string, hostname string, config *ssh.ClientConfig, duration time.Duration) {

//     timeout := time.After(duration)
//     results := make(chan string, 10)

//     client, err := ssh.Dial("tcp", hostname+":22", config)
//         if err != nil {
//             log.Fatalf("unable to connect: %v", err)
//         }
//     defer client.Close()

//     //new session
//     ss, err := client.NewSession()
//         if err != nil {
//             log.Fatal("unable to create SSH session: ", err)
//         }
//     // defer ss.Close()

//     go func(hostname string) {
//         results <- executeCmd(cmd, ss)
//     }(hostname)

//     time.Sleep(60 * time.Second)
//     ss.Close()

//     // ss.Signal(ssh.SIGINT)

//     select {
//         //case res := <-results:
//     case <- results:
//             // fmt.Print(res)
//             log.Println("Session killed")
//         case <-timeout:
//             fmt.Println("Timed out!")
//             return
//     }

// }

// func runAndStopCmd(cmd string, hostname string, client *ssh.Client, duration time.Duration) { //, results chan string) {

//     timeout := time.After(duration)
//     results := make(chan string, 10)

//     //new session
//     ss, err := client.NewSession()
//         if err != nil {
//             log.Fatal("unable to create SSH session: ", err)
//         }
//     // defer ss.Close()

//     go func(hostname string) {
//         results <- executeCmd(cmd, ss)
//     }(hostname)

//     time.Sleep(100 * time.Second)

//     ss.Signal(ssh.SIGINT)

//     select {
//         case res := <-results:
//             fmt.Print(res)
//         case <-timeout:
//             fmt.Println("Timed out!")
//             ss.Close()
//             return
//     }

// }


 // func runParallel(cmd string, hosts []string, config *ssh.ClientConfig, duration time.Duration, tp string) {

 //    results := make(chan string, 10)
 //    timeout := time.After(duration)

 //    for _, hostname := range hosts {
 //        go func(hostname string) {
 //            results <- executeCmd(cmd, hostname, config, tp)
 //        }(hostname)
 //    }

 //    for i := 0; i < len(hosts); i++ {
 //        select {
 //        case res := <-results:
 //            fmt.Print(res)
 //        case <-timeout:
 //            fmt.Println("Timed out!")
 //            return
 //        }
 //    }