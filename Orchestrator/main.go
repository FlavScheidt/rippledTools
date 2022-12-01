package main

import (
    // "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "time"
    "log"
    "flag"
    "strings"


	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
 )


// -----------------------------------------
//      Set paths
// -----------------------------------------
var PATH="/root/rippledTools/"
var RIPPLED_PATH="/root/sntrippled/my_build/"
var RIPPLED_CONFIG="/root/config/rippled.cfg"
var RIPPLED_QUORUM="15"

var GOSSIPSUB_PATH="/root/gossipGoSnt/"
var GOPATH="/usr/local/go/bin/"
var GOSSIPSUB_PARAMETERS=PATH+"/Orchestrator/parameters.csv"

var NODES_CONFIG=PATH+"/ConfigCluster/ClusterConfig.csv"

var PUPPET="liberty"

// var experiment="unl"

func main() {

	//------------------------------------------
	//	Proccess flags
	//------------------------------------------
	machineFlag := flag.String("machine", "master", "Is this machine a master or a puppet? Deafult is master")
  	experimentType := flag.String("type", "unl", "Type of experiment. Default is unl")

    d := flag.String("d", "8", "")
    dlo := flag.String("dlo", "6", "")
    dhi := flag.String("dhi", "12", "")
    dscore := flag.String("dscore", "4", "")
    dlazy := flag.String("dlazy", "8", "")
    dout := flag.String("dout", "2", "")
    gossipFactor := flag.String("gossipFactor", "0.25", "")

    // InitialDelay := flag.Duration("InitialDelay", 100 * time.Millisecond, "")
    // Interval := flag.Duration("Interval", 1 * time.Second, "")

    flag.Parse()

	machine := strings.ToLower(*machineFlag)
	experiment := strings.ToLower(*experimentType)

    // -----------------------------------------
    //      Set log file
    //			Just the go logging feature, nothing special
    // -----------------------------------------
    currentTime := time.Now()
    LOG_FILE := "./log_"+currentTime.Format("01022006_15_04_05")+"_"+experiment+".out"
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
    //		SSH config
    // -----------------------------------------
    user := "root"
    timeout := 4800 * time.Second

	// key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// hostKeyCallback, err := kh.New("/root/.ssh/known_hosts")
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

	if machine == "master" {
		// -----------------------------------------
	    // 		Parameters for GossipSub
	    // -----------------------------------------
	    // Read nodes name from config file
	    paramsList, er := readParamsFile(GOSSIPSUB_PARAMETERS)
	    if er != nil {
	        log.Panic(error)
	    }

	    fmt.Printf("%+v\n", paramsList)

		for _, param := range paramsList {
		    for _, hostname := range hosts {
		    	log.Println(hostname+" Starting rippled")
			    start := "nohup " + RIPPLED_PATH+"rippled --silent --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM+" & \n"
			    go remoteShell(start, hostname, config)
		    }
		    time.Sleep(60 * time.Second)
		    log.Println("Connecting to ", PUPPET)
    		go runPuppet(experiment, config, timeout, param)
		}
	} else if machine == "puppet" {

		//Get parameters from command line
		param := OverlayParams{
	        d:            *d,
	        dlo:          *dlo,
	        dhi:          *dhi,
	        dscore:       *dscore,
	        dlazy:        *dlazy,
	        dout:         *dout,
	        gossipFactor: *gossipFactor,
	    }

		//Connect and start gossipsub
		gossipsub := "cd "+GOSSIPSUB_PATH+" && go run . -type="+experiment+" -d="+param.d+" -dlo="+param.dlo+" -dhi="+param.dhi+" -dscore="+param.dscore+" -dlazy="+param.dlazy+" -dout="+param.dout+"\n"
		for _, hostname := range hosts {
			go executeCmd(gossipsub, hostname, config)
		}

	}

	// time.Sleep(100 * time.Second)
	select {}
 }