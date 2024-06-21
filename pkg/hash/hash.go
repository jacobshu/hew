package hash

import (
  "crypto/sha256"
  "math/rand"
  "fmt"
  "io"
  "log"
  "os"
)

func Hash(file string, outputFile string) {
  f, err := os.Open(file)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  h := sha256.New()
  if _, err := io.Copy(h, f); err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%x", h.Sum(nil))
}

func createTestFiles(num int, outDir string) {
  for i := 0; i < num; i++ {
    f := fmt.Sprintf("%04d", i) 
    c := randomString(16)

    err := os.WriteFile(f, []byte(c), 0644)
    if err != nil {
      panic(err)
    }
  }
}

func randomString(n int) string {
  const letters = "abcdefghijklmnopqrstuvwxyz"
  b := make([]byte, n)
  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }

  return string(b)
}

