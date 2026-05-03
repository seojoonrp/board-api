//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

baseURL := "http://localhost:8080"

type loginResp struct {
	User  struct{ ID string } `json:"user"`
	Token string              `json:"token"`
}

type postResp struct {
	ID   string `json:"id"`
	View int    `json:"view"`
}

func login(t *testing.T) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"username": fmt.Sprintf("loadtest-%d", time.Now().UnixNano()),
	})
	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("login status=%d body=%s", resp.StatusCode, b)
	}
	var lr loginResp
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	return lr.Token
}

func createPost(t *testing.T, token string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"title": "concurrency test",
		"body":  "increment view race check",
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("create status=%d body=%s", resp.StatusCode, b)
	}
	var p postResp
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatalf("decode post: %v", err)
	}
	return p.ID
}

func getPost(t *testing.T, id string) postResp {
	t.Helper()
	resp, err := http.Get(baseURL + "/posts/" + id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("get status=%d body=%s", resp.StatusCode, b)
	}
	var p postResp
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	return p
}

func TestIncrementView_Concurrent(t *testing.T) {
	const N = 200

	token := login(t)
	postID := createPost(t, token)

	var (
		wg     sync.WaitGroup
		failed atomic.Int64
		start  = make(chan struct{})
	)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			resp, err := http.Get(baseURL + "/posts/" + postID)
			if err != nil {
				failed.Add(1)
				return
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				failed.Add(1)
			}
		}()
	}

	close(start)
	wg.Wait()

	if f := failed.Load(); f > 0 {
		t.Fatalf("%d/%d requests failed", f, N)
	}

	final := getPost(t, postID)
	expected := N + 1
	if final.View != expected {
		t.Errorf("view=%d, want %d — lost updates detected", final.View, expected)
	} else {
		t.Logf("view=%d (atomic OK)", final.View)
	}
}
