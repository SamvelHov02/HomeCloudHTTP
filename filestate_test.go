package httphelper

import (
	"fmt"
	"reflect"
	"testing"
)

const TestVaultPath = "/home/samo/dev/HomeCloud/server/"

func TestTreeBuilder(t *testing.T) {
	tr := Tree{}
	tr.Init(TestVaultPath + "Vault")
	tr.Build()

	expected := Tree{
		Root: "/home/samo/dev/HomeCloud/server/Vault",
		Children: []*Leaf{
			{Category: "dir", Name: TestVaultPath + "Vault/TestFolder", Children: []*Leaf{
				{Category: "file", Name: TestVaultPath + "Vault/TestFolder/test3.md"},
			}},
			{Category: "file", Name: TestVaultPath + "Vault/test1.md"},
			{Category: "file", Name: TestVaultPath + "Vault/test2.md"},
		},
	}

	if !reflect.DeepEqual(tr, expected) {
		fmt.Printf("Expected %+v\n", expected.Children[1])
		fmt.Println("--------------------------")
		fmt.Printf("Expected %+v\n", tr.Children[1])
		t.Errorf("Failed")
	}
}

func TestComputeHashesFiles(t *testing.T) {
	l := &Leaf{
		Category: "file",
		Name:     TestVaultPath + "Vault/test1.md",
	}

	l.ComputeHash()

	expected := "b0681fc23ba56e1f40a0f471ab814b7ed05539edd486debad6dc85de8128b0ae"

	if l.Hash != expected {
		t.Errorf("Expected %s got %s", expected, l.Hash)
	}
}

func TestComputeHashes(t *testing.T) {
	expected := Tree{
		Root: "/home/samo/dev/HomeCloud/server/Vault",
		Children: []*Leaf{
			{Category: "dir", Hash: "09e5daa976edc8ccc26dd708cd8f8a7f6a1f521de8b92e2a18cec54a8f1a5a02", Name: TestVaultPath + "/TestFolder", Children: []*Leaf{
				{Category: "file", Name: TestVaultPath + "/TestFolder/test3.md", Hash: "970bf3582ad35d2661b54c0b942dba8c72c69e755a0a6ff51a38f466d32924d1"},
				{Category: "file", Name: TestVaultPath + "/TestFolder/test4.md", Hash: "625c2a35c58defbe8cfad1cfd7c2c342c38c6cd50f7e0ac8bc992fe288610de0"},
			}},
			{Category: "file", Name: TestVaultPath + "/test1.md", Hash: "b0681fc23ba56e1f40a0f471ab814b7ed05539edd486debad6dc85de8128b0ae"},
			{Category: "file", Name: TestVaultPath + "/test2.md", Hash: "2876743a49514f4a250d7b0b466a0fe8e4c6e471752268112530cf89d6c7cb8a"},
		},
		RootHash: "5b4470e8e14a148a61268df0498a18be9e3367c56451786e9db882aad1938bff",
	}

	actual := Tree{}
	actual.Init(TestVaultPath + "Vault")
	actual.Build()
	actual.ComputeHash()

	if actual.RootHash != expected.RootHash {
		t.Errorf("Expected root hash %s got %s", expected.RootHash, actual.RootHash)
	}
}

func TestComputeHashDir(t *testing.T) {
	l := &Leaf{
		Category: "dir",
		Name:     TestVaultPath + "Vault/TestFolder",
		Children: []*Leaf{
			{Category: "file", Name: TestVaultPath + "Vault/TestFolder/test3.md"},
			{Category: "file", Name: TestVaultPath + "Vault/TestFolder/test4.md"},
		},
	}

	actualHash := l.ComputeHash()
	fmt.Println("Actual hash:", actualHash)

	expecteHash := "09e5daa976edc8ccc26dd708cd8f8a7f6a1f521de8b92e2a18cec54a8f1a5a02"

	if l.Hash != expecteHash {
		t.Errorf("Expected %s got %s", expecteHash, l.Hash)
	}
}
