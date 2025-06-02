package model

type RoomMember struct {
	ID       string `gorm:"primaryKey;type:varchar(255)"`
	UserID   string `gorm:"type:varchar(255);not null"`
	RoomID   string `gorm:"type:varchar(255);not null"`
	JoinedAt int64  `gorm:"not null"`
	User     User   `gorm:"foreignKey:UserID"`
	Room     Room   `gorm:"foreignKey:RoomID"`
}
