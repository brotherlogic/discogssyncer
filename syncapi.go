package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"

	"os"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	"google.golang.org/grpc"
)

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	saveLocation string
	token        string
	retr         saver
	collection   *pb.RecordCollection
}

var (
	syncTime int64
)

// This is the only method that interacts with disk
func (s *Syncer) readRecordCollection() {
	log.Printf("Reading collection")
	savePath := s.saveLocation + "/collection"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		log.Printf("No collection exists!")
	} else {
		data, _ := ioutil.ReadFile(savePath)
		proto.Unmarshal(data, s.collection)
	}

	log.Printf("READ %v", s.collection)
}

func (s *Syncer) saveCollection() {
	log.Printf("SAVE: %v", s.collection)

	//Place holder for collection save
	savePath := s.saveLocation + "/collection"
	if _, err := os.Stat(s.saveLocation); os.IsNotExist(err) {
		os.MkdirAll(s.saveLocation, 0777)
	}
	data, _ := proto.Marshal(s.collection)
	ioutil.WriteFile(savePath, data, 0644)
}

func (s *Syncer) deleteRelease(rel *pbd.Release, folder int32) {
	index := -1
	for _, f := range s.collection.Folders {
		if f.Folder.Id == folder {
			for i, r := range f.Releases.Releases {
				if r.Id == rel.Id {
					index = i
				}

			}
			if index >= 0 {
				f.Releases.Releases = append(f.Releases.Releases[:index], f.Releases.Releases[index+1:]...)
			}
		}
	}

}

// DoRegister does RPC registration
func (s Syncer) DoRegister(server *grpc.Server) {
	pb.RegisterDiscogsServiceServer(server, &s)
}

// InitServer builds an initial server
func InitServer(token *string, folder *string, retr saver) Syncer {
	syncer := Syncer{&goserver.GoServer{}, *folder, *token, retr, &pb.RecordCollection{}}
	syncer.Register = syncer

	return syncer
}

// Mote promotes/demotes this server
func (s Syncer) Mote(master bool) error {
	return nil
}

// ReportHealth alerts if we're not healthy
func (s Syncer) ReportHealth() bool {
	return true
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	var quiet = flag.Bool("quiet", true, "Show all output")
	flag.Parse()
	retr := pbd.NewDiscogsRetriever(*token)
	syncer := InitServer(token, folder, retr)

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	syncer.PrepServer()
	syncer.RegisterServer("discogssyncer", false)
	syncer.Serve()
}
