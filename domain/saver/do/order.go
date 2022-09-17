package do

const (
	Normal  = 1 // 正常
	Deleted = 2 // 删除
)

type Order struct {
	OrderNO    string
	ModifyTime int64
	Version    int64
	Status     int
}

func (o *Order) IsEqualForReqInfo(n *Order) bool {
	if n == nil {
		return false
	}
	if o.OrderNO != n.OrderNO {
		return false
	}
	if o.ModifyTime != n.ModifyTime {
		return false
	}
	if o.Version != n.Version {
		return false
	}

	return true
}
func (o *Order) IsNewerTo(n *Order) bool {
	if n == nil {
		return false
	}
	if o.OrderNO != n.OrderNO {
		return false
	}
	if o.Version <= n.Version {
		return false
	}
	return true
}
