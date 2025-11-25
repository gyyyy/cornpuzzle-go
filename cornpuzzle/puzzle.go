package cornpuzzle

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var Verbose = true

// Corn 表示玉米拼图的二维映射板子（网格），支持环绕
type Corn struct {
	x      int     // 板子的宽度（x轴长度）
	y      int     // 板子的高度（y轴长度）
	remain int     // 剩余未填充的空格数量
	Groove [][]int // 二维数组表示板子的状态，0表示空，>0表示块编号
}

// next 返回下一个可填充的位置（x, y），从左到右，从上到下扫描
func (c *Corn) next() (int, int) {
	if c.remain == 0 {
		return -1, -1
	}
	for i := 0; i < c.y; i++ {
		for j := 0; j < c.x; j++ {
			if c.Groove[j][i] <= 0 {
				return j + 1, i + 1
			}
		}
	}
	return -1, -1
}

// fill 在指定位置 (x, y) 填充值 val
// x 和 y 从1开始，支持环绕（x超出范围时环绕到另一侧）
func (c *Corn) fill(x, y, val int) bool {
	if val <= 0 || y <= 0 || y > c.y {
		return false
	}
	if x > 0 {
		if x = x % c.x; x == 0 {
			x = c.x
		}
	} else if x <= 0 {
		if x = x % c.x; x == 0 {
			x = 1
		} else {
			x += c.x + 1
		}
	}
	empty := c.Groove[x-1][y-1] <= 0
	if empty {
		c.Groove[x-1][y-1] = val
		c.remain--
	}
	return empty
}

// remove 移除所有值为 val 的块，恢复为空格
func (c *Corn) remove(val int) {
	if val <= 0 {
		return
	}
	for i, g := range c.Groove {
		for j, v := range g {
			if v == val {
				c.Groove[i][j] = 0
				c.remain++
			}
		}
	}
}

// Set 尝试将拼图块放置到板子上，从下一个可用位置开始
// 参数：blk - 要放置的块
// 如果：true 如果放置成功，false 如果放置失败（回滚）
func (c *Corn) Set(blk *Block) bool {
	baseX, baseY := c.next()
	if baseY+blk.height-1 > c.y {
		return false
	}
	for i, item := range blk.item {
		for j, v := range item {
			if v <= 0 {
				continue
			}
			x := baseX + (j - blk.left + 1)
			if x < blk.left {
				x -= 1
			}
			if !c.fill(x, baseY+i, v) {
				c.remove(blk.no)
				return false
			}
		}
	}
	return true
}

// Unset 从板子上移除块
// 参数：blk - 要移除的块
func (c *Corn) Unset(blk *Block) {
	c.remove(blk.no)
}

// Size 返回板子的尺寸
// 返回：板子的宽度和高度
func (c *Corn) Size() (int, int) {
	return c.x, c.y
}

// String 返回板子的字符串表示，用于打印
func (c *Corn) String() string {
	var s string
	for i := c.y - 1; i >= 0; i-- {
		for j := 0; j < c.x; j++ {
			if v := c.Groove[j][i]; v <= 0 {
				s += "##"
			} else {
				s += fmt.Sprintf("%02d", v)
			}
		}
		s += "\n"
	}
	return s
}

// newCorn 创建一个新的 Corn 实例，尺寸为 x*y，所有位置初始为空
func newCorn(x, y int) *Corn {
	groove := make([][]int, x)
	for i := 0; i < x; i++ {
		groove[i] = make([]int, y)
	}
	return &Corn{
		x:      x,
		y:      y,
		remain: x * y,
		Groove: groove,
	}
}

// Block 表示一个拼图块
// 块有编号、左偏移、高度和形状（二维数组）
type Block struct {
	no     int     // 块的编号
	left   int     // 块的左偏移（从哪列开始）
	height int     // 块的高度
	item   [][]int // 二维数组表示块的形状，>0表示占用，0表示空
}

