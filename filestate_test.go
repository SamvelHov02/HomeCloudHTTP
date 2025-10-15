package httphelper

import (
	"fmt"
	"reflect"
	"testing"
)

const TestVaultPath = "/Users/samvelhovhannisyan/Documents/dev/Personal/HomeCloud/httphelper/"

func TestTreeBuilder(t *testing.T) {
	tr := Tree{}
	tr.Init(TestVaultPath + "Vault")
	tr.Build()

	expected := Tree{
		Root: "/Users/samvelhovhannisyan/Documents/dev/Personal/HomeCloud/httphelper/Vault",
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

func TestComputeHashes(t *testing.T) {
	expected := Tree{
		Root: "/Users/samvelhovhannisyan/Documents/dev/Personal/HomeCloud/server/Vault",
		Children: []*Leaf{
			{Category: "dir", Hash: "cd372fb85148700fa88095e3492d3f9f5beb43e555e5ff26d95f5a6adc36f8e6", Name: TestVaultPath + "/TestFolder", Children: []*Leaf{
				{Category: "file", Name: TestVaultPath + "/TestFolder/test3.md", Hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			}},
			{Category: "file", Name: TestVaultPath + "/test1.md", Hash: "44db641ce3deea2461f7110d99643c7df0b9c4e4d1f56883df96f52666cd1ba6"},
			{Category: "file", Name: TestVaultPath + "/test2.md", Hash: "0ef8d27f630a07b45d00f8e3c0018dc7232bc12707831d8c3c76e2364332dbdb"},
		},
		RootHash: "d8e9274b667892d4d9ff0c6a3089ab7978264242e897746171dc82c6fbb7b43a",
	}

	actual := Tree{}
	actual.Init(TestVaultPath + "Vault")
	actual.Build()
	actual.ComputeHash()

	if !reflect.DeepEqual(actual, expected) {
		fmt.Printf("Expected %+v\n", expected)
		fmt.Println("--------------------------")
		fmt.Printf("Expected %+v\n", actual)
	}
}
