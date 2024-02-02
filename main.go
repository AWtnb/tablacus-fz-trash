package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	if err := newDir(dest); err != nil {
		report(err.Error())
		return 1
	}
	copied, err := sfs.CopyFiles(dest)
	if err != nil {
		report(err.Error())
		return 1
	}
	if len(copied) < 1 {
		return 0
	}

	disposals := filesys.Files{Paths: copied}
	disposals.Show()
	fmt.Printf("\neverything successfully copied to '%s'.\nDeleting left files ===> ", trashname)
	if err := disposals.RemoveFiles(); err != nil {
		report(err.Error())
		return 1
	}
	fmt.Printf("[FINISHED]\n\n")
	fmt.Scanln()
	return 0
}
