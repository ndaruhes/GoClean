package requests

type PaginationRequest struct {
	Page    int    `json:"page" form:"page" validate:"required"`
	Size    int    `json:"size" form:"size" validate:"required,gte=0,lte=25"`
	OrderBy string `json:"orderBy" form:"orderBy"`
	Order   string `json:"order" form:"order"`
}
