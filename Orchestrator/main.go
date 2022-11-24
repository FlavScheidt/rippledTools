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
    //      Set paths
    // -----------------------------------------
	PATH:="/home/xrpl/rippledTools/"
	RIPPLED_PATH:="./sntrippled/my_build/"
	RIPPLED_CONFIG:="/root/config/rippled.cfg"
	RIPPLED_QUORUM:="15"

	GOSSIPSUB_PATH:="/root/gossipGoSnt/"
	GOPATH:="/usr/local/go/bin/"
	GOSSIPSUB_PARAMETERS:=PATH+"/Orchestrator/parameters.csv"

	NODES_CONFIG:=PATH+"/ConfigCluster/ClusterConfigSmall.csv"

	experiment:="unl"

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
    //		Parameters for GossipSub
    // -----------------------------------------
    //Read nodes name from config file
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

	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := kh.New("/root/.ssh/known_hosts")
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

	//Dial
	var client [len(hosts)]sshClient
	err, client = DialParallel(hosts, config)

	for _, param := range paramsList {
		// -----------------------------------------
	    //		Start rippled
	    // -----------------------------------------
	    log.Println("Starting Rippled")
	    cmd := RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" &"
	    go runParallel(hosts, config, timeout, client)
	    time.Sleep(10 * time.Second)

		// -----------------------------------------
	 	//    		Start gossipsub
	 	//    -----------------------------------------
	 	log.Println("Starting GossipSub")
	    cmd = "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment+" -d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout
	    go runParallel(cmd, hosts, config, timeout, client)
	}

	time.Sleep(60 * time.Second)
	select {}
 }