package export

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Eden/go-gin-example/pkg/setting"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"time"
)

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

func ExcelImport(r io.Reader, sheetName string) ([][]string, error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	rows := xlsx.GetRows(sheetName)
	return rows, nil
}

func ExcelExport(filename string, sheetName string, titles []string, keys []string, setDatas []map[string]string) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return "", err
	}

	//titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, maps := range setDatas {
		values := []string{}
		//keys := []string{"id", "name", "created_by", "created_on", "modified_by", "modified_on"}
		for _, key := range keys {
			values = append(values, maps[key])
		}
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename = filename + "-" + time + ".xlsx"

	fullPath := GetExcelFullPath() + filename
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return filename, nil
}
