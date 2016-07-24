package main

import "flag"
import "log"
import "os"
import "path/filepath"
import "strings"
import "time"

import "github.com/brotherlogic/godiscogs"
import "github.com/brotherlogic/goserver"
import "google.golang.org/grpc"


import pb "github.com/brotherlogic/discogssyncer/server"

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	saveLocation string
	token        string
	retr	     saver
}

var (
    syncTime int64
)

// DoRegister does RPC registration
func (s Syncer) DoRegister(server *grpc.Server) {
	pb.RegisterDiscogsServiceServer(server, &s)
}

func doDelete(path string, f os.FileInfo, err error) error {
     if !strings.Contains(path, "metadata/") && !f.IsDir() && f.ModTime().Unix() < syncTime {
     	return os.Remove(path)
     }
     return nil
}

func (s Syncer) clean() {
     filepath.Walk(s.saveLocation, doDelete)
}

// InitServer builds an initial server
func InitServer(token *string, folder *string, retr saver) Syncer{
	syncer := Syncer{&goserver.GoServer{}, *folder, *token, retr}
	syncer.Register = syncer
	return syncer
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	var sync = flag.Bool("sync", true, "Flag to serve rather than sync")
	flag.Parse()
	retr := godiscogs.NewDiscogsRetriever(*token)
	syncer := InitServer(token, folder, retr)

	log.Printf("HERE = %v", *sync)

	if *sync {
		syncTime = time.Now().Unix()
		syncer.SaveCollection(retr)
		syncer.clean()
	} else {
		syncer.PrepServer()
		syncer.RegisterServer("discogssyncer", false)
		syncer.Serve()
	}
}
