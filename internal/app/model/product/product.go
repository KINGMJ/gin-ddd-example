package product

import (
	"gin-ddd-example/internal/app/model"
)

type Product struct {
	model.BaseModel
	Unit        string `gorm:"column:unit;not null;default:''" json:"unit"`                // 单位
	EAN13       string `gorm:"column:ean13;not null" json:"ean13"`                         // 条形码，不限定 ean13
	Name        string `gorm:"column:name;not null" json:"name"`                           // 商品名称
	MerchantID  int64  `gorm:"column:merchant_id;not null" json:"merchant_id"`             // 商户id
	StoreID     int64  `gorm:"column:store_id;not null;default:0" json:"store_id"`         // 门店id
	Type        int64  `gorm:"column:type;not null;default:0" json:"type"`                 // 商品类型id
	Brand       string `gorm:"column:brand" json:"brand"`                                  // 商品品牌
	BuyPrice    int64  `gorm:"column:buy_price;not null;default:0" json:"buy_price"`       // 采购价，单位为分
	Price       int64  `gorm:"column:price;not null;default:0" json:"price"`               // 零售价，单位为分
	SettlePrice int64  `gorm:"column:settle_price;not null;default:0" json:"settle_price"` // 结算价，单位为分
	Discount    int8   `gorm:"column:discount;not null;default:100" json:"discount"`       // 折扣率：0-100
	Specs       string `gorm:"column:specs;not null" json:"specs"`                         // 规格
	ExpireTime  string `gorm:"column:expire_time;not null;default:''" json:"expire_time"`  // 过期时间
	Measurement int8   `gorm:"column:measurement;not null;default:1" json:"measurement"`   // 计量单位：1 计件；2 称重
	Source      int8   `gorm:"column:source;not null;default:1" json:"source"`             // 商品来源：1 自建商品；2 采销同步
	Status      int8   `gorm:"column:status;not null;default:3" json:"status"`             // 商品状态：1 待提交；2 待审核；3 已审核
}

func (t *Product) TableName() string {
	return "product"
}
