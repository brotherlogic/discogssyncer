package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brotherlogic/godiscogs"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
)

// GetRelease Gets the release and metadata for the release
func (syncer *Syncer) GetRelease(id int, folder int) (*pbd.Release, *pb.ReleaseMetadata) {
	filename := syncer.saveLocation + "/" + strconv.Itoa(folder) + "/" + strconv.Itoa(id) + ".release"
	releaseData, err := ioutil.ReadFile(filename)
	var release *pbd.Release
	if err != nil {
		log.Printf("Failing to read file: %v", err)
	} else {
		syncer.cache[int32(id)] = filename
		release = &pbd.Release{}
		proto.Unmarshal(releaseData, release)
	}
	metadataData, err := ioutil.ReadFile(syncer.saveLocation + "/static-metadata/" + strconv.Itoa(id) + ".metadata")
	if err == nil {
		metadata := &pb.ReleaseMetadata{}
		proto.Unmarshal(metadataData, metadata)
		return release, metadata
	}

	log.Printf("Error in reading metadata: %v", err)

	// We have no metadata for this release
	return release, nil
}

func match(query string, str string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(query))
}

// Search performs a search of the database
func (syncer *Syncer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.ReleaseList, error) {
	all, _ := syncer.GetCollection(ctx, &pb.Empty{})
	fil := &pb.ReleaseList{}
	for _, rel := range all.Releases {
		if match(req.Query, rel.Title) || match(req.Query, pbd.GetReleaseArtist(*rel)) {
			fil.Releases = append(fil.Releases, rel)
		}
	}
	return fil, nil
}

// GetSpend gets the spend
func (syncer *Syncer) GetSpend(ctx context.Context, req *pb.SpendRequest) (*pb.SpendResponse, error) {
	spend := 0
	var updates []*pb.MetadataUpdate
	col, _ := syncer.GetCollection(ctx, &pb.Empty{})
	for _, rel := range col.Releases {
		_, metadata := syncer.GetRelease(int(rel.Id), int(rel.FolderId))
		datev := time.Unix(metadata.DateAdded, 0)
		log.Printf("WE ARE HERE %v", req.Month)
		log.Printf("%v -> %v", req.Lower, req.Upper)
		log.Printf("%v", metadata.DateAdded)
		if (req.Year <= 0 || datev.Year() == int(req.Year)) && (req.Month <= 0 || int32(datev.Month()) == req.Month) && (req.Lower <= 0 || (metadata.DateAdded >= req.Lower && metadata.DateAdded <= req.Upper)) {
			spend += int(metadata.Cost)
			updates = append(updates, &pb.MetadataUpdate{Release: rel, Update: metadata})
		}
	}

	return &pb.SpendResponse{TotalSpend: int32(spend), Spends: updates}, nil
}

// AddWant adds a want to our list
func (syncer *Syncer) AddWant(ctx context.Context, req *pb.Want) (*pb.Empty, error) {
	//Add the want to discogs
	syncer.retr.AddToWantlist(int(req.ReleaseId))

	//Save and store the want
	release, _ := syncer.retr.GetRelease(int(req.ReleaseId))
	syncer.saveRelease(&release, -5)

	//Add the want internally
	syncer.wants.Want = append(syncer.wants.Want, req)
	syncer.saveWantList()

	return &pb.Empty{}, nil
}

func (syncer *Syncer) saveMetadata(rel *godiscogs.Release) {
	log.Printf("SAVING METADATA: %v", rel)
	metadataRoot := syncer.saveLocation + "/static-metadata/"
	metadataPath := metadataRoot + strconv.Itoa(int(rel.Id)) + ".metadata"
	if _, err := os.Stat(metadataRoot); os.IsNotExist(err) {
		os.MkdirAll(metadataRoot, 0777)
	}

	metadata := &pb.ReleaseMetadata{}
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		// Only set the date added if this isn't a want
		if rel.FolderId >= 0 {
			metadata.DateAdded = time.Now().Unix()
		}
		metadata.DateRefreshed = time.Now().Unix()
	} else {
		data, _ := ioutil.ReadFile(metadataPath)
		proto.Unmarshal(data, metadata)
		metadata.DateRefreshed = time.Now().Unix()

		//Set the data added if this is not a want
		if rel.FolderId >= 0 && metadata.DateAdded <= 0 {
			metadata.DateAdded = time.Now().Unix()
		}
	}
	log.Printf("SAVING %v", metadata)
	data, _ := proto.Marshal(metadata)
	ioutil.WriteFile(metadataPath, data, 0644)
}

