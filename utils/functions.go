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