// same 检查两个块是否形状相同（不考虑翻转）
func (b *Block) same(blk *Block) bool {
	if blk.left != b.left || blk.height != b.height {
		return false
	}
	for i, litem := range b.item {
		ritem := blk.item[i]
		if len(ritem) != len(litem) {
			return false
		}
		for j, lv := range litem {
			if rv := ritem[j]; (rv > 0 && lv == 0) || (rv == 0 && lv > 0) {
				return false
			}
		}
	}
	return true
}

// absame 检查两个块是否形状相同，包括翻转后的情况
func (b *Block) absame(blk *Block) bool {
	return b.same(blk) || b.same(blk.Reverse())
}

// Count 返回块中非空单元格的数量
// Count 返回块中非空单元格的数量
func (b *Block) Count() int {
	var count int
	for _, item := range b.item {
		for _, v := range item {
			if v > 0 {
				count++
			}
		}
	}
	return count
}

// Width 返回块的最大宽度（最宽行的长度）
func (b *Block) Width() int {
	maxWidth := 0
	for _, item := range b.item {
		if len(item) > maxWidth {
			maxWidth = len(item)
		}
	}
	return maxWidth
}

// Reverse 返回块的翻转版本（上下翻转）
// 返回：翻转后的块
func (b *Block) Reverse() *Block {
	rb := &Block{
		no:     b.no,
		left:   0,
		height: b.height,
		item:   make([][]int, 0, len(b.item)),
	}
	for i := len(b.item) - 1; i >= 0; i-- {
		var (
			curr = b.item[i]
			line = make([]int, 0, len(curr))
		)
		for j := len(curr) - 1; j >= 0; j-- {
			v := curr[j]
			if line = append(line, v); i == len(b.item)-1 && v > 0 && rb.left == 0 {
				rb.left = len(curr) - j
			}
		}
		rb.item = append(rb.item, line)
	}
	return rb
}

// String 返回块的字符串表示，用于打印
func (b *Block) String() string {
	var s string
	for i := len(b.item) - 1; i >= 0; i-- {
		for _, v := range b.item[i] {
			if v <= 0 {
				for n := b.no; n > 0; n = n / 10 {
					s += " "
				}
			} else {
				s += strconv.Itoa(v)
			}
		}
		s += "\n"
	}
	return s[:len(s)-1]
}

// newBlock 从字符串 s 创建一个新的 Block 实例，编号为 no
// 字符串格式：行用 | 或 , 分隔，空格表示空，#表示占用
func newBlock(no int, s string) *Block {
	sep := "|"
	if strings.Contains(s, ",") {
		sep = ","
	}
	var (
		line = strings.Split(s, sep)
		blk  = &Block{
			no:   no,
			item: make([][]int, len(line)),
		}
	)
	for i := 0; i < len(line); i++ {
		l := line[len(line)-i-1]
		item := make([]int, len(l))
		for j, r := range l {
			if r == ' ' {
				continue
			}
			if item[j] = no; i == 0 && blk.left == 0 {
				blk.left = j + 1
			}
		}
		blk.item[i] = item
	}
	blk.height = len(blk.item)
	return blk
}

// Puzzle 表示整个玉米拼图，包括板子和所有块
type Puzzle struct {
	Corn  *Corn    // 拼图板
	Block []*Block // 所有拼图块的列表
}

// Create 创建一个新的 Puzzle 实例
// 参数：x 和 y 是板子的尺寸，blk 是块的字符串列表
// 返回：创建的拼图或错误
func Create(x, y int, blk []string) (*Puzzle, error) {
	if x <= 0 || y <= 0 {
		return nil, fmt.Errorf("无效的玉米拼图尺寸 [%d,%d]", x, y)
	}
	if len(blk) == 0 {
		return nil, errors.New("拼图块不能为空")
	}
	blks := make([]*Block, 0, len(blk))
	for i, s := range blk {
		if strings.TrimSpace(s) == "" {
			continue
		}
		blks = append(blks, newBlock(i+1, s))
	}
	if len(blks) == 0 {
		return nil, errors.New("无效的拼图块数量")
	}
	if Verbose {
		log.Printf("创建玉米拼图: 尺寸 [%d,%d], 拼图块数量 %d\n", x, y, len(blks))
	}
	return &Puzzle{
		Corn:  newCorn(x, y),
		Block: blks,
	}, nil
}
