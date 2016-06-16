package main

import "flag"
import "fmt"
import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "log"
import "math/rand"
import "time"

import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"

func getRelease(folderName string, host string, port string) *pbd.Release {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folder := &pbd.Folder{Name: folderName}

	r, err := client.GetReleasesInFolder(context.Background(), folder)
	if err != nil {
		log.Fatal("Problem getting releases %v", err)
	}

	log.Printf("RELEASES = %v from %v", rand.Intn(len(r.Releases)), len(r.Releases))
	return r.Releases[rand.Intn(len(r.Releases))]
}

func main() {
	var folder = flag.String("foldername", "", "Folder to retrieve from.")
	var host = flag.String("host", "10.0.1.35", "Hostname of server.")
	var port = flag.String("port", "50051", "Port number of server")
	flag.Parse()

	rel := getRelease(*folder, *host, *port)
	log.Printf("HERE %v", rel)
	fmt.Printf(pbd.GetReleaseArtist(*rel) + " - " + rel.Title)
}
