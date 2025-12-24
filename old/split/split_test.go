package splitString

// https://www.cnblogs.com/sss4/p/12859027.html

import (
	"reflect"
	"testing"
)

// 测试用例1：以字符分割
func TestSplit(t *testing.T) {
	got := Newsplit("123N456", "N")
	want := []string{"123", "456"}
	//DeepEqual比较底层数组
	if !reflect.DeepEqual(got, want) {
		//如果got和want不一致说明你写得代码有问题
		t.Errorf("The values of %v is not %v\n", got, want)
	}

}

// 测试用例2：以标点符号分割
func TestPunctuationSplit(t *testing.T) {
	got := Newsplit("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.FailNow() //出错就stop别往下测了！
	}

}

// 测试用例3：增加分隔符的长度
func TestMultipleChartSplit(t *testing.T) {
	got := Newsplit("hellowbsdjshdworld", "bsdjshd")
	want := []string{"hellow", "world"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("无法通过多字符分隔符的测试！got: %v want:%v\n", got, want) //出错就stop别往下测了！
	}

}

//测试组:在1个函数中写多个测试用例，切支持灵活扩展！

type testCase struct {
	str      string
	separate string
	want     []string
}

var testGroup = []testCase{
	//测试用例1:单个英文字母
	testCase{
		str:      "123N456",
		separate: "N",
		want:     []string{"123", "456"},
	},
	//测试用例2：符号
	testCase{
		str:      "a:b:c",
		separate: ":",
		want:     []string{"a", "b", "c"},
	},
	//测试用例3：多个英文字母
	testCase{
		str:      "hellowbsdjshdworld",
		separate: "bsdjshd",
		want:     []string{"hellow", "world"},
	},
	//测试用例4:单个汉字
	testCase{
		str:      "山西运煤车煤运西山",
		separate: "山",
		want:     []string{"西运煤车煤运西"},
	},

	//测试用例4：多个汉字
	testCase{
		str:      "京北北京之北",
		separate: "北京",
		want:     []string{"京北", "之北"},
	},
}

func TestSplitOne(t *testing.T) {
	for _, test := range testGroup {
		got := Newsplit(test.str, test.separate)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("失败！got:%#v want:%#v\n", got, test.want)
		}
	}

}

// 子测试
type testCaseTwo struct {
	str      string
	separate string
	want     []string
}

var testGroupTwo = map[string]testCaseTwo{
	"punctuation": testCaseTwo{
		str:      "a:b:c",
		separate: ":",
		want:     []string{"a", "b", "c"},
	},
	"sigalLetter": testCaseTwo{
		str:      "123N456",
		separate: "N",
		want:     []string{"123", "456"},
	},

	"MultipleLetter": testCaseTwo{
		str:      "hellowbsdjshdworld",
		separate: "bsdjshd",
		want:     []string{"hellow", "world"},
	},
	"singalRune": testCaseTwo{
		str:      "山西运煤车煤运西山",
		separate: "山",
		want:     []string{"西运煤车煤运西"},
	},
	"multiplRune": testCaseTwo{
		str:      "京北北京之北",
		separate: "北京",
		want:     []string{"京北", "之北"},
	},
}

//测试用例函数
// 针对某1个测试用例进行单独测试
// go test -run=TestSplitTwo/punctuation

// 把测试覆盖率的详细信息输出到文件
// go test -cover -coverprofile=测试报告文件

//把测试报告输出到文件，就是为了分析测试结果，go内置的工具支持以HTML的方式打开测试报告文件！
// go tool cover -html=测试报告文件
// 打开的页面中每个用绿色标记的语句块表示被覆盖了，而红色的表示没有被覆盖。

func TestSplitTwo(t *testing.T) {
	for name, test := range testGroupTwo {
		//使用t参数的run方法
		t.Run(name, func(t *testing.T) {
			got := Newsplit(test.str, test.separate)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("失败！got:%#v want:%#v\n", got, test.want)
			}
		})
	}
}
