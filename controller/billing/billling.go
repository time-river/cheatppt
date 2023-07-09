package billing

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

type BillingItem struct {
	UUID          string  `json:"uuid"`
	Credit        float32 `json:"credit"` // RMB, unit: yuan
	Status        string  `json:"status"`
	PaymentMethod string  `json:"paymentMethod"`
	CreateAt      int64   `json:"createAt"`
	PaiedAt       int64   `json:"paiedAt"`
	Comment       string  `json:"comment"`
}

type billingStatus int

const (
	billingProcess  = 0
	billingComplete = 1
	billingCancel   = 2
)

func (b billingStatus) String() string {
	switch b {
	case billingProcess:
		return "进行中"
	case billingComplete:
		return "已支付"
	case billingCancel:
		return "已取消"
	default:
		return "未知状态"
	}
}

func GetBillings(rng *Range) ([]BillingItem, error) {
	var records []sql.UserBilling
	var data = make([]BillingItem, rng.Length)

	db := sql.NewSQLClient()

	if err := db.Offset(rng.From).Limit(rng.Length).Find(&records).Error; err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("内部错误")
	}

	for i, record := range records {
		data[i] = BillingItem{
			UUID:          record.UUID.String(),
			Credit:        sql.Coins2RMB(record.Coins),
			Status:        billingStatus(record.Status).String(),
			PaymentMethod: record.PaymentMethod,
			CreateAt:      record.CreatedAt.Unix(),
			PaiedAt:       record.CreatedAt.Unix(),
			Comment:       record.Comment,
		}
	}

	return data[0:], nil
}
