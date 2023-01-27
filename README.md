# ymlio

This is a command line tool that works with yaml files and is capable of the following: 

## 1. split 

### splits a multipe-file-yaml-file into single-yaml-files.

This is also capable of extracting only some files using the --only flag. It handles also other types of files (currently handles in addition to yaml also normal text files if they have the string **__RAW** and it can import content of other files if they have the string **__IMPORT**)

## 2. combine

### combines single-yaml-files into a multipe-file-yaml-file

This can also handle other types of files mentioned under 1. split