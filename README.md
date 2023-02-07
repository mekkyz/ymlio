# ymlio

This is a command line tool that works with yaml files and is capable of the following: 

## 1. split 

### splits a multipe-file-yaml-file into single-yaml-files.

This is also capable of extracting only some files using the --only flag. It handles also other types of files (currently handles in addition to yaml also normal text files if they have the string **__RAW** and it can import content of other files if they have the string **__IMPORT**)

## 2. combine

### combines single-yaml-files into a multipe-file-yaml-file

This can also handle other types of files mentioned under 1. split

# ----------------------------------------------------------------

# This CLI-Tool is still in development:

# ------------Testing------------

## You need to follow the following steps to proparly test the tool:

#### 1. Clone the repository to your local machine.

`git clone https://github.com/mekkyz/ymlio.git`

#### 2. Go to the project folder.

`cd PATH-TO-/ymlio`

#### 3. Build.

`go build main.go`

#### 4. To avoid file clutter. Move the executable to /Testing and go to /Testing
`mv main Testing && cd Testing`

#### 5. Testing: Split.

- `./main split bigyaml.yml`
This should split **all the files** in the bigyaml file.
- `./main split bigyaml.yml --only converter.yml database.yml` this should **only split** the specified files after the `--only`  tag.
    - **Special case:** if there is only one file specified after `--only`   -> it will only export this file and also **print it out to the terminal.**
    Example: `./main split bigyaml.yml --only converter.yml`
- The tool also handles `__RAW` Tags while splitting.
    - If there is a file in the bigyaml.yml file with the tag `__RAW`  The program will handle it as plain text. This is also handled in `combine` see 6.

- The tool also handles `__IMPORT` Tags while splitting.
    - If there is a file in the bigyaml.yml file with the tag `__IMPORT`  The program will import the contect of the specified file under `__IMPORT`.
    
    - In the bigyaml.yml file there are two files with `__IMPORT` which import their content from the `m1.md` & `m2.md` files. **(Please keep this files for testing purposes)**

- The tool also handles `anchors`
    - To know more about yaml anchors simply follow this link https://support.atlassian.com/bitbucket-cloud/docs/yaml-anchors/

#### 6. Testing: Combine.

- `./main combine converter.yml database.yml combined.yml`
This should combine **all the files before the last argument** and put them together in the last argument. The command here will combine `converter.yml` and `database.yml` and adds them to `combined.yml` (If the file `combined.yml` not there, it will be created)

- The combine also handles `__RAW`
    - If the program detects a file with a **text content** (not a yaml format). It will put the tag `__RAW` to it.