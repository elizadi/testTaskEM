package domain

type GetUsersReq struct {
	ID         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        uint8
	Gender     string
	Country    string
	Pag        Pagination
}

type Pagination struct {
	Page    int
	PerPage int
}

func (p *Pagination) Offset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 1
	}
	return (p.Page - 1) * p.PerPage
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) Limit() int {
	if p.PerPage <= 0 {
		p.PerPage = 1
	}
	return p.PerPage
}
