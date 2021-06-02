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
	"github.com/shopspring/decimal"
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
	hfValueMap := make(map[float64]string, 0)
	hfValueSlice := make([]float64, 0)
	//文件编号的slice
	locationSlice := make([]int, 0)
	//所有CH的值
	allResultSlice := make([]schema.Result, 0)
	//C的数值slice
	cResultSlice := make([]int, 0)
	//H的数值slice
	hResultSlice := make([]int, 0)

	files, _ := ioutil.ReadDir(PATH)
	for _, f := range files {
		fmt.Println("正在处理" + f.Name())
		hfValue, location, resultSlice, cSlice, hSlice := GetValueByFileName(f.Name())

		//处理hf值
		r, _ := decimal.NewFromString(hfValue)
		hfFloat, _ := r.Round(6).Float64()
		//有重复的则直接舍弃
		if _, ok := hfValueMap[hfFloat]; !ok {
			locationSlice = append(locationSlice, cast.ToInt(location))
			hfValueMap[hfFloat] = location
			hfValueSlice = append(hfValueSlice, hfFloat)
			//所有结果
			allResultSlice = append(allResultSlice, resultSlice...)
			cResultSlice = append(cResultSlice, cSlice...)
			hResultSlice = append(hResultSlice, hSlice...)
		}

	}

	//排序文件编号，写入hf值
	sort.Ints(locationSlice)
	sort.Float64s(hfValueSlice)
	//排序CH顺序
	uniqueCSlice := utils.RemoveDuplicate(cResultSlice)
	uniqueHSlice := utils.RemoveDuplicate(hResultSlice)
	sort.Ints(uniqueHSlice)
	sort.Ints(uniqueCSlice)

	//fmt.Println(hfValueSlice,hfValueMap)

	WriteExcel(locationSlice, hfValueSlice, hfValueMap, uniqueCSlice, uniqueHSlice, allResultSlice)

}

func GetValueByFileName(fileName string) (hfValue, location string, resultContentSlice []schema.Result, cResultSlice, hResultSlice []int) {
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
				if resultSlice[1] == "H" {
					hResultSlice = append(hResultSlice, cast.ToInt(resultSlice[0]))
				}

				if resultSlice[1] == "C" {
					cResultSlice = append(cResultSlice, cast.ToInt(resultSlice[0]))
				}
				resultContentSlice = append(resultContentSlice, result)
			}

		}

	}

	return

}

func WriteExcel(locationSlice []int, hfSlice []float64, hfValueMap map[float64]string, uniqueCSlice, uniqueHSlice []int, allResultSlice []schema.Result) {

	fileName := time.Now().Format("20060102150405") + ".xlsx"
	fmt.Println("数据处理中....")
	f := excelize.NewFile()
	Sheet2 := "Sheet2"
	Sheet1 := "Sheet1"
	index := f.NewSheet(Sheet1)
	f.NewSheet(Sheet2)
	f.SetCellValue(Sheet2, "C1", "HF")
	countHF := len(locationSlice)
	f.SetCellFormula(Sheet2, "G"+cast.ToString(countHF+2), "")
	sumFormula := fmt.Sprintf("SUM(G2:%s)", "G"+cast.ToString(countHF+1))
	f.SetCellFormula(Sheet2, "G"+cast.ToString(countHF+2), sumFormula)
	//fmt.Println(hfValueMap, hfSlice)
	for key, val := range hfSlice {
		yAxis := cast.ToString(key + 2)
		f.SetCellValue(Sheet2, "B"+yAxis, hfValueMap[val])
		f.SetCellValue(Sheet2, "C"+yAxis, val)
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

	fmt.Println("HF处理完毕")
	//处理sheet1的值
	for key, val := range locationSlice {
		coordinate, _ := excelize.CoordinatesToCellName(key+2, 1)
		f.SetCellValue(Sheet1, coordinate, val)
	}
	//第一行最后一列
	lastCoordinate, _ := excelize.CoordinatesToCellName(len(locationSlice)+2, 1)
	f.SetCellValue(Sheet1, lastCoordinate, "加权")

	locationMap := utils.TransferSliceToMap(locationSlice)
	cMap := utils.TransferSliceToMap(uniqueCSlice)
	hMap := utils.TransferSliceToMap(uniqueHSlice)
	cCount := len(uniqueCSlice)
	//hCount := len(uniqueHSlice)
	var x, y int
	for _, value := range allResultSlice {
		y = locationMap[cast.ToInt(value.Location)] + 2
		if value.Element == "C" {
			x = cMap[cast.ToInt(value.Sequence)] + 2
		} else {
			x = hMap[cast.ToInt(value.Sequence)] + cCount + 3
		}
		f.SetCellValue(Sheet1, "A"+cast.ToString(x), value.Element+value.Sequence)
		coordinate, _ := excelize.CoordinatesToCellName(y, x)
		f.SetCellValue(Sheet1, coordinate, value.Value)
	}

	fmt.Println("CH数据数据处理完毕...")
	//通过按照float排序后hfslice去找出每个location所在位置
	hfCalcMap := make(map[string]string, len(hfSlice))
	for fkey, fval := range hfSlice {
		hfCalcMap[hfValueMap[fval]] = cast.ToString(fkey + 2)
	}
	//计算最后的加权值，分别通过C,H

	for key, _ := range uniqueCSlice {
		sumFormulaSlice := make([]string, 0)
		for lkey, lval := range locationSlice {
			location, _ := excelize.CoordinatesToCellName(lkey+2, key+2)
			lvalStr := cast.ToString(lval)
			temp := fmt.Sprintf("%s*Sheet2!H%s", location, hfCalcMap[lvalStr])
			sumFormulaSlice = append(sumFormulaSlice, temp)
		}
		resultFormula := fmt.Sprintf("SUM(%s)", strings.Join(sumFormulaSlice, ","))
		resultLocation, _ := excelize.CoordinatesToCellName(len(locationSlice)+2, key+2)
		f.SetCellFormula("Sheet1", resultLocation, resultFormula)
	}

	//H的行值要加上C所有行值
	for key, _ := range uniqueHSlice {
		sumFormulaSlice := make([]string, 0)
		for lkey, lval := range locationSlice {
			location, _ := excelize.CoordinatesToCellName(lkey+2, key+2+cCount+1)
			lvalStr := cast.ToString(lval)
			temp := fmt.Sprintf("%s*Sheet2!H%s", location, hfCalcMap[lvalStr])
			sumFormulaSlice = append(sumFormulaSlice, temp)
		}
		resultFormula := fmt.Sprintf("SUM(%s)", strings.Join(sumFormulaSlice, ","))
		resultLocation, _ := excelize.CoordinatesToCellName(len(locationSlice)+2, key+2+cCount+1)
		f.SetCellFormula("Sheet1", resultLocation, resultFormula)
	}

	fmt.Println("CH加权数据处理完毕...")

	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Excel写入数据完毕...")
	return

}
