package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	_ "image/gif"
	_ "image/webp"
)

// 图片缓存项
type ImageCacheItem struct {
	Data        []byte    // 图片二进制数据
	ContentType string    // 内容类型
	Size        int64     // 图片大小
	LastAccess  time.Time // 最后访问时间
	AccessCount int       // 访问次数
	Key         string    // 缓存键
}

// 缓存策略接口
type CacheStrategy interface {
	Get(key string) (*ImageCacheItem, bool)
	Set(key string, item *ImageCacheItem)
	Remove(key string)
	Cleanup()
	Len() int
	Stats() CacheStats
}

// 缓存统计
type CacheStats struct {
	HitCount   int64
	MissCount  int64
	Size       int64
	ItemCount  int
	Evictions  int64
	MemoryUsed int64
}

// ==================== 实现1: LRU缓存策略 ====================
type LRUCache struct {
	capacity   int64                // 最大容量（字节）
	maxItems   int                  // 最大项目数
	current    int64                // 当前使用量
	items      map[string]*listNode // 存储节点
	list       *doublyLinkedList    // 双向链表
	mu         sync.RWMutex
	stats      CacheStats
	expiration time.Duration // 过期时间
}

type listNode struct {
	key   string
	item  *ImageCacheItem
	prev  *listNode
	next  *listNode
	added time.Time
}

type doublyLinkedList struct {
	head *listNode
	tail *listNode
	size int
}

func NewLRUCache(capacity int64, maxItems int, expiration time.Duration) *LRUCache {
	return &LRUCache{
		capacity:   capacity,
		maxItems:   maxItems,
		items:      make(map[string]*listNode),
		list:       newDoublyLinkedList(),
		expiration: expiration,
	}
}

func (c *LRUCache) Get(key string) (*ImageCacheItem, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		// 检查是否过期
		if c.expiration > 0 && time.Since(node.added) > c.expiration {
			c.removeNode(node)
			c.stats.MissCount++
			return nil, false
		}

		// 移动到链表头部（最近使用）
		c.list.moveToHead(node)
		node.item.LastAccess = time.Now()
		node.item.AccessCount++
		c.stats.HitCount++
		return node.item, true
	}

	c.stats.MissCount++
	return nil, false
}

func (c *LRUCache) Set(key string, item *ImageCacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	itemSize := int64(len(item.Data))

	// 如果单个项目超过容量，不缓存
	if itemSize > c.capacity {
		return
	}

	// 如果已存在，先移除
	if existing, exists := c.items[key]; exists {
		c.removeNode(existing)
	}

	// 确保有足够空间
	for (c.current+itemSize > c.capacity || len(c.items) >= c.maxItems) && len(c.items) > 0 {
		c.evict()
	}

	// 添加到链表头部
	node := &listNode{
		key:   key,
		item:  item,
		added: time.Now(),
	}
	c.items[key] = node
	c.list.addToHead(node)
	c.current += itemSize

	c.stats.ItemCount = len(c.items)
	c.stats.Size = c.current
}

func (c *LRUCache) evict() {
	if c.list.tail == nil {
		return
	}

	tail := c.list.tail
	c.removeNode(tail)
	c.stats.Evictions++
}

func (c *LRUCache) removeNode(node *listNode) {
	c.list.remove(node)
	delete(c.items, node.key)
	c.current -= int64(len(node.item.Data))
}

func (c *LRUCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		c.removeNode(node)
	}
}

