package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-asker"
	"github.com/AWtnb/tablacus-fz-trash/dir"
	"github.com/AWtnb/tablacus-fz-trash/filesys"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur       string
		focus     string
		trashname string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.StringVar(&focus, "focus", "", "selected item path")
	flag.StringVar(&trashname, "trashname", "_obsolete", "trash folder name")
	flag.Parse()
	os.Exit(run(cur, focus, trashname))
}

func report(s string) {
	fmt.Printf("ERROR: %s\n", s)
	fmt.Scanln()
}

func showLabel(heading string, s string) {
	fmt.Printf("\n\n[%s] %s:\n\n", strings.ToUpper(heading), s)
}

func run(cur string, focus string, trashname string) int {
	dest := filepath.Join(cur, trashname)
	var d dir.Dir
	d.Init(cur)
	d.Except(dest)
	var q string
	if len(focus) < 1 {
		q = ""
	} else {
		q = filepath.Base(focus)
	}
	selected, err := d.SelectItems(q)
	if err != nil {
		if err != fuzzyfinder.ErrAbort {
			report(err.Error())
		}
		return 1
	}

	var fes filesys.Entries
	fes.RegisterMulti(selected)
	dupls := fes.UnMovable(dest)
	if 0 < len(dupls) {
		for _, dp := range dupls {
			a := asker.Asker{Accept: "y", Reject: "n"}
			a.Ask(fmt.Sprintf("Name duplicated: '%s'\noverwrite?", filepath.Base(dp)))
			if !a.Accepted() {
				fmt.Printf("==> skipped\n")
				fes.Exclude(dp)
			}
		}
	}

	if fes.Size() < 1 {
		return 0
	}
	if err := dir.Create(dest); err != nil {
		report(err.Error())
		return 1
	}

	if err := fes.CopyTo(dest); err != nil {
		report(err.Error())
		return 1
	}
	showLabel("done", "successfully copied everything")
	fes.Show()
	fmt.Printf("\nDeleting left files ==>")
	if err := fes.Remove(); err != nil {
		report(err.Error())
		return 1
	}
	fmt.Println("FINISHED")
	fmt.Scanln()
	return 0
}
