/*
zlog provides logging features.

zlog provides build-time logging by [BuildLog] function.
It emits log recordes when the following build tag is provided.
Log records are output through the [BuildLogFunc]. It can be replaced if necessary.

	Build Tags:
		- zlogbuildlog  : Enables [BuildLog] function to work.
*/
package zlog
