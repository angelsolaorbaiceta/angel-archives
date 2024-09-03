# Angel Archives

An archiving tool that xz-compresses and bundles files together into an archive.

## Usage

Creating an archive:

```bash
$ aar -c -f archive.aarch file1.txt file2.txt file3.txt
```

Extracting all files an archive:

```bash
$ aar -x -f archive.aarch
```

Extracting a single file by name from an archive:

```bash
$ aar -x -f archive.aarch -n file2.txt
```

Listing the contents of an archive:

```bash
$ aar -l -f archive.aarch
```

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
```
