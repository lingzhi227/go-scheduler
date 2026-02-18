package memory_test

import (
	"context"
	"testing"

	"github.com/lingzhi227/go-scheduler/pkg/memory"
	"github.com/lingzhi227/go-scheduler/pkg/memory/backends"
)

func TestSharedSpaceACL(t *testing.T) {
	backend := backends.NewInMemory()
	ss := memory.NewSharedSpace(backend)

	// Create space
	if err := ss.CreateSpace("workspace", "admin_agent"); err != nil {
		t.Fatal(err)
	}

	// Admin can write
	if err := ss.Write(context.Background(), "workspace", "admin_agent", "key1", "value1"); err != nil {
		t.Fatalf("admin write failed: %v", err)
	}

	// Non-member cannot read
	_, _, err := ss.Read(context.Background(), "workspace", "outsider", "key1")
	if err == nil {
		t.Error("non-member should not be able to read")
	}

	// Join as reader
	ss.Join("workspace", "reader_agent", memory.PermReader)

	// Reader can read
	val, ok, err := ss.Read(context.Background(), "workspace", "reader_agent", "key1")
	if err != nil {
		t.Fatalf("reader read failed: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	// Reader cannot write
	err = ss.Write(context.Background(), "workspace", "reader_agent", "key2", "value2")
	if err == nil {
		t.Error("reader should not be able to write")
	}

	// Writer can write
	ss.Join("workspace", "writer_agent", memory.PermWriter)
	if err := ss.Write(context.Background(), "workspace", "writer_agent", "key2", "value2"); err != nil {
		t.Fatalf("writer write failed: %v", err)
	}
}

func TestSharedSpaceDuplicateCreate(t *testing.T) {
	backend := backends.NewInMemory()
	ss := memory.NewSharedSpace(backend)

	ss.CreateSpace("workspace", "a1")
	err := ss.CreateSpace("workspace", "a2")
	if err == nil {
		t.Error("expected error for duplicate space creation")
	}
}

func TestSharedSpaceLeave(t *testing.T) {
	backend := backends.NewInMemory()
	ss := memory.NewSharedSpace(backend)

	ss.CreateSpace("workspace", "a1")
	ss.Join("workspace", "a2", memory.PermReader)

	// Can read
	ss.Write(context.Background(), "workspace", "a1", "k", "v")
	_, _, err := ss.Read(context.Background(), "workspace", "a2", "k")
	if err != nil {
		t.Fatal("should be able to read after join")
	}

	// After leave, cannot read
	ss.Leave("workspace", "a2")
	_, _, err = ss.Read(context.Background(), "workspace", "a2", "k")
	if err == nil {
		t.Error("should not be able to read after leave")
	}
}
