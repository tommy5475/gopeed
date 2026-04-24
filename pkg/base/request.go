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
	// I bumped the default to 8 in the server config; this field just holds the override
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
		// Increased from 30s to 60s to better handle slow or congested servers
		// on my home network this was timing out too often on large file servers
		ResponseHeaderTimeout: 60 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}
