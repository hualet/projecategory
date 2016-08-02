package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// const alphabeta = "abcdefghijklmnopqrstuvwxyz"

func startsWithLetter(str string) bool {
	if len(str) == 0 {
		return false
	}

	return 'a' <= str[0] && str[0] <= 'z'
}

func category(str string) string {
	if len(str) == 0 {
		return ""
	}

	return strings.ToLower(string(str[0]))
}

func isProject(dir string) bool {
	path := filepath.Join(dir, ".git")
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func processDir(parent, root string) error {
	fmt.Printf("processing directory %v\n", parent)
	subDires, err := ioutil.ReadDir(parent)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	for _, dir := range subDires {
		name := dir.Name()
		cate := category(name)
		if cate == "" {
			fmt.Printf("name doesn't start with alphabeta : %v\n", name)
			continue
		} else if cate == name {
			fmt.Printf("skip directory %v\n", name)
			continue
		} else if isProject(filepath.Join(parent, dir.Name())) {
			target := filepath.Join(root, cate)
			_, err := os.Stat(target)
			if err != nil && os.IsNotExist(err) {
				err := os.Mkdir(target, 0755)
				if err != nil {
					return err
				}
			}

			from := filepath.Join(parent, dir.Name())
			to := filepath.Join(target, dir.Name())

			fmt.Printf("moving %v to %v \n", from, to)
			err = os.Rename(from, to)
			if err != nil {
				fmt.Println(err)
			}

			continue
		}

		processDir(filepath.Join(parent, dir.Name()), root)
	}

	return nil
}

func main() {
	flag.Parse()

	var root = flag.Arg(0)
	if root == "" {
		root = "."
	}

	processDir(root, root)
}
