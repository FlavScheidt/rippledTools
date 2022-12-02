package main

import (
    "encoding/csv"
    // "io"
    // "io/ioutil"
    "os"
    "log"
    "time"
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

func readData(fileName string) ([][]string, error) {

    f, err := os.Open(fileName)
    if err != nil {
        return [][]string{}, err
    }

    defer f.Close()

    r := csv.NewReader(f)
    // skip first line
    // if _, err := r.Read(); err != nil {
    //     return [][]string{}, err
    // }

    records, err := r.ReadAll()

    if err != nil {
        return [][]string{}, err
    }

    return records, nil
}
