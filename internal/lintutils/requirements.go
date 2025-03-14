package lintutils

import "golang.org/x/tools/go/analysis"

// ResultOf returns requirement result by given name
func ResultOf(pass *analysis.Pass, name string) interface{} {
	for an, req := range pass.ResultOf {
		if an.Name == name {
			return req
		}
	}
	return nil
}
