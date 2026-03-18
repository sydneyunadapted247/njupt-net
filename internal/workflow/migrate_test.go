package workflow

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hicancan/njupt-net-cli/internal/core"
)

type mockSessionClient struct {
	getFn      func(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error)
	postFormFn func(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error)
}

func (m *mockSessionClient) Get(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
	if m.getFn != nil {
		return m.getFn(ctx, path, opts)
	}
	return nil, errors.New("mock get not implemented")
}

func (m *mockSessionClient) PostForm(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
	if m.postFormFn != nil {
		return m.postFormFn(ctx, path, opts)
	}
	return nil, errors.New("mock post not implemented")
}

func TestMigrateBroadband_RejectsNilSession(t *testing.T) {
	err := MigrateBroadband(context.Background(), nil, "u1", "p1", "u2", "p2", map[string]string{"FLDEXTRA1": "x"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "source session is nil") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMigrateBroadband_RejectsInvalidArgs(t *testing.T) {
	mock := &mockSessionClient{}

	err := MigrateBroadband(context.Background(), mock, "", "p1", "u2", "p2", map[string]string{"FLDEXTRA1": "x"})
	if err == nil || !strings.Contains(err.Error(), "from/to credentials are required") {
		t.Fatalf("expected credential validation error, got: %v", err)
	}

	err = MigrateBroadband(context.Background(), mock, "u1", "p1", "u2", "p2", map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "target fields are required") {
		t.Fatalf("expected target field validation error, got: %v", err)
	}
}

func TestMigrateBroadband_FailFastOnSourceLogin(t *testing.T) {
	mock := &mockSessionClient{
		getFn: func(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
			_ = ctx
			_ = path
			_ = opts
			return nil, errors.New("network down")
		},
	}

	err := MigrateBroadband(context.Background(), mock, "from", "frompwd", "to", "topwd", map[string]string{"FLDEXTRA1": "x"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "source login failed") {
		t.Fatalf("expected source login failure, got: %v", err)
	}
}

func TestMigrateBroadband_FailFastOnSourceUnbindReadback(t *testing.T) {
	operatorReadCount := 0
	mock := &mockSessionClient{
		getFn: func(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
			_ = ctx
			_ = opts
			switch path {
			case "/Self/login/?302=LI":
				return &core.SessionResponse{Body: []byte(`<html><input name="checkcode" value="abcd"></html>`)}, nil
			case "/Self/login/randomCode":
				return &core.SessionResponse{Body: []byte("ok")}, nil
			case "/Self/service/operatorId":
				operatorReadCount++
				if operatorReadCount == 1 {
					return &core.SessionResponse{Body: []byte(`<html><input name="csrftoken" value="token1"><input name="FLDEXTRA1" value="old1"><input name="FLDEXTRA2" value="old2"><input name="FLDEXTRA3" value="old3"><input name="FLDEXTRA4" value="old4"></html>`)}, nil
				}
				return &core.SessionResponse{Body: []byte(`<html><input name="csrftoken" value="token2"><input name="FLDEXTRA1" value="still-old"><input name="FLDEXTRA2" value="old2"><input name="FLDEXTRA3" value="old3"><input name="FLDEXTRA4" value="old4"></html>`)}, nil
			default:
				return nil, errors.New("unexpected get path: " + path)
			}
		},
		postFormFn: func(ctx context.Context, path string, opts core.RequestOptions) (*core.SessionResponse, error) {
			_ = ctx
			_ = opts
			switch path {
			case "/Self/login/verify":
				return &core.SessionResponse{FinalURL: "/Self/dashboard"}, nil
			case "/Self/service/bind-operator":
				return &core.SessionResponse{StatusCode: 200, Body: []byte("ignored")}, nil
			default:
				return nil, errors.New("unexpected post path: " + path)
			}
		},
	}

	err := MigrateBroadband(context.Background(), mock, "from", "frompwd", "to", "topwd", map[string]string{"FLDEXTRA1": "target"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "source unbind failed") {
		t.Fatalf("expected source unbind failure, got: %v", err)
	}
}
