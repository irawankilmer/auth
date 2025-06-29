package model

type ProfileModel struct {
	ID       string
	UserID   string
	FullName string
	Address  *string
	Gender   *string
	User     UserModel
}
