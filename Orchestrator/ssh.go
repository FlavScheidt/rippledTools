package main

import (
    "bytes"
    // "fmt"
    // // "io"
    // // "io/ioutil"
    // // "os"
    "time"
    "log"
    "os"


	"golang.org/x/crypto/ssh"
	// kh "golang.org/x/crypto/ssh/knownhosts"
 )

func executeCmd(cmd string, hostname string, config *ssh.ClientConfig) string {//, client *ssh.Client) string {
    
    var stdoutBuf bytes.Buffer

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
    ss.Setenv("GOPATH", GOPATH)
    
    // Creating the buffer which will hold the remotly executed command's output.
    // ss.Stdout = &stdoutBuf
    ss.Stdout = os.Stdout
    ss.Stderr = os.Stderr
    ss.Run(cmd)

    // Let's print out the result of command.
    // fmt.Println(stdoutBuf.String())

    return hostname + ": " + stdoutBuf.String()
}

func remoteShell(cmd string, hostname string, config *ssh.ClientConfig,) {

    client, err := ssh.Dial("tcp", hostname+":22", config)
        if err != nil {
            log.Fatalf("unable to connect: %v", err)
        }
    defer client.Close()

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
    // var b bytes.Buffer
    // sess.Stdout = &b
    // sess.Stderr = &b

    // Enable system stdout
    // Comment these if you uncomment to store in variable
    // sess.Stdout = os.Stdout
    sess.Stderr = os.Stderr

    // Start remote shell
    err = sess.Shell()
    if err != nil {
        log.Fatal(err)
    }

    log.Println(hostname, ": Running command | ", cmd)
    stdinBuf.Write([]byte(cmd))
    time.Sleep(60 * time.Second)

    disown := "disown -h %1\n"
    log.Println(hostname, ": disown")
    stdinBuf.Write([]byte(disown))

    //Wait for sess to finish
    err = sess.Wait()
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

    // Uncomment to store in variable
    // log.Println(hostname, ": ", b.String())

}

// func runParallel(cmd string, hosts []string, config *ssh.ClientConfig, duration time.Duration, tp string) {

//     results := make(chan string, 10)
//     timeout := time.After(duration)

//     for _, hostname := range hosts {
//         go func(hostname string) {
//             results <- executeCmd(cmd, hostname, config, tp)
//         }(hostname)
//     }

//     for i := 0; i < len(hosts); i++ {
//         select {
//         case res := <-results:
//             fmt.Print(res)
//         case <-timeout:
//             fmt.Println("Timed out!")
//             return
//         }
//     }
// }
