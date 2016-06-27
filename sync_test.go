package main

import "golang.org/x/net/context"
import "os"
import "testing"
import "time"
import pbd "github.com/brotherlogic/godiscogs"

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

func GetTestSyncer(foldername string) Syncer {
	syncer := Syncer{
		saveLocation: foldername,
	}
	return syncer
}

func TestGetFolders(t *testing.T) {
	syncer := GetTestSyncer(".testgetfolders/")
	var folders []*pbd.Folder
	folders = append(folders, &pbd.Folder{Name: "TestOne", Id: 1234})
	syncer.SaveFolders(folders)

	release := &pbd.Release{Id: 1234}
	syncer.saveRelease(release, 1234)

	releases, err := syncer.GetReleasesInFolder(context.Background(), &pbd.Folder{Name: "TestOne"})

	if err != nil {
		t.Errorf("Error retrieveing releases: %v", err)
	}

	if len(releases.Releases) == 0 {
		t.Errorf("GetReleasesInFolder came back empty")
	}
}

func TestSaveFolderMetaata(t *testing.T) {
	syncer := GetTestSyncer(".testSaveFolderMetadata/")
	var folders []*pbd.Folder
	folders = append(folders, &pbd.Folder{Name: "TestOne", Id: 1234})
	folders = append(folders, &pbd.Folder{Name: "TestTwo", Id: 1232})

	syncer.SaveFolders(folders)

	if _, err := os.Stat(".testSaveFolderMetadata/metadata/folders"); os.IsNotExist(err) {
		t.Errorf("Folder metedata has not been save")
	}
}
