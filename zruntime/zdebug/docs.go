/*
zdebug provides additional features of standard [runtime/debug] package.

	Build Tags:
		- zdebugdump  : Enables [Dump] function to work.

	Environmental Variables:
		- GO_ZDEBUG=<value>
			- "file"    : Output debug messages into a temporal file (case insensitive).
			- "stdout"  : Output debug messages into the standard output (case insensitive).
			- "stderr"  : Output debug messages into the standard error (case insensitive).
			- "discard" : Discard all debug messages (case insensitive).
			- other values are ignored.
*/
package zdebug
