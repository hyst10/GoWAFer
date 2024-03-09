package pagination

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	DefaultPageSize = 15        // 默认页面大小
	MaxPageSize     = 100       // 最大页面大小
	PageQuery       = "page"    // 页码参数名
	PageSizeQuery   = "perPage" // 页面大小参数名
)

type Pages struct {
	Items   interface{} `json:"items"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"perPage"`
}

// NewPage 实例化分页
func NewPage(page, pageSize, total int) *Pages {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	if page <= 0 {
		page = 1
	}
	return &Pages{
		Page:    page,
		PerPage: pageSize,
		Total:   total,
	}
}

// 将string类型转换为int，转换失败返回默认值
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

// NewFromRequest 从请求中实例化分页
func NewFromRequest(c *gin.Context) *Pages {
	page := parseInt(c.Query(PageQuery), 1)
	pageSize := parseInt(c.Query(PageSizeQuery), DefaultPageSize)
	return NewPage(page, pageSize, -1)
}
