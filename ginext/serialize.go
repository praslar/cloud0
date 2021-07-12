package ginext

import (
	"encoding/json"
	"math"
)

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	Pages    int `json:"pages"`
	Total    int `json:"total"`
}

func (p *Pagination) MarshalJSON() ([]byte, error) {
	type Alias Pagination
	if p.Pages == 0 && p.Total != 0 && p.PageSize != 0 {
		p.Pages = TotalItemsToPages(p.Total, p.PageSize)
	}
	return json.Marshal((*Alias)(p))
}

// ResponseMeta is a general response metadata
type ResponseMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type GeneralResponse struct {
	Data  interface{}   `json:"data,omitempty"`
	Meta  *ResponseMeta `json:"meta,omitempty"`
	Error interface{}   `json:"error,omitempty"`
}

func NewResponseWithMeta(data interface{}, meta *ResponseMeta) *GeneralResponse {
	return &GeneralResponse{
		Data: data,
		Meta: meta,
	}
}

func NewPaginatedResponse(data interface{}, pager *Pager) *GeneralResponse {
	return &GeneralResponse{
		Data: data,
		Meta: &ResponseMeta{
			Page:       pager.GetPage(),
			TotalPages: pager.GetTotalPages(),
			PageSize:   pager.GetPageSize(),
			TotalRows:  pager.TotalRows,
		},
	}
}

// TotalItemsToPages converts to items to pages by page size
func TotalItemsToPages(totalItems, pageSize int) int {
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}
