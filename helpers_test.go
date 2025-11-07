package httphelper

import (
	"testing"
)

func TestPathToResource(t *testing.T) {
	path := "/home/lmao/Vault/Artificial Inteligence/Lecture 5.md"
	actual := PathToResource(path)
	expected := "/home/lmao/Vault/Artificial+Inteligence/Lecture+5.md"

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestResourceToPath(t *testing.T) {
	resource := "/home/lmao/Vault/Artificial+Inteligence/Lecture+5.md"
	actual := ResourceToPath(resource)
	expected := "/home/lmao/Vault/Artificial Inteligence/Lecture 5.md"

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestPathToResourceNoSpace(t *testing.T) {
	path := "/home/lmao/Vault/ArtificialInteligence/Lecture5.md"
	actual := PathToResource(path)

	if actual != path {
		t.Errorf("Expected %s Got %s", path, actual)
	}
}

func TestResourceToPathNoSpace(t *testing.T) {
	resource := "/Vault/lecture-4.md"
	actual := ResourceToPath(resource)

	if actual != resource {
		t.Errorf("Expected %s Got %s", resource, actual)
	}
}
