package cr

type Criteria struct {
	Filter Filter
	SortBy SortBy
	Offset *int
	Size   *int
}

func New(args ...string) *Criteria {
	var (
		filter Filter
		sortBy SortBy
	)

	if len(args) > 0 {
		filter = ParseFilter(args[0], false)
	}

	if len(args) > 1 {
		sortBy = ParseSort(args[1])
	}

	return &Criteria{Filter: filter, SortBy: sortBy}
}

func (c *Criteria) SetFilter(filter Filter) *Criteria {
	c.Filter = filter
	return c
}

func (c *Criteria) SetSortBy(sort ...Sort) *Criteria {
	c.SortBy = sort
	return c
}

func (c *Criteria) SetOffset(index int) *Criteria {
	if index < 0 {
		index = 0
	}
	c.Offset = &index
	return c
}

func (c *Criteria) GetOffset() int {
	if c.Offset == nil || *c.Offset < 0 {
		return 0
	}
	return *c.Offset
}

func (c *Criteria) SetSize(size int) *Criteria {
	if size < 1 {
		size = 1
	}
	c.Size = &size
	return c
}

func (c *Criteria) GetSize(def int) int {
	if c.Size == nil || *c.Size <= 0 {
		return def
	}
	return *c.Size
}
