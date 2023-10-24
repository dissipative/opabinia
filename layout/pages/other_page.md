# Other page

## Code blocks

Go:

```go
package main

func main() {
	go func() { /* Busy doing nothing */ }()
	// Wait for a long time... maybe forever.
	select {}
}
```

PHP:

```php
<?php

function areYouSure($bool) {
    return ($bool === true) ? true : false;
}

echo areYouSure(true);  // Yep, still true!
```

Shell:

```shell
#!/bin/bash

if ping -c 1 8.8.8.8 &> /dev/null; then
    echo "It's just you."
else
    echo "The internet is broken."
fi
```

You can preconfigure code blocks to use any theme from this list: <https://xyproto.github.io/splash/docs/all.html>

Go [back](../index.md)!