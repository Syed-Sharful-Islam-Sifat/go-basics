package user

import (
	"context"
	"errors"
	"testing"
)

// fakeRepo is a THIRD implementation of Repository — just a map in memory.
// No MySQL, no Postgres, no Docker. Tests run in milliseconds.
type fakeRepo struct {
	users  map[int]User
	nextID int
}

func newFakeRepo() *fakeRepo { return &fakeRepo{users: map[int]User{}} }

func (f *fakeRepo) CreateTable(ctx context.Context) error { return nil }

func (f *fakeRepo) Create(ctx context.Context, name, email string) (int, error) {
	f.nextID++
	f.users[f.nextID] = User{ID: f.nextID, Name: name, Email: email}
	return f.nextID, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id int) (User, error) {
	u, ok := f.users[id]
	if !ok {
		return User{}, errors.New("not found")
	}
	return u, nil
}

func TestRegisterAndGet(t *testing.T) {
	svc := NewService(newFakeRepo()) // same Service code that runs on MySQL/Postgres
	ctx := context.Background()

	id, err := svc.Register(ctx, "  tm ", "TM@Example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := svc.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "tm" || u.Email != "tm@example.com" {
		t.Errorf("got %+v, want name=tm email=tm@example.com (trimmed + lowercased)", u)
	}
}

func TestRegisterRejectsBadInput(t *testing.T) {
	svc := NewService(newFakeRepo())
	ctx := context.Background()

	if _, err := svc.Register(ctx, "", "a@b.com"); !errors.Is(err, ErrInvalidName) {
		t.Errorf("empty name: got %v, want ErrInvalidName", err)
	}
	if _, err := svc.Register(ctx, "tm", "not-an-email"); !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("bad email: got %v, want ErrInvalidEmail", err)
	}
}

func TestGetMissingUser(t *testing.T) {
	svc := NewService(newFakeRepo())
	if _, err := svc.Get(context.Background(), 999); err == nil {
		t.Error("expected an error for a missing user")
	}
}
