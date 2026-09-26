package pdf417

import "testing"

// ISO/IEC 15438: in cluster 0 the left row indicator carries (rows-1)/3,
// in cluster 1 the right one does. They must agree, or scanners reject the
// symbol (the upstream bug fixed in this copy, see doc.go).
func TestRowIndicatorsAgreeOnRowCount(t *testing.T) {
	for rows := 3; rows <= 30; rows++ {
		left := getLeftCodeWord(0, rows, 5, 2) % 30
		right := getRightCodeWord(1, rows, 5, 2) % 30
		if left != (rows-1)/3 || right != (rows-1)/3 {
			t.Errorf("rows=%d: left %d, right %d, want %d", rows, left, right, (rows-1)/3)
		}
	}
}
