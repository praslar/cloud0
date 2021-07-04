package ginext

import (
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	defaultPageSize = 20
	maxPageSize     = 500
)

// Pager represents a object that support paginate data in DB
// also parse request from client via gin.Context
type Pager struct {
	Page           int    `json:"page" form:"page"`
	PageSize       int    `json:"page_size" form:"page_size"`
	Sort           string `json:"sort" form:"sort"`
	TotalRows      int64  `json:"total_rows"`
	SortableFields []string
}

// SortableFieldsGetter represent a contract to all models that can give us a list of fields
//should be validated before sorting
//	type UserModel struct {
//		ID intName
//		String
//	}
//
//	func (u *UserModel) GetSortableFields() []string {
//		return []string{"id", "name" }
//	}
//
// the returned fields names should be database columns name instead of struct field names
type SortableFieldsGetter interface {
	GetSortableFields() []string
}

// NewPagerWithGinCtx initializes a new Pager from gin context by reading in order query, body request
func NewPagerWithGinCtx(c *gin.Context, logger *logrus.Entry) *Pager {
	pg := &Pager{}
	if err := c.ShouldBind(pg); err != nil && logger != nil {
		logger.WithError(err).Error("failed to parse pager request")
	}
	return pg
}

func (p *Pager) GetPage() int {
	if p.Page == 0 {
		return 1
	}
	return p.Page
}

func (p *Pager) GetOffset() int {
	return (p.GetPage() - 1) * p.PageSize
}

func (p *Pager) GetPageSize() int {
	if p.PageSize == 0 {
		return defaultPageSize
	}
	if p.PageSize > maxPageSize {
		return maxPageSize
	}
	return p.PageSize
}

// zerost is a empty struct (zero memory allowed), used for indexing map
type zerost struct{}

// GetOrder parses the order field then return Gorm order format
// 	eg. "name,-age" => "name asc, age desc"
func (p *Pager) GetOrder() string {
	if p.Sort == "" {
		return "id asc"
	}

	rawSegments := strings.Split(p.Sort, ",")
	var finalSortFields []string

	// sortable fields index
	var sortableFields = map[string]zerost{}
	for _, field := range p.SortableFields {
		sortableFields[field] = zerost{}
	}

	for _, segment := range rawSegments {
		segment = strings.TrimSpace(segment)

		var (
			fieldName string
			direction = "asc"
		)

		// convert :
		// 	-field -> field desc
		//	field -> field asc
		if strings.HasPrefix(segment, "-") {
			fieldName = segment[1:]
			direction = "desc"
		} else {
			fieldName = segment
		}

		if _, ok := sortableFields[fieldName]; ok {
			finalSortFields = append(finalSortFields, fieldName+" "+direction)
		}
	}

	return strings.Join(finalSortFields, ", ")
}

func (p *Pager) GetTotalPages() int {
	return int(math.Ceil(float64(p.TotalRows) / float64(p.GetPageSize())))
}

// DoQuery The execution will stop on count error then return that transaction
func (p *Pager) DoQuery(value interface{}, db *gorm.DB) *gorm.DB {
	var (
		totalRows int64
		tx        *gorm.DB
	)
	if tx = db.Count(&totalRows); tx.Error != nil {
		return tx
	}
	p.TotalRows = totalRows

	// get sortable fields if this model supports (just call in case empty page.SortableFields
	if len(p.SortableFields) == 0 {
		if getter, ok := value.(SortableFieldsGetter); ok {
			p.SortableFields = getter.GetSortableFields()
		}
	}

	return db.Offset(p.GetOffset()).Limit(p.GetPageSize()).Order(p.GetOrder()).Find(value)
}