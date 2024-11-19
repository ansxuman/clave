package objects

import (
	"time"

	"github.com/google/uuid"
)

type TotpSecretObject struct {
	Id           string `json:"id"`
	Timestamp    int64  `json:"timestamp"`
	Issuer       string `json:"issuer"`
	Secret       string `json:"secret"`
	LastCopiedAt int64  `json:"lastCopiedAt"`
}

func CreateNewTotpSecretObject(issuer, secret string) TotpSecretObject {
	return TotpSecretObject{
		Id:        uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
		Issuer:    issuer,
		Secret:    secret,
	}
}

type TotpSecretObjectList []TotpSecretObject

func (list TotpSecretObjectList) Len() int {
	return len(list)
}

func (list TotpSecretObjectList) Less(i, j int) bool {
	return list[i].Timestamp < list[j].Timestamp
}

func (list TotpSecretObjectList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
