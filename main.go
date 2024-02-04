package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-trash/filesys"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur       string
		trashname string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.StringVar(&trashname, "trashname", "_obsolete", "trash folder name")
	flag.Parse()
	os.Exit(run(cur, trashname))
}

func report(s string) {
	fmt.Printf("ERROR: %s\n", s)
	fmt.Scanln()
}

func border(s string) {
	fmt.Printf("\n======================================\n %s\n======================================\n", strings.ToUpper(s))
}

func newDir(path string) error {
	if f, err := os.Stat(path); err == nil && f.IsDir() {
		return nil
	}
	return os.Mkdir(path, os.ModePerm)
}

func run(c string, trashname string) int {
	d := filesys.Dir{Path: c}
	selected, err := d.SelectItems(true, false)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	sfs := filesys.Files{Paths: selected}
	dest := filepath.Join(c, trashname)
	targets := sfs.GetNonDuplicates(dest)
	dupls := sfs.GetDuplicates(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			pr := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp))
			a := Asker{Prompt: pr, Accept: "y", Reject: "n"}
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
			} else {
				targets = append(targets, dp)
			}
		}
	}

	if len(targets) < 1 {
		return 0
	}
	if err := newDir(dest); err != nil {
		report(err.Error())
		return 1
	}

	t := filesys.Files{Paths: targets}
	if err := t.CopyFiles(dest); err != nil {
		report(err.Error())
		return 1
	}
	border("successfully copied everything")
	t.Show()
	fmt.Printf("Deleting left files ==>")
	if err := t.RemoveFiles(); err != nil {
		report(err.Error())
		return 1
	}
	fmt.Println("FINISHED")
	fmt.Scanln()
	return 0
}
