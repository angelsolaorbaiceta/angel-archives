# Angel Archives

An archiving tool that xz-compresses and bundles files together into an archive.
Archives can be encrypted and decrypted for maximum privacy.

## Installation

Make sure you’ve correctly installed Go v1.23 or greater and that the Go binaries are in your PATH (you want to append `$GOPATH/bin` to your `PATH` to access installed Go binaries).
Then:

```bash
$ go install https://github.com/angelsolaorbaiceta/angel-archives@latest
```

You should have the _aar_ utility in your path:

```bash
$ which aar
/Users/yourusername/go/bin/aar
```

Where in the example above, `/Users/yourusername/go` is the value of the `$GOPATH` variable.

### Installation from source

Alternatively, you can clone the repository:

```bash
$ git clone https://github.com/angelsolaorbaiceta/angel-archives
```

And install the binary in your `$GOPATH/bin` by simply doing:

```bash
$ make install
```

You should have the _aar_ utility in your path:

```bash
$ which aar
/Users/yourusername/go/bin/aar
```

## Usage

Creating an archive:

```bash
$ aar create -f archive.aarch file1.txt file2.txt file3.txt
```

Extracting all files an archive:

```bash
$ aar extract -f archive.aarch
```

Extracting a single file by name from an archive:

```bash
$ aar extract -f archive.aarch -n file2.txt
```

Listing the contents of an archive:

```bash
$ aar list -f archive.aarch
```

Encrypting an archive:

```bash
$ aar encrypt -f archive.aarch
Password: <password>
Confirm password: <password>
```

Where `<password>` is the password you want to use to encrypt the archive, with a minimum length of 8 characters.
It removes the original _.aarch_ file and creates a new one with the encrypted data, with extension _.aarch.enc_.

> [!NOTE]
> The encryption is done using the AES-256-GCM algorithm, and it only works for angel archives.

Decrypting an archive:

```bash
$ aar decrypt -f archive.aarch.enc
Password: <password>
Confirm password: <password>
```

Where `<password>` is the password you used to encrypt the archive.
It removes the encrypted _.aarch.enc_ file and creates a new one with the decrypted data, with extension _.aarch_.

> [!NOTE]
> The decryption is done using the AES-256-GCM algorithm, and it only works for encrypted angel archives.

## File Format

### Archive Header

The archive file starts with a header that contains the following:

- **Magic**: A 4-byte sequence that identifies the file as an Angel Archive. The sequence is "AAR?" (0x41 0x41 0x52 0x3F)."
- **Header length**: A 4-byte integer that specifies the length of the header in bytes.
- **Files**: A list of files that are included in the archive. Each file entry contains the following:
  - **Name length**: A 2-byte integer that specifies the length of the file name in bytes.
  - **Name**: The name of the file.
  - **Offset**: A 4-byte integer that specifies the offset of the file data in the archive.
  - **Length**: A 4-byte integer that specifies the length of the file data in bytes.

Example:

```
0x41 0x41 0x52 0x3F // Magic
0x1A 0x00 0x00 0x00 // Header length (26 bytes)
0x08 0x00           // Name length (8 bytes)
"test.txt"          // File name
0x00 0x00 0x00 0x1B // Offset (27 bytes)
0x0B 0x00 0x00 0x00 // Length (11 bytes)
...                 // Next go the file bytes
```

### Archive Files

The files are stored sequentially after the header.
Their raw bytes are xz-compressed before being saved to disk.
