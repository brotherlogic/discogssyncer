package main

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

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

func processMetadataFile(id int32, f string) {

	log.Printf("READING %v", f)

	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	metadata := &pb.ReleaseMetadata{}
	err = proto.Unmarshal(data, metadata)
	if err != nil {
		log.Fatalf("Unable to unmarshall data %v", err)
	}

	host, port := findServer("discogssyncer")
	if port <= 0 {
		log.Fatalf("Unable to find server")
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to dial server: %v", err)
	}

	client := pb.NewDiscogsServiceClient(conn)
	rs, err := client.UpdateMetadata(context.Background(), &pb.MetadataUpdate{Release: &pbd.Release{Id: id}, Update: metadata})
	if err != nil {
		log.Printf("Update error: %v", err)
	} else {
		log.Printf("Update success: %v", rs)
	}
}

func main() {
	dir, _ := ioutil.ReadDir("data")
	for _, f := range dir {
		parts := strings.Split(f.Name(), ".")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("Error proc: %v", err)
		}
		processMetadataFile(int32(id), "data/"+f.Name())
	}
}
