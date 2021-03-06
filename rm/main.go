package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"strings"
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

var (
	force      *bool
	prompteach *bool
	promptonce *bool
)

func setPrompt(when string) error {
	when = strings.ToUpper(when)
	if when == "NEVER" {
		*prompteach = false
		*promptonce = false
		*force = true
	} else if when == "ALWAYS" {
		*prompteach = true
		*promptonce = false
		*force = false
	} else if when == "ONCE" {
		*prompteach = false
		*promptonce = true
		*force = false
	}
	return nil
}

func promptBeforeRemove(filename string, remove bool) bool {
	var prompt string
	if remove {
		prompt = "Remove " + filename + "?"
	} else {
		prompt = "Recurse into " + filename + "?"
	}
	var response string
	trueresponse := "yes"
	falseresponse := "no"
	for {
		fmt.Print(prompt)
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if strings.Contains(trueresponse, response) {
			return true
		} else if strings.Contains(falseresponse, response) || response == "" {
			return false
		}
	}
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Rm v0.1"
	goopt.Summary = "Remove each FILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... FILE...\n", os.Args[0]) + goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	force = goopt.Flag([]string{"-f", "--force"}, nil, "Ignore nonexistent files, don't prompt user", "")
	prompteach = goopt.Flag([]string{"-i"}, nil, "Prompt before each removal", "")
	promptonce = goopt.Flag([]string{"-I"}, nil, "Prompt before removing multiple files at once", "")
	goopt.OptArg([]string{"--interactive"}, "WHEN", "Prompt according to WHEN", setPrompt)
	/*onefs := goopt.Flag([]string{"--one-file-system"}, nil, "When -r is specified, skip directories on different filesystems", "")*/
	preserveroot := goopt.Flag([]string{"--no-preserve-root"}, []string{"--preserve-root"}, "Do not treat '/' specially", "Do not remove '/' (This is default)")
	recurse := goopt.Flag([]string{"-r", "-R", "--recursive"}, nil, "Recursively remove directories and their contents", "")
	emptydir := goopt.Flag([]string{"-d", "--dir"}, nil, "Remove empty directories", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	doubledash := false
	promptno := true
	var filenames []string
	var dirnames []string
	for i := range os.Args[1:] {
		if !doubledash && os.Args[i+1][0] == '-' {
			if os.Args[i+1] == "--" {
				doubledash = true
			}
			continue
		}
		fileinfo, err := os.Lstat(os.Args[i+1])
		if err != nil {
			fmt.Println("Error getting file info,", err)
		} else {
			if fileinfo.IsDir() {
				dirnames = append(dirnames, os.Args[i+1])
			}
		}
		filenames = append(filenames, os.Args[i+1])
	}
	i := 0
	l := len(filenames)
	for *recurse {
		*recurse = false
		for j := range filenames[i:] {
			fileinfo, err := os.Lstat(filenames[i+j])
			if err != nil {
				fmt.Println("Error getting file info,", err)
			} else {
				if fileinfo.IsDir() {
					dirnames = append(dirnames, filenames[i+j])
					if !*preserveroot || filenames[i+j] != "/" {
						if *prompteach || *promptonce {
							promptno = promptBeforeRemove(filenames[i+j], false)
						}
						filelisting, err := ioutil.ReadDir(filenames[i+j])
						if err != nil && !*force {
							fmt.Println("Could not recurse into", filenames[i+j], ":", err)
						} else if len(filelisting) > 0 {
							*recurse = true
							for h := range filelisting {
								filenames = append(filenames, filenames[i+j]+string(os.PathSeparator)+filelisting[h].Name())
							}
						}
					}
				}
			}
		}
		i = l
		l = len(filenames)
	}
	/* REVERSE FILENAMES HERE AND REPLACE l-i WITH i
	filenames = filenames.Reverse()*/
	l--
	isadir := false
	for i := range filenames {
		if *prompteach || *promptonce && (l-i)%3 == 1 {
			promptno = promptBeforeRemove(filenames[l-i], true)
		}
		for j := range dirnames {
			if filenames[l-i] == dirnames[j] {
				isadir = true
				break
			}
		}
		if promptno {
			if *emptydir || !isadir {
				if *verbose {
					fmt.Println("Removing", filenames[l-i])
				}
				err := os.Remove(filenames[l-i])
				if err != nil && !*force {
					fmt.Println("Could not remove", filenames[l-i], ":", err)
				}
			} else {
				fmt.Println("Could not remove", filenames[l-i], ": Is a directory")
			}
		}
		isadir = false
	}
	return
}
