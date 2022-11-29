package main

import (
    "bytes"
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
    
    // results := make(chan string, 10)

    //Dial
    client, err := ssh.Dial("tcp", hostname+":22", config)
        if err != nil {
            log.Fatalf("unable to connect: %v", err)
        }
    // defer client.Close()

    start := RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" &"
    // stop  := RIPPLED_PATH+"rippled stop"
    // status := RIPPLED_PATH+"rippled server_info"

    //Rippled Start
    log.Println(hostname+" Starting Rippled")
    go runUnblockingCmd(start, hostname, client, duration)//, results)
    time.Sleep(30 * time.Second)

    // Check rippled status
    // log.Println(hostname+" status Rippled")
    // go runAtomicCmd(status, hostname, client, duration)//, results)

    //GossipSubStart
    gossipsub := "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment + "&"//+"-d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout
    // go runAndStopCmd(gossipsub, hostname, client, duration)//, results)


    // Check rippled status
    log.Println(hostname+" GossipSub")
    go runAtomicCmd(gossipsub, hostname, client, duration)//, results)
}

func executeCmd(cmd string, ss *ssh.Session) {//} string { //, tp string) {

    // var stdoutBuf bytes.Buffer
    // ss.Stdout = &stdoutBuf
        
    err := ss.Run(cmd)
    if err != nil {
        log.Fatal(err)
    }

    // Let's print out the result of command.
    // fmt.Println(stdoutBuf.String())

    // return stdoutBuf.String()
}


//Run cmd and interprets output until a segmentation fault is detected
func runUnblockingCmd(cmd string, hostname string, client *ssh.Client, duration time.Duration) { //, results chan string) {

    // timeout := time.After(duration)
    // results := make(chan string, 10)

    //new session
    log.Println(hostname+": New session")
    ss, err := client.NewSession()
    if err != nil {
        log.Fatal("unable to create SSH session: ", err)
    }

    var stdoutBuf bytes.Buffer
    ss.Stdout = &stdoutBuf

    // go func(hostname string) {
    //     results <- executeCmd(cmd, ss)
    // }(hostname)

    go executeCmd(cmd, ss)

    // fmt.Println(stdoutBuf.String())

    // select {
    //     case res := <-results:
    //         log.Println(hostname+" result")
    //         fmt.Print(res)
    //     case <-timeout:
    //         fmt.Println("Timed out!")
    //         // ss.Close()
    //         return
    // }

}


//Run an atomic ssh command, usually just a status call
//Stdout goes directly to the stdout
func runAtomicCmd(cmd string, hostname string, client *ssh.Client, duration time.Duration) { //, results chan string) {

    // timeout := time.After(duration)
    // results := make(chan string, 10)

    //new session
    log.Println(hostname+": New session")
    sse, err := client.NewSession()
    if err != nil {
        log.Fatal("unable to create SSH session: ", err)
    }

    sse.Stdout = os.Stdout
    // ss.Stderr = os.Stderr

    // go func(hostname string) {
    //     results <- executeCmd(cmd, ss)
    // }(hostname)

    go executeCmd(cmd, sse)

    // select {
    //     case res := <-results:
    //         log.Println(hostname+" result")
    //         fmt.Print(res)
    //     case <-timeout:
    //         fmt.Println("Timed out!")
    //         // ss.Close()
    //         return
    // }

}

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