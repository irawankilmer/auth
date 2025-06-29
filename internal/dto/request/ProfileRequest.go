package request

type ProfileRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Address  string `json:"address"`
	Gender   string `json:"gender" binding:"required,oneof=laki-laki perempuan"`
}

func (p *ProfileRequest) Sanitize() map[string]any {
	return map[string]any{
		"full_name": p.FullName,
		"address":   p.Address,
		"gender":    p.Gender,
	}
}
