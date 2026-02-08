// Package utils 提供JavaScript安全的雪花ID生成器
//
// 功能特性：
// 1. 生成的ID在JavaScript Number.MAX_SAFE_INTEGER (2^53-1) 范围内，前端可安全使用
// 2. 全局单例模式，开箱即用，无需手动初始化
// 3. 基于秒级时间戳，支持每秒生成65536个唯一ID
// 4. 自动识别机器ID，支持最多64台机器分布式部署
// 5. 时钟回拨保护，小幅回拨自动等待，大幅回拨触发panic
//
// ID结构（53位以内）：
// +------------------+------------+----------------+
// |   时间戳(31位)    | 机器ID(6位) |   序列号(16位)  |
// +------------------+------------+----------------+
// | 从2026-01-01开始  |   0-63     |    0-65535     |
// +------------------+------------+----------------+
//
// 性能指标：
// - 单机QPS: 65536 (每秒)
// - 分布式节点: 64台
// - 理论寿命: ~68年 (2026-2094)
//
// 使用示例：
//
//	id := utils.NewJSSnowId()
//	// 生成的ID可直接传递给前端JavaScript使用，无精度丢失
package utils

import (
	"fmt"
	"hash/fnv"
	"net"
	"os"
	"sync"
	"time"
)

const (
	epoch        int64 = 1767225600       // 2026-01-01 00:00:00 UTC
	workerIDBits uint8 = 6                // 64 台机器 (2^6)
	sequenceBits uint8 = 16               // 65536 QPS (2^16，每秒最多生成65536个ID)
	maxSafeInt   int64 = 9007199254740991 // JavaScript Number.MAX_SAFE_INTEGER (2^53 - 1)
)

type worker struct {
	mu        sync.Mutex
	lastStamp int64
	workerID  int64
	sequence  int64
}

var globalWorker *worker

func init() {
	id := autoGenerateWorkerID()
	globalWorker = &worker{
		workerID:  id,
		lastStamp: -1,
	}
	// 启动打印，方便调试
	fmt.Printf("[IDGen] Initialized Machine ID: %d (JS-Safe Mode)\n", id)
}

func ID() int64 {
	return globalWorker.nextID()
}

func (w *worker) nextID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().Unix()

	// 1. 时钟回拨保护
	if now < w.lastStamp {
		diff := w.lastStamp - now
		if diff > 5 {
			// 如果回拨超过5秒，认为是严重问题
			panic(fmt.Sprintf("Clock moved backwards seriously! Rejecting requests for %d seconds", diff))
		}
		// 小幅回拨，等待追上
		time.Sleep(time.Duration(diff) * time.Second)
		now = time.Now().Unix()
	}

	// 2. 相同秒内生成
	if now == w.lastStamp {
		w.sequence = (w.sequence + 1) & (-1 ^ (-1 << sequenceBits))
		if w.sequence == 0 {
			// 序列号耗尽，休眠等待下一秒（避免CPU空转）
			time.Sleep(time.Millisecond * 100)
			for now <= w.lastStamp {
				time.Sleep(time.Millisecond * 10)
				now = time.Now().Unix()
			}
		}
	} else {
		// 3. 进入新的一秒，序列号重置
		w.sequence = 0
	}

	w.lastStamp = now

	// 拼接逻辑 (确保结果在 53 位以内)
	id := (now-epoch)<<(workerIDBits+sequenceBits) | (w.workerID << sequenceBits) | w.sequence

	// 检查是否超出JavaScript安全整数范围
	if id > maxSafeInt {
		panic(fmt.Sprintf("Generated ID %d exceeds JavaScript MAX_SAFE_INTEGER (%d)", id, maxSafeInt))
	}

	return id
}

func autoGenerateWorkerID() int64 {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return int64(ipnet.IP.To4()[3]) % 64
		}
	}
	name, _ := os.Hostname()
	h := fnv.New32a()
	h.Write([]byte(name))
	return int64(h.Sum32() % 64)
}
