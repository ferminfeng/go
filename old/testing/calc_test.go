package main

import "testing"

func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}

	// 错误的测试用例
	//if ans := Add(4, 9); ans != 3 {
	//	t.Errorf("4 + 9 expected be 13, but %d got", ans)
	//}

	if ans := Add(-10, -20); ans != -30 {
		t.Errorf("-10 + -20 expected be -30, but %d got", ans)
	}

	// 错误的测试用例
	//if ans := Add(-30, -40); ans != -30 {
	//	t.Errorf("-30 + -40 expected be -70, but %d got", ans)
	//}
}

func TestMul(t *testing.T) {
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}

	})
	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}
