package synctesting

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestTTLExpires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New()
		c.SetWithTTL("session", "active", time.Hour)

		if _, ok := c.Get("session"); !ok {
			t.Fatal("session should exist before its deadline")
		}

		synctest.Sleep(time.Hour)

		if _, ok := c.Get("session"); ok {
			t.Fatal("session should be expired")
		}
	})
}

func TestJanitorRemovesExpiredEntries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		c := New()
		done := c.StartJanitor(ctx, time.Minute)
		c.SetWithTTL("short", "gone", 30*time.Second)
		c.SetWithTTL("long", "present", 90*time.Second)

		synctest.Sleep(time.Minute)

		if got := c.Size(); got != 1 {
			t.Fatalf("Size() after first janitor cycle = %d, want 1", got)
		}

		cancel()
		synctest.Wait()
		select {
		case <-done:
		default:
			t.Fatal("janitor did not stop")
		}
	})
}

func TestHTTPEntryExpires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := New()
		server := httptest.NewTestServer(t, NewHandler(c))
		client := server.Client()

		putReq, err := http.NewRequest(http.MethodPut, "http://cache.test/cache/session?ttl=1h", strings.NewReader("active"))
		if err != nil {
			t.Fatal(err)
		}
		putResp, err := client.Do(putReq)
		if err != nil {
			t.Fatal(err)
		}
		putResp.Body.Close()
		if putResp.StatusCode != http.StatusNoContent {
			t.Fatalf("PUT status = %d, want 204", putResp.StatusCode)
		}

		getResp := get(t, client, "http://cache.test/cache/session")
		body, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if getResp.StatusCode != http.StatusOK || string(body) != "active" {
			t.Fatalf("GET before expiry = (%d, %q), want (200, active)", getResp.StatusCode, body)
		}

		synctest.Sleep(time.Hour)

		getResp = get(t, client, "http://cache.test/cache/session")
		getResp.Body.Close()
		if getResp.StatusCode != http.StatusNotFound {
			t.Fatalf("GET after expiry status = %d, want 404", getResp.StatusCode)
		}
	})
}

func get(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
