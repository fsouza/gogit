package git

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func createRepository() string {
	tmpdir := os.TempDir()
	p := path.Join(tmpdir, "gitrepo")
	err := os.MkdirAll(p, 0700)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("git", "init", "-q", p)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	configPath := path.Join(p, ".git", "config")
	cmd = exec.Command("cp", "testdata/config", configPath)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	return p
}

func removeRepository(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func TestConfigGetBool(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := GetRepository(p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ignorecase, err := config.GetBool("core.ignorecase")
	if err != nil {
		t.Error(err)
	} else if !ignorecase {
		t.Error("Failed to get core.ignorecase. Want true, got false.")
	}
}

func TestConfigGetString(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := GetRepository(p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user, err := config.GetString("github.user")
	if err != nil {
		t.Error(err)
	} else if user != "fsouza" {
		t.Errorf("Failed to get github.user. Want fsouza, got %s.", user)
	}
}

func TestConfigGetInt(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := GetRepository(p)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	commits, err := config.GetInt64("section.commits")
	if err != nil {
		t.Error(err)
	} else if commits != 800 {
		t.Errorf("Failed to get section.commits. Want 800, got %d.", commits)
	}
}
