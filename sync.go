package main

import "flag"
import "github.com/brotherlogic/godiscogs"
import "github.com/golang/protobuf/proto"
import "io/ioutil"
import "log"
import "os"
import "strconv"

// Syncer the configuration for the syncer
type Syncer struct {
	saveLocation string
	token        string
}

func (syncer *Syncer) saveRelease(rel *godiscogs.Release) {
	//Check that the save folder exists
	if _, err := os.Stat(syncer.saveLocation); os.IsNotExist(err) {
		os.Mkdir(syncer.saveLocation, 0777)
	}

	data, err := proto.Marshal(rel)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	ioutil.WriteFile(syncer.saveLocation+strconv.Itoa(int(rel.Id))+".release", data, 0644)
}

// SaveCollection writes out the full collection to files.
func (syncer *Syncer) SaveCollection(retr *godiscogs.DiscogsRetriever) {
	releases := retr.GetCollection()
	for _, release := range releases {
		syncer.saveRelease(&release)
	}
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	flag.Parse()

	syncer := Syncer{token: *token, saveLocation: *folder}

	retr := godiscogs.NewDiscogsRetriever(*token)
	syncer.SaveCollection(retr)
}
