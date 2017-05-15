package cli

type stringValue struct {
	val *string
}

func (v *stringValue) Set(value string) error {
	*v.val = value
	return nil
}

func (v *stringValue) String() string {
	return *v.val
}
