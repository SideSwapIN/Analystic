package model

import (
	"time"

	"github.com/SideSwapIN/Analystic/internal/db"
	"gorm.io/plugin/soft_delete"
)

type RouterMethod int

const (
	AddLiquidityRouterMethod RouterMethod = iota + 1
	RemoveLiquidityRouterMethod
	SwapRouterMethod
)

type BaseModel struct {
	ID        int64                 `gorm:"column:id;primaryKey;comment:自增ID"`
	CreatedAt time.Time             `gorm:"column:created_at;type:timestamp;<-:false;comment:创建时间"`
	UpdatedAt time.Time             `gorm:"column:updated_at;type:timestamp;<-:false;comment:更新时间"`
	DeletedAt int64                 `gorm:"column:deleted_at;default:NULL;comment:软删除时间"` // Unix time
	IsDelete  soft_delete.DeletedAt `gorm:"column:is_delete;softDelete:flag,DeletedAtField:DeletedAt;comment:是否删除"`
}

type SenderOperation struct {
	BaseModel
	From        string       `gorm:"type:varchar(255);column:from;not null;default:''"`
	To          string       `gorm:"type:varchar(255);column:to"`
	TxHash      string       `gorm:"type:varchar(255);column:tx_hash"`
	ChainID     int64        `gorm:"type:int;column:chain_id"`
	BlockNumber uint64       `gorm:"type:int;column:block_number"`
	BlockTime   uint64       `gorm:"type:int;column:block_time"`
	Type        RouterMethod `gorm:"type:int;column:type"`
}

func (SenderOperation) TableName() string {
	return "sender_operations"
}

func CreateSenderOperations(data []SenderOperation) error {
	return db.MysqlDB.Create(&data).Error
}
