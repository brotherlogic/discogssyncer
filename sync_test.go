package main

import "os"
import "testing"
import pb "github.com/brotherlogic/godiscogs"

func TestSaveLocation(t *testing.T) {
	syncer := Syncer{saveLocation: ".testfolder/"}
	release := &pb.Release{Id: 1234}
	syncer.saveRelease(release)

	//Check that the file is in the right location
	if _, err := os.Stat(".testfolder/1234.release"); os.IsNotExist(err) {
		t.Errorf("File does not exists")
	}
}
