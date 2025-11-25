package cornpuzzle

import (
	"fmt"
	"log"
	"slices"
)

var (
	dup = map[int][]int{} // 记录重复的块编号，用于优化搜索（跳过相同的块）
)

// sort 对拼图的块进行排序，按高度降序、宽度降序、单元格数量降序、编号升序
func sort(pzl *Puzzle) []int {
	var (
		n   = len(pzl.Block)
		sli = make([]int, n)
	)
	for i, blk := range pzl.Block {
		sli[i] = blk.no
	}
	slices.SortFunc(sli, func(a, b int) int {
		var (
			left  = pzl.Block[a-1]
			right = pzl.Block[b-1]
			cmp   = left.height - right.height
		)
		if cmp != 0 {
			return cmp
		}
		if cmp = left.width - right.width; cmp != 0 {
			return cmp
		}
		if cmp = left.Count() - right.Count(); cmp != 0 {
			return cmp
		}
		return left.no - right.no
	})
	slices.Reverse(sli)
	return sli
}

// check 检查拼图是否已完全解决（所有空格都被填满）
func check(pzl *Puzzle) bool {
	if pzl.Corn.remain > 0 {
		return false
	}
	x, y := pzl.Corn.next()
	return x == -1 && y == -1
}

// backtrack 使用回溯法尝试放置块
func backtrack(pzl *Puzzle, sorted []int, n int) bool {
	if n == 0 {
		return true
	}
	// 剪枝：计算剩余块的总单元格数，如果超过剩余空格，提前失败
	remainingCells := 0
	for i := 0; i < n; i++ {
		remainingCells += pzl.Block[sorted[i]-1].Count()
	}
	if remainingCells > pzl.Corn.remain {
		return false
	}
	skip := map[int]bool{}
	for i := range n {
		blk := pzl.Block[sorted[i]-1]
		if skip[blk.no] {
			continue
		}
		for _, b := range dup[blk.no] {
			skip[b] = true
		}
		if !pzl.Corn.Set(blk) {
			if rblk := blk.Reverse(); rblk.same(blk) || !pzl.Corn.Set(rblk) {
				continue
			}
		}
		if Verbose {
			fmt.Printf("%d ", blk.no)
		}
		if sorted[i], sorted[n-1] = sorted[n-1], sorted[i]; backtrack(pzl, sorted, n-1) {
			if Verbose && n == 1 {
				fmt.Println()
			}
			return true
		}
		sorted[i], sorted[n-1] = sorted[n-1], sorted[i]
		if pzl.Corn.remove(blk.no); Verbose {
			if n == len(pzl.Block) {
				fmt.Println()
			} else {
				j := 2
				if blk.no >= 10 {
					j++
				}
				for _, s := range []string{"\x08", " ", "\x08"} {
					for range j {
						fmt.Print(s)
					}
				}
			}
		}
	}
	return false
}

// Resolve 求解给定的拼图
// 参数：pzl - 要解决的拼图实例
// 返回：true 如果成功解决，false 如果无法解决或超时
func Resolve(pzl *Puzzle) bool {
	sorted := sort(pzl)
	for i := range sorted {
		left := pzl.Block[sorted[i]-1]
		if _, ok := dup[left.no]; !ok {
			dup[left.no] = []int{}
		}
		for j := i + 1; j < len(sorted); j++ {
			right := pzl.Block[sorted[j]-1]
			if !left.absame(right) {
				continue
			}
			dup[left.no] = append(dup[left.no], right.no)
			if _, ok := dup[right.no]; !ok {
				dup[right.no] = []int{}
			}
			dup[right.no] = append(dup[right.no], left.no)
		}
		if len(dup[left.no]) == 0 {
			delete(dup, left.no)
		}
	}
	if Verbose {
		log.Println("=== 开始求解玉米拼图 ===")
	}
	return backtrack(pzl, sorted, len(sorted)) && check(pzl)
}
