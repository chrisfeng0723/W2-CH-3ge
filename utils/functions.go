/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/5/28 14:50
 */

package utils

import (
	"regexp"
	"strings"
	"unicode"
)

//根据一个文件名字符串获取其中的数字
//eg:W-2_1.gjf.gjf.gjf.log
//_-和点号之间的数字
func GetFileNumber(fileName string) string {
	str := `[-|_]0*(\d*)\.`
	Regexp := regexp.MustCompile(str)
	params := Regexp.FindStringSubmatch(fileName)
	return params[1]
}

//获取文件的HF值
func GetFileHF(fileContent string) string {

	//str := `HF=(-?\d+.\d+)\\`
	str := `HF=(-?\d+.\d+)\\`
	Regexp := regexp.MustCompile(str)
	//去除空白字符
	temp := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, fileContent)
	params := Regexp.FindStringSubmatch(temp)
	if len(params) > 0 {
		return params[1]
	}
	return ""
}

//删除slice中的重复值
func RemoveDuplicate(slc []int) []int {
	result := make([]int, 0)         //存放返回的不重复切片
	tempMap := make(map[int]bool, 0) // 存放不重复主键
	for _, val := range slc {
		if _, ok := tempMap[val]; !ok {
			tempMap[val] = true
			result = append(result, val)
		}

	}
	return result
}

//将一个slice转换成一个map,value作为key，key作为value，仅适用于值不重复的slice

func TransferSliceToMap(slice []int) map[int]int {
	result := make(map[int]int, len(slice))
	for key, value := range slice {
		result[value] = key
	}
	return result
}
