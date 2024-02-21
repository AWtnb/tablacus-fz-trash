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

func showLabel(heading string, s string) {
	fmt.Printf("\n\n[%s] %s:\n\n", strings.ToUpper(heading), s)
}

func run(cur string, trashname string) int {
	d := filesys.Dir{Path: cur, Trashname: trashname}
	selected, err := d.SelectItems()
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	scs := filesys.Children{Paths: selected}
	dest := filepath.Join(cur, trashname)
	dupls := scs.Dupls(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			pr := fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp))
			a := Asker{Prompt: pr, Accept: "y", Reject: "n"}
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
				scs.Drop(dp)
			}
		}
	}

	if len(scs.Paths) < 1 {
		return 0
	}
	if err := filesys.MakeDir(dest); err != nil {
		report(err.Error())
		return 1
	}

	if err := scs.CopyTo(dest); err != nil {
		report(err.Error())
		return 1
	}
	showLabel("done", "successfully copied everything")
	scs.Show()
	fmt.Printf("\nDeleting left files ==>")
	if err := scs.Remove(); err != nil {
		report(err.Error())
		return 1
	}
	fmt.Println("FINISHED")
	fmt.Scanln()
	return 0
}
