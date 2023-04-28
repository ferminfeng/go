package main

import "fmt"

func main() {

	//list := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}

	list := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	fmt.Println("3列的板块，两个连续格子 在同一行：")
	for _, i := range list {
		for _, j := range list {
			if i != j && i < j {
				// 3列的板块
				if checkSameRowForThreeToTwo(i, j) {
					fmt.Println(i, "、", j)
				}
			}
		}
	}

	fmt.Println("\n4列的板块，两个连续格子 在同一行：")
	for _, i := range list {
		for _, j := range list {
			if i != j && i < j {
				// 4列的板块
				if checkSameRowForFourToTwo(i, j) {
					fmt.Println(i, "、", j)
				}
			}
		}
	}

	//fmt.Println("\n3列的板块，两个连续格子 在同一列：")
	//for _, i := range list {
	//	for _, j := range list {
	//		if i != j && i < j {
	//			//3列的板块
	//			if checkSameColForThreeToTwo(i, j) {
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
	//			if checkSameColForFourToTwo(i, j) {
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
	//					// 3列的板块
	//					if checkSameRowForThreeToThree(i, j, k) {
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
	//					// 4列的板块
	//					if checkSameRowForFourToThree(i, j, k) {
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
	//					// 3列的板块
	//					if checkSameColForThreeToThree(i, j, k) {
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
	//					// 4列的板块
	//					if checkSameColForFourToThree(i, j, k) {
	//						fmt.Println(i, "、", j, "、", k)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}
}

// 3*3、3*4板块判断其中两个格子是否连续且在同一行 (i<j)
func checkSameRowForThreeToTwo(i, j int32) bool {
	if i+1 != j {
		return false
	}

	if ((i+j)%3 == 0 && i/3 == j/3) || j%3 == 0 {
		return true
	}

	return false
}

// 4*3、板块判断其中两个格子是否连续且在同一行 (i<j)
func checkSameRowForFourToTwo(i, j int32) bool {
	if i+1 != j {
		return false
	}

	if (i%4 != 0 && i/4 == j/4) || (j%4 == 0) {
		return true
	}

	return false
}

// 3*3、3*4板块判断其中三个格子是否连续且在同一行 (i<j<k)
func checkSameRowForThreeToThree(i, j, k int32) bool {
	if !(i+1 == j && j+1 == k) {
		return false
	}

	if k%3 == 0 && i/3 == j/3 && j+1 == k {
		return true
	}

	return false
}

// 4*3板块判断其中三个格子是否连续且在同一行 (i<j<k)
func checkSameRowForFourToThree(i, j, k int32) bool {
	if !(i+1 == j && j+1 == k) {
		return false
	}

	if (i%4 != 0 && j%4 != 0 && k%4 != 0) || k%4 == 0 {
		return true
	}

	return false
}

// 3*3、3*4板块判断其中两个格子是否连续且在同一列 (i<j)
func checkSameColForThreeToTwo(i, j int32) bool {
	if i+3 == j && i%3 == j%3 {
		return true
	} else {
		return false
	}
}

// 4*3板块判断其中两个格子是否连续且在同一列 (i<j)
func checkSameColForFourToTwo(i, j int32) bool {
	if i+4 == j && i%4 == j%4 {
		return true
	} else {
		return false
	}
}

// 3*3、3*4板块判断其中三个格子是否连续且在同一列 (i<j<k)
func checkSameColForThreeToThree(i, j, k int32) bool {
	if i+3 == j && j+3 == k && i%3 == j%3 && j%3 == k%3 {
		return true
	} else {
		return false
	}
}

// 4*3板块判断其中三个格子是否连续且在同一列 (i<j<k)
func checkSameColForFourToThree(i, j, k int32) bool {
	if i+4 == j && j+4 == k && i%4 == j%4 && j%4 == k%4 {
		return true
	} else {
		return false
	}
}
