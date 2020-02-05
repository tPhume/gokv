package btree

import (
	"log"
	"testing"
)

func TestBtree_Simple(t *testing.T) {
	tree := NewBtree(3)

	// search empty tree
	item := tree.Search("example")
	if item != nil {
		t.Fatalf("search empty tree, expected nil, got = [%v]", item)
	}

	// insert first value
	aValue := "Hi, I am A"
	err := tree.Insert("A", map[string]string{
		"val": aValue,
	})
	if err != nil {
		t.Fatalf("insert first value, got error = [%v]", err)
	}

	// search value after insert
	item = tree.Search("A")
	if item["val"] != aValue {
		t.Fatalf("search value after insert, expected [%v], got = [%v]", aValue, item["val"])
	}
}

func TestBtree_Update(t *testing.T) {
	tree := NewBtree(3)

	// empty tree update
	err := tree.Update("DoesNotExist", map[string]string{"something": "Hi"})
	if err != KeyDoesNotExist {
		log.Fatalf("expected error = [KeyDoesNotExist], got = [%v]", err)
	}

	err = nil

	err = tree.Insert("DoesExist", map[string]string{"something": "Hi"})
	if err != nil {
		t.Fatal(err)
	}

	newValue := map[string]string{"something": "New Hi"}
	err = tree.Update("DoesExist", newValue)
	if err != nil {
		t.Fatal(err)
	}

	if v := tree.Search("DoesExist")["something"]; newValue["something"] != v {
		t.Fatalf("expected [NewHi], got = [%v]", v)
	}
}

func Test_SplitTree(t *testing.T) {
	tree := NewBtree(3)
	createTestTree(t, tree)

	if f := tree.Search("F"); f != nil {
		t.Fatalf("expected nil, got = [%v]", f)
	}

	err := tree.Insert("F", map[string]string{"val": "F"})
	if err != nil {
		t.Fatal(err)
	}

	if tree.root.items[0].getKey() != "C" {
		t.Fatalf("expected [%v], got = [%v]", "C", tree.root.items[0].getKey())
	}

	if !tree.root.node[0].leaf {
		t.Fatalf("node should be leaf")
	}

	if !tree.root.node[1].leaf {
		t.Fatalf("node should be leaf")
	}

	if tree.root.node[1].currKey != 3 {
		t.Fatalf("expected [3], got = [%v]", tree.root.node[1].currKey)
	}
}

func TestBtree_Remove(t *testing.T) {
	tree := NewBtree(3)

	err := tree.Insert("A", map[string]string{"val": "A"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Remove("A")
	if err != nil {
		t.Fatal(err)
	}

	if tree.root.node == nil {
		t.Fatalf("expected empty node, got == nil")
	}

	val := tree.Search("A")
	if val != nil {
		log.Fatalf("expected nil, got = [%v]", val)
	}

	createTestTree(t, tree)

	err = tree.Remove("J")
	t.Log(err)

	err = tree.Insert("F", map[string]string{"val": "F"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Remove("A")
	if err != nil {
		t.Fatal(err)
	}

	val = tree.Search("A")
	if val != nil {
		log.Fatalf("expected [nil], got = [%v]", err)
	}
}

func TestBtree_Remove2(t *testing.T) {
	tree := NewBtree(3)
	createTestTree(t, tree)
	err := tree.Insert("F", map[string]string{"val": "F"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Remove("C")
	if err != nil {
		t.Fatal(err)
	}
}

func createTestTree(t *testing.T, tree *Btree) {
	err := tree.Insert("A", map[string]string{"val": "A"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert("B", map[string]string{"val": "B"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert("C", map[string]string{"val": "C"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert("D", map[string]string{"val": "D"})
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert("E", map[string]string{"val": "E"})
	if err != nil {
		t.Fatal(err)
	}
}
