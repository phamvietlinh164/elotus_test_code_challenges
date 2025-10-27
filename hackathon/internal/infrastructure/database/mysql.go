package database

import (
	"context"
	"fmt"
	"time"

	"hackathon/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg config.Database) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type result struct {
		db  *gorm.DB
		err error
	}

	ch := make(chan result, 1)
	go func() {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		ch <- result{db, err}
	}()

	select {
	case <-ctx.Done():
		panic("database connection timeout after 10s")
	case r := <-ch:
		if r.err != nil {
			panic(fmt.Sprintf("failed to connect database: %v", r.err))
		}
		return r.db
	}
}
