package cache_penetration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-ddd-example/internal/app/model/product"
	"gin-ddd-example/internal/app/repo"
	"gin-ddd-example/pkg/container"
	"hash/fnv"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

// 缓存穿透的应对方案
type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*product.Product, error)
}

type productService struct {
	container.Container
	productRepo repo.ProductRepo
	bloomFilter *BloomFilter // 布隆过滤器
	sf          *singleflight.Group
}

func NewProductService(
	container container.Container,
	productRepo repo.ProductRepo,
) ProductService {
	// 初始化布隆过滤器
	bf := NewBloomFilter(container.Rdb, "product_filter", 100000, 0.01)

	// 立即同步一次数据库数据
	if err := bf.SyncWithDB(container.Db); err != nil {
		container.Logs.Error("初始化布隆过滤器失败", zap.Error(err))
	}
	// 定期更新布隆过滤器
	// 如果是旁路缓存，需要使用kafka或者mysql binlog同步数据
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := bf.SyncWithDB(container.Db); err != nil {
				container.Logs.Error("同步布隆过滤器失败", zap.Error(err))
			}
		}
	}()

	return &productService{
		productRepo: productRepo,
		Container:   container,
		bloomFilter: bf,
		sf:          &singleflight.Group{},
	}
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*product.Product, error) {
	// 使用布隆过滤器检查是否存在
	cacheKey := fmt.Sprintf("product:%d", id)
	// 记录请求信息
	var requestID string
	if id, ok := ctx.Value("request_id").(string); ok {
		requestID = id
	} else {
		requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	goroutineID := getGoroutineID() // 获取当前的 goroutine id

	exists, err := s.bloomFilter.Exists(fmt.Sprintf("%d", id))
	if err != nil {
		s.Logs.Error(fmt.Sprintf("布隆过滤器检查失败 - RequestID: %s, GoroutineID: %d, ProductID: %d, Error: %v", requestID, goroutineID, id, err))
	} else if !exists {
		// 如果布隆过滤器中不存在，则直接返回
		return nil, ErrProductNotFound
	}
	productValue, err := s.getFromLocalCache(ctx, cacheKey)
	if err == nil {
		return productValue, nil
	}
	if !errors.Is(err, redis.Nil) {
		s.Logs.Error(fmt.Sprintf("读取缓存错误 - RequestID: %s, GoroutineID: %d, ProductID: %d, Error: %v", requestID, goroutineID, id, err))
	}

	// 3. 使用 singleflight 合并请求
	value, err, _ := s.sf.Do(cacheKey, func() (interface{}, error) {
		// 3.1 再次检查缓存
		if product, err := s.getFromLocalCache(ctx, cacheKey); err == nil {
			return product, nil
		}

		// 3.2 查询数据库
		product, err := s.productRepo.FindById(s.Db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 3.3 对空值也进行缓存，避免缓存穿透
				s.cacheNilValue(ctx, cacheKey)
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		// 3.4 更新缓存
		go func() {
			if err := s.updateCache(context.Background(), cacheKey, product); err != nil {
				s.Logs.Error("更新缓存失败", zap.Error(err))
			}
		}()

		return product, nil
	})

	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, ErrProductNotFound
	}
	return value.(*product.Product), nil
}

type BloomFilter struct {
	rdb    *redis.Client
	filter *bloom.BloomFilter
	key    string
}

func NewBloomFilter(rdb *redis.Client, key string, size uint, fp float64) *BloomFilter {
	filter := bloom.NewWithEstimates(size, fp)
	bf := &BloomFilter{
		rdb:    rdb,
		filter: filter,
		key:    key,
	}

	// 从 Redis 加载已存在的数据
	bf.loadFromRedis()
	return bf
}

func (bf *BloomFilter) Add(item string) error {
	// 添加到本地布隆过滤器
	bf.filter.AddString(item)
	// 同步到 Redis
	h := fnv.New64()
	h.Write([]byte(item))
	hash := h.Sum64()
	capValue := uint64(bf.filter.Cap())
	return bf.rdb.SetBit(context.Background(), bf.key, int64(hash%capValue), 1).Err()
}

func (bf *BloomFilter) Exists(item string) (bool, error) {
	return bf.filter.TestString(item), nil
}

func (bf *BloomFilter) loadFromRedis() error {
	// 从 Redis 加载数据到本地布隆过滤器
	// 实际实现可能需要根据具体需求调整
	return nil
}

func getGoroutineID() interface{} {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// getFromLocalCache 从本地缓存获取数据
func (s *productService) getFromLocalCache(ctx context.Context, key string) (*product.Product, error) {
	// 从 Redis 获取数据
	data, err := s.Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var product product.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

// cacheNilValue 缓存空值，设置较短的过期时间
func (s *productService) cacheNilValue(ctx context.Context, key string) {
	err := s.Rdb.Set(ctx, key, "nil", 5*time.Minute).Err()
	if err != nil {
		s.Logs.Error("缓存空值失败", zap.Error(err))
	}
}

// updateCache 更新缓存
func (s *productService) updateCache(ctx context.Context, key string, product *product.Product) error {
	// 序列化数据
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	// 设置缓存，使用随机过期时间避免缓存雪崩
	expiration := getRandomExpiration()
	return s.Rdb.Set(ctx, key, data, expiration).Err()
}

// 过期时间随机化
func getRandomExpiration() time.Duration {
	baseExpiration := 30 * time.Minute
	// 随机添加0-5分钟
	random := time.Duration(rand.Int63n(300)) * time.Second
	return baseExpiration + random
}

func (bf *BloomFilter) SyncWithDB(db *gorm.DB) error {
	var ids []int64
	if err := db.Model(&product.Product{}).Pluck("id", &ids).Error; err != nil {
		return err
	}

	for _, id := range ids {
		if err := bf.Add(fmt.Sprintf("%d", id)); err != nil {
			return err
		}
	}
	return nil
}
