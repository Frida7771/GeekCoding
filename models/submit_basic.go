package models

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Submit_Basic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`
	UserBasic       *User_Basic   `gorm:"foreignKey:identity;references:user_identity"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	Path            string        `gorm:"column:path;type:varchar(255)" json:"path"`
	Status          int           `gorm:"column:status;type:tinyint(1)" json:"status"`
}

func (table *Submit_Basic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(Submit_Basic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")

	if problemIdentity != "" {
		tx = tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx = tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}

	return tx
}

// SubmitStatusInfo 提交状态信息（用于 Redis 存储）
type SubmitStatusInfo struct {
	Identity        string `json:"identity"`
	ProblemIdentity string `json:"problem_identity"`
	UserIdentity    string `json:"user_identity"`
	Status          int    `json:"status"` // -1待判断 1答案正确 2答案错误 3运行超时 4运行超内存 5编译错误
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// SaveSubmitStatusToRedis 保存提交状态到 Redis
func SaveSubmitStatusToRedis(submit *Submit_Basic) error {
	ctx := context.Background()
	key := "submit_status:" + submit.Identity

	now := time.Now()
	var createdAt int64
	if !submit.CreatedAt.IsZero() {
		createdAt = submit.CreatedAt.Unix()
	} else {
		createdAt = now.Unix()
	}

	statusInfo := SubmitStatusInfo{
		Identity:        submit.Identity,
		ProblemIdentity: submit.ProblemIdentity,
		UserIdentity:    submit.UserIdentity,
		Status:          submit.Status,
		CreatedAt:       createdAt,
		UpdatedAt:       now.Unix(),
	}

	data, err := json.Marshal(statusInfo)
	if err != nil {
		return err
	}

	// 保存到 Redis，过期时间 24 小时
	return RDB.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetSubmitStatusFromRedis 从 Redis 获取提交状态
func GetSubmitStatusFromRedis(identity string) (*SubmitStatusInfo, error) {
	ctx := context.Background()
	key := "submit_status:" + identity

	data, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var statusInfo SubmitStatusInfo
	if err := json.Unmarshal([]byte(data), &statusInfo); err != nil {
		return nil, err
	}

	return &statusInfo, nil
}

// UpdateSubmitStatusInRedis 更新 Redis 中的提交状态（原子操作）
func UpdateSubmitStatusInRedis(identity string, status int) error {
	ctx := context.Background()
	key := "submit_status:" + identity
	updatedAt := time.Now().Unix()

	// 使用 Lua 脚本保证原子性：获取、解析、更新、保存
	luaScript := `
		local key = KEYS[1]
		local new_status = tonumber(ARGV[1])
		local updated_at = tonumber(ARGV[2])
		local expire_seconds = tonumber(ARGV[3])
		
		-- 获取现有数据
		local data = redis.call('GET', key)
		if not data then
			-- 如果 Redis 中没有，返回 0（不更新）
			return 0
		end
		
		-- 解析 JSON（Redis Lua 支持 cjson）
		local cjson = cjson or require('cjson')
		local status_info = cjson.decode(data)
		
		-- 更新状态和时间戳
		status_info.status = new_status
		status_info.updated_at = updated_at
		
		-- 序列化回 JSON
		local new_data = cjson.encode(status_info)
		
		-- 保存到 Redis
		redis.call('SET', key, new_data, 'EX', expire_seconds)
		
		-- 返回 1（成功）
		return 1
	`

	result, err := RDB.Eval(ctx, luaScript, []string{key},
		status, updatedAt, 24*3600).Result()
	if err != nil {
		// Lua 脚本执行失败，可能是 JSON 解析错误或 Redis 中没有数据
		// 返回 nil 表示不更新（降级策略）
		return nil
	}

	// Lua 脚本返回 1 表示成功，0 表示数据不存在
	if result.(int64) == 0 {
		return nil
	}

	return nil
}
