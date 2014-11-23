# crypt

crypt is a small app to encrypt and decrypt files with aes encryption.  It is 
portable and works on OSX, Linux, and Windows along with any other OS that 
supports go.

A key can be specified with the `--key` flag or via stdin if `--key` is 
not specified.  You may want to provide your key via stdin for encryption
and decryption so that it is not saved in your bash history.

### Examples

**Encrypt the contents from STDIN to a file**
```bash
> echo "hello world" | crypt -e --key test -i encrypted-text
```

**Decrypt the contents from STDIN to STDOUT**
```bash
> cat encrypted-text | crypt -d --key test -i -o
hello world
```

**Encrypt a large tar file to a new file**
```bash
> crypt -e docker-image.tar encrypted-image.xxx
please enter your key:
> secret
1.24 GB / 1.25 GB [=========================================================================] 99.70 % 58.72 MB/s
# test to see if I can read the file
> tar -tvf encrypted-image.xxx
tar: This does not look like a tar archive
tar: Skipping to next header
tar: Exiting with failure status due to previous errors
```

**Decrypt a large file back to the original contents**
```bash
> crypt -d encrypted-image.xxx docker-image-unencrypted.tar
please enter your key:
> secret
1.24 GB / 1.25 GB [=========================================================================] 99.45 % 70.80 MB/s
# look at the tar headers
> tar -tvf docker-image-unencrypted.tar
drwxr-xr-x 0/0               0 2014-11-22 21:14 05c6d847812ccad7f327f12ec2404dd972fc64e65c6c8a996b402c8e3f990d7c/
-rw-r--r-- 0/0               3 2014-11-22 21:14 05c6d847812ccad7f327f12ec2404dd972fc64e65c6c8a996b402c8e3f990d7c/
VERSION
```

```bash
NAME:
   crypt - encrypt and decrypt files easily

USAGE:
   crypt [global options] command [command options] [arguments...]

VERSION:
   1

AUTHOR:
  @crosbymichael

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key                key to use for the encryption algo
   --encrypt, -e        encrypt a file
   --decrypt, -d        decrypt a file
   --stdin, -i          accept input for STDIN
   --stdout, -o         return output to STDOUT
   --help, -h           show help
   --version, -v        print the version
```

### TODO:
* Keep it simple stupid

## License - MIT
