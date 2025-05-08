package suprsend

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type AttachmentOption struct {
	// overrides filename of attachment, otherwise filename is picked from the filepath
	FileName string
	// ignore attachment if there is issue while accessing/downloading attachment from url
	// applicable when filepath is a url.
	IgnoreIfError bool
}

func GetAttachmentJson(filePath string, ao *AttachmentOption) (map[string]any, error) {
	fileName, ignoreIfError := "", false
	if ao != nil {
		fileName, ignoreIfError = ao.FileName, ao.IgnoreIfError
	}
	//
	if checkIsUrl(filePath) {
		return getAttachmentJsonForUrl(filePath, fileName, ignoreIfError)
	} else {
		return getAttachmentJsonForFile(filePath, fileName, ignoreIfError)
	}
}

func checkIsUrl(filePath string) bool {
	for _, scheme := range []string{"https://", "http://"} {
		if strings.HasPrefix(filePath, scheme) {
			return true
		}
	}
	return false
}

func getAttachmentJsonForUrl(fileUrl string, fileName string, ignoreIfError bool) (map[string]any, error) {
	return map[string]any{
		"filename":        fileName,
		"url":             fileUrl,
		"ignore_if_error": ignoreIfError,
	}, nil
}

func getAttachmentJsonForFile(filePath string, fileName string, ignoreIfError bool) (map[string]any, error) {
	// Get absolute path
	absPath, err := expandHomeDir(filePath)
	if err != nil {
		if ignoreIfError {
			log.Println("WARNING: ignoring error while processing attachment file.", err)
			return nil, nil
		}
		return nil, err
	}
	// Finalize file name
	finalFileName := filepath.Base(absPath)
	fileName = strings.TrimSpace(fileName)
	if fileName != "" {
		finalFileName = fileName
	}
	// extract content and mime-type
	content, err := os.ReadFile(absPath)
	if err != nil {
		if ignoreIfError {
			log.Println("WARNING: ignoring error while processing attachment file.", err)
			return nil, nil
		}
		return nil, err
	}
	b64Str := base64.StdEncoding.EncodeToString(content)
	mimeType := mimetype.Detect(content).String()
	//
	return map[string]any{
		"filename":        finalFileName,
		"contentType":     mimeType,
		"data":            b64Str,
		"ignore_if_error": ignoreIfError,
	}, nil
}

// expandHomeDir expands file paths relative to the user's home directory (~) into absolute paths.
func expandHomeDir(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if path == "~" {
		return homeDir, nil
	}
	if !strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		return "", &Error{Message: "cannot expand user-specific home dir"}
	}
	return filepath.Join(homeDir, strings.TrimPrefix(path, "~")), nil
}
