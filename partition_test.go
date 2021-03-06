package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPartitionToMaxSize(t *testing.T) {
	cmdArgs := []string{"/usr/bin/foo", "-c"}
	paths := []string{"one", "two", "three", "four"}

	parts := partitionToMaxSize(cmdArgs, paths, 24)
	expected := [][]string{
		append(cmdArgs, "one", "two"),
		append(cmdArgs, "three"),
		append(cmdArgs, "four"),
	}
	assert.Equal(t, expected, parts)
}

func TestPartitionToPackageFileGlobs(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-expand-paths")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cmdArgs := []string{"/usr/bin/foo", "-c"}
	paths := []string{
		filepath.Join(tmpdir, "one"),
		filepath.Join(tmpdir, "two"),
	}
	for _, dir := range paths {
		mkDir(t, dir)
		mkGoFile(t, dir, "other.go")
	}

	parts, err := partitionPathsAsFilesGroupedByPackage(cmdArgs, paths)
	require.NoError(t, err)
	expected := [][]string{
		append(cmdArgs, packagePaths(paths[0], "file.go", "other.go")...),
		append(cmdArgs, packagePaths(paths[1], "file.go", "other.go")...),
	}
	assert.Equal(t, expected, parts)
}

func packagePaths(dir string, filenames ...string) []string {
	paths := []string{}
	for _, filename := range filenames {
		paths = append(paths, filepath.Join(dir, filename))
	}
	return paths
}

func TestPartitionToPackageFileGlobsNoFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-expand-paths")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cmdArgs := []string{"/usr/bin/foo", "-c"}
	paths := []string{filepath.Join(tmpdir, "one"), filepath.Join(tmpdir, "two")}
	parts, err := partitionPathsAsFilesGroupedByPackage(cmdArgs, paths)
	require.NoError(t, err)
	assert.Len(t, parts, 0)
}

func TestPartitionToMaxArgSizeWithFileGlobsNoFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "test-expand-paths")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cmdArgs := []string{"/usr/bin/foo", "-c"}
	paths := []string{filepath.Join(tmpdir, "one"), filepath.Join(tmpdir, "two")}
	parts, err := partitionPathsAsFiles(cmdArgs, paths)
	require.NoError(t, err)
	assert.Len(t, parts, 0)
}

func TestPathsToPackagePaths(t *testing.T) {
	root1 := "/fake/root1"
	root2 := "/fake/root2"
	root3 := "/fake/root3"
	gopath := root1 + ":" + root2 + ":" + root3
	defer fakeGoPath(t, gopath)()

	wd, err := os.Getwd()
	require.NoError(t, err)

	outsidePath := "/fake/outside/example4.com/foo4"
	outsidePackage, _ := filepath.Rel(wd, outsidePath)

	packagePaths, err := pathsToPackagePaths([]string{
		filepath.Join(root1, "src", "example1.com", "foo1"),
		filepath.Join(root2, "src", "example2.com", "foo2"),
		filepath.Join(root3, "src", "example3.com", "foo3"),
		outsidePath,
		"./relative/package",
	})
	require.NoError(t, err)

	expected := []string{
		"example1.com/foo1",
		"example2.com/foo2",
		"example3.com/foo3",
		outsidePackage,
		"./relative/package",
	}
	assert.Equal(t, expected, packagePaths)
}

func fakeGoPath(t *testing.T, path string) func() {
	oldpath := os.Getenv("GOPATH")
	require.NoError(t, os.Setenv("GOPATH", path))
	return func() { require.NoError(t, os.Setenv("GOPATH", oldpath)) }
}

func TestPartitionPathsByDirectory(t *testing.T) {
	cmdArgs := []string{"/usr/bin/foo", "-c"}
	paths := []string{"one", "two", "three"}

	parts, err := partitionPathsByDirectory(cmdArgs, paths)
	require.NoError(t, err)
	expected := [][]string{
		append(cmdArgs, "one"),
		append(cmdArgs, "two"),
		append(cmdArgs, "three"),
	}
	assert.Equal(t, expected, parts)

}
