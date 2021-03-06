package pagination

import (
	"errors"
	"gorm.io/gorm"
	"math"
)

type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy string
	ShowSQL bool
}

type Paginator struct {
	TotalRecord int64       `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Data        interface{} `json:"data"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

func Paging(p *Param, result interface{}) (*Paginator, error) {
	db := p.DB.Session(&gorm.Session{})

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int64
	var offset int

	go getCounts(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	if errGet := db.Scopes(paginate(offset, p.Limit, p.OrderBy)).Find(result).Error; errGet != nil &&
		!errors.Is(errGet, gorm.ErrRecordNotFound) {
		return nil, errGet
	}
	<-done

	paginator.TotalRecord = count
	paginator.Data = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator, nil
}

func getCounts(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}
func paginate(offset int, limit int, orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orderBy == "" {
			return db.Offset(offset).Limit(limit)
		} else {
			return db.Offset(offset).Limit(limit).Order(orderBy)
		}
	}
}
