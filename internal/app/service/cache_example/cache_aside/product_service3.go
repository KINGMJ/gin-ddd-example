package cache_aside

import (
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"log"
	"time"
)

// 先更新数据库，再删除缓存
func (s *productService) UpdateProduct4(product *product.Product) error {
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	err := s.rdb.Del(s.ctx, cacheKey).Err()
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
	}
	// 模拟延迟，线程B读取到旧值
	time.Sleep(100 * time.Millisecond)
	// 更新数据库
	if err := s.productRepo.Update(s.db, product); err != nil {
		return err
	}
	// 再次删除
	err = s.rdb.Del(s.ctx, cacheKey).Err()
	if err != nil {
		log.Printf("删除缓存失败: %v", err)
	}
	return nil
}
