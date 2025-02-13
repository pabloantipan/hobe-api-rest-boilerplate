package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	maxTotalSize    = 100 * 1024 // 100KB total log size
	maxFieldSize    = 20 * 1024  // 20KB per field
	truncatedSuffix = "... (truncated)"
	maxCharsOnTrunc = 100 // Maximum characters to show when truncated
)

func ExtractBody(c *gin.Context) (map[string]interface{}, error) {
	contentType := c.ContentType()
	// result := make(map[string]interface{})

	switch {
	case strings.Contains(contentType, "application/json"):
		return ExtractJSONBody(c)

	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		return ExtractFormBody(c)

	case strings.Contains(contentType, "multipart/form-data"):
		return ExtractMultipartBody(c)

	case strings.Contains(contentType, "text/plain"):
		return ExtractTextBody(c)

	default:
		return ExtractRawBody(c)
	}
}

// ExtractJSONBody handles JSON request bodies
func ExtractJSONBody(c *gin.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Read the body
	bodyBytes, err := readBody(c)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, err
		}

		result = truncateMap(result)
	}

	// Restore the body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return result, nil
}

// ExtractFormBody handles URL-encoded form data
func ExtractFormBody(c *gin.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if err := c.Request.ParseForm(); err != nil {
		return nil, err
	}

	// Convert form values to map with size limits
	currentSize := 0
	for key, values := range c.Request.PostForm {
		valueSize := len(key)
		for _, v := range values {
			valueSize += len(v)
		}

		if currentSize+valueSize > maxTotalSize {
			result["_truncated_info"] = "Additional fields omitted due to size limits"
			break
		}

		if len(values) == 1 {
			result[key] = truncateValue(values[0])
		} else {
			result[key] = truncateValue(values)
		}
		currentSize += valueSize
	}

	return result, nil
}

// ExtractMultipartBody handles multipart form data
func ExtractMultipartBody(c *gin.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Parse multipart form with a reasonable max memory
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	// Get form values
	for key, values := range c.Request.MultipartForm.Value {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}

	// Add file information
	files := make(map[string][]string)
	for key, fileHeaders := range c.Request.MultipartForm.File {
		fileNames := make([]string, len(fileHeaders))
		for i, header := range fileHeaders {
			fileNames[i] = header.Filename
		}
		files[key] = fileNames
	}

	if len(files) > 0 {
		result["_files"] = files
	}

	return result, nil
}

// ExtractTextBody handles plain text bodies
func ExtractTextBody(c *gin.Context) (map[string]interface{}, error) {
	bodyBytes, err := readBody(c)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"content": string(bodyBytes),
	}

	// Restore the body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return result, nil
}

// ExtractRawBody handles any other type of body
func ExtractRawBody(c *gin.Context) (map[string]interface{}, error) {
	bodyBytes, err := readBody(c)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"raw_content": string(bodyBytes),
	}

	// Restore the body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return result, nil
}

// readBody is a helper function to read the request body
func readBody(c *gin.Context) ([]byte, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	// Restore the body for later middleware/handlers
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

func truncateValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if len(v) > maxFieldSize {
			return v[:maxCharsOnTrunc] + truncatedSuffix
		}
		return v
	case []interface{}:
		result := make([]interface{}, 0, len(v))
		currentSize := 0
		for _, item := range v {
			itemSize := approximateSize(item)
			if currentSize+itemSize > maxFieldSize {
				result = append(result, "(remaining items truncated)")
				break
			}
			result = append(result, truncateValue(item))
			currentSize += itemSize
		}
		return result
	case map[string]interface{}:
		return truncateMap(v)
	default:
		return v
	}
}

// truncateMap handles map size limiting
func truncateMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	currentSize := 0

	for k, v := range m {
		valueSize := approximateSize(v)
		if currentSize+valueSize > maxTotalSize {
			result["_truncated_info"] = "Additional fields omitted due to size limits"
			break
		}
		result[k] = truncateValue(v)
		currentSize += valueSize
	}
	return result
}

// approximateSize estimates the size of a value in bytes
func approximateSize(v interface{}) int {
	switch val := v.(type) {
	case string:
		return len(val)
	case []interface{}:
		size := 0
		for _, item := range val {
			size += approximateSize(item)
		}
		return size
	case map[string]interface{}:
		size := 0
		for k, v := range val {
			size += len(k) + approximateSize(v)
		}
		return size
	default:
		// Rough approximation for numbers, booleans, etc.
		return 8
	}
}
