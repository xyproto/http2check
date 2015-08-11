# http2check

Utility for checking if a given server supports HTTP/2.

Installation
------------

* `go get github.com/xyproto/http2check`
* Add `$GOPATH/bin` to the PATH (optional)

Example usage
-------------

`http2check twitter.com`

Output:

~~~
GET https://twitter.com
[protocol] HTTP/2.0
[status] 200 OK
~~~

General information
-------------------

* Version: 0.5
* License: MIT
* Alexander F Rødseth

