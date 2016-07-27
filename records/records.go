package main

import "flag"
import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "strconv"

import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"
import pbdi "github.com/brotherlogic/discovery/proto"

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, _ := registry.Discover(context.Background(), &entry)
	return r.Ip, int(r.Port)
}

func main() {
	var folder = flag.Int("folderid", 812802, "Folder to add to.")
	var host = flag.String("host", "10.0.1.17", "Hostname of server.")
	var port = flag.String("port", "50055", "Port number of server")
	var id = flag.Int("id", 0, "ID of record to add")

	flag.Parse()
	portVal, _ := strconv.Atoi(*port)

	dServer, dPort := getIP("discogssyncer", *host, portVal)

	//Move the previous record down to uncategorized
	dConn, err := grpc.Dial(dServer+":"+strconv.Itoa(dPort), grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	defer dConn.Close()
	dClient := pb.NewDiscogsServiceClient(dConn)

	release := &pbd.Release{Id: int32(*id)}
	folderAdd := &pb.ReleaseMove{Release: release, NewFolderId: int32(*folder)}
	_, err = dClient.AddToFolder(context.Background(), folderAdd)
	if err != nil {
		panic(err)
	}
}
