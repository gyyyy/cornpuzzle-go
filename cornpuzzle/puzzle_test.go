package cornpuzzle

import (
	"testing"
)

func TestResolve(t *testing.T) {
	// 测试用的拼图块
	blocks := []string{"  #|###", "  #|###", " ##|###", " ##|###", "## |###", "  #|# #|###", "  #|# #|###", "# #|# #|###", "  #|  #|###", "#  |## |###", " # |###| ##", " ##|###|## ", "  #|###|## ", "#  |#  |## |###", " # | # | # |###", " ##| # | # |## ", "  #|  #|###| ##"}
	pzl, err := Create(14, 7, blocks)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !Resolve(pzl) {
		t.Fatalf("Resolve() failed")
	}
}
