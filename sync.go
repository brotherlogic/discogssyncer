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
	log.Printf("READING: %v", syncer.saveLocation+"/"+strconv.Itoa(folder)+"/"+strconv.Itoa(id)+".release")
	releaseData, err := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(id) + ".release")
	if err != nil {
		return nil, nil
	}

	metadataData, _ := ioutil.ReadFile(syncer.saveLocation + "/static-metadata/" + strconv.Itoa(id) + ".metadata")
	release := &pbd.Release{}
	metadata := &pb.ReleaseMetadata{}

	proto.Unmarshal(releaseData, release)
	proto.Unmarshal(metadataData, metadata)

	log.Printf("Unmarshalled to %v", release)

	log.Printf("Reading %v from %v", metadataData, syncer.saveLocation+"/static-metadata/"+strconv.Itoa(id)+".metadata")

	return release, metadata
}

func (syncer *Syncer) saveMetadata(rel *godiscogs.Release) {
	metadataRoot := syncer.saveLocation + "/static-metadata/"
	metadataPath := metadataRoot + strconv.Itoa(int(rel.Id)) + ".metadata"
	if _, err := os.Stat(metadataRoot); os.IsNotExist(err) {
		os.MkdirAll(metadataRoot, 0777)
	}

	metadata := &pb.ReleaseMetadata{}
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		metadata.DateAdded = time.Now().Unix()
		metadata.DateRefreshed = time.Now().Unix()
	} else {
		data, _ := ioutil.ReadFile(metadataPath)
		proto.Unmarshal(data, metadata)
		metadata.DateRefreshed = time.Now().Unix()
	}
	log.Printf("Writing out %v", metadata)
	data, _ := proto.Marshal(metadata)
	ioutil.WriteFile(metadataPath, data, 0644)
}

func (syncer *Syncer) deleteRelease(rel *godiscogs.Release, folder int) {
	os.Remove(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(int(rel.Id)) + ".release")
}

func (syncer *Syncer) saveRelease(rel *godiscogs.Release, folder int) {
	//Check that the save folder exists
	savePath := syncer.saveLocation + "/" + strconv.Itoa(folder) + "/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

	log.Printf("Saving %v to %v", rel, savePath+strconv.Itoa(int(rel.Id))+".release")
	data, _ := proto.Marshal(rel)
	ioutil.WriteFile(savePath+strconv.Itoa(int(rel.Id))+".release", data, 0644)
	syncer.saveMetadata(rel)
}

type saver interface {
	GetCollection() []godiscogs.Release
	GetFolders() []godiscogs.Folder
	GetRelease(id int) (godiscogs.Release, error)
	MoveToFolder(folderID int, releaseID int, instanceID int, newFolderID int)
	AddToFolder(folderID int, releaseID int)
	SetRating(folderID int, releaseID int, instanceID int, rating int)
}

// SaveCollection writes out the full collection to files.
func (syncer *Syncer) SaveCollection(retr saver) {
	releases := retr.GetCollection()
	for _, release := range releases {
		fullRelease, _ := retr.GetRelease(int(release.Id))
		fullRelease.InstanceId = release.InstanceId
		fullRelease.FolderId = release.FolderId
		fullRelease.Rating = release.Rating
		syncer.saveRelease(&fullRelease, int(release.FolderId))
	}
	folders := retr.GetFolders()
	folderList := pb.FolderList{}
	for i := range folders {
		folder := folders[i]
		folderList.Folders = append(folderList.Folders, &folder)
	}
	syncer.SaveFolders(&folderList)
}

func (syncer *Syncer) getFolders() *pb.FolderList {
	data, _ := ioutil.ReadFile(syncer.saveLocation + "/metadata/folders")
	folderData := &pb.FolderList{}
	proto.Unmarshal(data, folderData)
	return folderData
}

// MoveToFolder moves a release to the specified folder
func (syncer *Syncer) MoveToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	syncer.retr.MoveToFolder(int(in.Release.FolderId), int(in.Release.Id), int(in.Release.InstanceId), int(in.NewFolderId))
	oldFolder := int(in.Release.FolderId)
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.saveRelease(&fullRelease, int(in.NewFolderId))
	syncer.deleteRelease(&fullRelease, oldFolder)
	return &pb.Empty{}, nil
}

// AddToFolder adds a release to the specified folder
func (syncer *Syncer) AddToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	log.Printf("Adding releases %v", in)
	syncer.retr.AddToFolder(int(in.NewFolderId), int(in.Release.Id))
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.saveRelease(&fullRelease, int(in.NewFolderId))
	return &pb.Empty{}, nil
}

// UpdateRating updates the rating of a release
func (syncer *Syncer) UpdateRating(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	log.Printf("Updating rating %v", in)
	syncer.retr.SetRating(int(in.FolderId), int(in.Id), int(in.InstanceId), int(in.Rating))
	fullRelease, _ := syncer.GetRelease(int(in.Id), int(in.FolderId))
	fullRelease.Rating = int32(in.Rating)
	log.Printf("SAVING %v", fullRelease)
	syncer.saveRelease(fullRelease, int(fullRelease.FolderId))
	return &pb.Empty{}, nil
}

// UpdateMetadata updates the metadata of a given record
func (syncer *Syncer) UpdateMetadata(ctx context.Context, in *pb.MetadataUpdate) (*pb.ReleaseMetadata, error) {
	release, metadata := syncer.GetRelease(int(in.Release.Id), int(in.Release.FolderId))
	if release == nil {
		return nil, errors.New("Unable to locate release")
	}
	proto.Merge(metadata, in.Update)

	metadataRoot := syncer.saveLocation + "/static-metadata/"
	metadataPath := metadataRoot + strconv.Itoa(int(in.Release.Id)) + ".metadata"
	data, _ := proto.Marshal(metadata)
	ioutil.WriteFile(metadataPath, data, 0644)

	return metadata, nil
}

// GetMetadata gets the metadata for a given release
func (syncer *Syncer) GetMetadata(ctx context.Context, in *pbd.Release) (*pb.ReleaseMetadata, error) {
	_, metadata := syncer.GetRelease(int(in.Id), int(in.FolderId))
	return metadata, nil
}

// GetReleasesInFolder serves up the releases in a given folder
func (syncer *Syncer) GetReleasesInFolder(ctx context.Context, in *pb.FolderList) (*pb.ReleaseList, error) {

	releases := pb.ReleaseList{}
	for _, folderSpec := range in.Folders {
		folders := syncer.getFolders()
		for _, folder := range folders.Folders {
			if folder.Name == folderSpec.Name || folder.Id == folderSpec.Id {
				innerReleases := syncer.getReleases(int(folder.Id))
				releases.Releases = append(releases.Releases, innerReleases.Releases...)
			}
		}
	}

	if len(releases.Releases) > 0 {
		return &releases, nil
	}

	return &pb.ReleaseList{}, errors.New("Folder does not exist in collection")
}

func (syncer *Syncer) getReleases(folderID int) *pb.ReleaseList {
	releases := pb.ReleaseList{}
	files, _ := ioutil.ReadDir(syncer.saveLocation + "/" + strconv.Itoa(folderID) + "/")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".release") {
			data, _ := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folderID) + "/" + file.Name())
			release := &pbd.Release{}
			proto.Unmarshal(data, release)
			releases.Releases = append(releases.Releases, release)
		}
	}
	return &releases
}

// SaveFolders saves out the list of folders
func (syncer *Syncer) SaveFolders(list *pb.FolderList) {
	savePath := syncer.saveLocation + "/metadata/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

	data, _ := proto.Marshal(list)
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
