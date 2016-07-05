package main

import "flag"
import "log"
import "github.com/brotherlogic/godiscogs"
import "github.com/brotherlogic/goserver"
import "google.golang.org/grpc"

import pb "github.com/brotherlogic/discogssyncer/server"

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	saveLocation string
	token        string
}

// DoRegister does RPC registration
func (s Syncer) DoRegister(server *grpc.Server) {
	pb.RegisterDiscogsServiceServer(server, &s)
}

// InitServer builds an initial server
func InitServer(token *string, folder *string) Syncer {
	syncer := Syncer{&goserver.GoServer{}, *folder, *token}
	syncer.Register = syncer
	return syncer
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	var sync = flag.Bool("sync", true, "Flag to serve rather than sync")
	flag.Parse()

	syncer := InitServer(token, folder)

	log.Printf("HERE = %v", *sync)

	if *sync {
		retr := godiscogs.NewDiscogsRetriever(*token)
		syncer.SaveCollection(retr)
	} else {
		syncer.PrepServer()
		syncer.RegisterServer("discogssyncer", false)
		syncer.Serve()
	}
}
