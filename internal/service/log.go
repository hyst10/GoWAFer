package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_helper"
	"time"
)

type LogService struct {
	logRepository *repository.LogRepository
}

func NewLogService(r *repository.LogRepository) *LogService {
	return &LogService{logRepository: r}
}

func (c *LogService) FindLogs(days, hours int) api_helper.LogStatsResponse {
	items := c.logRepository.FindLog(days, hours)

	// 初始化返回的结构
	response := api_helper.LogStatsResponse{}

	// 初始化开始时间和结束时间、日期分钟格式、区间
	endTime := time.Now()
	var startTime time.Time
	var timeFormat string
	var interval time.Duration

	if hours > 0 {
		startTime = endTime.Add(time.Duration(-hours) * time.Hour)
		endTime = endTime.Add(1 * time.Minute)
		timeFormat = "15:04"
		interval = 1 * time.Minute
	} else if days == 1 {
		startTime = endTime.Add(-24 * time.Hour)
		if endTime.Minute()%5 != 0 {
			endTime = endTime.Add(time.Duration(5-endTime.Minute()%5) * time.Minute)
		}
		timeFormat = "15:04"
		interval = 5 * time.Minute
	} else {
		startTime = endTime.AddDate(0, 0, -days)
		// 将 endTime 向前扩展到下一个天开始的时间
		if endTime.Hour() != 0 || endTime.Minute() != 0 || endTime.Second() != 0 {
			endTime = endTime.AddDate(0, 0, 1)                                                                  // 先加一天
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location()) // 然后调整到当天开始
		}
		timeFormat = "2006-01-02"
		interval = 24 * time.Hour
	}

	// 按小时或天分组的键值映射
	counts200 := make(map[string]int)
	counts403 := make(map[string]int)
	for t := startTime; t.Before(endTime); t = t.Add(interval) {
		key := t.Format(timeFormat)
		response.TimeSegments = append(response.TimeSegments, key)
		counts200[key] = 0
		counts403[key] = 0
	}

	// 填充数据到时间分组
	for _, item := range items {
		key := item.CreatedAt.Format(timeFormat)
		// 根据状态码计数
		switch item.Status {
		case 200:
			counts200[key]++
		case 403:
			counts403[key]++
		}
	}

	// 根据排序后的时间段填充计数数据
	for _, segment := range response.TimeSegments {
		response.Count200s = append(response.Count200s, counts200[segment])
		response.Count403s = append(response.Count403s, counts403[segment])
	}

	return response
}

func (c *LogService) FindPaginatedLogs(page *pagination.Pages, keyword string) *pagination.Pages {
	items, count := c.logRepository.FindPaginated(page.Page, page.PerPage, keyword)
	page.Items = items
	page.Total = count
	return page
}
