package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"golang.org/x/net/context"
)

func logRequest(ctx context.Context, w io.Writer, req *http.Request) {
	fmt.Fprintln(w, "")
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOISILON HTTP REQUEST")
	fmt.Fprintln(w, " -------------------------")
	buf, err := httputil.DumpRequest(req, true)
	if err != nil {
		return
	}
	WriteIndented(w, buf)
	fmt.Fprintln(w)
}

func logResponse(ctx context.Context, res *http.Response) {
	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOISILON HTTP RESPONSE")
	fmt.Fprintln(w, " -------------------------")

	buf, err := httputil.DumpResponse(
		res,
		res.Header.Get("Content-Type") != "application/octet-stream")
	if err != nil {
		return
	}

	bw := &bytes.Buffer{}
	WriteIndented(bw, buf)

	scanner := bufio.NewScanner(bw)
	for {
		if !scanner.Scan() {
			break
		}
		fmt.Fprintln(w, scanner.Text())
	}

	io.Copy(os.Stdout, w)
}

// WriteIndentedN indents all lines n spaces.
func WriteIndentedN(w io.Writer, b []byte, n int) error {
	s := bufio.NewScanner(bytes.NewReader(b))
	if !s.Scan() {
		return nil
	}
	l := s.Text()
	for {
		for x := 0; x < n; x++ {
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, l); err != nil {
			return err
		}
		if !s.Scan() {
			break
		}
		l = s.Text()
		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

// WriteIndented indents all lines four spaces.
func WriteIndented(w io.Writer, b []byte) error {
	return WriteIndentedN(w, b, 4)
}
