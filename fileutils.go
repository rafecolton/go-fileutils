// Package fileutils wraps or implements common file operations with familiar function names.
package fileutils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// ChmodR is like `chmod -R`
func ChmodR(name string, mode os.FileMode) error {
	return filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chmod(path, mode)
		}
		return err
	})
}

// ChownR is like `chown -R`
func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

// Cp is like `cp`
func Cp(src, dest string) (err error) {
	//open source
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	//create dest
	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	//copy to dest from source
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	//duplicate source permissions on dest
	si, err := os.Stat(src)
	if err != nil {
		return
	}

	if err = out.Chmod(si.Mode()); err != nil {
		return
	}

	//sync dest to disk
	err = out.Sync()

	return
}

/*
CpR is like `cp -R`
*/
func CpR(source, dest string) (err error) {
	// get properties of source dir
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return
	}

	if !sourceInfo.IsDir() {
		return errors.New("source is not a directory")
	}

	// ensure dest dir does not already exist
	if _, err = os.Open(dest); !os.IsNotExist(err) {
		return errors.New("destination already exists")
	}

	// create dest dir
	if err = os.MkdirAll(dest, sourceInfo.Mode()); err != nil {
		return
	}

	files, err := ioutil.ReadDir(source)

	for _, file := range files {
		sourceFilePath := fmt.Sprintf("%s/%s", source, file.Name())
		destFilePath := fmt.Sprintf("%s/%s", dest, file.Name())

		if file.IsDir() {
			if err = CpR(sourceFilePath, destFilePath); err != nil {
				return
			}
		} else {
			if err = Cp(sourceFilePath, destFilePath); err != nil {
				return
			}
		}

	}
	return
}

// MkdirP is `mkdir -p` / os.MkdirAll
func MkdirP(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Mv is `mv` / os.Rename
func Mv(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Rm is `rm` / os.Remove
func Rm(name string) error {
	return os.Remove(name)
}

// RmRF is `rm -rf` / os.RemoveAll
func RmRF(path string) error {
	return os.RemoveAll(path)
}

// Which is `which` / exec.LookPath
func Which(file string) (string, error) {
	return exec.LookPath(file)
}
