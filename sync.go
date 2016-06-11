package main

import "flag"
import "github.com/brotherlogic/godiscogs"
import "github.com/golang/protobuf/proto"
import "io/ioutil"
import "log"
import "os"
import "strconv"
import "strings"

import "golang.org/x/net/context"
import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"

// Syncer the configuration for the syncer
type Syncer struct {
	saveLocation string
	token        string
	host         string
	port         string
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

// GetCollection serves up the whole of the collection
func (syncer *Syncer) GetCollection(ctx context.Context, in *pb.Empty) (*pb.ReleaseList, error) {
	releases := &pb.ReleaseList{}
	files, _ := ioutil.ReadDir(syncer.saveLocation)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".release") {
			data, err1 := ioutil.ReadFile(syncer.saveLocation + "/" + file.Name())
			if err1 != nil {
				log.Printf("Error reading file: %v", err1)
			}
			release := &pbd.Release{}
			err2 := proto.Unmarshal(data, release)
			if err2 != nil {
				log.Printf("Error unmarshalling data: %v", err2)
			}
			releases.Releases = append(releases.Releases, release)
		}
	}
	return releases, nil
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	flag.Parse()

	syncer := Syncer{token: *token, saveLocation: *folder}

	retr := godiscogs.NewDiscogsRetriever(*token)
	syncer.SaveCollection(retr)
}
