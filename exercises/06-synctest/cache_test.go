package synctesting

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestTTLExpires(t *testing.T) {
	const ttl = 20 * time.Millisecond

	c := New()
	c.SetWithTTL("session", "active", ttl)

	if _, ok := c.Get("session"); !ok {
		t.Fatal("session should exist before its deadline")
	}

	started := time.Now()
	time.Sleep(ttl)
	t.Logf("waited %v of wall-clock time for a %v TTL", time.Since(started), ttl)

	if _, ok := c.Get("session"); ok {
		t.Fatal("session should be expired")
	}

	// TODO: reescreva o teste com synctest.Test, TTL de uma hora e tempo virtual.
}

func TestJanitorRemovesExpiredEntries(t *testing.T) {
	t.Skip("TODO: remova o skip e sincronize deadline mais reação do janitor")

	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		c := New()
		done := c.StartJanitor(ctx, time.Minute)
		c.SetWithTTL("short", "gone", 30*time.Second)
		c.SetWithTTL("long", "present", 90*time.Second)

		time.Sleep(time.Minute)
		// TODO: espere a reação do janitor; depois use a forma composta.

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

func get(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestHTTPEntryExpires(t *testing.T) {
	t.Skip("TODO: remova o skip e use httptest.NewTestServer dentro da bubble")

	synctest.Test(t, func(t *testing.T) {
		c := New()
		server := httptest.NewTestServer(t, NewHandler(c))
		client := server.Client()

		req, err := http.NewRequest(http.MethodPut, "http://cache.test/cache/session?ttl=1h", strings.NewReader("active"))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("PUT status = %d, want 204", resp.StatusCode)
		}

		// TODO: faça GET antes, avance uma hora, faça GET depois.
	})
}
