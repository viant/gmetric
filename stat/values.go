package stat

//Values represents stats values
type Values []interface{}

//Append appends values to a slice
func (v *Values) Append(value interface{}) {
	*v = append(*v, value)
}

//AppendAll appends all values
func (v *Values) AppendAll(values []interface{}) {
	*v = append(*v, values...)
}

//Values returns all values
func (v *Values) Values() []interface{} {
	return *v
}

//New creates stat values
func New() *Values {
	return &Values{}
}