func (c *LRUCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.expiration <= 0 {
		return
	}

	now := time.Now()
	for key, node := range c.items {
		if now.Sub(node.added) > c.expiration {
			c.removeNode(node)
			delete(c.items, key)
		}
	}
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *LRUCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// ==================== 实现2: 使用go-cache库（简单） ====================
type GoCacheWrapper struct {
	cache      *cache.Cache
	defaultTTL time.Duration
	stats      CacheStats
	mu         sync.RWMutex
}

func NewGoCacheWrapper(defaultTTL time.Duration, cleanupInterval time.Duration) *GoCacheWrapper {
	return &GoCacheWrapper{
		cache:      cache.New(defaultTTL, cleanupInterval),
		defaultTTL: defaultTTL,
	}
}

func (g *GoCacheWrapper) Get(key string) (*ImageCacheItem, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if x, found := g.cache.Get(key); found {
		if item, ok := x.(*ImageCacheItem); ok {
			item.LastAccess = time.Now()
			item.AccessCount++
			g.stats.HitCount++
			return item, true
		}
	}
	g.stats.MissCount++
	return nil, false
}

func (g *GoCacheWrapper) Set(key string, item *ImageCacheItem) {
	g.mu.Lock()
	defer g.mu.Unlock()

	item.Key = key
	g.cache.Set(key, item, g.defaultTTL)
	g.stats.ItemCount = g.cache.ItemCount()
	g.stats.MemoryUsed += int64(len(item.Data))
}

func (g *GoCacheWrapper) Remove(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cache.Delete(key)
}

func (g *GoCacheWrapper) Cleanup() {
	// go-cache会自动清理
}

func (g *GoCacheWrapper) Len() int {
	return g.cache.ItemCount()
}

func (g *GoCacheWrapper) Stats() CacheStats {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.stats
}

// ==================== 图片缓存管理器 ====================
type ImageCacheManager struct {
	strategy      CacheStrategy
	httpClient    *http.Client
	cloudEndpoint string // 云存储地址
	mu            sync.RWMutex

	// 统计信息
	totalRequests  int64
	cacheHits      int64
	cacheMisses    int64
	bandwidthSaved int64 // 节省的带宽（字节）
}

func NewImageCacheManager(strategy CacheStrategy, cloudEndpoint string) *ImageCacheManager {
	return &ImageCacheManager{
		strategy:      strategy,
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		cloudEndpoint: cloudEndpoint,
	}
}

// 获取图片（核心方法）
func (m *ImageCacheManager) GetImage(ctx context.Context, imagePath string) ([]byte, string, error) {
	m.mu.Lock()
	m.totalRequests++
	m.mu.Unlock()

	// 1. 尝试从缓存获取
	if item, found := m.strategy.Get(imagePath); found {
		m.mu.Lock()
		m.cacheHits++
		m.mu.Unlock()
		return item.Data, item.ContentType, nil
	}

	// 2. 缓存未命中，从云存储获取
	m.mu.Lock()
	m.cacheMisses++
	m.mu.Unlock()

	data, contentType, err := m.fetchFromCloud(ctx, imagePath)
	if err != nil {
		return nil, "", err
	}

	// 3. 存入缓存
	item := &ImageCacheItem{
		Data:        data,
		ContentType: contentType,
		Size:        int64(len(data)),
		LastAccess:  time.Now(),
		AccessCount: 1,
		Key:         imagePath,
	}
	m.strategy.Set(imagePath, item)

	// 4. 记录节省的带宽
	m.mu.Lock()
	m.bandwidthSaved += int64(len(data))
	m.mu.Unlock()

	return data, contentType, nil
}

// 从云存储获取图片
func (m *ImageCacheManager) fetchFromCloud(ctx context.Context, imagePath string) ([]byte, string, error) {
	url := m.cloudEndpoint + "/" + imagePath

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// 读取数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// 尝试检测图片类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	return data, contentType, nil
}

// 批量预加载热门图片
func (m *ImageCacheManager) PreloadImages(ctx context.Context, imagePaths []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 限制并发数

	for _, path := range imagePaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 检查是否已在缓存中
			if _, found := m.strategy.Get(p); !found {
				// 异步加载
				go func() {
					m.GetImage(ctx, p)
				}()
			}
		}(path)
	}
	wg.Wait()
}

// 清理过期缓存
func (m *ImageCacheManager) Cleanup() {
	m.strategy.Cleanup()
}

// 获取缓存统计
func (m *ImageCacheManager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := m.strategy.Stats()
	return map[string]interface{}{
		"total_requests":  m.totalRequests,
		"cache_hits":      m.cacheHits,
		"cache_misses":    m.cacheMisses,
		"hit_rate":        float64(m.cacheHits) / float64(m.totalRequests),
		"bandwidth_saved": m.bandwidthSaved,
		"cache_stats": map[string]interface{}{
			"item_count":  stats.ItemCount,
			"total_size":  stats.Size,
			"hit_count":   stats.HitCount,
			"miss_count":  stats.MissCount,
			"evictions":   stats.Evictions,
			"memory_used": stats.MemoryUsed,
		},
	}
}

