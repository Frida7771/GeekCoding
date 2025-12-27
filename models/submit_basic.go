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

// UpdateSubmitStatusInRedis 更新 Redis 中的提交状态
func UpdateSubmitStatusInRedis(identity string, status int) error {
	ctx := context.Background()
	key := "submit_status:" + identity

	// 先获取现有数据
	data, err := RDB.Get(ctx, key).Result()
	if err != nil {
		// 如果 Redis 中没有，就不更新了
		return nil
	}

	var statusInfo SubmitStatusInfo
	if err := json.Unmarshal([]byte(data), &statusInfo); err != nil {
		return err
	}

	// 更新状态和时间
	statusInfo.Status = status
	statusInfo.UpdatedAt = time.Now().Unix()

	// 保存回 Redis
	newData, err := json.Marshal(statusInfo)
	if err != nil {
		return err
	}

	return RDB.Set(ctx, key, newData, 24*time.Hour).Err()
}
