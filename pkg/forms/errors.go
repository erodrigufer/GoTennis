package forms

// New errors type, will hold the validation error messages for forms.
// Each form can instantiate an errors type object for the specific errors of
// its particular form
// The name of the form field will be used as the key in this map.
// Store string slices inside the map for a given key
type errors map[string][]string

// Add() method adds error messages for a given field to the map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get() method retrieves the first error message for a given field from the map
func (e errors) Get(field string) string {
	es := e[field]
	// No error messages for a given field
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