// ==================== HTTP服务器封装 ====================
type ImageServer struct {
	cacheManager *ImageCacheManager
}

func NewImageServer(cacheManager *ImageCacheManager) *ImageServer {
	return &ImageServer{
		cacheManager: cacheManager,
	}
}

func (s *ImageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imagePath := r.URL.Path[len("/image/"):]
	if imagePath == "" {
		http.Error(w, "Image path required", http.StatusBadRequest)
		return
	}

	// 从缓存获取图片
	data, contentType, err := s.cacheManager.GetImage(r.Context(), imagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("X-Cache-Hit", "true") // 实际应该根据是否缓存命中设置

	// 返回图片数据
	w.Write(data)
}

// ==================== 辅助函数 ====================
func newDoublyLinkedList() *doublyLinkedList {
	return &doublyLinkedList{}
}

func (d *doublyLinkedList) addToHead(node *listNode) {
	node.next = d.head
	node.prev = nil

	if d.head != nil {
		d.head.prev = node
	}

	d.head = node
	if d.tail == nil {
		d.tail = node
	}
	d.size++
}

func (d *doublyLinkedList) remove(node *listNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		d.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		d.tail = node.prev
	}

	node.prev = nil
	node.next = nil
	d.size--
}

func (d *doublyLinkedList) moveToHead(node *listNode) {
	d.remove(node)
	d.addToHead(node)
}

// ==================== 使用示例 ====================
func main() {
	// 方式1：使用自定义LRU缓存
	lruCache := NewLRUCache(
		100*1024*1024,  // 100MB
		1000,           // 最多1000个图片
		30*time.Minute, // 30分钟过期
	)

	// 方式2：使用go-cache（更简单）
	// goCache := NewGoCacheWrapper(30*time.Minute, 5*time.Minute)

	// 创建缓存管理器
	cacheManager := NewImageCacheManager(
		lruCache, // 或 goCache
		"https://your-cloud-storage.com",
	)

	// 创建HTTP服务器
	server := NewImageServer(cacheManager)

	// 设置路由
	http.Handle("/image/", server)
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := cacheManager.GetStats()
		fmt.Fprintf(w, "Cache Stats: %+v\n", stats)
	})

	// 定期清理过期缓存
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cacheManager.Cleanup()
		}
	}()

	// 预加载热门图片
	ctx := context.Background()
	hotImages := []string{
		"products/featured.jpg",
		"users/avatar/default.png",
		"banners/main.webp",
	}
	cacheManager.PreloadImages(ctx, hotImages)

	fmt.Println("Image cache server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

// 图片处理扩展：压缩和格式转换
func (m *ImageCacheManager) GetOptimizedImage(ctx context.Context, imagePath string, width, height int, format string) ([]byte, string, error) {
	cacheKey := fmt.Sprintf("%s_%dx%d_%s", imagePath, width, height, format)

	// 尝试从缓存获取优化版本
	if item, found := m.strategy.Get(cacheKey); found {
		return item.Data, item.ContentType, nil
	}

	// 获取原始图片
	originalData, contentType, err := m.GetImage(ctx, imagePath)
	if err != nil {
		return nil, "", err
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(originalData))
	if err != nil {
		return nil, "", err
	}

	// 这里可以添加图片缩放、压缩等处理
	// 示例：简单的格式转换

	var processedData []byte
	var processedType string

	buf := new(bytes.Buffer)
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
		processedType = "image/jpeg"
	case "png":
		err = png.Encode(buf, img)
		processedType = "image/png"
	default:
		// 保持原格式
		processedData = originalData
		processedType = contentType
	}

	if err == nil && processedData == nil {
		processedData = buf.Bytes()
	}

	// 缓存优化版本
	item := &ImageCacheItem{
		Data:        processedData,
		ContentType: processedType,
		Size:        int64(len(processedData)),
		LastAccess:  time.Now(),
		AccessCount: 1,
		Key:         cacheKey,
	}
	m.strategy.Set(cacheKey, item)

	return processedData, processedType, nil
}
