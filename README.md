# dl_cli (Downloader CLI)
** Downloads files from list of source URLs at a well defined specified path

# Pre-requisites
The CLI has been tested on ubuntu18 (64 bit), golang1.17.6.

# How to run?
Unzip the code and place it in $GO_PATH/src/github.com
Then on terminal, run

```sh
make clean; make build
```

This will create a binary in $GO_PATH/src/github.com/dl_cli/bin

# Usage

## dl_cli Commands
Below are supported commands for `dl_cli`

```sh
% ./dl_cli --help
CLI to download files.

Usage:
  dl_cli [command]

Available Commands:
  download    download files from list of urls
  help        Help about any command
  version     Current version of file downloader CLI being used

Flags:
  -h, --help   help for dl_cli

Use "dl_cli [command] --help" for more information about a command.
```

## Examples

``` sh
% ./dl_cli version
dl_cli version: v1.0

# Download the files represented by comma separated urls at a specified path.
% ./dl_cli download -p <path to download> -u <comma separated list of urls>

# Sample command to download files,
% ./dl_cli download -p /tmp/test -u http://my.file.com/file,https://you.file.com/test
```

# Notes about the implementation

- The code is written in Golang.
- The code is multi-threaded.
- Used waitgroups to make sure all the file downloads are completed before the main process exits.
- Currently supports http and https protocols.
- Used protocol interface to extend it to other protocols.
- Used combination of http.Get and io.copy to make sure that process does not go out-of-memory.
- This is achieved due to the fact that io.Copy copies data in chunks of 32k.
- If download fails, partial file is deleted.
- In case of failures, rollbacks are carried out.
- The code creates timestamp based directory, followed by sub-dir structure in order to avoid collisions.

# Drawbacks

- Not able to get automated unit tests implemented. I ran out of time due to my job commitments.
- Looking back, I could have structured the code a bit differently, creating different interfaces.
- Not tested with directory urls.
- Could have added checksumming logic.

# Manual Testing

- Tested with real urls from S3
- Tested for larget files
- Tested for out of memory on local file system
- Tested for protocols that are not http or https

Example output of testing

% ./bin/dl_cli download -p /tmp -u https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/macos/appctl,https://pmkft-assets.s3.us-west-1.amazonaws.com/non-existing-file,https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/windows/appctl,ftp://my.file.com/file,https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/linux/appctl
***************************
Source URL   : https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/macos/appctl
Destination  : /tmp/2022-03-29-16-21-06/pmkft-assets.s3.us-west-1.amazonaws.com/appctl/macos/appctl
Status       : Success
***************************
Source URL   : https://pmkft-assets.s3.us-west-1.amazonaws.com/non-existing-file
Destination  : /tmp/2022-03-29-16-21-06/pmkft-assets.s3.us-west-1.amazonaws.com/non-existing-file
Status       : Failed, Reason : Download failed for URL: https://pmkft-assets.s3.us-west-1.amazonaws.com/non-existing-file
***************************
Source URL   : https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/windows/appctl
Destination  : /tmp/2022-03-29-16-21-06/pmkft-assets.s3.us-west-1.amazonaws.com/appctl/windows/appctl
Status       : Success
***************************
Source URL   : ftp://my.file.com/file
Destination  : /tmp/2022-03-29-16-21-06/my.file.com/file
Status       : Failed, Reason : Error: Protocol Unsupported: ftp, URL: ftp://my.file.com/file
***************************
Source URL   : https://pmkft-assets.s3.us-west-1.amazonaws.com/appctl/linux/appctl
Destination  : /tmp/2022-03-29-16-21-06/pmkft-assets.s3.us-west-1.amazonaws.com/appctl/linux/appctl
Status       : Success
***************************
Download operation complete, all URLs processed.

% tree /tmp/2022-03-29-16-21-06/
/tmp/2022-03-29-16-21-06/
`-- pmkft-assets.s3.us-west-1.amazonaws.com
    `-- appctl
        |-- linux
        |   `-- appctl
        |-- macos
        |   `-- appctl
        `-- windows
            `-- appctl

5 directories, 3 files
