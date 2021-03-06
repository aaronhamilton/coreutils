package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var License = `License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law`

func version() error {
	fmt.Println(goopt.Suite + " " + goopt.Version)
	fmt.Println()
	fmt.Println("Copyright (C) 2013 " + goopt.Author)
	fmt.Println(License)
	os.Exit(0)
	return nil
}

func removeEmptyParents(dir string, verbose, ignorefail bool) bool {
	error := false
	for {
		dir = filepath.Dir(dir)
		if len(dir) == 1 {
			break
		}
		if verbose {
			fmt.Printf("Removing directory %s\n", dir)
		}
		filelisting, err := ioutil.ReadDir(dir)
		if len(filelisting) == 0 {
			err = os.Remove(dir)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", dir, err)
				error = true
			}
		} else {
			if !ignorefail {
				fmt.Println("Failed to remove 'test': directory not empty\n", dir)
			}
			return true
		}
	}
	return error
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Rmdir v0.1"
	goopt.Summary = "Remove each DIRECTORY if it is empty"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... DIRECTORY...\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	ignorefail := goopt.Flag([]string{"--ignore-fail-on-non-empty"}, nil,
		"Ignore each failure that is from a directory not being empty", "")
	parents := goopt.Flag([]string{"-p", "--parents"}, nil, "Remove DIRECTORY and ancestors if ancestors become empty", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each directory as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	error := false
	for i := range os.Args[1:] {
		if os.Args[i+1][0] == '-' {
			continue
		}
		if *verbose {
			fmt.Printf("Removing directory %s\n", os.Args[i+1])
		}
		filelisting, err := ioutil.ReadDir(os.Args[i+1])
		if err != nil {
			fmt.Printf("Failed to remove %s: %v\n", os.Args[i+1], err)
			error = true
		}
		if len(filelisting) == 0 {
			err = os.Remove(os.Args[i+1])
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", os.Args[i+1], err)
				error = true
			} else if *parents {
				dir := os.Args[i+1]
				if dir[len(dir)-1] == '/' {
					dir = filepath.Dir(dir)
				}
				if removeEmptyParents(dir, *verbose, *ignorefail) {
					error = true
				}
			}
		} else if !*ignorefail {
			fmt.Println("Failed to remove", os.Args[i+1], "directory is non-empty")
			error = true
		}
	}
	if error {
		os.Exit(1)
	}
	return
}
