# Ymlio 

Ymlio is a cli tool that works with yaml files and is capable of the following: 

## 1. Splitting 

### splits a multipe-file-yaml-file into single-yaml-files.

This is also capable of extracting only some files using the --only flag. It handles also other types of files (currently handles in addition to yaml also normal text files if they have the string **__RAW** and it can import content of other files if they have the string **__IMPORT**)

## 2. Combining

### combines single-yaml-files into a multipe-file-yaml-file

This can also handle other types of files mentioned under **1. Splitting**

# Ymlio is still in development:

# Testing

## You need to follow these steps to proparly test the tool:

#### 1. Clone the repository to your local machine.

`git clone https://github.com/mekkyz/ymlio.git`

#### 2. Go to the project folder.

`cd PATH-TO-/ymlio`

#### 3. Build.

`go build main.go`

#### 4. To avoid file clutter, move the executable to /Testing and go to /Testing
`mv main Testing && cd Testing`

#### 5. Test: Splitting.

- `./main split bigyaml.yml` This should split **all the files** in the bigyaml file.
- `./main split bigyaml.yml --only converter.yml database.yml` this should **only split** the specified files after the `--only` tag.
    - **Special case:** if there is only one file specified after `--only`   -> it will export this file and also **print it out to the terminal.**
    Example: `./main split bigyaml.yml --only converter.yml`
- The tool also handles `__RAW` Tags while splitting.
    - If there is a file in the bigyaml.yml file with the tag `__RAW`  The program will handle it as plain text. This is also handled in `combine` see **6. Test: Combining**

- The tool also handles `__IMPORT` Tags while splitting.
    - If there is a file in the bigyaml.yml file with the tag `__IMPORT`  The program will import the content of the specified file under `__IMPORT`.
    
    - In the bigyaml.yml file there are two files with `__IMPORT` which import their content from the `m1.md` & `m2.md` files. **(Please keep these files for testing purposes)**

- The tool also handles `anchors`
    - To know more about yaml anchors simply follow this **link:** https://support.atlassian.com/bitbucket-cloud/docs/yaml-anchors/

#### 6. Test: Combining.

- `./main combine converter.yml database.yml combined.yml`
This should combine **all the files before the last argument** and put them together in the last argument. The command here will combine `converter.yml` and `database.yml` and adds them to `combined.yml` (If the file `combined.yml` not there, it will be created)

- The combine also handles `__RAW`
    - If the program detects a file with a **text content** (not a yaml format). It will put the tag `__RAW` to it.