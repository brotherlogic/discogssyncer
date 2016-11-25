package main

import (
	"flag"
	"fmt"
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

// MoveToFolder moves a release to the specified folder
func (s *Syncer) MoveToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	s.retr.MoveToFolder(int(in.Release.FolderId), int(in.Release.Id), int(in.Release.InstanceId), int(in.NewFolderId))
	oldFolder := int(in.Release.FolderId)
	fullRelease, _ := s.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	s.relMap[fullRelease.Id] = &fullRelease

	s.Log(fmt.Sprintf("Moving %v from %v to %v", in.Release.Id, in.Release.FolderId, in.NewFolderId))
	s.saveRelease(&fullRelease, int(in.NewFolderId))
	s.deleteRelease(&fullRelease, oldFolder)
	return &pb.Empty{}, nil
}

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
	s := Syncer{&goserver.GoServer{}, *folder, *token, retr, make(map[int32]*godiscogs.Release), pb.Wantlist{}}
	s.relMap = make(map[int32]*godiscogs.Release)

	//Build out the release map
	releases, _ := s.GetCollection(context.Background(), &pb.Empty{})
	for _, release := range releases.Releases {
		s.relMap[release.Id] = release
	}

	s.initWantlist()
	s.Register = s

	return s
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

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	var token = flag.String("token", "", "Discogs Token")
	var sync = flag.Bool("sync", true, "Flag to serve rather than sync")
	flag.Parse()
	retr := godiscogs.NewDiscogsRetriever(*token)
	s := InitServer(token, folder, retr)

	if *sync {
		syncTime = time.Now().Unix()
		s.SaveCollection(retr)
		s.SyncWantlist()
		s.clean()
	} else {
		s.PrepServer()
		s.RegisterServer("discogss", false)
		s.Serve()
	}
}
