/*
zerrors provides additional features to standard [errors] package.

	Build Tags:
		- zerrorstrace  : Enables error tracing to work.

	Environmental Variables:
		- GO_ZERRORS=<value>
			- "file"    : Output trace messages into a temporary file (case insensitive).
			- "stdout"  : Output trace messages into the standard output (case insensitive).
			- "stderr"  : Output trace messages into the standard error (case insensitive).
			- "discard" : Discard all trace messages (case insensitive).
			- other values are ignored.
*/
package zerrors
