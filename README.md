# progstack ssg

The progstack SSG emphasises extreme simplicity.

It works by crawling the provided directory recursively, treating each Markdown
file it finds as a page.
The directory structure is mirrored across exactly.

## Setting up for private importing

First you need to ensure that your Git-config is set to use SSH:

```
[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
```

Then use this command:

```bash
GOPRIVATE=github.com/xr0-org/progstack-ssg go get github.com/xr0-org/progstack-ssg@feat/basic-generation
```

After this you should be able to import [ssg](pkg/ssg/ssg.go) like so:

```Go
package main

import (
	"log"

	"github.com/xr0-org/progstack-ssg/pkg/ssg"
)

func main() {
	if err := ssg.Generate("hello", "output", ""); err != nil {
		log.Fatal(err)
	}
}
```
