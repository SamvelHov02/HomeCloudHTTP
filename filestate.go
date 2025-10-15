package httphelper

import (
	"crypto/sha256"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const VaultPath = "/Users/samvelhovhannisyan/Documents/dev/Personal/HomeCloud/server/"

type Tree struct {
	Root     string  `json:"root"`
	RootHash string  `json:"hash"`
	Children []*Leaf `json:"children"`
}

type Leaf struct {
	Category string  `json:"type"`
	Children []*Leaf `json:"children"`
	Name     string  `json:"name"`
	Hash     string  `json:"hash"`
}

func (t *Tree) Init(root string) {
	t.Root = root
}

func (t *Tree) Build() {
	// Want to recurse the Vault for each level create a Leaf
	t.Children = []*Leaf{}

	contents, err := os.ReadDir(t.Root)

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range contents {
		path := filepath.Join(t.Root, entry.Name())

		if entry.IsDir() {
			l := &Leaf{
				Category: "dir",
				Name:     path,
				Children: []*Leaf{},
			}
			t.Children = append(t.Children, GetAllNames(path, l))
		} else {
			t.Children = append(t.Children, &Leaf{
				Category: "file",
				Name:     path,
			})
		}
	}
}

func (t *Tree) ComputeHash() {
	l := &Leaf{
		Category: "dir",
		Name:     t.Root,
		Children: t.Children,
	}

	t.RootHash = l.ComputeHash()
}

func (l *Leaf) ComputeHash() string {
	if l.Category == "file" {
		data, err := os.ReadFile(l.Name)
		if err != nil {
			log.Fatal(err)
		}
		hashData := sha256.Sum256(data)
		hashStr := string(hashData[:])
		return hashStr
	}

	ChildHashes := []byte{}
	for _, child := range l.Children {
		ChildHashes = append(ChildHashes, []byte(child.ComputeHash())...)
	}

	hash := sha256.Sum256(ChildHashes)

	return string(hash[:])
}

func (l Leaf) Equal(l2 Leaf) bool {
	if RelativePath(l.Name) == RelativePath(l2.Name) && l.Hash == l2.Hash {
		return true
	}
	return false
}

func GetAllNames(root string, l *Leaf) *Leaf {
	l.Children = []*Leaf{}
	contents, err := os.ReadDir(root)

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range contents {
		path := filepath.Join(root, entry.Name())

		if entry.IsDir() {
			l2 := &Leaf{
				Category: "dir",
				Name:     path,
				Children: []*Leaf{},
			}
			l.Children = append(l.Children, GetAllNames(path, l2))
		} else {
			l.Children = append(l.Children, &Leaf{
				Category: "file",
				Name:     path,
			})
		}
	}
	return l
}

func TurnFilesToLeafs(files []string) []Leaf {
	var Leaves []Leaf

	for _, file := range files {
		Leaves = append(Leaves, Leaf{Category: "file", Name: file})
	}
	return Leaves
}

func RelativePath(path string) string {
	labels := strings.Split(path, "/")

	for i, l := range labels {
		if l == "Vault" {
			return strings.Join(labels[i:], "/")
		}
	}
	return ""
}
