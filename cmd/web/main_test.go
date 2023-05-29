package main

import "testing"

func testRun(t *testing.T) {
	_, err := run()

	if err != nil {
		t.Errorf("Failed run()")
	}
}
