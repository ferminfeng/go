package splitString

import "strings"

//Newsplit 切割字符串
//example:
//abc,b=>[ac]
func Newsplit(str, sep string) (des []string) {
	index := strings.Index(str, sep)
	for index > -1 {
		sectionBefor := str[:index]
		if len(sectionBefor) >= 1 {
			des = append(des, sectionBefor)
		}
		str = str[index+len(sep):]
		index = strings.Index(str, sep)
	}
	//最后1
	if len(str) >= 1 {
		des = append(des, str)
	}

	return
}
