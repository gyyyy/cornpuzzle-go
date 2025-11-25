package cornpuzzle

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var Verbose = true

type Corn struct {
	x      int
	y      int
	remain int
	Groove [][]int
}

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

func (c *Corn) Unset(blk *Block) {
	c.remove(blk.no)
}

func (c *Corn) Size() (int, int) {
	return c.x, c.y
}

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

type Block struct {
	no     int
	left   int
	height int
	item   [][]int
}

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

func (b *Block) absame(blk *Block) bool {
	return b.same(blk) || b.same(blk.Reverse())
}

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

type Puzzle struct {
	Corn  *Corn
	Block []*Block
}

// corn 3d jigsaw puzzle
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
