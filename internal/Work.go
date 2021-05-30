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
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/spf13/cast"
	"io/ioutil"
	"strings"
	"time"
)

const PATH = "./data/"

func Worker() {
	//resultSlice := make([]schema.Result,0)
	hfValueMap := make(map[string]string, 0)
	locationSlice := make([]int, 0)
	allResultSlice := make([]schema.Result, 0)
	files, _ := ioutil.ReadDir(PATH)
	for _, f := range files {
		fmt.Println("正在处理" + f.Name())
		hfValue, location, resultSlice := GetValueByFileName(f.Name())
		locationSlice = append(locationSlice, cast.ToInt(location))
		hfValueMap[location] = hfValue
		allResultSlice = append(allResultSlice, resultSlice...)
	}

	fmt.Println(locationSlice, hfValueMap)
	WriteExcel(locationSlice,hfValueMap)

	/**
	locationSlice := make([]int, 0)

	//处理hf的值
	locationSlice = append(locationSlice, cast.ToInt(location))
	hfValueMap[location] = hfValue
	fmt.Println("111")
	fmt.Println(locationSlice, hfValueMap)
	fmt.Println("222")
	*/

}

func GetValueByFileName(fileName string) (hfValue, location string, resultContentSlice []schema.Result) {
	fileName = PATH + fileName
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("open %s failed:%s", fileName, err)
	}
	//获取文件location
	location = utils.GetFileNumber(fileName)

	//获取hf值
	hfValue = utils.GetFileHF(string(file))

	for _, line := range strings.Split(string(file), "\n") {
		if strings.Contains(line, "Isotropic") && strings.Contains(line, "Anisotropy") {
			resultSlice := strings.Fields(line)
			if resultSlice[1] == "H" || resultSlice[1] == "C" {
				result := schema.Result{
					Location: location,
					Sequence: resultSlice[0],
					Element:  resultSlice[1],
					Value:    resultSlice[4],
				}
				resultContentSlice = append(resultContentSlice, result)
			}

		}

	}

	return

}

func WriteExcel(locationSlice []int, hfValueMap map[string]string) {

	fileName := time.Now().Format("20060102150405") + ".xlsx"
	f := excelize.NewFile()
	Sheet2 := "Sheet2"
	index := f.NewSheet(Sheet2)
	f.SetCellValue(Sheet2,"C1","HF")
	countHF :=len(locationSlice)
	f.SetCellFormula(Sheet2,"G"+cast.ToString(countHF+2),"")
	sumFormula := fmt.Sprintf("SUM(G2:%s)","G"+cast.ToString(countHF+1))
	f.SetCellFormula(Sheet2,"G"+cast.ToString(countHF+2),sumFormula)
	for key, val := range locationSlice {
		yAxis := cast.ToString(key+2)
		f.SetCellValue(Sheet2,"B"+yAxis,val)
		f.SetCellValue(Sheet2,"C"+yAxis,hfValueMap[cast.ToString(val)])
		formulaD :=     "C"+yAxis+"-C2"
		f.SetCellFormula(Sheet2,"D"+yAxis,formulaD)
		formulaE := "D"+yAxis+"*627.5"
		f.SetCellFormula(Sheet2,"E"+yAxis,formulaE)
		formulaF :="-E"+yAxis+"/(0.0019858955*298.15)"
		f.SetCellFormula(Sheet2,"F"+yAxis,formulaF)
		formulaG :="EXP(F"+yAxis+")"
		f.SetCellFormula(Sheet2,"G"+yAxis,formulaG)
		formulaH :="G"+yAxis+"/G"+cast.ToString(countHF+2)
		f.SetCellFormula(Sheet2,"H"+yAxis,formulaH)

	}
	f.SetActiveSheet(index)
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	return

}

/**
func WriteExcel(resultMap chan schema.Result,hfResultMap chan schema.HFResult){

	fileName := "test.xlsx"
	f, err := excelize.OpenFile(fileName)
	defer f.Save()
	if err != nil {
		fmt.Println(err)
		return
	}
	index := f.NewSheet("Sheet2")

	for {
		hf := <- hfResultMap
		f.SetCellValue("Sheet2", "A"+CurrentLine,Content[val])
		Column2 :="B"+CurrentLine
		f.SetCellValue("Sheet2",Column2,val)
		formula3 :=	Column2+"-B1"
		f.SetCellFormula("Sheet2","C"+CurrentLine,formula3)
		formula4 :=	"C"+CurrentLine+"*627.5"
		f.SetCellFormula("Sheet2","D"+CurrentLine,formula4)
		formula5 :=	"-D"+CurrentLine+"/(0.0019858955*298.15)"
		f.SetCellFormula("Sheet2","E"+CurrentLine,formula5)
		formula6 :=	"EXP(E"+CurrentLine+")"
		f.SetCellFormula("Sheet2","F"+CurrentLine,formula6)
		formula7 :=	"F"+CurrentLine+"/"+"F"+sumLine
		f.SetCellFormula("Sheet2","G"+CurrentLine,formula7)

		numberLocation[Content[val]] = "G"+CurrentLine
		start++

	}

	sumFormula := fmt.Sprintf("SUM(F1:%s)","F"+cast.ToString(totalLine))
	f.SetCellFormula("Sheet2","F"+sumLine,sumFormula)
	f.SetActiveSheet(index)
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}
	return

}
*/
