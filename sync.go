package main

import (
	"errors"
	"fmt"
	"log"
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
func (syncer *Syncer) GetRelease(id int32, folder int32) (*pbd.Release, *pb.ReleaseMetadata) {
	var release *pbd.Release
	var metadata *pb.ReleaseMetadata
	for _, f := range syncer.collection.Folders {
		if f.Folder.Id == folder {
			for _, r := range f.Releases.Releases {
				if r.Id == id {
					release = r
				}
			}
		}
	}
	for _, m := range syncer.collection.Metadata {
		if m.Id == id {
			metadata = m
		}
	}
	return release, metadata
}

// MoveToFolder moves a release to the specified folder
func (syncer *Syncer) MoveToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {

	log.Printf("MOVE TO FOLDER: %v", in)
	log.Printf("Current folders: %v", syncer.collection.Folders)

	//Before doing anything check that the new folder exists
	legit := false
	for _, f := range syncer.getFolders().Folders {
		log.Printf("FOLDER = %v", f)
		if f.Id == in.NewFolderId {
			legit = true
		}
	}

	if !legit {
		return nil, errors.New("Unable to locate folder with id " + strconv.Itoa(int(in.NewFolderId)))
	}

	syncer.retr.MoveToFolder(int(in.Release.FolderId), int(in.Release.Id), int(in.Release.InstanceId), int(in.NewFolderId))
	oldFolder := in.Release.FolderId
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)

	log.Printf("Moving %v from %v to %v", in.Release.Id, in.Release.FolderId, in.NewFolderId)
	syncer.Log(fmt.Sprintf("Moving %v from %v to %v", in.Release.Id, in.Release.FolderId, in.NewFolderId))
	syncer.saveRelease(&fullRelease, in.NewFolderId)
	syncer.deleteRelease(&fullRelease, oldFolder)
	return &pb.Empty{}, nil
}

func match(query string, str string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(query))
}

