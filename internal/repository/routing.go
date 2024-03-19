package repository

import (
	"GoWAFer/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type RoutingRepository struct {
	db *gorm.DB
}

func NewRoutingRepository(db *gorm.DB) *RoutingRepository {
	return &RoutingRepository{db: db}
}

func (r *RoutingRepository) Create(route *model.Routing) error {
	return r.db.Create(route).Error
}

func (r *RoutingRepository) Update(route *model.Routing) error {
	return r.db.Save(route).Error
}

func (r *RoutingRepository) Delete(route *model.Routing) error {
	return r.db.Delete(route).Error
}

func (r *RoutingRepository) FindPaginated(pageIndex, pageSize int, routeType, keyword string) ([]model.Routing, int) {
	var routes []model.Routing
	var count int64
	query := r.db.Model(&model.Routing{}).Where("Type = ?", routeType)

	if keyword != "" {
		query = query.Where("Route LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	query.Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&routes)

	return routes, int(count)
}

func (r *RoutingRepository) FindByID(id uint) (*model.Routing, error) {
	var route *model.Routing
	err := r.db.First(&route, id).Error
	if err != nil {
		return nil, err
	}
	return route, err
}

func (r *RoutingRepository) IsExist(route, method string) (*model.Routing, error) {
	var current model.Routing
	if err := r.db.Where("Route = ? AND Method = ?", route, method).First(&current).Error; err != nil {
		return nil, err
	}
	return &current, nil
}
