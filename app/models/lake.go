package models

type Lake struct {
	Name     string `json:"name"`
	Season   string `json:"season"`
	IsFrozen bool   `json:"is_frozen"`
}

func (lake Lake) Frozen() string {
	if lake.IsFrozen {
		return "Yes"
	} else {
		return "No"
	}
}
