package mysql

import (
	"context"

	"gorm.io/gorm"
)

type VIP struct {
	gorm.Model
	UserID   uint `gorm:"user_id"`  //VIP用户id
	Duration uint `gorm:"duration"` //VIP有效时间,单位是天
}

func (VIP) TableName() string {
	return "vip"
}

// 新增VIP
func CreateVIP(ctx context.Context, vip *VIP) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(vip).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 从数据库中查询VIP信息
func GetVIPByUserID(ctx context.Context, userID int64) (*VIP, error) {
	vip := new(VIP)
	err := GetDB().WithContext(ctx).Where("user_id = ?", userID).Find(&vip).Error
	if err != nil {
		zapLogger.Error(err.Error())
		return nil, err
	}
	return vip, nil
}
