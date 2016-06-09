package main

import "flag"
import "github.com/brotherlogic/godiscogs"
import "github.com/golang/protobuf/proto"
import "io/ioutil"
import "log"

// SaveRelease writes a single release out to a file.
func SaveRelease(rel *godiscogs.Release) {
	data, err := proto.Marshal(rel)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	ioutil.WriteFile(string(rel.Id)+".release", data, 0644)
}

// SaveCollection writes out the full collection to files.
func SaveCollection(retr *godiscogs.DiscogsRetriever) {
	releases := retr.GetCollection()
	for _, release := range releases {
		SaveRelease(&release)
	}
}

func main() {

	var token = flag.String("token", "", "Discogs Token")
	flag.Parse()

	retr := godiscogs.NewDiscogsRetriever(*token)
	SaveCollection(retr)
}
