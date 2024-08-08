package model

import "time"

type StorePo struct {
	Store
}

type Store struct {
	ID                      int64     `gorm:"primaryKey,autoIncrement" json:"id"`
	Name                    string    `json:"name"`
	MerchantID              int64     `json:"merchantId"`
	HuifuId                 string    `json:"huifuId"` // 汇付商户号
	Location                string    `json:"location"`
	Count                   int64     `json:"count"` // 门店人数
	ContactName             string    `json:"contactName"`
	ContactPhone            string    `json:"contactPhone"`
	AuditId                 int64     `json:"auditId"`
	Type                    int       `json:"type"`                    // 1.门店 2.仓库
	Sort                    int64     `json:"sort"`                    // 编号
	Status                  int64     `json:"status"`                  // 状态(1:开启;0:关闭)
	JoinStatus              int64     `json:"joinStatus"`              // 加盟状态(0:未加盟,1:加盟,2:自营)
	SeparatePointSettlement bool      `json:"separatePointSettlement"` // 是否为独立积分结算
	PointSettlementRatio    int64     `json:"point_settlement_ratio"`  // 积分结算抽佣比例
	CreatedAt               time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (table *StorePo) TableName() string {
	return "store"
}
