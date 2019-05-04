# http2check

Utility for checking if a given server supports HTTP/2.

Installation
------------

Optionally add `$GOPATH/bin` to the PATH, then:

    go get -u github.com/xyproto/http2check

Example usage
-------------

    http2check twitter.com

Output:

~~~
GET https://twitter.com
[protocol] HTTP/2.0
[status] 200 OK
~~~

Limitations
-----------

* IPv6 addresses are not supported.

General information
-------------------

* Version: 0.6
* License: MIT
* Alexander F Rødseth
