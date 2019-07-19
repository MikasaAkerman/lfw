package writer

// OptionKey define the type of key
type OptionKey string

// define the key of ops
const (
	OptionDuration   = "duration"
	OptionPrefix     = "prifix"
	OptionSuffix     = "suffix"
	OptionTimeLayout = "timeLayout"
	OptionMaxFiles   = "maxFiles"
	OptionMaxSize    = "maxSize"
	OptionLogdir     = "logDir"
)

// Option the writer init ops
type Option struct {
	Key   OptionKey
	Value interface{}
}

// NewOption create a new option with given key,value
func NewOption(key OptionKey, value interface{}) *Option {
	return &Option{key, value}
}

// GetKey get key of option
func (o *Option) GetKey() OptionKey {
	return o.Key
}

// GetValue get value of option
func (o *Option) GetValue() interface{} {
	return o.Value
}

// SetKey ...
func (o *Option) SetKey(key OptionKey) {
	o.Key = key
}

// SetValue ...
func (o *Option) SetValue(value interface{}) {
	o.Value = value
}
