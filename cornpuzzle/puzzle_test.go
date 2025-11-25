package cornpuzzle

import (
	"testing"
)

func TestResolve(t *testing.T) {
	// 测试用的拼图块
	blocks := []string{"  #|###", "  #|###", " ##|###", " ##|###", "## |###", "  #|# #|###", "  #|# #|###", "# #|# #|###", "  #|  #|###", "#  |## |###", " # |###| ##", " ##|###|## ", "  #|###|## ", "#  |#  |## |###", " # | # | # |###", " ##| # | # |## ", "  #|  #|###| ##"}
	pzl, err := Create(14, 7, blocks)
	if err != nil {
		t.Fatalf("创建拼图失败: %v", err)
	}
	result := Resolve(pzl)
	t.Logf("求解结果: %v", result)
	if result {
		t.Log("求解成功")
	} else {
		t.Log("求解失败或超时")
	}
}
