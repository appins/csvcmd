# `csvcmd` Example Usage
As `csvcmd` is a tool that many people could benefit from, I have decided that it would be worthwile task to document an example use case of this tool.

### Prerequisites:
Before getting started you should have a recent version of Go installed. This library was tested with Go version 1.11.4, which was released a little over two years ago. You will also want to have some CSV file(s) to work with.

### Installation:
As README.md states, you'll want to run `go get github.com/appins/csvcmd` to install the tool. Additionally, you'll want to add it to your `$PATH` environment variable. If you're on Linux or MacOS, you can do this by adding some line in your `.bash_profile` with a command like `echo "PATH=$PATH:$GOPATH/bin" >> .bash_profile`. 

For Windows, consider just `cd`'ing to the folder (typically `%homepath%/go/src/github.com/appins/csvcmd`) and running `csvcmd.exe`from there. Note that while this tool works on Windows, it's designed primarily for Linux or MacOS users.

### Viewing files:
Calling csvcmd without any options prints a csv file without any changes. This is functionally equivalent to `cat`'ing a file. Example: `csvcmd data.csv`.

To view a formatted (fixed column width) view, you can use the `-h` option. Example: `csvcmd -h data.csv`.
