// common - generate utiles
package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ErrCommonInvalidMedia - wrong media type in Request
var ErrCommonInvalidMedia = errors.New("unexpected media type")

// Message - body format for response
type Message map[string]any

func (m Message) String() string {
	buff := &bytes.Buffer{}
	for k, v := range m {
		_, _ = fmt.Fprintf(buff, `{%s : %v},`, k, v)
	}
	if n := buff.Len(); n != 0 {
		buff.Truncate(n - 1)
	}
	return buff.String()
}

// ScanSQL - generic template for scan object from sql database
//
// example: func scanSomeObject[T ScanSQL](r T) (obj,err)
// T - can be '*sql.Row', '*sql.Rows' ets.
type ScanSQL interface {
	Scan(dest ...any) error
}

// CreatePathWithFile - create a file and all directories to it if they do not exist
func CreatePathWithFile(partOfFilePath string) error {
	fileName := filepath.Base(partOfFilePath)
	if fileName == "" {
		return errors.New("common: incorrect path of file")
	}
	if fileExtension := filepath.Ext(fileName); fileExtension != ".db" {
		return errors.New("common: invalid file extension")
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(currentDir, partOfFilePath)
	onlyDir := strings.Replace(fullPath, fileName, "", 1)
	if err := os.MkdirAll(onlyDir, 0o755); err != nil {
		return err
	}
	_, err = os.Create(fullPath)
	return err
}

// Abs - absolute value
func Abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

// DecodeJSON - get object from 'Request'
// check: media type,
// struct fields -> json.DisallowUnknownFields()
func DecodeJSON(r *http.Request, obj any) error {
	media := r.Header.Get("Content-Type")
	parse, _, err := mime.ParseMediaType(media)
	if err != nil || parse != "application/json" {
		return ErrCommonInvalidMedia
	}
	reqBody := r.Body
	defer func() {
		if err := reqBody.Close(); err != nil {
			log.Printf("common: r.Body.Close error - %v", err)
		}
	}()
	dec := json.NewDecoder(reqBody)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

// EncodeJSON - we write the status and the object type of 'json' to 'ResponseWriter'
//
// context Deadline not null - set status 408
func EncodeJSON(ctx context.Context, w http.ResponseWriter, httpCode int, obj any) {
	if ctx.Err() == context.DeadlineExceeded {
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("common: json.Encode error - %v", err)
	}
}

func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ReduceTimeToDay - yaer,month,day
func ReduceTimeToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
