package response

type ProfileResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"-"`
	FullName string  `json:"full_name"`
	Address  *string `json:"address"`
	Gender   *string `json:"gender"`
}
