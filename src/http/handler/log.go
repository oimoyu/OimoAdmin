package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oimoyu/OimoAdmin/src/utils/g"
	"github.com/oimoyu/OimoAdmin/src/utils/restful"
	"io"
	"os"
	"strings"
)

func readLastNLines(filePath string, n int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}
	fileSize := fileStat.Size()

	var lines []string
	var currentPos int64 = fileSize
	var lineBuffer []byte

	buf := make([]byte, 1) // Buffer to store file content read in reverse

	for currentPos > 0 && len(lines) < n {
		currentPos--

		if _, err := file.Seek(currentPos, io.SeekStart); err != nil { // Seek to the current position from the end of the file
			return nil, fmt.Errorf("failed to seek log file: %w", err)
		}

		if _, err := file.Read(buf); err != nil { // Read one byte at a time
			return nil, fmt.Errorf("failed to read log file: %w", err)
		}

		if buf[0] == '\n' { // If we encounter a newline, prepend the current line to lines slice
			lines = append([]string{string(lineBuffer)}, lines...)
			lineBuffer = nil
		} else {

			lineBuffer = append([]byte{buf[0]}, lineBuffer...) // Otherwise, prepend the byte to the line buffer
		}
	}

	if len(lineBuffer) > 0 {
		lines = append([]string{string(lineBuffer)}, lines...) // Add the last line if there is no trailing newline
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:] // Trim lines to the last n lines
	}

	// reverse
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines, nil
}

func Log(c *gin.Context) {
	logLength := 100

	logFilePath := g.OimoAdmin.Logger.FileLogPath
	lines, err := readLastNLines(logFilePath, logLength)
	if err != nil {
		restful.ParamErr(c, fmt.Sprintf("read log file failed: %v", err))
		return
	}
	items := make([]map[string]interface{}, 0)
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) < 2 {
			continue
		}
		items = append(items, map[string]interface{}{
			"create_time": parts[0],
			"msg":         parts[1],
		})
	}

	if len(lines) >= logLength-1 {
		items = append(items, map[string]interface{}{
			"create_time": "...",
			"msg":         fmt.Sprintf("more log check: %s", g.OimoAdmin.Logger.FileLogPath),
		})
	}

	returnData := map[string]interface{}{
		"items": items,
	}

	restful.Ok(c, returnData)
}
