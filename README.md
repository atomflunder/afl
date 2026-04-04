# Atomflunder's Programming Language

Work in progress of a very simple scripting language that will bring literally nothing new to the table.  
Even the name is a work in progress.

## Examples

```
// Comment

/* 
    Multi-line comment
*/

x = 5           // Type inference
y: int = 6      // Explicit type annotation
z: float = 7.0  // Explicit type annotation

z = 8.5 // Forbidden

a? = 10 // Mutable
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
} else {
    print("x is less than 3")
}

for (i in 0->5) {
    print(i)
}
```

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