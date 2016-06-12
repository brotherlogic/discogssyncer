package main

import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "os"
import "testing"
import pbd "github.com/brotherlogic/godiscogs"
import pb "github.com/brotherlogic/discogssyncer/server"

func TestSaveLocation(t *testing.T) {
	syncer := Syncer{saveLocation: ".testfolder/"}
	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release)

	//Check that the file is in the right location
	if _, err := os.Stat(".testfolder/1234.release"); os.IsNotExist(err) {
		t.Errorf("File does not exists")
	}
}

func GetTestSyncer() Syncer {
	syncer := Syncer{
		saveLocation: ".testfolder/",
		host:         "localhost",
		port:         "12345",
	}
	return syncer
}

func RunServer() {
	syncer := GetTestSyncer()
	syncer.Serve()
}

func TestServer(t *testing.T) {
	RunServer()
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure())
	if err != nil {
		t.Errorf("Error connecting to server: %v", err)
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)

	r, err := client.GetCollection(context.Background(), &pb.Empty{})
	if err != nil {
		t.Errorf("Error getting collection: %v", err)
	}
	if len(r.Releases) == 0 {
		t.Errorf("Collection has come back empty")
	}
}
