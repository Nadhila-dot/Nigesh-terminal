package server

import (
    "log"
    "net/http"
    "os"
    "path/filepath"
   
)

func StartFileServer() {
    port := os.Getenv("PORT_NIGESH")
    if port == "" {
        port = "8087"
    }

    // Serve private workspace files
    http.Handle("/private/", http.StripPrefix("/private/", http.FileServer(http.Dir(".nigesh/workspace"))))

    // Serve public files
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("nigesh-public"))))

    //log.Printf("Nigesh file server running on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Utility: Move files from workspace to public
func MoveToPublic(files []string) error {
    publicDir := "nigesh-public"
    if _, err := os.Stat(publicDir); os.IsNotExist(err) {
        if err := os.Mkdir(publicDir, 0755); err != nil {
            return err
        }
    }
    for _, file := range files {
        src := filepath.Join(".nigesh/workspace", file)
        dst := filepath.Join(publicDir, filepath.Base(file))
        input, err := os.ReadFile(src)
        if err != nil {
            return err
        }
        if err := os.WriteFile(dst, input, 0644); err != nil {
            return err
        }
    }
    return nil
}