package mitm

import (
	"bytes"
	"io"
	"testing"
	"time"
)

// This test assumes a mitm-cache is running at http://localhost:8000 with key "abc"
func TestRequest(t *testing.T) {
	c := New("http://localhost:8000", "abc")

	r1, err1 := c.Request("https://httpbin.org/uuid", 0) // ensure r1 is new uuid
	r2, err2 := c.Request("https://httpbin.org/uuid", 2) // r2 should equal r1 (as long as no one else invalidated)

	time.Sleep(time.Second * 2)
	r3, err3 := c.Request("https://httpbin.org/uuid", 1) // r2 is now 2s old, r3 must be different

	for i, err := range []error{err1, err2, err3} {
		if err != nil {
			t.Fatalf("Request %d failed with message: %v", i, err)
		}
	}

	r := make([]string, 3)

	for i, reader := range []io.ReadCloser{r1, r2, r3} {
		b := new(bytes.Buffer)
		io.Copy(b, reader)
		r[i] = b.String()
		t.Logf("Request %d has content %s", i, r[i])
	}

	if r[0] != r[1] {
		t.Error("In-date cache not used?")
	}

	if r[0] == r[2] || r[1] == r[2] {
		t.Error("Cache kept beyond maxage")
	}
}
