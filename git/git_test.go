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

func TestConfigGetInt64(t *testing.T) {
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

func TestConfigSetBool(t *testing.T) {
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
	err = config.SetBool("core.ignorecase", false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = config.SetBool("github.login", true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ignorecase, _ := config.GetBool("core.ignorecase")
	if ignorecase {
		t.Error("Failed to set core.ignorecase to false.")
	}
	login, err := config.GetBool("github.login")
	if err != nil {
		t.Error(err)
	} else if !login {
		t.Error("Set github.login to false instead of setting it to true.")
	}
}

func TestConfigSetInt64(t *testing.T) {
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
	err = config.SetInt64("section.commits", 300)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = config.SetInt64("section.errors", -10)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	commits, _ := config.GetInt64("section.commits")
	if commits != 300 {
		t.Errorf("Failed to get the right value for commits. Want 300, got %d.", commits)
	}
	errors, _ := config.GetInt64("section.errors")
	if errors != -10 {
		t.Errorf("Failed to errors. Want -10, got %d.", errors)
	}
}

func TestConfigSetString(t *testing.T) {
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
	err = config.SetString("github.user", "franciscosouza")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user, _ := config.GetString("github.user")
	if user != "franciscosouza" {
		t.Errorf("Failed to set github.user value, it's %s.", user)
	}
}
