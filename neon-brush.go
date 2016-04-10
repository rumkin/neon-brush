package main

import (
  "fmt"
  "encoding/hex"
  "crypto/sha256"
  "github.com/agl/ed25519"
  "github.com/agl/ed25519/edwards25519"
  "os"
  "golang.org/x/crypto/ssh/terminal"
  "bufio"
  "path"
  "io/ioutil"
  "flag"
)

// Read line from stdin
func ReadLine(question string) string {
  fmt.Print(question)
  reader := bufio.NewReader(os.Stdin)
  answer, _ := reader.ReadString('\n')

  return answer[0:len(answer) - 1]
}

// Read password from stdin
func ReadPassword(question string) string {
  fmt.Print(question)
  password, err := terminal.ReadPassword(0)
  if err != nil {
    panic(err)
  }

  return string(password)
}

// Generate key from custom bytes
func GenerateKey(bytes []byte) (publicKey *[ed25519.PublicKeySize]byte, privateKey *[ed25519.PrivateKeySize]byte, err error) {
	privateKey = new([64]byte)
	publicKey = new([32]byte)
	h := sha256.New()
	h.Write(bytes)
	digest := h.Sum(nil)
  copy(privateKey[:32], digest)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest)
	edwards25519.GeScalarMultBase(&A, &hBytes)
	A.ToBytes(publicKey)

	copy(privateKey[32:], publicKey[:])
	return
}

func writeFile(filepath string, content []byte) (error) {
  f, err := os.Create(filepath)
  if err != nil {
    return err
  }

  defer f.Close()

  err = ioutil.WriteFile(
    filepath,
    []byte(hex.EncodeToString(content)),
    0600,
  )

  if err != nil {
    return err
  }

  return nil
}


func main() {
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n\n1. Run key generator.\n2. Enter key file location (or skip empty for default).\n3. Enter username.\n4. Enter password\n5. Get key\n", path.Base(os.Args[0]))
    flag.PrintDefaults()
  }

  flag.Parse()

  filename := ReadLine("Key file path (default: .nbs/key): ")
  if len(filename) < 1 {
    filename = ".nbs/key"
  }

  username := ReadLine("Enter username: ")
  if len(username) < 1 {
    fmt.Fprintln(os.Stderr, "Username is empty")
    os.Exit(1)
  }

  password := ReadPassword("Enter password: ")
  fmt.Println("")

  str := fmt.Sprintf("%s|%s|%d", username,  password, len(username) + len(password))
  public, private, err := GenerateKey([]byte(str))

  if err != nil {
    panic(err)
  }

  home := os.Getenv("HOME")
  dir := path.Dir(path.Join(home, filename))
  err = os.MkdirAll(dir, 0700)
  if err != nil {
    panic(err)
  }

  err = writeFile(
    path.Join(home, filename + ".pub"),
    public[:],
  )

  if err != nil {
    panic(err)
  }

  err = writeFile(
    path.Join(home, filename),
    private[:],
  )

  if err != nil {
    panic(err)
  }
}
