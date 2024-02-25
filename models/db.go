package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=efosa.me user=postgres password=5005227a52c02361b7e95a1f5acfc7f0 dbname=jobs_db port=44553 sslmode=disable TimeZone=America/Los_Angeles"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if cache, err := cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,  // when u create/update/delete objects, invalidate cache
		CacheTTL:             10000, // 5000 ms
		// if length of objects retrieved one single time
		// exceeds this number, don't cache
		CacheMaxItemCnt: 20,
	}); err != nil {
		panic("failed to setup db cache")
	} else {
		conn.Use(cache)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(30)
	DB = db
}
