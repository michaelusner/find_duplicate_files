// Copyright Michael Usner - 2017
// Yes, comments are scarce - I'm still learning go.
// The basic idea is:
// Iterate all files in the folder and map by size.
// Then, if there are multiple files that map to the same size, calculate the sha256 hash.
// Any matching hashes are duplicates (with pretty high certainty).
// We then ask the user if they want to keep the dupe.
// TODO - I think I need to consider the possiblity of more than one dupe in my loop.
package main

import (
    "github.com/minio/sha256-simd"
    "path/filepath"
    "os"
    "log"
    "io"
    "flag"
    "fmt"
    "encoding/hex"
)

type FileInfo struct {
    path string
    hash string
}

var files map[int64][]string

func GetHash(path string) string {
    file, err := os.Open(path)
    if (err != nil) {
        log.Printf("Failed to open %s", path)
    }
    defer file.Close()
    shaWriter := sha256.New()
    io.Copy(shaWriter, file)
    return hex.EncodeToString(shaWriter.Sum(nil))
}

func visit(path string, f os.FileInfo, err error) error {
    if (!f.IsDir()) {
        files[f.Size()] = append(files[f.Size()], path)
    } else {
        fmt.Printf("%s\n", path)
    }
    return nil
} 


func main() {
    files = make(map[int64][]string)
        
    flag.Parse()
    root := flag.Arg(0)
    err := filepath.Walk(root, visit)
    if (err != nil) {
        fmt.Print("Failed to walk path")
    }

    for _, paths := range files {
        if (len(paths) > 1) {
            hashes := make(map[string]string)
            for _, path := range(paths) {
                //fmt.Println(path)
                hash := GetHash(path)
                if val, ok := hashes[hash]; ok {
                    fmt.Println("a: ", val)
                    fmt.Println("b: ", path)
                    fmt.Println("Keep A or B")
                    var input string
                    fmt.Scanln(&input)
                    if (input == "a") {
                        fmt.Println("Keeping a")
                        os.Remove(path)
                    } else if (input == "b") {
                        fmt.Println("Keeping b")
                        os.Remove(val)
                    }
                } else {
                    hashes[hash] = path
                }
            }
            
            fmt.Println()
        }
    }

}
