// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git

import (
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

func createRepository() string {
	tmpdir := os.TempDir()
	p := path.Join(tmpdir, "gitrepo")
	_, err := os.Stat(p)
	if !os.IsNotExist(err) {
		os.RemoveAll(p)
	}
	err = os.MkdirAll(p, 0700)
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
	err = exec.Command("touch", path.Join(p, "README")).Run()
	if err != nil {
		panic(err)
	}
	err = exec.Command("git", "--work-tree="+p, "--git-dir="+p+"/.git", "add", ".").Run()
	if err != nil {
		panic(err)
	}
	out, err := exec.Command("git", "--work-tree="+p, "--git-dir="+p+"/.git", "commit", "-m", "initial commit").CombinedOutput()
	if err != nil {
		log.Panicf("Failed to commit (%s): %s", err, out)
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
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	ignorecase, err := config.GetBool("core.ignorecase")
	if err != nil {
		t.Fatal(err)
	} else if !ignorecase {
		t.Error("Failed to get core.ignorecase. Want true, got false.")
	}
}

func TestConfigGetString(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	user, err := config.GetString("github.user")
	if err != nil {
		t.Fatal(err)
	} else if user != "fsouza" {
		t.Errorf("Failed to get github.user. Want fsouza, got %s.", user)
	}
}

func TestConfigGetInt64(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	commits, err := config.GetInt64("section.commits")
	if err != nil {
		t.Fatal(err)
	} else if commits != 800 {
		t.Errorf("Failed to get section.commits. Want 800, got %d.", commits)
	}
}

func TestConfigSetBool(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	err = config.SetBool("core.ignorecase", false)
	if err != nil {
		t.Fatal(err)
	}
	err = config.SetBool("github.login", true)
	if err != nil {
		t.Fatal(err)
	}
	ignorecase, _ := config.GetBool("core.ignorecase")
	if ignorecase {
		t.Error("Failed to set core.ignorecase to false.")
	}
	login, err := config.GetBool("github.login")
	if err != nil {
		t.Fatal(err)
	} else if !login {
		t.Error("Set github.login to false instead of setting it to true.")
	}
}

func TestConfigSetInt64(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	err = config.SetInt64("section.commits", 300)
	if err != nil {
		t.Fatal(err)
	}
	err = config.SetInt64("section.errors", -10)
	if err != nil {
		t.Fatal(err)
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
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	config, err := r.Config()
	if err != nil {
		t.Fatal(err)
	}
	err = config.SetString("github.user", "franciscosouza")
	if err != nil {
		t.Fatal(err)
	}
	user, _ := config.GetString("github.user")
	if user != "franciscosouza" {
		t.Errorf("Failed to set github.user value, it's %s.", user)
	}
}

func TestHead(t *testing.T) {
	p := createRepository()
	defer removeRepository(p)
	r, err := OpenRepository(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()
	last, err := exec.Command("git", "--work-tree="+p, "--git-dir="+p+"/.git", "log", "-1", "--format=format:%H").CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to get last commit (%s): %s", err, last)
	}
	commit, err := r.Head()
	if err != nil {
		t.Fatal(err)
	}
	if commit.Id() != string(last) {
		t.Errorf("Failed to get head. Want %s, got %s.", last, commit.Id())
	}
}

func TestGitErrorIsAnError(t *testing.T) {
	var _ error = GitError("")
}

func TestGitError(t *testing.T) {
	expected := "Failed to do git stuff."
	err := GitError(expected)
	got := err.Error()
	if got != expected {
		t.Errorf("GitError failed on checking message. Want %s, got %s.", expected, got)
	}
}
