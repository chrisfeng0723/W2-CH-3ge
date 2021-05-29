/**
 * @Author: fxl
 * @Description: 
 * @File:  Result.go
 * @Version: 1.0.0
 * @Date: 2021/5/29 22:06
 */
package schema

type Result struct {
	Location string
	Sequence string
	Element  string
	Value    string
}

type HFResult struct {
	Location string
	Value string
}
