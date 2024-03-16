package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

type BlockLogService struct {
	blockLogRepository *repository.BlockLogRepository
}

func NewBlockLogService(r *repository.BlockLogRepository) *BlockLogService {
	return &BlockLogService{blockLogRepository: r}
}

func (c *BlockLogService) FindPaginatedLogs(page *pagination.Pages, keyword string) *pagination.Pages {
	items, count := c.blockLogRepository.FindPaginated(page.Page, page.PerPage, keyword)
	page.Items = items
	page.Total = count
	return page
}
