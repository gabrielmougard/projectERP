package utils

type Paginator struct {
	CurrentPage     int64 `json:"currentPage"`     
	NextPage        int64 `json:"nextPage"`        
	PrePage         int64 `json:"prePage"`        
	PageSize        int64 `json:"pageSize"`        
	CurrentPageSize int64 `json:"currentPageSize"` 
	TotalPage       int64 `json:"totalPage"`       
	TotalCount      int64 `json:"totalCount"`      
	FirstPage       bool  `json:"firstPage"`      
	LastPage        bool  `json:"lastPage"`       
	// PageList        []int64 `json:"pageList"`       
	Max int64
}

func GenPaginator(limit, offset, count int64) Paginator {
	var paginator Paginator
	paginator.TotalCount = count
	paginator.TotalPage = (count + limit - 1) / limit
	paginator.PageSize = limit
	if offset == 0 {
		paginator.FirstPage = true
	} else {
		paginator.FirstPage = false
	}
	if offset == paginator.TotalPage {
		paginator.LastPage = true
	} else {
		paginator.LastPage = false
	}
	if paginator.TotalCount > 0 && paginator.CurrentPage > 0 {
		paginator.Max = paginator.TotalCount / paginator.CurrentPage
	} else {
		paginator.Max = 0
	}
	return paginator

}