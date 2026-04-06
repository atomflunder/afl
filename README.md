# afl

Very simple scripting language that will bring literally nothing new to the table.

## Examples

```
// Comment

/* 
    Multi-line comment
*/

x = 5           // Type inference
z: float = 7.0  // Explicit type annotation
z = 8.5         // Will not do anything, not mutable

a? = 10 // Declared as mutable
a = 15  // Allowed

fn a_function() {
    print("Hello, World!")

    return 42
}

a_function(); // Semicolon is optional

if (x > 3) {
    print("x is greater than 3 but actually {x}")
} elseif (x == 3) {
    print("x is equal to 3")
} else { print("x is less than 3") } // Whitespace does not matter

for (i in 0->5) {
    print(i)
}
```

The script will produce the following output:

```sh
Hello, World!
x is greater than 3 but actually 5
0
1
2
3
4
```

For more examples, check out the [`./examples`](./examples) directory.

## Usage

Build the compiler:

```sh
go build -o afl ./src/main.go   # manually
./scripts/build.sh              # or just via the build script
```

Run the compiler on your script:

```sh
./afl your_script.afl       # Linux / MacOS
afl.exe your_script.afl     # Windows
```