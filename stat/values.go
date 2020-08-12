package stat

//Values represents stats values
type Values []interface{}

func (v *Values) Append(value interface{}) {
	*v = append(*v, value)
}


//AppendAll appends all values
func (v *Values) AppendAll(values [] interface{}) {
	*v = append(*v, values...)
}


func (v *Values) Values() []interface{} {
	return *v
}

//New creates stat values
func New() *Values {
	return &Values{}
}
