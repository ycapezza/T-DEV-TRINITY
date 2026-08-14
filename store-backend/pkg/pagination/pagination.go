package pagination

type Pagination struct {
    Page     int   `json:"page"`
    PageSize int   `json:"page_size"`
    Total    int64 `json:"total"`
}

func NewPagination(page, pageSize int) *Pagination {
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }
    return &Pagination{
        Page:     page,
        PageSize: pageSize,
    }
}

func (p *Pagination) GetOffset() int {
    return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
    return p.PageSize
}