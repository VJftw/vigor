package html

import "syscall/js"

// BindValue two-way binds a text input, textarea or select to a signal: the
// element shows get(), and typing (or choosing) calls set. Use it as an
// element's option: html.El("input", html.BindValue(name, setName)).
func BindValue(get func() string, set func(string)) *elementGroup {
	return Group(
		Property("value", get),
		On("input", func(this js.Value, _ []js.Value) { set(this.Get("value").String()) }),
	)
}

// BindChecked two-way binds a checkbox (or radio) to a signal.
func BindChecked(get func() bool, set func(bool)) *elementGroup {
	return Group(
		Property("checked", get),
		On("change", func(this js.Value, _ []js.Value) { set(this.Get("checked").Bool()) }),
	)
}
