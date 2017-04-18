package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"os"
	"strings"

	pb "github.com/brotherlogic/discogssyncer/server"
	"github.com/brotherlogic/godiscogs"
	"google.golang.org/grpc"
)

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	saveLocation string
	token        string
	retr         saver
	wants        pb.Wantlist
	cache        map[int32]string
}

var (
	syncTime int64
)

func (s *Syncer) initWantlist() {
	wldata, _ := ioutil.ReadFile(s.saveLocation + "/metadata/wantlist")
	proto.Unmarshal(wldata, &s.wants)

	for _, want := range s.wants.Want {
		rel, err := s.GetRelease(int(want.ReleaseId), -5)
		if err != nil && rel != nil {
			rel.FolderId = -5
		}
	}
}

func (s *Syncer) deleteRelease(rel *godiscogs.Release, folder int) {
	os.Remove(s.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(int(rel.Id)) + ".release")
}

// DoRegister does RPC registration
func (s Syncer) DoRegister(server *grpc.Server) {
	pb.RegisterDiscogsServiceServer(server, &s)
}

// MoveToFolder moves a release to the specified folder
func (s *Syncer) MoveToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	s.retr.MoveToFolder(int(in.Release.FolderId), int(in.Release.Id), int(in.Release.InstanceId), int(in.NewFolderId))
	oldFolder := int(in.Release.FolderId)
	fullRelease, _ := s.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)

	s.Log(fmt.Sprintf("Moving %v from %v to %v", in.Release.Id, in.Release.FolderId, in.NewFolderId))
	s.saveRelease(&fullRelease, int(in.NewFolderId))
	s.deleteRelease(&fullRelease, oldFolder)
	return &pb.Empty{}, nil
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
func InitServer(token *string, folder *string, retr saver) Syncer {
	syncer := Syncer{&goserver.GoServer{}, *folder, *token, retr, pb.Wantlist{}, make(map[int32]string)}
	syncer.initWantlist()
	syncer.Register = syncer

	return syncer
}

// ReportHealth alerts if we're not healthy
func (s Syncer) ReportHealth() bool {
	return true
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	var sync = flag.Bool("sync", true, "Flag to serve rather than sync")
	var verbose = flag.Bool("verbose", false, "Show all output")
	flag.Parse()
	retr := godiscogs.NewDiscogsRetriever(*token)
	syncer := InitServer(token, folder, retr)

	if *sync {
		//Turn off logging
		if !*verbose {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
		}
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
