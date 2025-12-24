package main

import "fmt"

func main() {

	//list := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//list := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	//fmt.Println("3列的板块，两个连续格子 在同一行：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			// 3列的板块
	//			if checkSameRowForTwoBox(i, j, 3) {
	//				fmt.Println(i, "、", j)
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Println("\n4列的板块，两个连续格子 在同一行：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			// 4列的板块
	//			if checkSameRowForTwoBox(i, j, 4) {
	//				fmt.Println(i, "、", j)
	//			}
	//		}
	//	}
	//}

	//fmt.Println("\n3列的板块，两个连续格子 在同一列：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			//3列的板块
	//			if checkSameColForTwoBox(i, j, 3) {
	//				fmt.Println(i, "、", j)
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Println("\n4列的板块，两个连续格子 在同一列：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			// 4列的板块
	//			if checkSameColForTwoBox(i, j, 4) {
	//				fmt.Println(i, "、", j)
	//			}
	//		}
	//	}
	//}

	//fmt.Println("\n3列的板块，三个连续格子 在同一行：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			for _, k := range list {
	//				if j != k && j < k {
	//					if checkSameRowForThreeBox(i, j, k, 3) {
	//						fmt.Println(i, "、", j, "、", k)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Println("\n4列的板块，三个连续格子 在同一行：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			for _, k := range list {
	//				if j != k && j < k {
	//					if checkSameRowForThreeBox(i, j, k, 4) {
	//						fmt.Println(i, "、", j, "、", k)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	//fmt.Println("\n3列的板块，三个连续格子 在同一列：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			for _, k := range list {
	//				if j != k && j < k {
	//					if checkSameColForThreeBox(i, j, k, 3) {
	//						fmt.Println(i, "、", j, "、", k)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Println("\n4列的板块，三个连续格子 在同一列：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			for _, k := range list {
	//				if j != k && j < k {
	//					if checkSameColForThreeBox(i, j, k, 4) {
	//						fmt.Println(i, "、", j, "、", k)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	// 是否存在2*2的格子
	columnsNum := int32(4)

	var list [][]int32

	boxList := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	// 随机选6块地
	for ik, iv := range boxList {
		for jk := ik + 1; jk < len(boxList); jk++ {
			for mk := jk + 1; mk < len(boxList); mk++ {
				for nk := mk + 1; nk < len(boxList); nk++ {
					for ak := nk + 1; ak < len(boxList); ak++ {
						for bk := ak + 1; bk < len(boxList); bk++ {
							list = append(list, []int32{iv, boxList[jk], boxList[mk], boxList[nk], boxList[ak], boxList[bk]})
						}
					}
				}
			}

		}
	}

	// 是否存在2*2的格子
	//for _, v := range list {
	//	returnData := checkTwoToTwoLandShape(v, columnsNum)
	//	if len(returnData) == 4 {
	//		fmt.Println(v, "存在2*2的格子", returnData)
	//	}
	//}

	// 是否存在2*3的格子
	//for _, v := range list {
	//	returnData := checkTwoToThreeLandShape(v, columnsNum)
	//	if len(returnData) > 0 {
	//		fmt.Println(v, "存在2*3的格子", returnData)
	//	}
	//}

	// 是否存在3*2的格子
	//for _, v := range list {
	//	returnData := checkThreeToTwoLandShape(v, columnsNum)
	//	if len(returnData) > 0 {
	//		fmt.Println(v, "存在3*2的格子", returnData)
	//	}
	//}

	// 是否存在1*2的格子
	//for _, v := range list {
	//	returnData := checkOneToTwoLandShape(v, columnsNum)
	//	if len(returnData) > 0 {
	//		fmt.Println(v, "存在1*2的格子", returnData)
	//	}
	//}

	// 是否存在2*1的格子
	for _, v := range list {
		returnData := checkTwoToOneLandShape(v, columnsNum)
		if len(returnData) > 0 {
			fmt.Println(v, "存在2*1的格子", returnData)
		}
	}
}

// 是否存在2*2的形状
func checkTwoToTwoLandShape(list []int32, columnsNum int32) []int32 {
	boxNum := len(list)
	returnData := make([]int32, 0)
	for ik, iv := range list {
		if ik+1 < boxNum {
			// 一行是否存在 连续两个格子
			if checkSameRowForTwoBox(iv, list[ik+1], columnsNum) {
				// 如果一行存在连续的两个格子 判断下一行是否存在 连续两个格子
				nextIv := iv + columnsNum
				nextJv := list[ik+1] + columnsNum

				isNextExist := 0
				for jk := ik + 2; jk < boxNum; jk++ {
					if list[jk] == nextIv || list[jk] == nextJv {
						isNextExist++
						if isNextExist == 2 {
							returnData = append(returnData, iv, list[ik+1], nextIv, nextJv)
							return returnData
						}
					}
				}
			}
		}
	}

	return returnData
}

// 是否存在2*3的形状
func checkTwoToThreeLandShape(list []int32, columnsNum int32) []int32 {
	boxNum := len(list)
	returnData := make([]int32, 0)
	for ik, iv := range list {
		if ik+1 < boxNum {
			// 一行是否存在 连续两个格子
			if checkSameRowForTwoBox(iv, list[ik+1], columnsNum) {
				// 如果一行存在连续的两个格子 判断下一行是否存在 连续两个格子
				nextIv := iv + columnsNum
				nextJv := list[ik+1] + columnsNum

				isNextExist := 0
				jk := ik + 2
				for ; jk < boxNum; jk++ {
					if list[jk] == nextIv || list[jk] == nextJv {
						isNextExist++
						if isNextExist == 2 {
							break
						}
					}
				}

				// 第三行是否也存在
				if isNextExist == 2 {
					lastIv := nextIv + columnsNum
					lastJv := nextJv + columnsNum

					isLastExist := 0
					for mk := jk + 1; mk < boxNum; mk++ {
						if list[mk] == lastIv || list[mk] == lastJv {
							isLastExist++
							if isLastExist == 2 {
								returnData = append(returnData, iv, list[ik+1], nextIv, nextJv, lastIv, lastJv)
								return returnData
							}
						}
					}
				}
			}
		}
	}

	return returnData
}

// 是否存在3*2的形状
func checkThreeToTwoLandShape(list []int32, columnsNum int32) []int32 {
	boxNum := len(list)
	returnData := make([]int32, 0)

	for ik, iv := range list {
		if ik+2 < boxNum {
			// 一行是否存在 连续的三个格子
			if checkSameRowForThreeBox(iv, list[ik+1], list[ik+2], columnsNum) {

				// 如果一行存在连续的三个格子 判断下一行是否存在 连续三个格子
				nextIv := iv + columnsNum
				nextJv := list[ik+1] + columnsNum
				nextKv := list[ik+2] + columnsNum
				isNextExist := 0
				jk := ik + 3
				for ; jk < boxNum; jk++ {
					if list[jk] == nextIv || list[jk] == nextJv || list[jk] == nextKv {
						isNextExist++
						if isNextExist == 3 {
							returnData = append(returnData, iv, list[ik+1], list[ik+2], nextIv, nextJv, nextKv)
							return returnData
						}
					}
				}
			}
		}
	}

	return returnData
}

// 是否存在1*2的形状
func checkOneToTwoLandShape(list []int32, columnsNum int32) []int32 {
	boxNum := len(list)
	returnData := make([]int32, 0)

	for ik, iv := range list {
		if ik+1 < boxNum {
			// 一列是否存在 连续的两个格子
			if checkSameColForTwoBox(iv, list[ik+1], columnsNum) {
				returnData = append(returnData, iv, list[ik+1])
				return returnData
			}
		}
	}

	return returnData
}

// 是否存在2*1的形状
func checkTwoToOneLandShape(list []int32, columnsNum int32) []int32 {
	boxNum := len(list)
	returnData := make([]int32, 0)

	for ik, iv := range list {
		if ik+1 < boxNum {
			// 一行是否存在 连续的两个格子
			if checkSameRowForTwoBox(iv, list[ik+1], columnsNum) {
				returnData = append(returnData, iv, list[ik+1])
				return returnData
			}
		}
	}

	return returnData
}

// 判断两个格子是否连续且在同一行 (i<j) columnsNum：列数
func checkSameRowForTwoBox(i, j, columnsNum int32) bool {
	if i+1 != j {
		return false
	}

	if (i%columnsNum != 0 && i/columnsNum == j/columnsNum) || (j%columnsNum == 0) {
		return true
	}

	return false
}

// 判断两个格子是否连续且在同一列 (i<j) rowNum：行数
func checkSameColForTwoBox(i, j, rowNum int32) bool {
	if i+rowNum == j && i%rowNum == j%rowNum {
		return true
	} else {
		return false
	}
}

// 判断三个格子是否连续且在同一行 (i<j<k) columnsNum：列数
func checkSameRowForThreeBox(i, j, k, columnsNum int32) bool {
	if !(i+1 == j && j+1 == k) {
		return false
	}

	if (i%columnsNum != 0 && j%columnsNum != 0 && k%columnsNum != 0) || k%columnsNum == 0 {
		return true
	}

	return false
}

// 判断三个格子是否连续且在同一列 (i<j<k) rowNum：行数
func checkSameColForThreeBox(i, j, k, rowNum int32) bool {
	if i+rowNum == j && j+rowNum == k && i%rowNum == j%rowNum && j%rowNum == k%rowNum {
		return true
	} else {
		return false
	}
}
