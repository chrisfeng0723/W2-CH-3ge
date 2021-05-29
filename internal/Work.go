/**
* @Author:fengxinlei
* @Description:
* @Version 1.0.0
* @Date: 2021/5/28 15:14
 */

package internal

import (
	"W2-CH-3ge/internal/schema"
	"W2-CH-3ge/utils"
	"fmt"
	"io/ioutil"
	"strings"
)

func Worker(fileName string, resultMap chan schema.Result, hfResultMap chan schema.HFResult) {
	fileName = "./data/" + fileName
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("open %s failed:%s", fileName, err)
	}
	//获取文件location
	location := utils.GetFileNumber(fileName)

	//获取hf值
	hfValue := utils.GetFileHF(string(file))
	//hf对应关系值获取
	hfResultMap <- schema.HFResult{
		Location: location,
		Value:    hfValue,
	}

	for _, line := range strings.Split(string(file), "\n") {
		if strings.Contains(line, "Isotropic") && strings.Contains(line, "Anisotropy") {
			resultSlice := strings.Fields(line)
			if resultSlice[1] == "H" || resultSlice[1] == "C" {
				resultMap <- schema.Result{
					Location: location,
					Sequence: resultSlice[0],
					Element:  resultSlice[1],
					Value:    resultSlice[4],
				}
			}

		}

	}

}


func WriteExcel(resultMap chan schema.Result,hfResultMap chan schema.HFResult){

}
