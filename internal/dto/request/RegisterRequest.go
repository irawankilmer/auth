package request

type RegisterRequest struct {
	Username string   `json:"username" binding:"required,excludesall= "`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Roles    []string `json:"roles" binding:"required"`
}

func (r *RegisterRequest) Sanitize() map[string]any {
	return map[string]any{
		"username": r.Username,
		"email":    r.Email,
		"roles":    r.Roles,
	}
}
