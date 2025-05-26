/*
zslog provides extensional logging features to the standard [log/slog] package.

zslog provides build-time logging by the [BuildLog] function.
It emits log records when the following build tag is provided.
Log records are output through the [BuildLogFunc]. It can be replaced if necessary.

	Build Tags:
		- zslogbuildlog  : Enables [BuildLog] function to work.
*/
package zslog
