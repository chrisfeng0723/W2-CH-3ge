/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/5/28 14:50
 */

package utils

import "regexp"

//根据一个文件名字符串获取其中的数字
//eg:W-2_1.gjf.gjf.gjf.log
//_-和点号之间的数字
func GetFileNumber(fileName string) string {
	str := `[-|_]0*(\d*)\.`
	Regexp := regexp.MustCompile(str)
	params := Regexp.FindStringSubmatch(fileName)
	return params[1]
}
