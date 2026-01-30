package server

import (
    "os"
    "path/filepath"
)

// Checks if .nigesh/workspace exists and has files
func HasWorkspaceFiles() bool {
    workspace := ".nigesh/workspace"
    files, err := os.ReadDir(workspace)
    if err != nil || len(files) == 0 {
        return false
    }
    return true
}

// Checks if nigesh-public exists in current directory and has files
func HasPublicFiles() bool {
    public := "nigesh-public"
    files, err := os.ReadDir(public)
    if err != nil || len(files) == 0 {
        return false
    }
    return true
}

// Returns all file paths in .nigesh/workspace as a []string
func ListWorkspaceFiles() ([]string, error) {
    return listFilesRecursive(".nigesh/workspace")
}

// Returns all file paths in nigesh-public as a []string
func ListPublicFiles() ([]string, error) {
    return listFilesRecursive("nigesh-public")
}

// Helper to recursively list all files in a directory
func listFilesRecursive(root string) ([]string, error) {
    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}