package wx

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
StreamingFile serves a file over HTTP with support for byte-range requests.

It performs the following:
- Opens the requested file and retrieves its size and MIME type.
- Sets appropriate HTTP headers:
  - Content-Type: the file's MIME type.
  - Cache-Control: allows caching by clients for 24 hours.
  - Accept-Ranges: indicates support for partial content requests.

- Checks if the client sent a "Range" header to request a partial file.
  - If no Range header is present, it streams the entire file with status 200 OK.
  - If a valid Range header is present, it parses the range and streams only
    the requested byte range with status 206 Partial Content.
  - Reads the file in 32KB chunks and writes them to the response,
    flushing the output buffer after each chunk to enable streaming.
  - Handles errors and properly closes the file.

This function enables efficient serving of large files and supports resumable downloads.
*/
func (ctx *httpContext) StreamingFile(fileName string) error {
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		http.Error(ctx.Res, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(ctx.Res, "Cannot read file info", http.StatusInternalServerError)
		return err
	}

	fileSize := stat.Size()

	ext := strings.ToLower(filepath.Ext(fileName))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType = http.DetectContentType(buf[:n])
		// Reset offset sau khi đọc
		_, _ = file.Seek(0, io.SeekStart)
	}

	ctx.Res.Header().Set("Content-Type", contentType)
	ctx.Res.Header().Set("Cache-Control", "public, max-age=86400")
	ctx.Res.Header().Set("Accept-Ranges", "bytes")

	rangeHeader := ctx.Req.Header.Get("Range")
	flusher, _ := ctx.Res.(http.Flusher)

	if rangeHeader == "" {
		// Trả toàn bộ file, status 200
		ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		ctx.Res.WriteHeader(http.StatusOK)

		buf := make([]byte, 32*1024) // 32KB
		for {
			n, err := file.Read(buf)
			if n > 0 {
				_, writeErr := ctx.Res.Write(buf[:n])
				if writeErr != nil {
					return writeErr
				}
				if flusher != nil {
					flusher.Flush()
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
		return nil
	}

	// Parse Range header
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		http.Error(ctx.Res, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return err
	}

	contentLength := end - start + 1

	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	ctx.Res.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	ctx.Res.Header().Set("ETag", fmt.Sprintf("%x-%d", stat.ModTime().Unix(), fileSize))
	ctx.Res.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))
	ctx.Res.WriteHeader(http.StatusPartialContent)

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, 32*1024) // 32KB
	remaining := contentLength
	for remaining > 0 {
		readSize := int64(len(buf))
		if remaining < readSize {
			readSize = remaining
		}

		n, err := file.Read(buf[:readSize])
		if n > 0 {
			_, writeErr := ctx.Res.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			if flusher != nil {
				flusher.Flush()
			}
			remaining -= int64(n)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

/*
parseRange parses a Range header of the form "bytes=start-end" and returns
the start and end byte positions as int64.

Note:
- The Range header might be like "bytes=0-" (from start to end of file)
- Or "bytes=12345-" (from position 12345 to end of file)
- The function validates the range according to the fileSize.
- It returns an error if the range is invalid or malformed.

Parameters:
- rangeHeader: the value of the Range header from the HTTP request.
- fileSize: the total size of the file in bytes.

Returns:
- start: the starting byte position of the requested range.
- end: the ending byte position of the requested range.
- error: non-nil if the range is invalid or cannot be parsed.
*/
func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	const prefix = "bytes="
	// Check if the header starts with "bytes="
	if !strings.HasPrefix(rangeHeader, prefix) {
		return 0, 0, fmt.Errorf("invalid range header")
	}

	// Remove the "bytes=" prefix
	r := strings.TrimPrefix(rangeHeader, prefix)
	// Split into start and end parts, e.g. "123-456" -> ["123", "456"]
	items := strings.SplitN(r, "-", 2)
	if len(items) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	// Parse start position
	start, errStart := strconv.ParseInt(strings.TrimSpace(items[0]), 10, 64)
	// Parse end position
	end, errEnd := strconv.ParseInt(strings.TrimSpace(items[1]), 10, 64)

	if errStart != nil {
		// If start is not a valid number, return error
		// (Note: suffix ranges like "bytes=-500" are not supported here)
		return 0, 0, fmt.Errorf("invalid start range")
	}

	if errEnd != nil {
		// If end is missing or invalid (e.g. "bytes=12345-"),
		// treat it as the last byte of the file
		end = fileSize - 1
	}

	// Validate range boundaries
	if start < 0 || start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("invalid range values")
	}

	return start, end, nil
}

/*
DownloadFile streams the specified file to the HTTP client as a downloadable attachment.

Parameters:
- w: http.ResponseWriter to write the HTTP response to the client.
- fileName: the path to the file on disk to be downloaded.

Behavior:
- Opens the specified file and returns an error if it cannot be opened.
- Determines the Content-Type based on the file extension.
- Sets headers to prompt the client to download the file with the original filename.
- Sets Content-Length header according to the file size.
- Reads the file in chunks (32KB) and writes to the response.
- Flushes the response buffer if supported to enable streaming.
- Automatically closes the file after finishing or upon error.

Returns:
- error: non-nil if any error occurs during file opening, reading, or writing.
*/
func (ctx *httpContext) DownloadFile(fileName string) error {
	// Open the file for reading
	file, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		http.Error(ctx.Res, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	// Get file information to obtain size
	stat, err := file.Stat()
	if err != nil {
		http.Error(ctx.Res, "Cannot read file info", http.StatusInternalServerError)
		return err
	}
	fileSize := stat.Size()

	// Determine Content-Type based on file extension
	ext := strings.ToLower(filepath.Ext(fileName))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// Fallback: read first 512 bytes to detect content type
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType = http.DetectContentType(buf[:n])
		_, _ = file.Seek(0, io.SeekStart) // Reset offset after reading
	}

	// Set response headers for download
	ctx.Res.Header().Set("Content-Type", contentType)
	ctx.Res.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(fileName)))
	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
	ctx.Res.Header().Set("Cache-Control", "no-cache")

	// Try to get Flusher interface to flush response buffer
	flusher, canFlush := ctx.Res.(http.Flusher)

	// Buffer size for streaming
	buf := make([]byte, 32*1024) // 32KB buffer

	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, writeErr := ctx.Res.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			if canFlush {
				flusher.Flush()
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}
func (h *httpContext) GetAbsRootUri() string {
	if h.rootAbsUrl != "" {
		return h.rootAbsUrl
	}
	scheme := h.GetScheme()

	h.rootAbsUrl = scheme + "://" + h.Req.Host
	return h.rootAbsUrl
}
func (h *httpContext) GetScheme() string {
	if h.schema != "" {
		return h.schema
	}
	h.schema = "http"

	// Trường hợp Go server trực tiếp nhận HTTPS
	if h.Req.TLS != nil {
		return "https"
	}

	// Trường hợp có reverse proxy thêm header
	if proto := h.Req.Header.Get("X-Forwarded-Proto"); proto != "" {
		h.schema = strings.ToLower(proto)
	}

	if forwarded := h.Req.Header.Get("Forwarded"); forwarded != "" {
		// Ví dụ: "for=192.0.2.43; proto=https; by=203.0.113.43"
		for _, part := range strings.Split(forwarded, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "proto=") {
				h.schema = strings.TrimPrefix(part, "proto=")
			}
		}
	}

	// Mặc định là http
	return h.schema
}
