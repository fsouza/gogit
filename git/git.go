// Copyright 2012 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git

// #cgo pkg-config: libgit2
// #include <git2.h>
import "C"

import (
	"unsafe"
)

func oidToString(oid *C.git_oid) string {
	id := C.git_oid_allocfmt(oid)
	defer C.free(unsafe.Pointer(id))
	return C.GoString(id)
}

// Commit represents a git commit.
type Commit struct {
	commit *C.struct_git_commit
}

// Free is used to deallocate a commit object.
func (c *Commit) Free() {
	C.git_commit_free(c.commit)
}

// Id returns the hash of the commit.
func (c *Commit) Id() string {
	oid := C.git_commit_id(c.commit)
	defer C.free(unsafe.Pointer(oid))
	return oidToString(oid)
}

// Tree returns the tree pointed by the commit.
func (c *Commit) Tree() (*Tree, error) {
	t := new(Tree)
	if C.git_commit_tree(&t.tree, c.commit) != C.GIT_OK {
		return nil, lastErr()
	}
	return t, nil
}

// Repository is the basic type of the git package, it represents a git
// repository.
type Repository struct {
	repository *C.struct_git_repository
}

// InitRepository inits a new repository.
//
// If the path does not exist, it will be created.
//
// Returns an instance of Repository, or an error in case of failure.
func InitRepository(path string, bare bool) (*Repository, error) {
	var cbare C.unsigned = 0
	if bare {
		cbare = 1
	}
	repo := new(Repository)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.git_repository_init(&repo.repository, cpath, cbare) != C.GIT_OK {
		return nil, lastErr()
	}
	return repo, nil
}

// OpenRepository opens a repository by its path.
//
// Returns an error in case of failure (e.g.: if the path does not exist; if
// it exists but is not a git repository or if the path exists, is a git
// repository, but the user does not have access to it).
func OpenRepository(path string) (*Repository, error) {
	repo := new(Repository)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.git_repository_open(&repo.repository, cpath) != C.GIT_OK {
		return nil, lastErr()
	}
	return repo, nil
}

// Config returns a Config instance, representing the configuration of the
// repository.
func (r *Repository) Config() (*Config, error) {
	conf := new(Config)
	if C.git_repository_config(&conf.config, r.repository) != C.GIT_OK {
		return nil, lastErr()
	}
	return conf, nil
}

// Free is used to deallocate the repository. It should be called to finish the
// repository. It's a good practice to use it with the defer statement:
//
//     repo, err := git.GetRepository("/path/to/repository")
//     // check error
//     defer repo.Free()
//     // use repo
func (r *Repository) Free() {
	C.git_repository_free(r.repository)
}

// Head returns the commit at the head of the repository.
func (r *Repository) Head() (*Commit, error) {
	var reference *C.struct_git_reference
	if C.git_repository_head(&reference, r.repository) != C.GIT_OK {
		return nil, lastErr()
	}
	defer C.git_reference_free(reference)
	c := new(Commit)
	if C.git_commit_lookup(&c.commit, r.repository, C.git_reference_oid(reference)) != C.GIT_OK {
		return nil, lastErr()
	}
	return c, nil
}

// Config represents the configuration of a git repository.
//
// You can use it to retrieve or to define settings on the repository.
type Config struct {
	config *C.struct_git_config
}

// Free is used to deallocate the Config instance. It should be called to
// finish the instance. You can use it with the defer statement:
//
//     // get repository instance
//     config, err := repo.Config()
//     // check error
//     defer config.Free()
func (c *Config) Free() {
	C.git_config_free(c.config)
}

// GetBool is used to get boolean config values.
//
// The dot notation is used for configuration parameters. Example:
//
//     v, err := config.GetBool("core.ignorecase")
//     // check errors and use v
func (c *Config) GetBool(name string) (bool, error) {
	var v C.int
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_bool(&v, c.config, cname) != C.GIT_OK {
		return false, lastErr()
	}
	return v == 1, nil
}

// SetBool is used to add a boolean setting to the configuration file.
//
// The format of the configuration parameter is the same as in GetBool. If the
// configuration parameter is not declared in the config file, it will be
// created. Example of use:
//
//     err := config.SetBool("core.ignorecase", true)
//     if err != nil {
//         panic(err)
//     }
func (c *Config) SetBool(name string, value bool) error {
	var v C.int = 0
	if value {
		v = 1
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_set_bool(c.config, cname, v) != C.GIT_OK {
		return lastErr()
	}
	return nil
}

// GetString is used to get string config values.
//
// The format of the configuration parameter is the same as in GetBool.
func (c *Config) GetString(name string) (string, error) {
	var v *C.char
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_string(&v, c.config, cname) != C.GIT_OK {
		return "", lastErr()
	}
	return C.GoString(v), nil
}

// SetString is used to add a string setting to the config file.
//
// The format of the configuration parameter is the same as in GetBool. If the
// parameter is not declared in the config file, it will be created.
func (c *Config) SetString(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	if C.git_config_set_string(c.config, cname, cvalue) != C.GIT_OK {
		return lastErr()
	}
	return nil
}

// GetInt64 is used to get int64 config values.
//
// The format of the configuration parameter is the same as in GetBool.
func (c *Config) GetInt64(name string) (int64, error) {
	var v C.int64_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_int64(&v, c.config, cname) != C.GIT_OK {
		return 0, lastErr()
	}
	return int64(v), nil
}

// SetInt64 is used to add a int64 setting to the config file.
//
// The format of the configuration parameter is the same as in GetBool. If the
// parameter is not declared in the config file, it will be created.
func (c *Config) SetInt64(name string, value int64) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_set_int64(c.config, cname, C.int64_t(value)) != C.GIT_OK {
		return lastErr()
	}
	return nil
}

// Tree represents a git tree.
type Tree struct {
	tree *C.struct_git_tree
}

// Free is used to deallocate a git tree.
func (t *Tree) Free() {
	C.git_tree_free(t.tree)
}

func (t *Tree) Id() string {
	oid := C.git_tree_id(t.tree)
	defer C.free(unsafe.Pointer(oid))
	return oidToString(oid)
}

// GitError is the type used for errors in this package.
type GitError string

func (err GitError) Error() string {
	return string(err)
}

func lastErr() GitError {
	err := C.giterr_last()
	return GitError(C.GoString(err.message))
}
