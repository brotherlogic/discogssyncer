package main

import "golang.org/x/net/context"
import "os"
import "testing"
import "time"
import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"

type testDiscogsRetriever struct{}

func (testDiscogsRetriever) GetCollection() []pbd.Release {
	var releases = make([]pbd.Release, 0)
	releases = append(releases, pbd.Release{FolderId: 23, Id: 25})
	releases = append(releases, pbd.Release{FolderId: 23, Id: 32})
	return releases
}

func (testDiscogsRetriever) GetRelease(id int) (pbd.Release, error) {
	return pbd.Release{}, nil
}

func (testDiscogsRetriever) GetFolders() []pbd.Folder {
	var folders = make([]pbd.Folder, 0)
	folders = append(folders, pbd.Folder{Id: 23, Name: "Testing"})
	folders = append(folders, pbd.Folder{Id: 25, Name: "TestingTwo"})
	return folders
}

func (testDiscogsRetriever) MoveToFolder(fodlerID int, releaseID int, instanceID int, newFolderID int) {
	// Do nothing
}

func (testDiscogsRetriever) AddToFolder(fodlerID int, releaseID int) {
	// Do nothing
}

func TestMoveToFolder(t *testing.T) {
	syncer := GetTestSyncer(".testMoveToFolder")
	release := &pbd.Release{FolderId: 23, Id: 25, InstanceId: 37}
	releaseMove := &pb.ReleaseMove{Release: release, NewFolderId: 20}
	_, err := syncer.MoveToFolder(context.Background(), releaseMove)
	if err != nil {
		t.Errorf("Move to uncat has returned error")
	}
}

func TestAddToFolder(t *testing.T) {
	syncer := GetTestSyncer(".testAddToFolder")
	release := &pbd.Release{Id: 25}
	releaseMove := &pb.ReleaseMove{Release: release, NewFolderId: 20}
	_, err := syncer.AddToFolder(context.Background(), releaseMove)
	if err != nil {
		t.Errorf("Move to uncat has returned error")
	}
}

func TestSaveCollection(t *testing.T) {
	syncer := Syncer{saveLocation: ".testcollectionsave/"}
	syncer.SaveCollection(&testDiscogsRetriever{})
}

func TestGetCollection(t *testing.T) {
	syncer := Syncer{saveLocation: ".testcollectionsave/"}
	syncer.SaveCollection(&testDiscogsRetriever{})

	releases, err := syncer.GetCollection(context.Background(), &pb.Empty{})
	if err != nil {
		t.Errorf("Error returned on Get Collection")
	}

	if len(releases.Releases) == 0 {
		t.Errorf("No releases have been returned")
	}

	folders := syncer.getFolders()
	if len(folders.Folders) != 2 {
		t.Errorf("Not enough folders: %v", folders)
	}

	if folders.Folders[0].Name == folders.Folders[1].Name {
		t.Errorf("FOlders have same name: %v", folders)
	}
}

func TestRetrieveEmptyCollection(t *testing.T) {
	syncer := Syncer{saveLocation: ".testemptyfolder/"}
	folderList := &pb.FolderList{}
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "TestOne", Id: 1234})
	_, err := syncer.GetReleasesInFolder(context.Background(), folderList)
	if err == nil {
		t.Errorf("Pull from empty folder returns no error!")
	}
}

func TestSaveLocation(t *testing.T) {
	syncer := Syncer{saveLocation: ".testfolder/"}
	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release, 12)

	//Check that the file is in the right location
	if _, err := os.Stat(".testfolder/12/1234.release"); os.IsNotExist(err) {
		t.Errorf("File does not exists")
	}
}

func TestSaveMetadata(t *testing.T) {
	now := time.Now()
	syncer := Syncer{saveLocation: ".testmetadatasave/"}
	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release, 12)

	_, metadata := syncer.GetRelease(1234, 12)
	if metadata.DateAdded > now.Unix() {
		t.Errorf("Metadata is prior to adding the release: %v (%v)", metadata.DateAdded, metadata.DateAdded-now.Unix())
	}
}

func TestSaveAndRefreshMetadata(t *testing.T) {
	now := time.Now()
	syncer := Syncer{saveLocation: ".testmetadatasave/"}
	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release, 12)

	_, metadata := syncer.GetRelease(1234, 12)
	if metadata.DateAdded > now.Unix() {
		t.Errorf("Metadata is prior to adding the release: %v (%v)", metadata.DateAdded, metadata.DateAdded-now.Unix())
	}

	time.Sleep(time.Second)
	syncer.saveRelease(release, 12)
	_, metadata2 := syncer.GetRelease(1234, 12)
	if metadata2.DateRefreshed == metadata.DateRefreshed {
		t.Errorf("Metadata has not been refreshed: %v and %v", metadata.DateRefreshed, metadata2.DateRefreshed)
	}
}

func GetTestSyncer(foldername string) Syncer {
	syncer := Syncer{
		saveLocation: foldername,
		retr:         testDiscogsRetriever{},
	}
	return syncer
}

func TestGetFolders(t *testing.T) {
	syncer := GetTestSyncer(".testgetfolders/")
	folders := &pb.FolderList{}
	folders.Folders = append(folders.Folders, &pbd.Folder{Name: "TestOne", Id: 1234})
	folders.Folders = append(folders.Folders, &pbd.Folder{Name: "TestTwo", Id: 1235})
	syncer.SaveFolders(folders)

	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release, 1234)

	folderList := &pb.FolderList{}
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "TestOne"})
	releases, err := syncer.GetReleasesInFolder(context.Background(), folderList)

	folderList2 := &pb.FolderList{}
	folderList2.Folders = append(folderList2.Folders, &pbd.Folder{Name: "TestTwo"})
	releases2, err2 := syncer.GetReleasesInFolder(context.Background(), folderList2)

	if err != nil {
		t.Errorf("Error retrieveing releases: %v", err)
	}

	if len(releases.Releases) == 0 {
		t.Errorf("GetReleasesInFolder came back empty")
	}

	if err2 == nil {
		t.Errorf("Error retrieving releases: %v", err2)
	}

	if len(releases2.Releases) != 0 {
		t.Errorf("Releases returned for folder 2: %v", releases2)
	}
}

func TestSaveFolderMetaata(t *testing.T) {
	syncer := GetTestSyncer(".testSaveFolderMetadata/")
	folderList := &pb.FolderList{}
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "TestOne", Id: 1234})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "TestTwo", Id: 1232})

	syncer.SaveFolders(folderList)

	if _, err := os.Stat(".testSaveFolderMetadata/metadata/folders"); os.IsNotExist(err) {
		t.Errorf("Folder metedata has not been save")
	}

}
