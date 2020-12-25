# csvcmd
A command line tool to filter, read, and modify CSV files

### Installation:
```bash
$ go get github.com/appins/csvcmd
$ # Make sure your $PATH contains $GOPATH/bin, you can
$ # do that with PATH=$PATH:$GOPATH/bin
```

### Usage:
Help:
```
csvcmd [OPTION]... [FILE]...
Print all lines that meet the paramters specified by OPTIONs
in each FILE (Or stdin, if no files specified).
Example: csvcmd --end=20 --filter="First Name=Alex" people.csv

Selection:
	--start=int	Specify the first line (inclusive) that should be read

	--end=int	Specify the last line (inclusive) that should be read
			Note: To control how many lines are read pipe command
			output to head or tail. try `| head -nXX` where XX=#lines

	--filter="..."	Specify a set of filters that a row needs to meet to
			be printed. See the filters section at
			github.com/appins/csvcmd

	--or		Set filters to be OR'd rather than AND'd together.
			Requires only 1 filter to be met to print a line.
	
	--shown="...	Specify columns which should be shown, seperated with a
			semicolon. Either use the header row's text or specify
			the number of the column with `_#`, like `_3`

Output:
	-h		Human readable output (default: regular CSV)
```

Example Usage:
```bash
$ csvcmd -h --end=5 somecsv.csv
> Sex Weight ... Weight ... BMI (Sep) BMI (Apr)
> M   159        130        22.02     18.14
> M   214        190        19.7      17.44
> M   163        152        24.09     22.43
> M   205        194        26.97     25.57
> F   150        141        21.51     20.1
$ csvcmd -h --end=5 --filter "Sex=F" somecsv.csv
> Sex Weight ... Weight ... BMI (Sep) BMI (Apr)
> F   150        141        21.51     20.1
```

### Features:

##### Filters:
Filters follow the following format: `filter1;filter2;filter3`.
As of right now, there is only 1 kind of filter (TODO: Add more).
Also, for a line to print, every filter must return true.

Note that for all filters, you can specify a column number rather
than exact name. To do this, just set the column name to be `_#`.
Ex: the third column can be refered to at `_3`. Counting starts at 1.

Filter types:
1. Equality: `column_name=value`: Passes when the cell under `column_name` is equal to
value

