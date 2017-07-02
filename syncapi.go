package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

// Syncer the configuration for the syncer
type Syncer struct {
	*goserver.GoServer
	token      string
	retr       saver
	collection *pb.RecordCollection
}

var (
	syncTime int64
)

const (
	//KEY under which we store the collection
	KEY = "/github.com/brotherlogic/discogssyncer/collection"

	//TOKEN for discogs
	TOKEN = "/github.com/brotherlogic/discogssyncer/token"
)

// This is the only method that interacts with disk
func (s *Syncer) readRecordCollection() error {
	log.Printf("Reading collection")
	collection := &pb.RecordCollection{}
	data, err := s.KSclient.Read(KEY, collection)

	log.Printf("READ: %v", data)

	if err != nil {
		log.Printf("Unable to read collection: %v", err)
		return err
	}

	s.collection = data.(*pb.RecordCollection)
	log.Printf("FOLDERS = %v", len(s.collection.Folders))
	return nil
}

func (s *Syncer) saveCollection() {
	log.Printf("Writing collection")
	s.KSclient.Save(KEY, s.collection)
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

func findServer(name string) (string, int) {
	conn, err := grpc.Dial("192.168.86.64:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot reach discover server: %v (trying to discover %v)", err, name)
	}
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	rs, err := registry.ListAllServices(context.Background(), &pbdi.Empty{})

	if err != nil {
		log.Fatalf("Failure to list: %v", err)
	}

	for _, r := range rs.Services {
		if r.Name == name {
			log.Printf("%v -> %v", name, r)
			return r.Ip, int(r.Port)
		}
	}

	log.Printf("No %v running", name)

	return "", -1
}

// InitServer builds an initial server
func InitServer() Syncer {
	syncer := Syncer{GoServer: &goserver.GoServer{}, collection: &pb.RecordCollection{Wantlist: &pb.Wantlist{}}}
	syncer.GoServer.KSclient = *keystoreclient.GetClient(findServer)
	err := syncer.readRecordCollection()
	if err != nil {
		log.Fatalf("Unable to read record collection")
	}

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
	var quiet = flag.Bool("quiet", true, "Show all output")
	var token = flag.String("token", "", "Discogs token")
	flag.Parse()

	syncer := InitServer()

	if len(*token) > 0 {
		syncer.KSclient.Save(TOKEN, &pb.Token{Token: *token})
	}

	tType := &pb.Token{}
	tResp, err := syncer.KSclient.Read(TOKEN, tType)

	if err != nil {
		log.Fatalf("Unable to read token: %v", err)
	}

	sToken := tResp.(*pb.Token).Token
	syncer.retr = pbd.NewDiscogsRetriever(sToken)
	syncer.token = sToken

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	syncer.Register = syncer
	syncer.PrepServer()
	syncer.RegisterServer("discogssyncer", false)

	log.Printf("PRESERVER %v", len(syncer.collection.Folders))

	syncer.Serve()
}
