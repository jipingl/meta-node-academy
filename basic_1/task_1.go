package main

// 136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
// 可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
func singleNumber(nums []int) int {
	// result := 0
	// for a := 0; a < len(nums); a++ {
	// 	result ^= nums[a]
	// }
	// return result

	// 使用 map 记录每个元素出现的次数
	countMap := make(map[int]int)
	for _, num := range nums {
		countMap[num]++
	}
	// 遍历 map 找到出现次数为1的元素
	for num, count := range countMap {
		if count == 1 {
			return num
		}
	}
	return 0
}

// 给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false
// 回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
func isPalindrome(x int) bool {
	// 负数不是回文数
	if x < 0 {
		return false
	}
	// 最后一位为0也不会是回文数
	if x%10 == 0 && x != 0 {
		return false
	}
	// 获取后边一半的数字
	halfNumber := 0
	for halfNumber < x {
		halfNumber = halfNumber*10 + x%10
		x /= 10
	}
	// 前后数字相等就死回文数字
	return halfNumber == x || halfNumber/10 == x
}

// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
// 有效字符串需满足：
// 1.左括号必须用相同类型的右括号闭合。
// 2.左括号必须以正确的顺序闭合。
// 3.每个右括号都有一个对应的相同类型的左括号。
func isValid(s string) bool {
	n := len(s)

	// 字符串长度如果为奇数直接返回
	if n%2 == 1 {
		return false
	}

	// 使用map保存括号的左右关系
	pairMap := map[byte]byte{
		')': '(',
		'}': '{',
		']': '[',
	}

	// 使用数组切片构造栈结构
	stack := []byte{}
	for i := 0; i < n; i++ {
		// 匹配括号右侧时需要与栈顶做对比
		if pairMap[s[i]] > 0 {
			if len(stack) == 0 || stack[len(stack)-1] != pairMap[s[i]] {
				return false
			}
			// 匹配成功移除栈顶元素
			stack = stack[:len(stack)-1]
		} else {
			// 左侧括号先入栈
			stack = append(stack, s[i])
		}
	}
	// 全部匹配后栈应为空
	return len(stack) == 0
}

// 编写一个函数来查找字符串数组中的最长公共前缀。
// 如果不存在公共前缀，返回空字符串 ""。
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	// 二分法查找最长前缀
	var lcp func(int, int) string
	lcp = func(start, end int) string {
		// 只有一个字符串直接返回
		if start == end {
			return strs[start]
		}
		mid := (start + end) / 2
		longL, longR := lcp(start, mid), lcp(mid+1, end)
		minLen := min(len(longL), len(longR))
		// 逐位比较
		for i := 0; i < minLen; i++ {
			if longL[i] != longR[i] {
				return longL[:i]
			}
		}
		return longL[:minLen]
	}
	// 调用二分法查找
	return lcp(0, len(strs)-1)
}

// 给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。
// 元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
func removeDuplicates(nums []int) int {
	if len(nums) <= 1 {
		return len(nums)
	}
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}

// 给定一个由 整数 组成的 非空 数组所表示的非负整数，在该数的基础上加一。
// 最高位数字存放在数组的首位， 数组中每个元素只存储单个数字。
// 你可以假设除了整数 0 之外，这个整数不会以零开头。
func plusOne(digits []int) []int {
	// 从后向前遍历
	for i := len(digits) - 1; i >= 0; i-- {
		// 不为9加一不会发生进位
		if digits[i] != 9 {
			digits[i]++
			// 发生进位后的位置值为0
			for j := i + 1; j < len(digits); j++ {
				digits[j] = 0
			}
			return digits
		}
	}
	// 如果经过第一步没有找到不为9的位置
	res := make([]int, len(digits)+1)
	res[0] = 1
	return res
}

// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间
func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}
	// 按区间起始值排序
	// 外部循环控制比较的轮次
	for i := 1; i < len(intervals); i++ {
		// 内部循环交换值的位置
		for j := i; j < len(intervals); j++ {
			front := intervals[j-1]
			after := intervals[j]
			if front[0] > after[0] {
				intervals[j-1] = after
				intervals[j] = front
			}
		}
	}
	// 排序结束后开始合并区间
	merged := intervals[:1]
	for i := 1; i < len(intervals); i++ {
		// 存在交集保留结束值中的最大值
		lastMerged := merged[len(merged)-1]
		if lastMerged[0] <= intervals[i][0] && intervals[i][0] <= lastMerged[1] {
			maxEnd := max(lastMerged[1], intervals[i][1])
			lastMerged[1] = maxEnd
			continue
		}
		// 不存在交集直接将区间存入最终结果
		merged = append(merged, intervals[i])
	}
	return merged
}

// 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
func twoSum(nums []int, target int) []int {
	// 采用映射表来实现
	numsMap := map[int]int{}
	for i, v := range nums {
		// 检查映射表是否存在和为目标值的元素
		if p, exist := numsMap[target-v]; exist {
			return []int{i, p}
		}
		numsMap[v] = i
	}
	return nil
}

func main() {
	// nums := []int{4, 1, 2, 1, 2}
	// fmt.Println(singleNumber(nums))

	// x := 100011
	// fmt.Println(isPalindrome(x))

	// fmt.Println(isValid("{[{()}]}"))

	// strs := [3]string{"gene", "generies", "gender"}
	// fmt.Printf(longestCommonPrefix(strs[:]))

	// tmp := [][]int{{1, 3}, {0, 5}, {8, 10}, {2, 6}, {15, 18}}
	// fmt.Printf("%v\n", merge(tmp))
}
