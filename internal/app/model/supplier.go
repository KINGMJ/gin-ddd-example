package model

import (
	"gin-ddd-example/internal/app/model/ctype"
)

type SupplierPo struct {
	Supplier
	Stores []StorePo `json:"stores" gorm:"many2many:supplier_store;foreignKey:ID;joinForeignKey:SupplierID;references:ID;joinReferences:StoreID"`
}

// db 模型
type Supplier struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`          // 供应商名称
	SType           int64          `json:"sType"`                         // 供应商类型
	Region          string         `json:"region"`                        // 供应商区域
	ComMobile       string         `json:"comMobile"`                     // 公司手机号
	Fax             string         `json:"fax"`                           // 公司传真
	BName           string         `json:"bName"`                         // 订货人名称
	BMobile         string         `json:"bMobile"`                       // 订货人手机号
	TaxesCard       string         `json:"taxesCard"`                     // 纳税人等级号
	Purchaser       string         `json:"purchaser"`                     // 采购员
	AdvanceType     string         `json:"advanceType"`                   // 进货方式
	OrderDate       ctype.NullTime `json:"orderDate"`                     // 订货日
	Cycle           int64          `json:"cycle"`                         // 到货周期
	PayWhere        int64          `json:"payWhere" gorm:"default:1"`     // 付款条件 1:货到付款 2:订单付款
	BankName        string         `json:"bankName"`                      // 银行名称
	OpenBank        string         `json:"openBank"`                      // 开户支行
	BankAccount     string         `json:"bankAccount"`                   // 银行账户/对公账户
	BankAccountName string         `json:"bankAccountName"`               // 银行卡户名/公司名称
	Address         string         `json:"address"`                       // 供应商地址
	MerchantId      int64          `json:"merchantId"`                    // 商户id
	Scale           int            `json:"scale"`                         // 个人还是企业 1:个人 2:企业
	Mode            int            `json:"mode"`                          // 经营方式 1:购消 2:联营
	Created         ctype.NullTime `json:"created" gorm:"autoCreateTime"` // 添加时间
	Updated         ctype.NullTime `json:"updated" gorm:"autoUpdateTime"` // 修改时间
}

func (table *SupplierPo) TableName() string {
	return "supplier"
}

// ----------- (●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●)(●'◡'●) ------------
const (
	_                    int = iota
	SupplierScalePerson      // 个人
	SupplierScaleCompany     // 企业
)

const (
	_                        int = iota
	SupplierModePurchaseSale     // 购销
	SupplierModeJoint            // 联营
)

// 其他自定义的类型
type (
	Suppliers []*Supplier
)
