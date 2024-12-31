package dto


type Sort struct {
	Sort string `json:"sort"`
	Column int `json:"column"`
}

type Filter struct {

	// text or number
	FilterType string `json:"filter_type"`
	
	//lessThan - greaterThan - inRange - equals - notEquals - contains - notContains ...
	Type string `json:"type"`
	From string `json:"from"`
	To string `json:"to"`
}

type DynamicFilter struct {
	Sort []Sort `json:"sort"`
	Filters map[string]Filter `json:"filters"`
}