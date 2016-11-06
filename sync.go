package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/brotherlogic/godiscogs"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

import pb "github.com/brotherlogic/discogssyncer/server"

// GetRelease Gets the release and metadata for the release
func (syncer *Syncer) GetRelease(id int, folder int) (*pbd.Release, *pb.ReleaseMetadata) {
	releaseData, err := ioutil.ReadFile(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(id) + ".release")
	if err != nil {
		return nil, nil
	}

	metadataData, _ := ioutil.ReadFile(syncer.saveLocation + "/static-metadata/" + strconv.Itoa(id) + ".metadata")
	release := &pbd.Release{}
	metadata := &pb.ReleaseMetadata{}

	proto.Unmarshal(releaseData, release)
	proto.Unmarshal(metadataData, metadata)

	return release, metadata
}

// GetMonthlySpend gets the monthly spend
func (syncer *Syncer) GetMonthlySpend(ctx context.Context, req *pb.SpendRequest) (*pb.SpendResponse, error) {
	spend := 0
	for _, rel := range syncer.relMap {
		_, metadata := syncer.GetRelease(int(rel.Id), int(rel.FolderId))
		datev := time.Unix(metadata.DateAdded, 0)
		if datev.Year() == int(req.Year) && int32(datev.Month()) == req.Month {
			spend += int(metadata.Cost)
		}
	}

	return &pb.SpendResponse{TotalSpend: int32(spend)}, nil
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
	data, _ := proto.Marshal(metadata)
	ioutil.WriteFile(metadataPath, data, 0644)
}

func (syncer *Syncer) deleteRelease(rel *godiscogs.Release, folder int) {
	os.Remove(syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(int(rel.Id)) + ".release")
}

func (syncer *Syncer) saveRelease(rel *godiscogs.Release, folder int) {
	log.Printf("SAVING %v -> %v", rel, folder)

	//Check that the save folder exists
	savePath := syncer.saveLocation + "/" + strconv.Itoa(folder) + "/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

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
	GetWantlist() ([]pbd.Release, error)
	RemoveFromWantlist(releaseID int)
	AddToWantlist(releaseID int)
}

// EditWant edits a want in the wantlist
func (syncer *Syncer) EditWant(ctx context.Context, wantIn *pb.Want) (*pb.Want, error) {
	for _, want := range syncer.wants.Want {
		if want.ReleaseId == wantIn.ReleaseId {
			want.Valued = wantIn.Valued
		}
	}
	syncer.saveWantList()

	return wantIn, nil
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

// SyncWantlist syncs the wantlist with the server
func (syncer *Syncer) SyncWantlist() {
	wants, _ := syncer.retr.GetWantlist()

	for _, want := range wants {
		seen := false
		var val *pb.Want
		for _, swant := range syncer.wants.Want {
			if swant.ReleaseId == want.Id {
				seen = true
				val = swant
			}
		}

		if seen {
			val.Wanted = true
		} else {
			syncer.wants.Want = append(syncer.wants.Want, &pb.Want{ReleaseId: want.Id, Valued: false, Wanted: true})
		}
	}

	// Cache the want list releases
	for _, want := range syncer.wants.Want {
		release, _ := syncer.retr.GetRelease(int(want.ReleaseId))
		syncer.saveRelease(&release, -5)
	}

	syncer.saveWantList()
}

func (syncer *Syncer) getFolders() *pb.FolderList {
	data, _ := ioutil.ReadFile(syncer.saveLocation + "/metadata/folders")
	folderData := &pb.FolderList{}
	proto.Unmarshal(data, folderData)
	return folderData
}

// GetSingleRelease gets a single release
func (syncer *Syncer) GetSingleRelease(ctx context.Context, in *pbd.Release) (*pbd.Release, error) {
	if val, ok := syncer.relMap[in.Id]; ok {
		return val, nil
	}
	return nil, errors.New("Unable to find release")
}

// CollapseWantlist collapses the wantlist
func (syncer *Syncer) CollapseWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	for _, want := range syncer.wants.Want {
		if !want.Valued {
			log.Printf("AVOIDING %v", want)
			syncer.retr.RemoveFromWantlist(int(want.ReleaseId))
			want.Wanted = false
		}
	}

	return &syncer.wants, nil
}

// RebuildWantlist rebuilds the wantlist
func (syncer *Syncer) RebuildWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	for _, want := range syncer.wants.Want {
		syncer.retr.AddToWantlist(int(want.ReleaseId))
		want.Wanted = true
	}

	return &syncer.wants, nil
}

// MoveToFolder moves a release to the specified folder
func (syncer *Syncer) MoveToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	syncer.retr.MoveToFolder(int(in.Release.FolderId), int(in.Release.Id), int(in.Release.InstanceId), int(in.NewFolderId))
	oldFolder := int(in.Release.FolderId)
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.relMap[fullRelease.Id] = &fullRelease

	syncer.Log(fmt.Sprintf("Moving %v from %v to %v", in.Release.Id, in.Release.FolderId, in.NewFolderId))
	syncer.saveRelease(&fullRelease, int(in.NewFolderId))
	syncer.deleteRelease(&fullRelease, oldFolder)
	return &pb.Empty{}, nil
}

// AddToFolder adds a release to the specified folder
func (syncer *Syncer) AddToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	syncer.retr.AddToFolder(int(in.NewFolderId), int(in.Release.Id))
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.saveRelease(&fullRelease, int(in.NewFolderId))
	syncer.relMap[in.Release.Id] = &fullRelease
	return &pb.Empty{}, nil
}

// UpdateRating updates the rating of a release
func (syncer *Syncer) UpdateRating(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	syncer.retr.SetRating(int(in.FolderId), int(in.Id), int(in.InstanceId), int(in.Rating))
	fullRelease, _ := syncer.GetRelease(int(in.Id), int(in.FolderId))
	fullRelease.Rating = int32(in.Rating)
	syncer.relMap[in.Id] = fullRelease
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

// GetWantlist gets the wantlist
func (syncer *Syncer) GetWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	return &syncer.wants, nil
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

	return &releases, nil
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

func (syncer *Syncer) saveWantList() {
	data, _ := proto.Marshal(&syncer.wants)
	savePath := syncer.saveLocation + "/metadata/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}
	ioutil.WriteFile(savePath+"wantlist", data, 0644)
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
		if bfile.IsDir() && bfile.Name() != "-5" {
			folderID, _ := strconv.Atoi(bfile.Name())
			for _, release := range syncer.getReleases(folderID).Releases {
				releases.Releases = append(releases.Releases, release)
			}
		}
	}
	return &releases, nil
}

// DeleteWant removes a want from the system
func (syncer *Syncer) DeleteWant(ctx context.Context, in *pb.Want) (*pb.Empty, error) {
	//Remove the want file and remove from
	index := -1
	for i, val := range syncer.wants.Want {
		if val.ReleaseId == in.ReleaseId {
			index = i
		}
	}

	if index >= 0 {
		syncer.wants.Want = append(syncer.wants.Want[:index], syncer.wants.Want[index+1:]...)
	}
	return &pb.Empty{}, nil
}
