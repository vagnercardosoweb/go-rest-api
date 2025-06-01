package utils

import (
	"testing"
)

func TestValidateSourceFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errType error
	}{
		{
			name:    "valid go file in project",
			path:    "/path/to/go-rest-api/main.go",
			wantErr: false,
		},
		{
			name:    "valid go runtime file",
			path:    "/usr/local/go/src/runtime/proc.go",
			wantErr: false,
		},
		{
			name:    "valid go mod file",
			path:    "/go/pkg/mod/github.com/gin-gonic/gin@v1.9.1/gin.go",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errType: ErrInvalidFilePath,
		},
		{
			name:    "path traversal attempt",
			path:    "../../../etc/passwd",
			wantErr: true,
			errType: ErrPathTraversal,
		},
		{
			name:    "non-go file",
			path:    "/path/to/file.txt",
			wantErr: true,
			errType: ErrInvalidFilePath,
		},
		{
			name:    "invalid path outside allowed directories",
			path:    "/etc/shadow.go",
			wantErr: true,
			errType: ErrInvalidFilePath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSourceFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSourceFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil && err != tt.errType {
				t.Errorf("ValidateSourceFilePath() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidateDownloadFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errType error
	}{
		{
			name:    "valid tmp path",
			path:    "/tmp/download.txt",
			wantErr: false,
		},
		{
			name:    "valid relative tmp path",
			path:    "./tmp/download.txt",
			wantErr: false,
		},
		{
			name:    "valid downloads path",
			path:    "./downloads/file.pdf",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errType: ErrInvalidFilePath,
		},
		{
			name:    "path traversal attempt",
			path:    "../../../etc/passwd",
			wantErr: true,
			errType: ErrPathTraversal,
		},
		{
			name:    "invalid path outside allowed directories",
			path:    "/etc/shadow",
			wantErr: true,
			errType: ErrInvalidFilePath,
		},
		{
			name:    "path traversal in allowed directory",
			path:    "/tmp/../etc/passwd",
			wantErr: true,
			errType: ErrPathTraversal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDownloadFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDownloadFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil && err != tt.errType {
				t.Errorf("ValidateDownloadFilePath() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "normal filename",
			filename: "document.pdf",
			want:     "document.pdf",
		},
		{
			name:     "filename with path traversal",
			filename: "../../../etc/passwd",
			want:     "______etc_passwd",
		},
		{
			name:     "filename with dangerous characters",
			filename: "file<>:\"|?*.txt",
			want:     "file_______.txt",
		},
		{
			name:     "filename with backslashes",
			filename: "path\\to\\file.txt",
			want:     "path_to_file.txt",
		},
		{
			name:     "very long filename",
			filename: string(make([]byte, 300)),
			want:     string(make([]byte, 255)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFileName(tt.filename)
			if got != tt.want {
				t.Errorf("SanitizeFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidGoSourcePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "go runtime path",
			path: "/usr/local/go/src/runtime/proc.go",
			want: true,
		},
		{
			name: "project path",
			path: "/path/to/go-rest-api/main.go",
			want: true,
		},
		{
			name: "go mod path",
			path: "/go/pkg/mod/github.com/gin-gonic/gin@v1.9.1/gin.go",
			want: true,
		},
		{
			name: "invalid system path",
			path: "/etc/passwd",
			want: false,
		},
		{
			name: "invalid user path",
			path: "/home/user/malicious.go",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidGoSourcePath(tt.path)
			if got != tt.want {
				t.Errorf("isValidGoSourcePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
