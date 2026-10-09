package domain

type CartItem struct {
	SKU   int64
	Name  string
	Count uint32
	Price uint32
}

type Cart struct {
	Items   []CartItem
	TotalPrice uint32
}

func (c *Cart) AddItem(item CartItem) {
	c.Items = append(c.Items, item)
}