# stacktrace

I know, I know. Any experienced gopher has undoubtedly rolled their eyes by now, but bear with me.

Logs or errors that are just printed to the screen are difficult to debug. I _suppose_ you could clone the repo and
search through it with Grep, Ripgrep, your IDE or even use Github to search for it like some kind of wild animal.

That being said, there's nothing wrong with logging or printing errors to the screen. In fact, you should do that
instead of attempting to abuse this package by throwing stacktraces in everyone's faces.

This package is a little bit different than some of the other stacktrace packages for Go. When you use this package, you
_should_ be using `stacktrace.Propogate(err, cause, args)` every time you return an error. It does _not_ mean that the
stacktrace will be "thrown" so to speak. If `stacktrace.Throw(err)` is never called, then this package behaves similarly
to `errors.Wrap(err, cause)` from `github.com/pkg/errors` in that the stack is still recorded, but if you just want to
print the errors you picked up along the way in that `fmt.Sprintf("msg: cause")` format, you can just by
printing/logging/returning `stacktrace.Error(err)`.
