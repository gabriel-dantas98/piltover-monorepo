package cli

import (
	"io"
	"os"
)

var stderrSinkFn = func() io.Writer { return os.Stderr }

func stderrSink() io.Writer { return stderrSinkFn() }
