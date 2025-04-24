// common - generate utiles
package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	// ErrCommonInvalidMedia - wrong media type in Request
	ErrCommonInvalidMedia = errors.New("unexpected media type")

	// ErrCookieEmptyKey - return if value by key in cookie is empty
	ErrCookieEmptyKey = errors.New("empty cookie key")

	// ErrCommonEmptyBody - use in DecodeJSON see below
	ErrCommonEmptyBody = errors.New("request body is empty")
)

// Message - body format for response
type Message map[string]any

//	String - many keys in message
//
// sort keys - result was predictable
func (m Message) String() string {
	lineMSG := make([]string, 0, len(m))
	for k, v := range m {
		lineMSG = append(lineMSG, fmt.Sprintf(`{%s:%v}`, k, v))
	}
	sort.Strings(lineMSG)
	return strings.Join(lineMSG, ",")
}

// MessageError - format of error for Response
type MessageError struct {
	ErrLine string `json:"error"`
}

func NewError(err error) MessageError {
	return MessageError{ErrLine: err.Error()}
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
	onlyDir := strings.Replace(fullPath, fileName, "", -1)
	if err := os.MkdirAll(onlyDir, 0o755); err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	return file.Close()
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
	if r.Body == nil {
		return ErrCommonEmptyBody
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("common: r.Body.Close error - %v", err)
		}
	}()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(obj)
}

// EncodeJSON - we write the status and the object type of 'json' to 'ResponseWriter'
func EncodeJSON(w http.ResponseWriter, httpCode int, obj any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		log.Printf("common: json.Encode error - %v", err)
	}
}

func BeginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// ReduceTimeToDay - yaer,month,day
func ReduceTimeToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// HashPData - use 'sha256.Sum256' for hashing string line
func HashData(line string) string {
	hashLine := sha256.Sum256([]byte(line))
	return hex.EncodeToString(hashLine[:])
}

// ReadCookie - return value from Cookie by key
func ReadCookie(r *http.Request, key string) (string, error) {
	if key == "" {
		return "", ErrCookieEmptyKey
	}
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(cookie.Value)
}

// CleanCookie - set all cookies -> MaxAge = -1
func CleanCookie(w http.ResponseWriter, r *http.Request) {
	for _, val := range r.Cookies() {
		c := http.Cookie{
			Value:  val.Name,
			MaxAge: -1,
		}
		http.SetCookie(w, &c)
	}
}
