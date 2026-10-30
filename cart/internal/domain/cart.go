package domain

import "sync/atomic"

type CartItem struct {
	SKU   int64
	Name  string
	Count uint32
	Price uint32
}

type Cart struct {
	Items      []CartItem
	TotalPrice uint32
}

func (c *Cart) AddItem(item CartItem) {
	for i := range c.Items {
		if c.Items[i].SKU == item.SKU {
			atomic.AddUint32(&c.Items[i].Count, item.Count)
			return
		}
	}
	c.Items = append(c.Items, item)
}
