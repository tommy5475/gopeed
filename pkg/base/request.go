package base

import (
	"net/http"
	"time"
)

// Request represents a download request with all necessary metadata.
type Request struct {
	// URL is the target download URL
	URL string `json:"url"`
	// Extra contains protocol-specific extra information
	Extra interface{} `json:"extra,omitempty"`
	// Labels are user-defined key-value pairs for categorization
	Labels map[string]string `json:"labels,omitempty"`
	// Headers are HTTP headers to include in the request
	Headers map[string]string `json:"headers,omitempty"`
}

// Resource represents a downloadable resource resolved from a Request.
type Resource struct {
	// Name is the suggested filename for the resource
	Name string `json:"name"`
	// Size is the total size in bytes; 0 means unknown
	Size int64 `json:"size"`
	// Range indicates whether the server supports range requests
	Range bool `json:"range"`
	// Files contains the list of files in this resource (for multi-file downloads)
	Files []*FileInfo `json:"files"`
	// Hash is the optional checksum of the resource
	Hash string `json:"hash,omitempty"`
}

// FileInfo holds metadata about a single file within a resource.
type FileInfo struct {
	// Name is the filename
	Name string `json:"name"`
	// Path is the relative directory path within the download folder
	Path string `json:"path"`
	// Size is the file size in bytes
	Size int64 `json:"size"`
	// Req is the underlying HTTP request used to fetch this file
	Req *Request `json:"req,omitempty"`
}

// Options holds configuration options for a download task.
type Options struct {
	// Name overrides the resource name if set
	Name string `json:"name,omitempty"`
	// Path is the local directory where files will be saved
	Path string `json:"path"`
	// SelectFiles is a list of file indices to download; empty means all files
	SelectFiles []int `json:"selectFiles,omitempty"`
	// Extra contains protocol-specific options
	Extra interface{} `json:"extra,omitempty"`
	// Connections is the number of concurrent connections per file
	Connections int `json:"connections,omitempty"`
}

// HTTPOptions holds HTTP-specific download options.
type HTTPOptions struct {
	// Method is the HTTP method (default: GET)
	Method string `json:"method,omitempty"`
	// Body is the request body for POST/PUT requests
	Body string `json:"body,omitempty"`
	// AutoTorrent enables automatic torrent detection and switching
	AutoTorrent bool `json:"autoTorrent,omitempty"`
}

// BuildHTTPClient creates an *http.Client configured with sensible defaults.
// The timeout parameter sets the overall request timeout.
func BuildHTTPClient(timeout time.Duration, proxy string) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if proxy != "" {
		// Proxy configuration would be applied here
		// Kept as a hook for future proxy support
		_ = proxy
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

// Status represents the current state of a download task.
type Status int

const (
	StatusReady   Status = iota // Task created but not started
	StatusRunning               // Task is actively downloading
	StatusPause                 // Task has been paused
	StatusWait                  // Task is waiting in queue
	StatusError                 // Task encountered an error
	StatusDone                  // Task completed successfully
)

// String returns a human-readable representation of the Status.
func (s Status) String() string {
	switch s {
	case StatusReady:
		return "ready"
	case StatusRunning:
		return "running"
	case StatusPause:
		return "pause"
	case StatusWait:
		return "wait"
	case StatusError:
		return "error"
	case StatusDone:
		return "done"
	default:
		return "unknown"
	}
}
