package repo

import (
	"fmt"
	"time"
)

// структуры, методы и интерфейсы для абстрагирования параметров поиска
type DateSearchParams struct {
	Date time.Time
}

func (dp *DateSearchParams) GetQueryData() *QueryData {
	return &QueryData{
		Param:     dp.Date.Format(timeTemplate),
		Condition: "WHERE date LIKE :search",
	}
}

type TextSearchParams struct {
	Text string
}

func (tp *TextSearchParams) GetQueryData() *QueryData {
	return &QueryData{
		Param:     fmt.Sprintf("%%%s%%", tp.Text),
		Condition: "WHERE title LIKE :search OR comment LIKE :search",
	}
}

type SearchQueryData interface {
	GetQueryData() *QueryData
}

func QueryDataFromString(search string) SearchQueryData {
	searchDate, err := time.Parse("02.01.2006", search)
	if err != nil {
		return &TextSearchParams{Text: search}
	} else {
		return &DateSearchParams{Date: searchDate}
	}
}

type QueryData struct {
	Param     string
	Condition string
}
