package parser_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// binaryName represents the name of the
// generated binary file used for testing
const binaryName = "logy"

// helperGenerateBinary is a test helper function
// that generates the binary file
func helperGenerateBinary(t *testing.T) func() {
	cmd := exec.Command("go", "build", "-o", binaryName)
	err := cmd.Run()
	if err != nil {
		t.Fatal("Could not generate binary")
	}
	return func() {
		if err := os.Remove(binaryName); err != nil {
			t.Fatal("Could not remove binary")
		}
	}
}

// TestNew tests if our parser can be instantiated
// successfully based upon the input command line flags
func TestNew(t *testing.T) {
	defer helperGenerateBinary(t)()

	tests := []struct {
		flags []string
		exit  bool
	}{
		{
			nil,
			true,
		},
		{
			[]string{},
			true,
		},
		{
			[]string{"-h"},
			true,
		},
		{
			[]string{"--help"},
			true,
		},
		{
			[]string{"-path"},
			false,
		},
		{
			[]string{"-path="},
			false,
		},
		{
			[]string{"-path=none"},
			false,
		},
		{
			[]string{"-lines"},
			false,
		},
		{
			[]string{"-lines="},
			false,
		},
		{
			[]string{"-lines=0"},
			false,
		},
		{
			[]string{"-lines=-1"},
			false,
		},
		{
			[]string{"-page"},
			false,
		},
		{
			[]string{"-page="},
			false,
		},
		{
			[]string{"-page=0"},
			false,
		},
		{
			[]string{"-page=-1"},
			false,
		},
		{
			[]string{"-text"},
			false,
		},
		{
			[]string{"-text="},
			false,
		},
		{
			[]string{"-text=none"},
			false,
		},
		{
			[]string{"-text=plain"},
			false,
		},
		{
			[]string{"-text=json"},
			false,
		},
		{
			[]string{"-filter"},
			false,
		},
		{
			[]string{"-filter="},
			false,
		},
		{
			[]string{"-filter=none"},
			false,
		},
		{
			[]string{"--with-regex"},
			false,
		},
		{
			[]string{"--with-regex="},
			false,
		},
		{
			[]string{"--with-regex=none"},
			false,
		},
		{
			[]string{"--no-color"},
			false,
		},
		{
			[]string{"--no-color="},
			false,
		},
		{
			[]string{"--no-color=none"},
			false,
		},
		{
			[]string{"-ext"},
			false,
		},
		{
			[]string{"-ext="},
			false,
		},
		{
			[]string{"-ext=none"},
			false,
		},
		{
			[]string{"-path=testdata"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden"},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden"))},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.none"))},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "--with-regex"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "--with-regex"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-text"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-text="},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-text=none"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-text=plain"},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-text=json"},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-lines"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-lines="},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-lines=0"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-lines=-1"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-lines=1"},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-page"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-page="},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-page=0"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-page=-1"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-page=1"},
			true,
		},
		{
			[]string{"-path=testdata", "-text"},
			false,
		},
		{
			[]string{"-path=testdata", "-text="},
			false,
		},
		{
			[]string{"-path=testdata", "-text=none"},
			false,
		},
		{
			[]string{"-path=testdata", "-text=plain"},
			false,
		},
		{
			[]string{"-path=testdata", "-text=json"},
			false,
		},
		{
			[]string{"-path=testdata", "-lines"},
			false,
		},
		{
			[]string{"-path=testdata", "-lines="},
			false,
		},
		{
			[]string{"-path=testdata", "-lines=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-lines=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-lines=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-page"},
			false,
		},
		{
			[]string{"-path=testdata", "-page="},
			false,
		},
		{
			[]string{"-path=testdata", "-page=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-page=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-page=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-text"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-text="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-text=none"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-text=plain"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-text=json"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-lines"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-lines="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-lines=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-lines=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-lines=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-page"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-page="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-page=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-page=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext", "-page=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-text"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-text="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-text=none"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-text=plain"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-text=json"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-lines"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-lines="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-lines=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-lines=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-lines=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-page"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-page="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-page=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-page=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=none", "-page=1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-text"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-text="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-text=none"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-text=plain"},
			true,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-text=json"},
			true,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-lines"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-lines="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-lines=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-lines=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-lines=1"},
			true,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-page"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-page="},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-page=0"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-page=-1"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-page=1"},
			true,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-filter=[>}", "--with-regex"},
			false,
		},
		{
			[]string{fmt.Sprintf("-path=%s", filepath.Join("testdata", "file.golden")), "-filter=[0-9][a-z]", "--with-regex"},
			true,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-filter=[>}", "--with-regex"},
			false,
		},
		{
			[]string{"-path=testdata", "-ext=golden", "-filter=[0-9][a-z]", "--with-regex"},
			true,
		},
	}

	for _, tc := range tests {
		cmd := exec.Command(binaryName, tc.flags...)
		err := cmd.Run()
		if err == nil && tc.exit == false {
			t.Fatalf("With flags: %v, exited with error", tc.flags)
		}
		if e, ok := err.(*exec.ExitError); ok && e.Success() != tc.exit {
			t.Fatalf("With flags: %v, expected %v; got %v", tc.flags, tc.exit, e.Success())
		}
	}
}
