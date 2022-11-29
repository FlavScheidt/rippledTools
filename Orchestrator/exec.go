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
    
    // results := make(chan string, 10)

    //Dial
    client, err := ssh.Dial("tcp", hostname+":22", config)
        if err != nil {
            log.Fatalf("unable to connect: %v", err)
        }
    // defer client.Close()

    // log.Println(hostname+": New session")
    // ss, err := client.NewSession()
    // if err != nil {
    //     log.Fatal("unable to create SSH session: ", err)
    // }
    // defer ss.Close()

    start := RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" &\n"
    // stop  := RIPPLED_PATH+"rippled stop"
    // status := RIPPLED_PATH+"rippled server_info"
    // gossipsub := "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment + "&\n"//+"-d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout
    gossipsub := "cd "+GOSSIPSUB_PATH+" && go run . -type="+experiment + "&\n"//+"-d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout


    //Rippled Start
    log.Println(hostname+" Starting Rippled")
    go remoteShell(start, gossipsub, 60, client)
    // go runUnblockingCmd(start, hostname, client, duration)//, results)
    // // Check rippled status
    // log.Println(hostname+" GossipSub")
    // go runAtomicCmd(gossipsub, hostname, client, duration)//, results)


    // // StdinPipe for commands
    // stdin, err := sess.StdinPipe()
    // if err != nil {
    //     log.Fatal(err)
    // }



}

// func executeCmd(cmd string, ss *ssh.Session) {//} string { //, tp string) {
        
//     err := ss.Run(cmd)
//     if err != nil {
//         log.Fatal(err)
//     }
// }

func remoteShell(cmd1 string, cmd2 string, sleep int,client *ssh.Client) {

  // Create sesssion
    sess, err := client.NewSession()
    if err != nil {
        log.Fatal("Failed to create session: ", err)
    }
    defer sess.Close()

    sess.Setenv("GOPATH", GOPATH)

    // StdinPipe for commands
    stdinBuf, err := sess.StdinPipe()
    if err != nil {
        log.Fatal(err)
    }

    // Uncomment to store output in variable
    //var b bytes.Buffer
    //sess.Stdout = &b
    //sess.Stderr = &b

    // Enable system stdout
    // Comment these if you uncomment to store in variable
    // sess.Stdout = os.Stdout
    sess.Stderr = os.Stderr

    // Start remote shell
    err = sess.Shell()
    if err != nil {
        log.Fatal(err)
    }

    // send the commands
    // commands := []string{
    //     "pwd",
    //     "whoami",
    //     "echo 'bye'",
    //     "exit",
    // }

    // for _, cmd := range commands {
    //     _, err = fmt.Fprintf(stdin, "%s\n", cmd)
    //     if err != nil {
    //         log.Fatal(err)
    //     }
    // }
    stdinBuf.Write([]byte(cmd1))
    time.Sleep(60 * time.Second)
    stdinBuf.Write([]byte(cmd2))


    // Wait for sess to finish
    err = sess.Wait()
    if err != nil {
        log.Fatal(err)
    }

    // Uncomment to store in variable
    //fmt.Println(b.String())

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