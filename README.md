# IRC Stress Test

Stress tests IRC servers and presents results on how long the tests took.

Unless you're writing an IRC server or something similar, this isn't that useful for you.

**Note:** Very pre-release. Very early. Does not yet work.


## Waiting

By default, we only wait for the final `QUIT` message to be processed (i.e. for an `ERROR` message to be returned to us). Passing the `--wait` flag makes us wait after every command we can wait after (i.e. channel joins, parts, etc).


## Recommendations

* Test over localhost.
* Disable ident lookup.
* Disable connection limits.
* Disable rate limiting.
