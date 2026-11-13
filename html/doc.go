// Package html builds a page from Go: elements, text and structural nodes
// (Each, When, Await...) that update from signals. It needs a browser, so
// its files are for js/wasm; natively the package is empty, which lets
// programs and build systems that cover every platform (Bazel's //...) build
// it without special cases.
package html
