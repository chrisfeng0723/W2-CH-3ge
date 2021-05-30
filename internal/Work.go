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
	_ "github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

const PATH = "./data/"

func Worker() {
	//resultSlice := make([]schema.Result,0)
	//hf的map值，key为所在的文件，value为hf值
	hfValueMap := make(map[string]string, 0)
	//文件编号的slice
	locationSlice := make([]int, 0)
	//所有CH的值
	allResultSlice := make([]schema.Result, 0)
	//C的数值slice
	cResultSlice :=make([]int,0)
	//H的数值slice
	hResultSlice :=make([]int,0)

	files, _ := ioutil.ReadDir(PATH)
	for _, f := range files {
		fmt.Println("正在处理" + f.Name())
		hfValue, location, resultSlice,cSlice,hSlice := GetValueByFileName(f.Name())
		locationSlice = append(locationSlice, cast.ToInt(location))
		hfValueMap[location] = hfValue
		allResultSlice = append(allResultSlice, resultSlice...)
		cResultSlice = append(cResultSlice,cSlice...)
		hResultSlice = append(hResultSlice,hSlice...)
	}

	//排序文件编号，写入hf值
	sort.Ints(locationSlice)
	//排序CH顺序
	uniqueCSlice :=utils.RemoveDuplicate(cResultSlice)
	uniqueHSlice :=utils.RemoveDuplicate(hResultSlice)
	sort.Ints(uniqueHSlice)
	sort.Ints(uniqueCSlice)

	fmt.Println(uniqueCSlice,uniqueHSlice)


	WriteExcel(locationSlice, hfValueMap,uniqueCSlice,uniqueHSlice,allResultSlice)



}


func GetValueByFileName(fileName string) (hfValue, location string, resultContentSlice []schema.Result,cResultSlice,hResultSlice []int) {
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
				if resultSlice[1] == "H"{
					hResultSlice = append(hResultSlice,cast.ToInt(resultSlice[0]))
				}

				if resultSlice[1] == "C"{
					cResultSlice = append(cResultSlice,cast.ToInt(resultSlice[0]))
				}
				resultContentSlice = append(resultContentSlice, result)
			}

		}

	}

	return

}

func WriteExcel(locationSlice []int, hfValueMap map[string]string,uniqueCSlice,uniqueHSlice []int,allResultSlice []schema.Result) {

	fileName := time.Now().Format("20060102150405") + ".xlsx"
	f := excelize.NewFile()
	Sheet2 := "Sheet2"
	Sheet1 := "Sheet1"
	index := f.NewSheet(Sheet1)
	f.SetCellValue(Sheet2, "C1", "HF")
	countHF := len(locationSlice)
	f.SetCellFormula(Sheet2, "G"+cast.ToString(countHF+2), "")
	sumFormula := fmt.Sprintf("SUM(G2:%s)", "G"+cast.ToString(countHF+1))
	f.SetCellFormula(Sheet2, "G"+cast.ToString(countHF+2), sumFormula)
	for key, val := range locationSlice {
		yAxis := cast.ToString(key + 2)
		f.SetCellValue(Sheet2, "B"+yAxis, val)
		f.SetCellValue(Sheet2, "C"+yAxis, hfValueMap[cast.ToString(val)])
		formulaD := "C" + yAxis + "-C2"
		f.SetCellFormula(Sheet2, "D"+yAxis, formulaD)
		formulaE := "D" + yAxis + "*627.5"
		f.SetCellFormula(Sheet2, "E"+yAxis, formulaE)
		formulaF := "-E" + yAxis + "/(0.0019858955*298.15)"
		f.SetCellFormula(Sheet2, "F"+yAxis, formulaF)
		formulaG := "EXP(F" + yAxis + ")"
		f.SetCellFormula(Sheet2, "G"+yAxis, formulaG)
		formulaH := "G" + yAxis + "/G" + cast.ToString(countHF+2)
		f.SetCellFormula(Sheet2, "H"+yAxis, formulaH)

	}
	f.SetActiveSheet(index)




	for key,val := range locationSlice{
		coordinate,_ :=excelize.CoordinatesToCellName(key+2,1)

		f.SetCellValue(Sheet1,coordinate,val)
	}
	//第一行最后一列
	lastCoordinate,_ :=excelize.CoordinatesToCellName(len(locationSlice)+2,1)
	f.SetCellValue(Sheet1,lastCoordinate,"加权")

	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	return

}

func GetCoorddinate(){

}