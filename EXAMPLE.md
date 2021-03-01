# `csvcmd` Example Usage
As `csvcmd` is a tool that many people could benefit from, I have decided that it would be worthwhile task to document an example use case of this tool.

### Prerequisites:
Before getting started you should have a recent version of Go installed. This library was tested with Go version 1.11.4, which was released a little over two years ago. You will also want to have some CSV file(s) to work with.

### Installation:
As README.md states, you'll want to run `go get github.com/appins/csvcmd` to install the tool. Additionally, you'll want to add it to your `$PATH` environment variable. If you're on Linux or MacOS, you can do this by adding some line in your `.bash_profile` with a command like `echo "PATH=$PATH:$GOPATH/bin" >> .bash_profile`. 

For Windows, consider just `cd`'ing to the folder (typically `%homepath%/go/src/github.com/appins/csvcmd`) and running `csvcmd.exe`from there. Note that while this tool works on Windows, it's designed primarily for Linux or MacOS users.

## Methods we'll be using:
#### Viewing files:
Calling csvcmd without any options prints a csv file without any changes. This is functionally equivalent to `cat`'ing a file. Example: `csvcmd data.csv`.

To view a formatted (fixed column width) view, you can use the `-h` option. Example: `csvcmd -h data.csv`.

(insert picture here)

#### Truncating files:
There are a few ways to truncate files in order to see only relevant data. They are as follows:
- filters: filters allow you to see only rows that meet certain criteria. Say you only would like to see contacts who have the first name Alex, or only view wrestlers that have a weight of 150 to 160lbs. To use a filter, use the `-filter` flag. For example, `-filter "name=Alex"`.
- start and end: if you decide that you only want to see rows in a portion of the file, you can use the `-start` and `-end` flags. For example, if we only wanted to see the rows that are between 50 and 100 inclusive, we would use the flags `-start 50 -end 100`. (Note that you cannot currently use these and split at the same time).
- shown columns: rather than just truncating rows, it is sometimes helpful to truncate columns. To select only relevant columns, use the `-shown` flag. For example, to only view the name and phone number, use `-shown "name;phone number"`.
- splits: oftentimes you'll want to split one CSV file into several. To do this, use the `-split` flag. For example, `-split "1/2"` and `-split "2/2"` will produce the first and second half of the file. The files will both have CSV headers, but will have no overlapping data.

## A use case: Split office workers into two groups.
We begin with the following data. How can we split these people into two groups with zero overlaps?

