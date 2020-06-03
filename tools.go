// +build tools

// Manage versions of tool dependencies as per:
//  https://github.com/go-modules-by-example/index/tree/master/010_tools

package tools

import (
	_ "github.com/kisielk/errcheck"
	_ "golang.org/x/lint/golint"
)
