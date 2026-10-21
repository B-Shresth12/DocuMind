package auth

type RegisterRequest struct {
	DisplayName   string `json:"display_name" binding:"required,min=2,max=20"`
	Email         string `json:"email" binding:"required,email,max=320"`
	Password      string `json:"password" binding:"required,min=12,max=128"`
	WorkspaceName string `json:"workspace_name" binding:"required,min=2,max=120"`
}

type RegisterResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}
