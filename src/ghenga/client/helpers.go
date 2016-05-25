package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
)

func dumpHTTP(wr io.Writer, req *http.Request, res *http.Response) {
	if wr == nil {
		return
	}

	fmt.Fprintf(wr, "====== REQUEST ====================================\n")
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintf(wr, "unable to dump http request: %v\n", err)
	}
	wr.Write(dump)

	fmt.Fprintf(wr, "====== RESPONSE ===================================\n")
	dump, err = httputil.DumpResponse(res, true)
	if err != nil {
		fmt.Fprintf(wr, "unable to dump http response: %v\n", err)
	}
	wr.Write(dump)
	fmt.Fprintf(wr, "====== END ========================================\n\n\n")
}

// TestClient returns a fully authenticated client against a server.
func TestClient(t *testing.T, url string, username, password string) *Client {
	client := New(url)

	_, err := client.Login(username, password)
	if err != nil {
		t.Fatalf("login with username %q and password %q failed: %v", username, password, err)
	}

	if os.Getenv("HTTPTRACE") != "" {
		client.TraceOn(os.Stderr)
	}

	return client
}
