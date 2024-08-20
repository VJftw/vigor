package html

/*
onCleanup design


Elements which may remove/switch elements must evaluate onCleanup methods for a
Node.

Options:

a). Create an OnCleanup wrapper which wraps the Node with a cleanup method that
	the remove/switch elements call.

b). Pass an owner/context/scope variable which vigor.OnCleanup acts as a
	subscriber on. e.g. `vigor.OnCleanup(scope, fn)`, where on a "cleanup"
	event, for the `scope`, call `fn`.
	- Could I elaborate/expand this? i.e. add a `scope` to `context.Context` to
	  then extract. This `scope` would set by the elements which call `Node.DOMObject()`?

*/
