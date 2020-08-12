package base

import (
	"github.com/viant/gmetric/stat"
	"github.com/viant/toolbox"
)

//Provider represents a base provider, note that this provider uses 0 index for an error type.
type Provider struct {
	keys []string
	aMap   map[interface{}]int
}

func (p *Provider) Map(key interface{}) int {
	if key == nil {
		return -1
	}
	if _, ok := key.(error); ok {
		return 0
	}
	index, ok := p.aMap[key]
	if ! ok {
		return -1
	}
	return index

}

//Keys  returns mapped keys
func (p *Provider) Keys() []string {
	return p.keys
}
//NewBaseProvider
func NewProvider(values ...interface{}) *Provider {
	values = append([]interface{}{stat.ErrorKey}, values...)
	var aMap = make(map[interface{}]int)
	var keys = make([]string, len(values))
	for i, item := range values {
		aMap[item] = i
		keys[i] = toolbox.AsString(item)
	}
	return &Provider{keys:keys, aMap:aMap}
}
