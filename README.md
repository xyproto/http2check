# http2check

Utility for checking if a given server supports HTTP/2.

Installation
------------

For Go 1.17 or later:

    go install github.com/xyproto/http2check@latest

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

* Version: 0.7.1
* License: BSD-3
* Alexander F. Rødseth
