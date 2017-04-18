// Copyright 2017 Pulumi, Inc. All rights reserved.

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := NewCocoCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(-1)
	}
}
