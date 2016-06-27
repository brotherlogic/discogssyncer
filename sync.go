package main

import "errors"
import "github.com/brotherlogic/godiscogs"
import "github.com/golang/protobuf/proto"
import "io/ioutil"
import "log"
import "os"
import "strconv"
import "strings"
import "time"

import "golang.org/x/net/context"
import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"

// GetRelease Gets the release and metadata
func (syncer *Syncer) GetRelease(id int, folder int) (*pbd.Release, *pb.ReleaseMetadata) {
	releaseData, _ := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(id) + ".release")
	metadataData, _ := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(id) + ".metadata")
	release := &pbd.Release{}
	metadata := &pb.ReleaseMetadata{}

	proto.Unmarshal(releaseData, release)
	proto.Unmarshal(metadataData, metadata)
	return release, metadata
}

func (syncer *Syncer) saveMetadata(rel *godiscogs.Release, folder int) {
	metadataPath := syncer.saveLocation + strconv.Itoa(folder) + "/" + strconv.Itoa(int(rel.Id)) + ".metadata"

	metadata := &pb.ReleaseMetadata{}
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		metadata.DateAdded = time.Now().Unix()
		metadata.DateRefreshed = time.Now().Unix()
	} else {
		data, _ := ioutil.ReadFile(metadataPath)
		proto.Unmarshal(data, metadata)
		metadata.DateRefreshed = time.Now().Unix()
	}
	data, _ := proto.Marshal(metadata)
	ioutil.WriteFile(metadataPath, data, 0644)
}

func (syncer *Syncer) saveRelease(rel *godiscogs.Release, folder int) {
	//Check that the save folder exists
	savePath := syncer.saveLocation + strconv.Itoa(folder) + "/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

	data, err := proto.Marshal(rel)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	ioutil.WriteFile(savePath+strconv.Itoa(int(rel.Id))+".release", data, 0644)

	syncer.saveMetadata(rel, folder)
}

// SaveCollection writes out the full collection to files.
func (syncer *Syncer) SaveCollection(retr *godiscogs.DiscogsRetriever) {
	releases := retr.GetCollection()
	for _, release := range releases {
		syncer.saveRelease(&release, int(release.FolderId))
	}
}

func (syncer *Syncer) getFolders() *pb.FolderList {
	data, err := ioutil.ReadFile(syncer.saveLocation + "/metadata/folders")
	if err != nil {
		log.Fatal("Error reading metadata folders file: %v", err)
	}

	folderData := &pb.FolderList{}
	err2 := proto.Unmarshal(data, folderData)
	if err2 != nil {
		log.Fatal("Error unmarshalling data: %v", err2)
	}
	return folderData
}

// GetReleasesInFolder serves up the releases in a given folder
func (syncer *Syncer) GetReleasesInFolder(ctx context.Context, in *godiscogs.Folder) (*pb.ReleaseList, error) {

	folders := syncer.getFolders()
	for _, folder := range folders.Folders {
		if folder.Name == in.Name {
			return syncer.getReleases(int(folder.Id)), nil
		}
	}

	return &pb.ReleaseList{}, errors.New("Folder does not exist in collection")
}

func (syncer *Syncer) getReleases(folderID int) *pb.ReleaseList {
	releases := pb.ReleaseList{}
	files, _ := ioutil.ReadDir(syncer.saveLocation + "/" + strconv.Itoa(folderID) + "/")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".release") {
			data, err1 := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folderID) + "/" + file.Name())
			if err1 != nil {
				log.Printf("Error reading file: %v", err1)
			}
			release := &pbd.Release{}
			err2 := proto.Unmarshal(data, release)
			if err2 != nil {
				log.Printf("Error unmarshalling data: %v", err2)
			}
			releases.Releases = append(releases.Releases, release)
		}
	}
	return &releases
}

// SaveFolders saves out the list of folders
func (syncer *Syncer) SaveFolders(folders []*pbd.Folder) {
	list := pb.FolderList{}
	list.Folders = append(list.Folders, folders...)

	savePath := syncer.saveLocation + "metadata/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

	data, err := proto.Marshal(&list)
	if err != nil {
		log.Fatal("Marshalling error: %v", err)
	}
	ioutil.WriteFile(savePath+"folders", data, 0644)
}

// GetCollection serves up the whole of the collection
func (syncer *Syncer) GetCollection(ctx context.Context, in *pb.Empty) (*pb.ReleaseList, error) {
	releases := pb.ReleaseList{}
	bfiles, _ := ioutil.ReadDir(syncer.saveLocation)
	for _, bfile := range bfiles {
		if bfile.IsDir() {
			log.Printf("Searching %v", bfile)
			folderID, _ := strconv.Atoi(bfile.Name())
			for _, release := range syncer.getReleases(folderID).Releases {
				releases.Releases = append(releases.Releases, release)
			}
		}
	}
	return &releases, nil
}
