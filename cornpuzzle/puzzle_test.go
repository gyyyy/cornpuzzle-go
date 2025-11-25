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

// TestPreCheck 测试预检查优化：单元格数不匹配时应快速失败
func TestPreCheck(t *testing.T) {
	// 创建一个单元格数少于板子的拼图
	blocks := []string{"  #|###", "  #|###"} // 只有2个块，总单元格少
	pzl, err := Create(14, 7, blocks)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	// 预检查应使Resolve快速返回false
	if Resolve(pzl) {
		t.Fatalf("Resolve() should fail due to insufficient cells")
	}
}
