#  Example Usage of `csvcmd`
As `csvcmd` is a tool that many people could benefit from, I have decided that it would be a worthwhile task to document a few examples use cases of this tool.

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

## A use case: Analyzing flight statistics.
We begin with the following data:

![Screenshot from 2021-03-01 19-12-11](https://user-images.githubusercontent.com/13598541/109591257-0d480580-7ac2-11eb-912d-2a320fb437c3.png)


How could you split this into two groups with zero overlaps? Before we attempt that feat, it may be useful to visualize the data. Do this by typing `csvcmd -h file.csv`. You'll get a much more readable output that looks like this:

![Screenshot from 2021-03-01 19-14-16](https://user-images.githubusercontent.com/13598541/109591585-a24afe80-7ac2-11eb-8e50-bdc014c37f29.png)

You may only be interested in one year, 1959. In that case, use the shown columns option to filter the columns to hide the other columns. Do this by typing `csvcmd -h -shown "Month;1959" file.csv`. You'll get an output that has the exact columns we want in a readable format:

![Screenshot from 2021-03-01 19-19-11](https://user-images.githubusercontent.com/13598541/109591842-21d8cd80-7ac3-11eb-8b8a-964e659b1278.png)

For splitting the files, use the split option. To retrieve the first third of the file, for example, use `-split 1/3`. Likewise this works for the second and third portions of the file by substituting for `-split 2/3` and `-split 3/3`. The full command now appears as `csvcmd -h -shown "Month;1959" -split 1/3 file.csv`. The output is exactly as we wanted:

![Screenshot from 2021-03-01 19-23-34](https://user-images.githubusercontent.com/13598541/109592143-a1ff3300-7ac3-11eb-9d9e-352292971e7b.png)

To store this in a file, remove the `-h` flag and output it to the desired file. The command will look something like `csvcmd -shown "Month;1959" -split 1/3 file.csv > output1.csv`. The file is ready for use and could be uploaded to a website, edited with excel, or used by another program. In a GUI spreadsheet program, the new CSV file will look something like:

![Screenshot from 2021-03-01 19-26-22](https://user-images.githubusercontent.com/13598541/109592373-05896080-7ac4-11eb-8cde-dc40408448a7.png)

## Another use case: Separating city data
We begin with the following data (truncated):

![Screenshot from 2021-03-01 19-29-13](https://user-images.githubusercontent.com/13598541/109592634-6d3fab80-7ac4-11eb-8f0e-aa43b106047a.png)


It contains a list of many cities around the US along with their latitude and longitude. We have the following question about the data: Which cities are in the West coast states?

First, filter the data with the filter option. In this case, our command should look something like `csvcmd -h -filter "State=OR;State=WA;State=CA" file.csv`. If we enter that in, we'll get our header row, but zero data back out.

![Screenshot from 2021-03-01 19-38-24](https://user-images.githubusercontent.com/13598541/109593403-bc3a1080-7ac5-11eb-8f40-2674605c6dcc.png)

This is because the `-or` flag was not set, meaning that for a row to show, every single option must be true. Instead, we would like a row to show up if any of the criteria if the filter is met. To do this, add `-or` to your command. It should now look like `csvcmd -h -filter "State=OR;State=WA;State=CA" -or file.csv`. The data you'll get in return will look similar to the following:

![Screenshot from 2021-03-01 19-45-49](https://user-images.githubusercontent.com/13598541/109593934-c0b2f900-7ac6-11eb-962c-d3e5a902cc54.png)


Now, filter the rows and output to a CSV file. The command I used is `csvcmd -filter "State=OR;State=WA;State=CA" -or -shown "State;City" file.csv > west_coast.csv`. 

![Screenshot from 2021-03-01 19-48-26](https://user-images.githubusercontent.com/13598541/109594126-18e9fb00-7ac7-11eb-80a2-e69dc4bbbf7f.png)

## Conclusion
`csvcmd` is a tool that can operate efficiently on CSV files by allowing the user to quickly hide irrelevant information or split the file into segments. The repository is actively being maintained and new features will be added periodically. If you feel that this project would be beneficial to your work, feel free to download it and give it a try. Thank you.


-Alex Anderson

alexandersonone@gmail.com
