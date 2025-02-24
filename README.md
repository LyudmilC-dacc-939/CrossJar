# CrossJar
A Cross-Platform JAR-to-Executable Converter

Transform Java JAR files into native executables that work seamlessly across Windows (.exe), Linux, and macOS. This tool embeds your JAR into a lightweight Go wrapper, creating a platform-specific binary that runs the JAR using the system’s Java runtime.

## Features:
✅ Cross-platform output (Windows, Linux, macOS)
✅ Architecture support (amd64, arm64, etc.)
✅ No external dependencies (except Java on target systems)
✅ Automatic cleanup of temporary files
✅ Simple CLI interface with build flags

## Usage:

crossjar -input app.jar -output app.exe -os windows -arch amd64

## Requirements:

Go 1.16+ (to build the converter)

Java (on systems running the generated executable)

## How It Works:

Embeds your JAR file into a Go binary.

When executed, extracts the JAR to a temp directory.

Runs java -jar on the extracted file.

Deletes temp files on exit.

## Ideal For:

Distributing JAR-based tools to non-technical users

Creating OS-specific launchers for Java applications

Simplifying deployment workflows

## Limitations:

Requires Java on target machines (does not bundle a JRE).
