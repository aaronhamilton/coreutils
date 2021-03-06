package main

import (
	"fmt"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
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
	uid, gid       int
	usingreference bool
)

func fromReference(rfile string) error {
	usingreference = true
	/*owner, err*/
	return nil
}

func main() {
	goopt.Suite = "XQZ coreutils"
	goopt.Author = "William Pearson"
	goopt.Version = "Chown v0.1"
	goopt.Summary = "Change owner or group of each FILE to OWNER or GROUP\nWith reference, change owner and group of each FILE to the owner and group of RFILE"
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage:\t%s [OPTION]... [OWNER][:[GROUP]] FILE...\n or:\t%s [OPTION]... --reference=RFILE FILE...\n", os.Args[0], os.Args[0], os.Args[0]) +
			goopt.Summary + "\n\n" + goopt.Help()
	}
	goopt.Description = func() string {
		return goopt.Summary + "\n\nUnless --help or --version is passed."
	}
	changes := goopt.Flag([]string{"-c", "--changes"}, nil, "Like verbose but only report changes", "")
	silent := goopt.Flag([]string{"-f", "--silent", "--quiet"}, nil, "Suppress most error messages", "")
	verbose := goopt.Flag([]string{"-v", "--verbose"}, nil, "Output each file as it is processed", "")
	nodereference := goopt.Flag([]string{"-h", "--no-dereference"}, []string{"--derference"}, "Affect symbolic links directly instead of dereferencing them", "Dereference symbolic links before operating on them (This is default)")
	preserveroot := goopt.Flag([]string{"--preserve-root"}, []string{"--no-preserve-root"}, "Don't recurse on '/'", "Treat '/' normally (This is default)")
	goopt.OptArg([]string{"--reference"}, "RFILE", "Use RFILE's owner and group", fromReference)
	recurse := goopt.Flag([]string{"-R", "--recursive"}, nil, "Operate recursively on files and directories", "")
	goopt.NoArg([]string{"--version"}, "outputs version information and exits", version)
	goopt.Parse(nil)
	i := -1
	if !usingreference {
		for j := range os.Args[1:] {
			if os.Args[j+1][0] != '-' {
				i = j + 1
				break
			}
		}
		usergroup := strings.Split(os.Args[i], ":")
		owner, err := user.Lookup(usergroup[0])
		if err != nil {
			uid = -1
		} else {
			uid, err = strconv.Atoi(owner.Uid)
			if err != nil && !*silent {

			}
			/* hack until Go matures */
			gid, err = strconv.Atoi(owner.Gid)
			if err != nil && !*silent {
				/* stuff */
			}
		}
	}
	for j := range os.Args[1:] {
		if j != i && os.Args[j+1][0] != '-' {
			filenames := []string{os.Args[j+1]}
			for h := 0; h < len(filenames); h++ {
				/* Fix to only print when changes occur for *changes */
				if *changes || *verbose {
					fmt.Printf("Modifying ownership of %s\n", filenames[h])
				}
				if *nodereference {
					os.Lchown(filenames[h], uid, gid)
				} else {
					os.Chown(filenames[h], uid, gid)
				}
				if *recurse && (!*preserveroot || filenames[h] != "/") {
					filelisting, err := ioutil.ReadDir(filenames[h])
					if err != nil {
						fmt.Println("Could not recurse into", filenames[h])
					} else {
						for g := range filelisting {
							filenames = append(filenames, filelisting[g].Name())
						}
					}
				}
			}
		}
	}
	return
}
