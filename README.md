# Go-kType

Go-kType is a little experimental, for-fun, tool that generates Kotlin classes based from Golang classes, allowing HTTP 
APIs to expose types that they want to expose, plus a recommended import map and generator, to generate types for their  
clients.

## Methodology

Go-kType works with the heavy use of Reflection, we have this little internal thing called a `runtime` (not really a runtime), 
which contains all the translated classes. In this `runtime`, we have private methods called `classFrom` that heavily uses 
reflection to build a little structure of the class, one that we can use to translate down to Kotlin. 

In addition, we have a type translator that basically turns Golang types to a more dumbed down type, for example, 
`uint` to `uint32` and even `int` to `int32`is dumbed down to `Int` while `uint64` and `int64` are dumbed down to `Long`. 
Furthermore, the type translator will dumb down special properties such as `time.Time` to `instant`, additionally configurable 
by using `runtime.BiasedTypeTranslator.Map(name, value)`, for example:
```go
runtime.BiasedTypeTranslator
    .Map("time.Time", "instant")
    .Map("sync.Mutex", "mutex")
```

This dumbed down types should still be valid types in the desired language, in this case, Kotlin, but Kotlin doesn't have 
an 'Instant' type, which is where import maps comes to play. Import maps are simply JSON files that maps these dumbed down 
types to specific imports, for example, the `instant` we had is mapped to `kotlinx.datetime.Instant`:
```json
{
  "instant": "kotlinx.datetime.Instant"
}
```

As a result of the import map, our little generator will see if any of the fields has a type that is mapped with a property 
in the import map and add that as an import, which later down, gets generated into something such as:
```kotlin
import kotlinx.datetime.Instant
```

Now, the last part of the sequence is the `Typist` which is simply templates for the generation. As a default, we have under 
[`generators/kotlin`](generators/kotlin) a Kotlin template called [`data_class.kt.typist`](generators/kotlin/data_class.kt.typist) which 
is simply a dumb template of how a Kotlin file with data class should be generated.

Combined, all of these eventually generate into "valid" Kotlin files, but it's not perfect and shouldn't be as this is simply a little 
experimental and for-fun project.

## I want to try it out!

Well, you can certainly try it out by downloading the package:
```shell
go get github.com/ShindouMihou/go-ktype
```

After downloading the package, you can copy the [`generators/kotlin`](generators/kotlin) folder under somewhere where inside your code, 
you can retrieve, then simply create your structs and then import the import map and typist file to create your little runtime. You can 
read the following example file:
- [`examples/sample.go`](examples/sample.go)