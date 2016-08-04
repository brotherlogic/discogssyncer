package main

import "flag"
import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "log"
import "math/rand"
import "strconv"
import "strings"
import "time"
import "io/ioutil"

import "github.com/golang/protobuf/proto"

import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"
import pbdi "github.com/brotherlogic/discovery/proto"
import pbc "github.com/brotherlogic/cardserver/card"

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, _ := registry.Discover(context.Background(), &entry)
	return r.Ip, int(r.Port)
}

func getRelease(folderName []string, host string, port string) *pbd.Release {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	for _, name := range folderName {
		folder := &pbd.Folder{Name: name}
		folderList.Folders = append(folderList.Folders, folder)
	}
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatal("Problem getting releases %v", err)
	}

	return r.Releases[rand.Intn(len(r.Releases))]
}

func main() {
	var folder = flag.String("foldername", "", "Folder to retrieve from.")
	var host = flag.String("host", "10.0.1.17", "Hostname of server.")
	var port = flag.String("port", "50055", "Port number of server")
	flag.Parse()
	portVal, _ := strconv.Atoi(*port)

	//Read the last written record
	data, _ := ioutil.ReadFile("last_written")
	lastWritten := &pbd.Release{}
	proto.Unmarshal(data, lastWritten)

	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver", *host, portVal)
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	cardList, erf := client.GetCards(context.Background(), &pbc.Empty{})
	found := false
	foundCard := false
	for _, card := range cardList.Cards {
		if lastWritten.Title != "" && card.Hash == "discogs" {
			if pbd.GetReleaseArtist(*lastWritten)+" - "+lastWritten.Title == card.Text {
				found = true
			}
			foundCard = true
		}
	}

	if !foundCard && !found {
		dServer, dPort := getIP("discogssyncer", *host, portVal)

		//Move the previous record down to uncategorized
		dConn, _ := grpc.Dial(dServer+":"+strconv.Itoa(dPort), grpc.WithInsecure())
		defer dConn.Close()
		dClient := pb.NewDiscogsServiceClient(dConn)
		folderMove := &pb.ReleaseMove{Release: lastWritten, NewFolderId: 673768}
		log.Printf("Moving to folder: %v from %v", folderMove, lastWritten)
		log.Printf("Cardlist was %v from error %v", cardList, erf)
		dClient.MoveToFolder(context.Background(), folderMove)

		rel := getRelease(strings.Split(*folder, ","), dServer, strconv.Itoa(dPort))

		cards := pbc.CardList{}

		imageURL := ""
		backupURL := ""
		for _, image := range rel.Images {
			if image.Type == "primary" {
				imageURL = image.Uri
			}
			backupURL = image.Uri
		}
		if imageURL == "" {
			imageURL = backupURL
		}

		// Write out the chosen record
		data, _ := proto.Marshal(rel)
		ioutil.WriteFile("last_written", data, 0644)

		card := pbc.Card{Text: pbd.GetReleaseArtist(*rel) + " - " + rel.Title, Hash: "discogs", Image: imageURL, Action: pbc.Card_DISMISS, Priority: 100}
		cards.Cards = append(cards.Cards, &card)
		_, err = client.AddCards(context.Background(), &cards)
		if err != nil {
			log.Printf("Problem adding cards %v", err)
		}
	}
}