// Search performs a search of the database
func (syncer *Syncer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.ReleaseList, error) {
	all, _ := syncer.GetCollection(ctx, &pb.Empty{})
	log.Printf("Running over collection: %v", all)
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
		_, metadata := syncer.GetRelease(rel.Id, rel.FolderId)
		datev := time.Unix(metadata.DateAdded, 0)
		log.Printf("WE ARE HERE %v", req.Month)
		log.Printf("%v -> %v", req.Lower, req.Upper)
		log.Printf("%v", metadata.DateAdded)
		if (req.Year <= 0 || datev.Year() == int(req.Year)) && (req.Month <= 0 || int32(datev.Month()) == req.Month) && (req.Lower <= 0 || (metadata.DateAdded >= req.Lower && metadata.DateAdded <= req.Upper)) {
			if metadata.Cost == 0 {
				spend += 3000
			} else {
				spend += int(metadata.Cost)
			}
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
	syncer.collection.Wantlist.Want = append(syncer.collection.Wantlist.Want, req)
	return &pb.Empty{}, nil
}

func (syncer *Syncer) saveMetadata(rel *godiscogs.Release) {
	log.Printf("SAVING METADATA: %v", rel)
	metadata := &pb.ReleaseMetadata{}

	index := -1
	for i, m := range syncer.collection.Metadata {
		if m.Id == rel.Id {
			metadata = m
			index = i
		}
	}

	log.Printf("Found metadata %v and %v", index, metadata)

	// Only set the date added if this isn't a want
	if rel.FolderId >= 0 && metadata.DateAdded <= 0 {
		metadata.DateAdded = time.Now().Unix()
	}

	metadata.DateRefreshed = time.Now().Unix()
	metadata.Id = rel.Id

	log.Printf("Updated %v", metadata)

	if index < 0 {
		syncer.collection.Metadata = append(syncer.collection.Metadata, metadata)
	}

	syncer.saveCollection()
}

func (syncer *Syncer) saveRelease(rel *pbd.Release, folder int32) {
	log.Printf("SAVING: %v into %v", rel, folder)

	foundFolder := false
	for _, f := range syncer.collection.Folders {
		if f.Folder.Id == folder {
			foundFolder = true
			found := false
			for i, r := range f.Releases.Releases {
				if r.Id == rel.Id {
					found = true
					f.Releases.Releases[i] = rel
				}
			}

			if !found {
				f.Releases.Releases = append(f.Releases.Releases, rel)
			}
		}
	}

	if !foundFolder {
		f := &pb.CollectionFolder{}
		f.Folder = &pbd.Folder{Id: folder}
		f.Releases = &pb.ReleaseList{Releases: []*pbd.Release{rel}}
		syncer.collection.Folders = append(syncer.collection.Folders, f)
	}

	log.Printf("FOLDERS: %v", syncer.collection.Folders)

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
	for _, want := range syncer.collection.Wantlist.Want {
		if want.ReleaseId == wantIn.ReleaseId {
			want.Valued = wantIn.Valued
		}
	}

	return wantIn, nil
}

// SaveCollection writes out the full collection to files.
func (syncer *Syncer) SaveCollection() {
	releases := syncer.retr.GetCollection()
	masterMap := make(map[int32][]int32)
	for _, release := range releases {
		fullRelease, err := syncer.retr.GetRelease(int(release.Id))
		log.Printf("PULL RELEASE %v from %v with %v", fullRelease, release.Id, err)
		fullRelease.InstanceId = release.InstanceId
		fullRelease.FolderId = release.FolderId
		fullRelease.Rating = release.Rating
		syncer.saveRelease(&fullRelease, release.FolderId)
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

	folders := syncer.retr.GetFolders()
	for _, f := range folders {
		found := false
		for _, f2 := range syncer.collection.Folders {
			log.Printf("COMPARING %v and %v (%v)", f, f2.Folder, found)
			log.Printf("%v and %v -> %v", f.Id, f2.Folder.Id, f.Id == f2.Folder.Id)
			if f.Id == f2.Folder.Id {
				f2.Folder.Name = f.Name
				found = true
			}
		}
		if !found {
			log.Printf("NOT FOUND FOLDER: %v", f)
			syncer.collection.Folders = append(syncer.collection.Folders, &pb.CollectionFolder{Folder: &f, Releases: &pb.ReleaseList{Releases: make([]*pbd.Release, 0)}})
		}
	}

	syncer.saveCollection()
}

// SyncWantlist syncs the wantlist with the server
func (syncer *Syncer) SyncWantlist() {
	wants, _ := syncer.retr.GetWantlist()

	for _, want := range wants {
		seen := false
		var val *pb.Want
		for _, swant := range syncer.collection.Wantlist.Want {
			if swant.ReleaseId == want.Id {
				seen = true
				val = swant
			}
		}

		if seen {
			val.Wanted = true
		} else {
			syncer.collection.Wantlist.Want = append(syncer.collection.Wantlist.Want, &pb.Want{ReleaseId: want.Id, Valued: false, Wanted: true})
		}
	}

	// Cache the want list releases
	for _, want := range syncer.collection.Wantlist.Want {
		release, _ := syncer.retr.GetRelease(int(want.ReleaseId))
		syncer.saveRelease(&release, -5)
	}

	syncer.saveCollection()
}

func (syncer *Syncer) getFolders() *pb.FolderList {
	folderList := &pb.FolderList{}
	for _, folder := range syncer.collection.Folders {
		folderList.Folders = append(folderList.Folders, folder.GetFolder())
	}
	return folderList
}

// GetSingleRelease gets a single release
func (syncer *Syncer) GetSingleRelease(ctx context.Context, in *pbd.Release) (*pbd.Release, error) {
	t1 := time.Now()
	log.Printf("HERE :%v -> %v", in, len(syncer.collection.Folders))
	col, _ := syncer.GetCollection(ctx, &pb.Empty{})
	for _, rel := range col.Releases {
		if rel.Id == in.Id {
			log.Printf("Returning %v", rel)
			syncer.LogFunction("GetSingleRelease-coll", int32(time.Now().Sub(t1).Nanoseconds()/1000000))
			return rel, nil
		}
	}

	// We might be asking for a want here
	rel, _ := syncer.GetRelease(in.Id, -5)
	if rel != nil {
		syncer.LogFunction("GetSingleRelease-want", int32(time.Now().Sub(t1).Nanoseconds()/1000000))
		return rel, nil
	}

	//Let's reach out to discogs and see if this is there
	frel, err := syncer.retr.GetRelease(int(in.Id))
	log.Printf("LOGGING")
	syncer.LogFunction("GetSingleRelease-discogs", int32(time.Now().Sub(t1).Nanoseconds()/1000000))
	return &frel, err
}

// CollapseWantlist collapses the wantlist
func (syncer *Syncer) CollapseWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	for _, want := range syncer.collection.Wantlist.Want {
		if !want.Valued {
			log.Printf("AVOIDING %v", want)
			syncer.retr.RemoveFromWantlist(int(want.ReleaseId))
			want.Wanted = false
		}
	}

	return syncer.collection.Wantlist, nil
}

// RebuildWantlist rebuilds the wantlist
func (syncer *Syncer) RebuildWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	for _, want := range syncer.collection.Wantlist.Want {
		syncer.retr.AddToWantlist(int(want.ReleaseId))
		want.Wanted = true
	}

	return syncer.collection.Wantlist, nil
}

