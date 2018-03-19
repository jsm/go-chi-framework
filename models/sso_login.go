package models

type SSOLogin struct {
	Base
	User    User
	UserID  uint
	SSOID   string `gorm:"column:sso_id"`
	SSOType string
}

var SSOLoginTypes = struct {
	Google   string
	LinkedIn string
}{
	Google:   "google",
	LinkedIn: "linkedin",
}
