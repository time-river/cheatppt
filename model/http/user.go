package msg

type CommonResponse struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Recaptcha string `json:"recaptcha"`
}

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Recaptcha string `json:"recaptcha"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
}

type ResetRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