// AddToFolder adds a release to the specified folder
func (syncer *Syncer) AddToFolder(ctx context.Context, in *pb.ReleaseMove) (*pb.Empty, error) {
	syncer.retr.AddToFolder(int(in.NewFolderId), int(in.Release.Id))
	fullRelease, _ := syncer.retr.GetRelease(int(in.Release.Id))
	fullRelease.FolderId = int32(in.NewFolderId)
	syncer.saveRelease(&fullRelease, in.NewFolderId)
	return &pb.Empty{}, nil
}

// UpdateRating updates the rating of a release
func (syncer *Syncer) UpdateRating(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	syncer.retr.SetRating(int(in.FolderId), int(in.Id), int(in.InstanceId), int(in.Rating))
	fullRelease, _ := syncer.GetRelease(in.Id, in.FolderId)
	fullRelease.Rating = int32(in.Rating)
	syncer.saveRelease(fullRelease, fullRelease.FolderId)
	return &pb.Empty{}, nil
}

// UpdateMetadata updates the metadata of a given record
func (syncer *Syncer) UpdateMetadata(ctx context.Context, in *pb.MetadataUpdate) (*pb.ReleaseMetadata, error) {
	t := time.Now()
	metadata, err := syncer.GetMetadata(ctx, in.Release)
	if err != nil {
		return nil, err
	}
	proto.Merge(metadata, in.Update)

	// Manual set of boolean fields
	if !in.Update.Others {
		metadata.Others = false
	}

	syncer.saveCollection()
	syncer.LogFunction("UpdateMetadata", int32(time.Now().Sub(t).Nanoseconds()/1000000))
	return metadata, nil
}

// GetWantlist gets the wantlist
func (syncer *Syncer) GetWantlist(ctx context.Context, in *pb.Empty) (*pb.Wantlist, error) {
	return syncer.collection.Wantlist, nil
}

// GetMetadata gets the metadata for a given release
func (syncer *Syncer) GetMetadata(ctx context.Context, in *pbd.Release) (*pb.ReleaseMetadata, error) {
	_, metadata := syncer.GetRelease(in.Id, in.FolderId)
	log.Printf("Getting Metadata for %v -> %v", in, metadata)
	if metadata == nil {
		return nil, errors.New("Failed to get metadata for release")
	}
	return metadata, nil
}

// GetReleasesInFolder serves up the releases in a given folder
func (syncer *Syncer) GetReleasesInFolder(ctx context.Context, in *pb.FolderList) (*pb.ReleaseList, error) {
	t := time.Now()
	releases := pb.ReleaseList{}
	for _, folderSpec := range in.Folders {
		folders := syncer.getFolders()
		for _, folder := range folders.Folders {
			if (len(folder.Name) > 0 && folder.Name == folderSpec.Name) || folder.Id == folderSpec.Id {
				innerReleases := syncer.getReleases(folder.Id)
				releases.Releases = append(releases.Releases, innerReleases.Releases...)
			}
		}
	}

	syncer.LogFunction("GetReleasesInFolder", int32(time.Now().Sub(t).Nanoseconds()/1000000))
	return &releases, nil
}

func (syncer *Syncer) getReleases(folderID int32) *pb.ReleaseList {
	for _, f := range syncer.collection.Folders {
		if f.Folder.Id == folderID {
			return f.Releases
		}
	}
	return nil
}

// GetCollection serves up the whole of the collection
func (syncer *Syncer) GetCollection(ctx context.Context, in *pb.Empty) (*pb.ReleaseList, error) {
	t1 := time.Now()
	releases := &pb.ReleaseList{}
	log.Printf("NOW: %v", len(syncer.collection.Folders))
	log.Printf("FOLDERS: %v", syncer.collection.Folders)
	for _, f := range syncer.collection.Folders {
		if f.Folder.Id != -5 {
			releases.Releases = append(releases.Releases, f.Releases.Releases...)
		}
	}
	syncer.LogFunction("GetCollection", int32(time.Now().Sub(t1).Nanoseconds()/1000000))
	return releases, nil
}

// DeleteWant removes a want from the system
func (syncer *Syncer) DeleteWant(ctx context.Context, in *pb.Want) (*pb.Wantlist, error) {
	//Remove the want file and remove from
	index := -1
	for i, val := range syncer.collection.Wantlist.Want {
		if val.ReleaseId == in.ReleaseId {
			index = i
		}
	}

	if index >= 0 {
		syncer.collection.Wantlist.Want = append(syncer.collection.Wantlist.Want[:index], syncer.collection.Wantlist.Want[index+1:]...)
	}

	syncer.retr.RemoveFromWantlist(int(in.ReleaseId))
	syncer.saveCollection()
	return syncer.collection.Wantlist, nil
}

//SyncWithDiscogs Syncs everything with discogs
func (syncer *Syncer) SyncWithDiscogs(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	syncer.SaveCollection()
	syncer.SyncWantlist()
	return &pb.Empty{}, nil
}
