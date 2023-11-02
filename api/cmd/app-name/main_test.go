package main

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/fredbi/go-cli/cli"

	"github.com/stretchr/testify/require"
)

const appName = "app-name" // CHANGE_ME

var argsMx sync.Mutex

func mockFatal(called *bool) func(string, ...any) {
	return func(_ string, _ ...any) {
		log.Println("Die called")
		*called = true
	}
}

func setMockFatal() (*bool, func()) {
	var called bool
	cli.SetDie(mockFatal(&called))
	revert := func() { cli.SetDie(log.Fatalf) }

	return &called, revert
}

func setMockArgs(args []string) func() {
	argsMx.Lock()
	defer argsMx.Unlock()

	original := os.Args
	os.Args = args
	revert := func() { os.Args = original }

	return revert
}

func TestMain(t *testing.T) {
	// warning: since this test alters globals, it is not compatible with parallel testing

	t.Run("with valid args", func(t *testing.T) {
		called, revertFatal := setMockFatal()
		defer revertFatal()

		revertArgs := setMockArgs(
			[]string{appName, "--version"},
		)
		defer revertArgs()

		main()
		require.False(t, *called)
	})

	t.Run("with invalid args", func(t *testing.T) {
		called, revertFatal := setMockFatal()
		defer revertFatal()

		revertArgs := setMockArgs(
			[]string{appName, "zorg"},
		)
		defer revertArgs()

		main()
		require.True(t, *called)
	})
}
