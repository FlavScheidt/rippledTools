package main

import (
    // "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "time"
    "log"


	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
 )


// -----------------------------------------
//      Set paths
// -----------------------------------------
var PATH="/home/xrpl/rippledTools/"
var RIPPLED_PATH="./sntrippled/my_build/"
var RIPPLED_CONFIG="/root/config/rippled.cfg"
var RIPPLED_QUORUM="15"

var GOSSIPSUB_PATH="/root/gossipGoSnt/"
var GOPATH="/usr/local/go/bin/"
var GOSSIPSUB_PARAMETERS=PATH+"/Orchestrator/parameters.csv"

var NODES_CONFIG=PATH+"/ConfigCluster/ClusterConfig.csv"

var experiment="general"

func readNodesFile(fileName string) ([]string, error) {

	var nodeList []string

	records, err := readData(fileName)
    if err != nil {
        log.Fatal(err)
        return nodeList, err
    }

    for _, record := range records {
    		nodeList = append(nodeList, record[1])
    }

    log.Println("[INFO] Nodes config file read")
    return nodeList, nil
}


func readParamsFile(fileName string) ([]OverlayParams, error) {

	var paramsList []OverlayParams

	records, err := readData(fileName)
    if err != nil {
        log.Fatal(err)
        return paramsList, err
    }

    for _, record := range records {
    	param := OverlayParams{
    		d:            record[0],
	        dlo:          record[1],
	        dhi:          record[2],
	        dscore:       record[3],
	        dlazy:        record[4],
	        dout:         record[5],
	        gossipFactor: record[6]}

    	paramsList = append(paramsList, param)
    }

    log.Println("[INFO] Parameters file read")
    return paramsList, nil
}

func main() {



    // -----------------------------------------
    //      Set log file
    //			Just the go logging feature, nothing special
    // -----------------------------------------
    currentTime := time.Now()
    LOG_FILE := "./log_"+currentTime.Format("01022006_15_04_05")+".out"
    // open log file
    logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        log.Panic(err)
    }
    defer logFile.Close()

    mw := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(mw)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // -----------------------------------------
    //		Nodes
    // -----------------------------------------
    //Read nodes name from config file
    hosts, error := readNodesFile(NODES_CONFIG)
    if err != nil {
        log.Panic(error)
    }

    fmt.Printf("%+v\n", hosts)

    // -----------------------------------------
    // 		Parameters for GossipSub
    // -----------------------------------------
    // Read nodes name from config file
    paramsList, er := readParamsFile(GOSSIPSUB_PARAMETERS)
    if er != nil {
        log.Panic(error)
    }

    fmt.Printf("%+v\n", paramsList)

    // -----------------------------------------
    //		SSH config
    // -----------------------------------------
    user := "root"
    timeout := 4800 * time.Second

	// key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	key, err := ioutil.ReadFile("/home/xrpl/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// hostKeyCallback, err := kh.New("/root/.ssh/known_hosts")
	hostKeyCallback, err := kh.New("/home/xrpl/.ssh/known_hosts")
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	// configG := &ssh.ClientConfig{
	// 	User: user,
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.PublicKeys(signer),
	// 	},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }

	for _, param := range paramsList {

	    for _, hostname := range hosts {

	    	go runParallel(hostname, config, timeout, param)


	    }






		// // -----------------------------------------
	    // //		Start rippled
	    // // -----------------------------------------
	    // log.Println("Starting Rippled")
	    // cmd := RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" &"
	    // go runParallel(cmd, hosts, config, timeout, "rippled")
	    // time.Sleep(10 * time.Second)



		// // -----------------------------------------
	 	// //    		Start gossipsub
	 	// //    -----------------------------------------
	 	// log.Println("Starting GossipSub")
	    // cmd = "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment+" -d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout
	    // go runParallel(cmd, hosts, configG, timeout, "gossip")


	    // time.Sleep(100 * time.Second)


	    // log.Println("Stopping rippled")
	   	// cmd = RIPPLED_PATH+"rippled stop"
	    // go runParallel(cmd, hosts, config, timeout, "rippledStop") //, client)

	}

	time.Sleep(100 * time.Second)
	// select {}
 }