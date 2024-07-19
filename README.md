# go-xml-patch

A golang implementation of [RFC 5261](https://www.rfc-editor.org/rfc/rfc5261.html): "An Extensible Markup Language (XML)
Patch Operations Framework Utilizing XML Path Language (XPath) Selectors".

# Status

This library is still in early stage. API may break. Missing functionality (see [progress](#progress)). Probably buggy
🐛.

# Usage

### Install

```shell
go get github.com/jfrog/go-xml-patch
```

### Code

```go
package main

import (
	"fmt"
	"github.com/jfrog/go-xml-patch"
	"os"
)

func main() {
	target, _ := os.ReadFile("target.xml")
	diff, _ := os.ReadFile("diff.xml")
	patch, err := xmlpatch.Patch(target, diff)
	if err != nil {
		panic(err)
	}
	fmt.Println(patch)
}
```

# Example [from RFC 5261](https://www.rfc-editor.org/rfc/rfc5261#appendix-A.6)

**An example target XML document:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<doc>
    <foo a="1">This is a sample document</foo>
</doc>
```

**An XML diff document:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<diff>
    <replace sel="doc/foo[@a='1']">
        <bar a="2"/>
    </replace>
</diff>
```

**A result XML document:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<doc>
    <bar a="2"/>
</doc>
```

# Progress

### [Specification](https://www.rfc-editor.org/rfc/rfc5261) Items

- [ ] 4.3. `<add>` Element
    - [x] 4.3.1. Adding an Element
    - [ ] 4.3.2. Adding an Attribute
    - [ ] 4.3.3. Adding a Prefixed Namespace Declaration
    - [ ] 4.3.4. Adding Node(s) with the 'pos' Attribute
    - [ ] 4.3.5. Adding Multiple Nodes
- [ ] 4.4. `<replace>` Element (with optional auto-create, see [#2](https://github.com/jfrog/go-xml-patch/issues/2))
    - [x] 4.4.1. Replacing an Element
    - [x] 4.4.2. Replacing an Attribute Value
    - [ ] 4.4.3. Replacing a Namespace Declaration URI
    - [ ] 4.4.4. Replacing a Comment Node
    - [ ] 4.4.5. Replacing a Processing Instruction Node
    - [ ] 4.4.6. Replacing a Text Node
- [ ] 4.5. `<remove>` Element
    - [ ] 4.5.1. Removing an Element
    - [ ] 4.5.2. Removing an Attribute
    - [ ] 4.5.3. Removing a Prefixed Namespace Declaration
    - [ ] 4.5.4. Removing a Comment Node
    - [ ] 4.5.5. Removing a Processing Instruction Node
    - [ ] 4.5.6. Removing a Text Node

### Other

- [ ] CI
- [ ] Release flow
- [ ] better docs
- [ ] coverage

# Release

(manual for now)

Follow [Go module publishing instructions](https://go.dev/doc/modules/publishing) and then create a new [Github release](https://go.dev/doc/modules/publishing).