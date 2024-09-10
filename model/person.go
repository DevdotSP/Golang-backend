package model

import "time"

type User struct {
	ID            uint          `gorm:"primaryKey;column:id" json:"id"`
	Name          string        `gorm:"column:name;not null" validate:"required,min=8,max=12" json:"name"`
	Age           int           `gorm:"column:age;not null" validate:"required,gte=18,lte=65" json:"age"`
	Email         string        `gorm:"column:email;unique;not null" validate:"required,email" json:"email"`
	Password      string        `gorm:"column:password;not null" validate:"required,min=8,max=12" json:"password"`
	AccountDetail AccountDetail `gorm:"foreignKey:UserID" json:"account_details"`
	History       History       `gorm:"foreignKey:UserId" json:"histories"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=12"`
}

type AccountDetail struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	UserID  uint    `gorm:"index" json:"user_id"` // Adding an index to the foreign key
	Balance float64 `json:"balance"`
}

type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserId    uint      `gorm:"index" json:"user_id"` // Adding an index to the foreign key
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}
