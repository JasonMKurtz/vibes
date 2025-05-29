package models

import "strings"

// DiffOp represents a line level diff operation.
type DiffOp struct {
	Type    string `json:"type"` // add, del, equal
	Text    string `json:"text"`
	OldLine int    `json:"old_line,omitempty"`
	NewLine int    `json:"new_line,omitempty"`
}

// DiffLines returns a simple line-based diff between before and after.
func DiffLines(before, after string) []DiffOp {
	a := strings.Split(before, "\n")
	b := strings.Split(after, "\n")
	// LCS dynamic programming
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}
	for i := len(a) - 1; i >= 0; i-- {
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}
	i, j := 0, 0
	var ops []DiffOp
	oldLine, newLine := 1, 1
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			ops = append(ops, DiffOp{Type: "equal", Text: a[i], OldLine: oldLine, NewLine: newLine})
			i++
			j++
			oldLine++
			newLine++
		} else if dp[i+1][j] >= dp[i][j+1] {
			ops = append(ops, DiffOp{Type: "del", Text: a[i], OldLine: oldLine})
			i++
			oldLine++
		} else {
			ops = append(ops, DiffOp{Type: "add", Text: b[j], NewLine: newLine})
			j++
			newLine++
		}
	}
	for i < len(a) {
		ops = append(ops, DiffOp{Type: "del", Text: a[i], OldLine: oldLine})
		i++
		oldLine++
	}
	for j < len(b) {
		ops = append(ops, DiffOp{Type: "add", Text: b[j], NewLine: newLine})
		j++
		newLine++
	}
	return ops
}