func (syncer *Syncer) saveRelease(rel *godiscogs.Release, folder int) {
	//Check that the save folder exists
	savePath := syncer.saveLocation + "/" + strconv.Itoa(folder) + "/"
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0777)
	}

	data, _ := proto.Marshal(rel)
	log.Printf("SAVING RELEASE %v", rel)
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
	masterMap := make(map[int32][]int32)
	for _, release := range releases {
		fullRelease, err := retr.GetRelease(int(release.Id))
		if err != nil {
			log.Printf("ERROR in SaveCollection: %v for release %v", err, release)
		}
		fullRelease.InstanceId = release.InstanceId
		fullRelease.FolderId = release.FolderId
		fullRelease.Rating = release.Rating
		syncer.saveRelease(&fullRelease, int(release.FolderId))
		if _, ok := masterMap[fullRelease.MasterId]; ok {
			masterMap[fullRelease.MasterId] = append(masterMap[fullRelease.MasterId], fullRelease.Id)
		} else {
			masterMap[fullRelease.MasterId] = []int32{fullRelease.Id}
		}
	}

	//Process out the multi release map
	log.Printf("META MAP: %v", masterMap)
	for key, value := range masterMap {
		for _, rel := range value {
			meta, _ := syncer.GetMetadata(context.Background(), &godiscogs.Release{Id: rel})
			if key != 0 && len(value) > 1 {
				meta.Others = true
			} else {
				meta.Others = false
			}
			log.Printf("Updating %v with %v", rel, meta)
			syncer.UpdateMetadata(context.Background(), &pb.MetadataUpdate{Release: &godiscogs.Release{Id: rel}, Update: meta})
		}
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
	log.Printf("Getting Single Release: %v", in)

	if val, ok := syncer.cache[in.Id]; ok {
		log.Printf("READING FROM CACHE: %v", val)
		releaseData, err := ioutil.ReadFile(val)
		if err == nil {
			release := &pbd.Release{}
			proto.Unmarshal(releaseData, release)
			return release, nil
		}
	}

	col, _ := syncer.GetCollection(ctx, &pb.Empty{})
	for _, rel := range col.Releases {
		if rel.Id == in.Id {
			log.Printf("Returning %v", rel)
			return rel, nil
		}
	}

	// We might be asking for a want here
	rel, _ := syncer.GetRelease(int(in.Id), -5)
	if rel != nil {
		return rel, nil
	}

	log.Printf("Returning nil!")
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

// AddToFolder adds a release to the specified folder
func (syncer *Syncer) AddToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	syncer.retr.AddToFolder(int(in.NewFolderId), int(in.Release.Id))
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.saveRelease(&fullRelease, int(in.NewFolderId))
	return &pb.Empty{}, nil
}

// UpdateRating updates the rating of a release
func (syncer *Syncer) UpdateRating(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	syncer.retr.SetRating(int(in.FolderId), int(in.Id), int(in.InstanceId), int(in.Rating))
	fullRelease, _ := syncer.GetRelease(int(in.Id), int(in.FolderId))
	fullRelease.Rating = int32(in.Rating)
	syncer.saveRelease(fullRelease, int(fullRelease.FolderId))
	return &pb.Empty{}, nil
}

// UpdateMetadata updates the metadata of a given record
func (syncer *Syncer) UpdateMetadata(ctx context.Context, in *pb.MetadataUpdate) (*pb.ReleaseMetadata, error) {
	metadata, err := syncer.GetMetadata(ctx, in.Release)
	if err != nil {
		return nil, err
	}

	proto.Merge(metadata, in.Update)

	// Manual set of boolean fields
	if !in.Update.Others {
		metadata.Others = false
	}

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
	log.Printf("Getting Metadata for %v -> %v", in, metadata)
	if metadata == nil {
		return nil, errors.New("Failed to get metadata for release")
	}
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
			filename := syncer.saveLocation + "/" + strconv.Itoa(folderID) + "/" + file.Name()
			data, _ := ioutil.ReadFile(filename)
			release := &pbd.Release{}
			proto.Unmarshal(data, release)
			syncer.cache[release.Id] = filename
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
func (syncer *Syncer) DeleteWant(ctx context.Context, in *pb.Want) (*pb.Wantlist, error) {
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

	syncer.retr.RemoveFromWantlist(int(in.ReleaseId))
	syncer.saveWantList()
	return &syncer.wants, nil
}
