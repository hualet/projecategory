package main

import (
	"errors"
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
	subDires, err := ioutil.ReadDir(parent)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	for _, dir := range subDires {
		fmt.Println(dir.Name())
		name := dir.Name()
		cate := category(name)
		if cate == "" {
			return errors.New("name doesn't start with alphabeta")
		} else if cate == name {
			return nil
		} else if isProject(filepath.Join(parent, dir.Name())) {
			target := filepath.Join(root, cate)
			_, err := os.Stat(target)
			if err != nil && os.IsNotExist(err) {
				err := os.Mkdir(target, 664)
				if err != nil {
					return err
				}
			}
			return os.Rename(filepath.Join(parent, dir.Name()), filepath.Join(target, dir.Name()))
		}

		return processDir(filepath.Join(parent, dir.Name()), root)
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
