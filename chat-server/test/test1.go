package main

import "fmt"

//func twoSum(nums []int, target int) []int {
//	hash := make(map[int32]uint16)
//	for i := range nums {
//		temp := int32(target - nums[i])
//		if j, ok := hash[temp]; ok {
//			fmt.Printf("%d", hash[temp])
//			return []int{int(j), i}
//		}
//		hash[int32(nums[i])] = uint16(i)
//	}
//	return nil
//}

func twoSum(nums []int, target int) []int {
	ans := make([]int, 2)
	mp := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if val, ok := mp[target-nums[i]]; ok {
			ans[0], ans[1] = val, i
			break
		}
		mp[nums[i]] = i
	}
	return ans
}

func main() {
	index := twoSum([]int{3, 2, 3}, 6)
	fmt.Println(index)
}
