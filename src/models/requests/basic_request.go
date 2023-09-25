package requests

type PaginationRequest struct {
	Page      int    `json:"page" form:"page" validate:"required"`
	Size      int    `json:"size" form:"size" validate:"required,gte=0,lte=25"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}
