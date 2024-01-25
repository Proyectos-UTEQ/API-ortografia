package types

// Paginated para utilizarlos en los endpoint de paginados.
type Paginated struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Query string `query:"query"`
	Sort  string `query:"sort"`
	Order string `query:"order"`
}

func (p *Paginated) Validate() error {

	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Sort == "" {
		p.Sort = "id"
	}

	if p.Order == "" {
		p.Order = "asc"
	}

	return nil
}

// Datos adecimados para la paginación.
type PagintaedDetails struct {
	Page         int   `json:"page"`
	TotalPage    int64 `json:"total_page"`
	TotalItems   int64 `json:"total_items"`
	ItemsPerPage int   `json:"items_per_page"`
}
