package functions

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func GenerateToken(dataMap map[string]interface{}, secretString string, expireTime time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range dataMap {
		claims[key] = value
	}
	claims["exp"] = time.Now().Add(expireTime).Unix()

	secretByte := []byte(secretString)
	tokenString, err := token.SignedString(secretByte)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func DecodeToken(tokenString string, secretString string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return []byte(secretString), nil })
	result := make(map[string]interface{})
	if err != nil {
		return result, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		for key, value := range claims {
			result[key] = value
		}
		return result, err
	} else {
		return result, err
	}
}

func GetRootDir() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	rootDir := filepath.Dir(executable)
	return rootDir
}

func GetRuntimeDir() string {
	rootDir := GetRootDir()
	runtimeDir := filepath.Join(rootDir, "admin_runtime")
	return runtimeDir
}
func GetEnvDir() string {
	rootDir := GetRuntimeDir()
	envDir := filepath.Join(rootDir, ".env")
	return envDir
}
func GetWWWRootDir() string {
	rootDir := GetRuntimeDir()
	wwwRootDir := filepath.Join(rootDir, "wwwroot")
	return wwwRootDir
}

func EnsureDirExists(dir string, options ...interface{}) error {
	perm := os.ModePerm // default permissions
	for _, option := range options {
		switch option.(type) {
		case os.FileMode:
			perm = option.(os.FileMode)
		default:
			return fmt.Errorf("input options is not os.FileMode, rather: %T", option)
		}
	}

	if isPathExists := IsPathExists(dir); !isPathExists {
		err := os.MkdirAll(dir, perm)
		if err != nil {
			return err
		}
	}
	return nil
}
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func DownloadAndExtract(url, destDir string, tempFilePath string) error {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Get File Size
	totalSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// using fix temp path avoid garbage
	tmpFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to create or open temp file: %w", err)
	}
	defer os.Remove(tempFilePath)
	defer tmpFile.Close()

	// Create a TeeReader to echo progress
	progress := &ProgressWriter{Total: int64(totalSize)}
	tee := io.TeeReader(resp.Body, progress)

	// Download Content to temp IO File
	if _, err := io.Copy(tmpFile, tee); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// ensure dist dir exist before unzip
	if err := EnsureDirExists(destDir); err != nil {
		return err
	}

	// Unzip gz
	if err := ExtractTarGz(tmpFile.Name(), destDir); err != nil {
		os.RemoveAll(destDir)
		return fmt.Errorf("failed to extract tar.gz file: %w", err)
	}

	return nil
}

type ProgressWriter struct {
	Total   int64
	Current int64
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Current += int64(n)
	pw.printProgress()
	return n, nil
}

func (pw *ProgressWriter) printProgress() {
	fmt.Printf("\rDownloading... %d%% \r", pw.Current*100/pw.Total)
}

func ExtractTarGz(gzPath string, destDir string) error {
	gzipStream, err := os.Open(gzPath)
	if err != nil {
		return fmt.Errorf("failed to open gz file: %w", err)
	}
	defer gzipStream.Close()

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			// unzip done
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// need to manually create sub dir
			if err := EnsureDirExists(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close() // close resource when err
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close() // close resource when no err

		default:
			//do nothing when unknown header
		}

	}

	return nil
}
func CreateFileWithContent(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := EnsureDirExists(dir); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
func ParseTemplate(templateContent string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	content := buf.String()

	return content, nil
}

func FormatJsonString(jsonString string) (string, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return "", err
	}

	formattedJSON, err := json.MarshalIndent(jsonData, "", "	")
	if err != nil {
		return "", err
	}
	return string(formattedJSON), nil
}
func IsMatchInSlice(totalString string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(totalString, substring) {
			return true
		}
	}
	return false
}
func IsStringUUID(s string) bool {
	// check is uuid without lib
	re := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	return re.MatchString(s)
}
func IsStringNumeric(s string) bool {
	re := regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)$`)
	return re.MatchString(s)
}
func IsStringInt(s string) bool {
	re := regexp.MustCompile(`^[+-]?\d+$`)
	return re.MatchString(s)
}
