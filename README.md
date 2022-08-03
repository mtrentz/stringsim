# stringsim
Calculate the similarity between at least two strings. Accepts input from file and export your results to json or csv.

## Installation
```
git clone https://github.com/mtrentz/stringsim
go get
go build -o stringsim .
./stringsim hello world
```

To access it everywhere
```
sudo mv ./stringsim /usr/bin
stringsim hello world
```

## With docker
```
docker run mtrentz/stringsim hello world
```

## Usage
```
stringsim <s1> <s2> [<s3> ...] [flags]
stringsim -h
```

## Examples
```
# Comparing s1 to s2
  stringsim adam adan

# Comparing s1 to s2 and s3, case insensitive, output result to file
  stringsim adam adan Aden -i -o output.csv

# Reading s2, s3, ..., from a txt file separated by newlines and comparing to 'adam' using Levenshtein as metric
  stringsim adam --f2 strings.txt -m Levenshtein

# Reading many words from a json file (formated as array of strings ["a", "b", ...])
# and comparing each to every word in a txt file separated by newlines.
  stringsim --f1 strings_one.json --f2 strings_two.txt
  
# Reading and writing to file when running it in docker
  docker run -v $PWD:/app -it mtrentz/stringsim adam --f2 strings.txt -o output.json
```

