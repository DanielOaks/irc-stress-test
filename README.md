# IRC Stress Test

Stress tests IRC servers and presents results on how long the tests took.

Unless you're writing an IRC server or something similar, this isn't that useful for you.

**Note:** Very pre-release. Very early. Does not yet work.


## Waiting

By default, we only wait for the final `QUIT` message to be processed (i.e. for an `ERROR` message to be returned to us). Passing the `--wait` flag makes us wait after every command we can wait after (i.e. channel joins, parts, etc).


## Recommendations

* Ensure that both the server and the stress test are allowed to open enough file descriptors to complete the test (check the output of `ulimit` or the contents of `/proc/${pid}/limits`).
* Test over localhost.
* Disable ident lookup.
* Disable connection limits.
* Disable rate limiting.
* Check `dmesg` for warnings about SYN flooding and adjust `net.ipv4.tcp_max_syn_backlog` as necessary
