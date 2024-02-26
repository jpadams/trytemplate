// The Foo module's short description.
//
// The Foo module's long description is here.
// It can span multiple lines and provide
// more detail about your module's usage.
//
// Here's another paragraph.

package main

import (
	"context"
	"fmt"
)

// Functions for working with Foo.
type Foo struct{}

// Returns a base Container with some packages installed.
func (m *Foo) Base() *Container {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "bash", "figlet"})
}

// Returns a container that echoes a message and its length using a Bash script
func (m *Foo) EchoLength(message string) *Container {
	bashScript := fmt.Sprintf(`
	declare -A myArray
	myArray["msg"]="%s"
	myArray["len"]=%d
	echo "Message: ${myArray["msg"]}"
	echo "Length: ${myArray["len"]}"
	`, message, len(message))
	return m.Base().WithExec([]string{"bash", "-c", bashScript})
}

// Returns a message as a large banner and optionally provides metadata
func (m *Foo) Big(
	ctx context.Context,
	// +optional
	// +default="Hello, Dagger!"
	message string,
	// +optional
	// +default=false
	verbose bool,
) string {
	output, _ := m.Base().WithExec([]string{"figlet", message}).Stdout(ctx)
	if verbose {
		metadata, _ := m.EchoLength(message).Stdout(ctx)
		output = output + metadata
	}
	return output
}

// Runs grep using pattern in a container over provided Directory
func (m *Foo) GrepDir(ctx context.Context, directory *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directory).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-Rn", pattern, "."}).
		Stdout(ctx)
}

// Takes a message and creates a Directory with a banner and metadata Files
func (m *Foo) Bundle(ctx context.Context, message string) *Directory {
	metadata, _ := m.EchoLength(message).Stdout(ctx)
	return dag.Directory().
		WithNewFile("/banner", m.Big(ctx, message, false)).
		WithNewFile("/metadata", metadata)
}

// Searches over a Directory of files created using message returning result
func (m *Foo) Search(ctx context.Context, message string) string {
	dir := m.Bundle(ctx, message)
	result, _ := m.GrepDir(ctx, dir, message)
	return result
}
