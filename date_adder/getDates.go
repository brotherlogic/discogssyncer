package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, _ := registry.Discover(context.Background(), &entry)
	return r.Ip, int(r.Port)
}

func setDate(releaseID int, folderID int, date string) {
	layout := "02-Jan-06 03:04 PM"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Fatal(err)
	}

	release := &pbd.Release{Id: int32(releaseID), FolderId: int32(folderID)}
	metadata := &pb.ReleaseMetadata{DateAdded: t.Unix()}
	update := &pb.MetadataUpdate{Release: release, Update: metadata}

	server, port := getIP("discogssyncer", "10.0.1.17", 50055)
	conn, err := grpc.Dial(server+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDiscogsServiceClient(conn)
	_, err = client.UpdateMetadata(context.Background(), update)
	if err != nil {
		log.Fatal(err)
	}
}

func processFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	regger, err := regexp.Compile("/release/(\\d*?)\"")
	if err != nil {
		log.Fatal(err)
	}
	reggerd, err := regexp.Compile("data-header=\"Added\".*span title=\"(.*?)\"")
	if err != nil {
		log.Fatal(err)
	}

	reggerf, err := regexp.Compile("\\?folder=(\\d+)")
	if err != nil {
		log.Fatal(err)
	}

	releaseID := ""
	folderID := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		results := regger.FindStringSubmatch(text)
		if results != nil {
			releaseID = results[1]
		}
		results3 := reggerf.FindStringSubmatch(text)
		if results3 != nil {
			folderID = results3[1]
		}
		results2 := reggerd.FindStringSubmatch(text)
		if results2 != nil {
			strDate := results2[1]
			relID, _ := strconv.Atoi(releaseID)
			folID, _ := strconv.Atoi(folderID)
			setDate(relID, folID, strDate)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	processFile("data/1.htm")
	processFile("data/2.htm")
	processFile("data/3.htm")
	processFile("data/4.htm")
	processFile("data/5.htm")
	processFile("data/6.htm")
	processFile("data/7.htm")
	processFile("data/8.htm")

}
