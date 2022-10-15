package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/git-starter-go/object"
)

func lsTreeCmd(args []string) (err error) {
	// Assuming that args is ["ls-tree", "--name-only", "hash"], just like os.Args

	if len(args) < 3 || args[1] != "--name-only" {
		fmt.Fprintf(os.Stderr, "usage: mygit ls-tree --name-only <hash>\n")

		return fmt.Errorf("bad usage")
	}

	hash, err := object.HashFromString(args[2])
	if err != nil {
		return fmt.Errorf("%w: %v", err, args[2])
	}

	typ, content, err := object.LoadByHash(hash)
	if err != nil {
		return fmt.Errorf("load by hash: %w", err)
	}

	if typ != "tree" {
		return fmt.Errorf("unsupported object type: %v", typ)
	}

	return lsTree(content)
}

func lsTree(d []byte) error {
	r := bufio.NewReader(bytes.NewReader(d))

	for {
		// row format: 100644<space>file_name.txt<null_byte><20_byte_sha1>

		_, err := r.ReadString(' ') // mode
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		name, err := r.ReadString('\000')
		if err != nil {
			return err
		}

		name = name[:len(name)-1] // cut delimiter

		_, err = r.Discard(sha1.Size) // sha1
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", name)
	}

	return nil
}
