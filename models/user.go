package models

type User struct {
	ID  int64 `gorm:"primary key" json:"id"`
	NamaLengkap string `gorm:"varchar(300)" json:"nama_lengkap" validate:"required"`
	UserName string `gorm:"varchar(300)" json:"username" validate:"required"`
	Password string `gorm:"varchar(300)" json:"Password" validate:"required,gte=8,alphanumunicode"`
}