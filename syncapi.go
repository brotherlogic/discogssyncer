package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/discogssyncer/server"
)

import "os"

import "strings"

import "github.com/brotherlogic/godiscogs"

import "google.golang.org/grpc"

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	saveLocation string
	token        string
	retr         saver
	relMap       map[int32]*godiscogs.Release
	wants        pb.Wantlist
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

func (s *Syncer) initWantlist() {
	wldata, _ := ioutil.ReadFile(s.saveLocation + "/metadata/wantlist")
	proto.Unmarshal(wldata, &s.wants)

	for _, want := range s.wants.Want {
		rel, _ := s.GetRelease(int(want.ReleaseId), -5)
		rel.FolderId = -5
		s.relMap[rel.Id] = rel
	}
}

// InitServer builds an initial server
func InitServer(token *string, folder *string, retr saver) Syncer {
	syncer := Syncer{&goserver.GoServer{}, *folder, *token, retr, make(map[int32]*godiscogs.Release), pb.Wantlist{}}
	syncer.relMap = make(map[int32]*godiscogs.Release)

	//Build out the release map
	releases, _ := syncer.GetCollection(context.Background(), &pb.Empty{})
	for _, release := range releases.Releases {
		syncer.relMap[release.Id] = release
	}

	syncer.initWantlist()
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

	if *sync {
		syncTime = time.Now().Unix()
		syncer.SaveCollection(retr)
		syncer.SyncWantlist()
		syncer.clean()
	} else {
		syncer.PrepServer()
		syncer.RegisterServer("discogssyncer", false)
		syncer.Serve()
	}
}
