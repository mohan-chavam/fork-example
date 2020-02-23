package service

import (
	"github.com/kuangcp/gobase/myth-bookkeeping/dal"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
)

func AutoMigrateAll() {
	db := dal.GetDB()
	db.AutoMigrate(&domain.Account{})
	db.AutoMigrate(&domain.Category{})
	db.AutoMigrate(&domain.Record{})
	db.AutoMigrate(&domain.BookKeeping{})
	db.AutoMigrate(&domain.TransferRecord{})
}