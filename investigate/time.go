package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbds "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

func findServer(name string) (string, int) {
	conn, err := grpc.Dial(utils.Discover, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot reach discover server: %v (trying to discover %v)", err, name)
	}
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	re := &pbdi.RegistryEntry{Name: name}
	r, err := registry.Discover(ctx, re)

	e, ok := status.FromError(err)
	if ok && e.Code() == codes.Unavailable {
		log.Printf("RETRY")
		r, err = registry.Discover(ctx, re)
	}

	if err != nil {
		return "", -1
	}
	return r.Ip, int(r.Port)
}

func run() {
	s := time.Now()
	host, port := findServer("discogssyncer")

	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot reach discover server: %v", err)
	}
	defer conn.Close()

	client := pbds.NewDiscogsServiceClient(conn)
	res, err := client.GetReleasesInFolder(context.Background(), &pbds.FolderList{Folders: []*pbd.Folder{&pbd.Folder{Id: 812802}}})
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	log.Printf("Found %v in %v", len(res.GetRecords()), time.Now().Sub(s))
}

func main() {
	run()
}
