package workspaces

type AuthWorkspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
type AuthWorkspaceList struct {
	Workspaces []AuthWorkspace `json:"workspaces"`
}

type WorkspaceRequest struct {
	Name string `json:"name" binding:"required,min=2,max=20"`
}
