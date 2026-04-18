package node

import "testing"

func TestCompatJarBaseNameKeepsSemanticVersion(t *testing.T) {
	got := compatJarBaseName("/tmp/platform-0.0.1")
	if got != "platform-0.0.1" {
		t.Fatalf("compatJarBaseName() = %q, want %q", got, "platform-0.0.1")
	}
}

func TestCompatJarBaseNameStripsExeExtension(t *testing.T) {
	got := compatJarBaseName("platform-0.0.1.exe")
	if got != "platform-0.0.1" {
		t.Fatalf("compatJarBaseName() = %q, want %q", got, "platform-0.0.1")
	}
}
