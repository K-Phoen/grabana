# Unmarshalling a YAML file

Grabana provides a way to unmarshal YAML. This can be done by giving a file or anything that
satisfies the `io.Reader` interface to the `decoder.UnmarshalYAML` function.

The result is a `dashboard.Builder` that can then be used by Grabana's client to upsert the dashboard.

```go
package main 

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana/decoder"
)

func main() {
	filePath := "some/awesome/dashboard.yaml"

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(1)
	}

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBuffer(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", err)
		os.Exit(1)
	}

	// do something with `dashboard`
}
```


## That was it!

[Return to the index to explore the other possibilities of the module](index.md)
