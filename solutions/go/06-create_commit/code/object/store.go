package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func StoreFromFile(name, typ string) (Hash, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return Hash{}, fmt.Errorf("read src file: %w", err)
	}

	return Store(bytes.NewReader(content), typ, int64(len(content)))
}

func Store(src io.Reader, typ string, size int64) (Hash, error) {
	var buf bytes.Buffer

	err := encodeObject(&buf, src, typ, size)

	fileContent, err := compress(buf.Bytes())
	if err != nil {
		return Hash{}, fmt.Errorf("compress: %w", err)
	}

	sum := sha1.Sum(buf.Bytes())
	name := hex.EncodeToString(sum[:])

	objPath := filepath.Join(".git", "objects", name[:2], name[2:])
	dirPath := filepath.Dir(objPath)

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return Hash{}, fmt.Errorf("mkdir: %w", err)
	}

	err = os.WriteFile(objPath, fileContent, 0644)
	if err != nil {
		return Hash{}, fmt.Errorf("write file: %w", err)
	}

	return Hash(sum), nil
}

func encodeObject(dst io.Writer, src io.Reader, typ string, size int64) error {
	_, err := fmt.Fprintf(dst, "%v %d\000", typ, size)
	if err != nil {
		return err
	}

	n, err := io.Copy(dst, src)
	if err != nil {
		return err
	}

	if n != size {
		return fmt.Errorf("file size mismatch, got %v, content %v", size, n)
	}

	return nil
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
