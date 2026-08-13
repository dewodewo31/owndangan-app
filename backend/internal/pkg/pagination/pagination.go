package pagination

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

type Params struct {
	Page    int
	PerPage int
}

func FromRequest(page, perPage int) Params {
	if page < 1 {
		page = defaultPage
	}
	if perPage < 1 {
		perPage = defaultPerPage
	}
	if perPage > maxPerPage {
		perPage = maxPerPage
	}
	return Params{Page: page, PerPage: perPage}
}

func Offset(params Params) int {
	return (params.Page - 1) * params.PerPage
}
