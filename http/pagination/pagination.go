package pagination

import (
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	Total        int64
	Query        func(page, pageSize, showPageNum int64) string
	Page         int64
	PageSize     int64
	SkipNum      int64
	ShowPageNum  int64
	PageTotal    int64
	ShowPageList []int64
}

func Parse(r *http.Request, cnt int64) Pagination {
	pageTag := "page"
	pageStr := r.URL.Query().Get(pageTag)

	pageSizeTag := "page_size"
	pageSizeStr := r.URL.Query().Get(pageSizeTag)
	if pageSizeStr == "" {
		pageSizeTag = "pageSize"
		pageSizeStr = r.URL.Query().Get(pageSizeTag)
	}

	showPageNumTag := "show_page_num"
	showPageNumStr := r.URL.Query().Get(showPageNumTag)

	query := r.URL.Query()
	pg := Pagination{
		ShowPageList: []int64{},
	}

	pg.Total = cnt
	pg.Query = func(page, pageSize, showPageNum int64) string {
		query.Set(pageTag, strconv.FormatInt(page, 10))
		query.Set(pageSizeTag, strconv.FormatInt(pageSize, 10))
		query.Set(showPageNumTag, strconv.FormatInt(showPageNum, 10))
		return query.Encode()
	}

	if i, e := strconv.ParseInt(pageStr, 10, 64); e == nil {
		pg.Page = i
	}

	if i, e := strconv.ParseInt(pageSizeStr, 10, 64); e == nil {
		if i > 0 {
			pg.PageSize = i
		} else {
			pg.PageSize = 10
		}
	}

	if i, e := strconv.ParseInt(showPageNumStr, 10, 64); e == nil {
		if i > 0 {
			pg.ShowPageNum = i
		} else {
			pg.ShowPageNum = 5
		}
	}

	if pg.Page > 0 {
		pg.SkipNum = (pg.Page - 1) * pg.PageSize
	} else {
		pg.Page = 1
	}

	pg.PageTotal = int64(math.Ceil(float64(cnt) / float64(pg.PageSize)))

	for i := pg.Page - pg.ShowPageNum; i <= pg.Page+pg.ShowPageNum; i++ {
		pg.ShowPageList = append(pg.ShowPageList, i)
	}

	return pg
}
