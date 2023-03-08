# Ymlio 

Ymlio is a cli tool that splits and combines yaml files and deals with multiple special cases (See below): 

## Installation

### Installation using golang

#### 1. Install golang on your system.

see https://go.dev/doc/install 

#### 2. Install ymlio using: 

`go install github.com/mekkyz/ymlio@latest`

### Installation without golang

#### 1. Get the binary

The ymlio tool is a binary file and needs no installation. The only prerequisite is that you have access as an administrator on the OS you are installing it on. Download the latest release for your OS from [here](https://github.com/mekkyz/ymlio/releases/)

#### 2. Make it an executable and run it.

On Linux, make this file executable by doing: `chmod u+x ymlio`.

On Windows, the file should be executable by default, i.e. do nothing.

On macOS, make this file executable by doing: `chmod u+x ymlio.amd.osx` or `chmod u+x ymlio.arm.osx`. If the there is a security pop-up when running the command, please also `Allow` the executable in `System Preferences > Security & Privacy`.

#### 3. For Linux users

If you want to use the tool globally just copy the `ymlio` file to your bin folder using `sudo cp ymlio /usr/local/bin/ymlio`

# Usage

## 1. Splitting 

### splits a multipe-file-yaml-file into single-yaml-files.

This is also capable of extracting only some files using the --only flag. It handles also other types of files (currently handles in addition to yaml also normal text files if they have the string **__RAW** and it can import content of other files if they have the string **__IMPORT**)

- `ymlio split file.yml` This should split **all the files** in the yaml file.
- `ymlio split file.yml --only key1.yml key2.yml` this should **only split** the specified files after the `--only` tag.
    - **Special case:** if there is only one file specified after `--only`   -> it will export this file and also **print it out to the terminal.** -> Stdout
    Example: `ymlio split file.yml --only key.yml`
- The tool also handles `__RAW` Tags while splitting.
    - If there is a file in the file.yml with the tag `__RAW`  The program will handle it as plain text. This is also handled in `combine`
- If the file name is `-` it will its content from **Stdin**
Example: `ymlio split -` or `ymlio split - --only key.yml`

- The tool also handles `__IMPORT` Tags while splitting.
    - If there is a file in the file.yml with the tag `__IMPORT`  The program will import the content of the specified file under `__IMPORT`.
    Example:
    ```yaml
        importtest
            __IMPORT: m1.md
    ```
    This will import the content of m1.md file. The m1.md must be in the same folder as the file.yml


- The tool also handles `anchors`
    - To know more about yaml anchors simply follow this **link:** https://support.atlassian.com/bitbucket-cloud/docs/yaml-anchors/



## 2. Combining

### combines single-yaml-files into a multipe-file-yaml-file

This can also handle other types of files mentioned under **1. Splitting**

- `ymlio combine file1.yml file2.yml combined.yml`
This should combine **all the files before the last argument** and put them together in the last argument. The command here will combine `file1.yml` and `file2.yml` and adds them to `combined.yml` (If the file `combined.yml` not there, it will be created)
    - If the `combined.yml` is already there. It will be overwritten
    - If you want to append to it using the flag `--extend`
    Example: `ymlio combine file1.yml file2.yml combined.yml --extend`

- The combine also handles `__RAW`
    - If the program detects a file with a **text content** (not a yaml format). It will put the tag `__RAW` to it.