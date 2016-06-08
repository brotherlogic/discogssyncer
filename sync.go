package main

import "github.com/brotherlogic/godiscogs"
import "github.com/golang/protobuf/proto"
import "io/ioutil"
import "log"

// SaveRelease writes a single release out to a file
func SaveRelease(rel *godiscogs.Release) {
data, err := proto.Marshal(rel)
        if err != nil {
	            log.Fatal("marshaling error: ", err)
		            }
	log.Printf("%v", data)

	ioutil.WriteFile("test.release", data, 0644)
}

func main() {
     retr := godiscogs.NewDiscogsRetriever()
     release, _ := retr.GetRelease(249504)
     SaveRelease(&release)
}