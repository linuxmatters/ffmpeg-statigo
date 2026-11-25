package main

import (
	"bufio"
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	arMagic              = "!<arch>\n"
	arHeaderSize         = 60
	coffHeaderSize       = 20
	coffBigObjHeaderSize = 56
	maxBSDNameProbeSize  = 4096
)

// coffBigObjClassID is the Windows GUID {D1BAA1C7-BAEE-4ba9-AF20-FAF66AA4DCB8}
// encoded in mixed-endian binary form, identifying ANON_OBJECT_HEADER_BIGOBJ.
var coffBigObjClassID = [16]byte{
	0xc7, 0xa1, 0xba, 0xd1,
	0xee, 0xba,
	0xa9, 0x4b,
	0xaf, 0x20, 0xfa, 0xf6, 0x6a, 0xa4, 0xdc, 0xb8,
}

type arMember struct {
	headerOffset int64
	size         int64
	keep         bool
}

// pruneZeroSectionCOFFMembers removes zero-section COFF object members and
// split DWARF (.dwo) members from a standard SVR4/GNU archive in two streaming
// passes without loading the full archive into memory. When no members are
// pruned, the archive is left untouched. When one or more members are pruned,
// symbol table members are omitted so the caller can rebuild the index with
// "ar s".
func pruneZeroSectionCOFFMembers(archivePath string) (int, error) {
	src, err := os.Open(archivePath)
	if err != nil {
		return 0, fmt.Errorf("open archive %s: %w", archivePath, err)
	}
	defer func() {
		if src != nil {
			_ = src.Close()
		}
	}()

	members, pruned, err := scanArchiveMembers(src)
	if err != nil {
		return 0, fmt.Errorf("scan archive %s: %w", archivePath, err)
	}
	if pruned == 0 {
		err = src.Close()
		src = nil
		return 0, err
	}

	tmp, err := os.CreateTemp(filepath.Dir(archivePath), ".filtered-*.a")
	if err != nil {
		return 0, fmt.Errorf("create filtered archive: %w", err)
	}
	tmpPath := tmp.Name()
	keepTmp := false
	defer func() {
		if tmp != nil {
			_ = tmp.Close()
		}
		if !keepTmp {
			_ = os.Remove(tmpPath)
		}
	}()

	if err := writeFilteredArchive(src, tmp, members); err != nil {
		return 0, fmt.Errorf("write filtered archive %s: %w", archivePath, err)
	}
	if err := tmp.Close(); err != nil {
		tmp = nil
		return 0, fmt.Errorf("close filtered archive: %w", err)
	}
	tmp = nil

	if err := src.Close(); err != nil {
		src = nil
		return 0, fmt.Errorf("close source archive %s: %w", archivePath, err)
	}
	src = nil

	if err := os.Rename(tmpPath, archivePath); err != nil {
		return 0, fmt.Errorf("replace archive %s: %w", archivePath, err)
	}
	keepTmp = true
	return pruned, nil
}

func scanArchiveMembers(src *os.File) ([]arMember, int, error) {
	info, err := src.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("stat archive: %w", err)
	}
	fileSize := info.Size()

	var magic [len(arMagic)]byte
	if _, err := io.ReadFull(src, magic[:]); err != nil {
		return nil, 0, fmt.Errorf("read archive magic: %w", err)
	}
	if string(magic[:]) != arMagic {
		return nil, 0, fmt.Errorf("unexpected archive magic %q", magic[:])
	}

	var hdr [arHeaderSize]byte
	var prefix [coffBigObjHeaderSize]byte
	var members []arMember
	var longNames []byte
	pruned := 0

	off := int64(len(arMagic))
	for off < fileSize {
		if off+arHeaderSize > fileSize {
			return nil, 0, fmt.Errorf("truncated archive header at offset %d", off)
		}
		if _, err := src.ReadAt(hdr[:], off); err != nil {
			return nil, 0, fmt.Errorf("read archive header at offset %d: %w", off, err)
		}
		if hdr[58] != '`' || hdr[59] != '\n' {
			return nil, 0, fmt.Errorf("malformed archive header at offset %d", off)
		}

		name := strings.TrimSpace(string(hdr[0:16]))
		sizeStr := strings.TrimSpace(string(hdr[48:58]))
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		payloadOff := off + arHeaderSize
		if err != nil || size < 0 || size > fileSize-payloadOff {
			return nil, 0, fmt.Errorf("malformed archive member size %q at offset %d", sizeStr, off)
		}
		pad := size & 1
		if pad > fileSize-payloadOff-size {
			return nil, 0, fmt.Errorf("truncated archive member padding at offset %d", off)
		}

		member := arMember{
			headerOffset: off,
			size:         size,
			keep:         true,
		}

		switch {
		case isArchiveSymbolTableName(name):
			member.keep = false
		case name == "//":
			// Preserve GNU long filename table verbatim and cache it to
			// resolve subsequent /<offset> member names.
			if size > 0 {
				longNames = make([]byte, size)
				if _, err := src.ReadAt(longNames, payloadOff); err != nil {
					return nil, 0, fmt.Errorf("read GNU long filename table at offset %d: %w", off, err)
				}
			} else {
				longNames = nil
			}
		default:
			dataOff := payloadOff
			payloadSize := size
			resolvedName := ""

			switch {
			case strings.HasPrefix(name, "#1/"):
				nameLen, err := strconv.ParseInt(name[3:], 10, 64)
				if err != nil || nameLen <= 0 || nameLen > size {
					return nil, 0, fmt.Errorf("malformed BSD archive member %q at offset %d", name, off)
				}
				if nameLen <= maxBSDNameProbeSize {
					bsdName := make([]byte, nameLen)
					if _, err := src.ReadAt(bsdName, payloadOff); err != nil {
						return nil, 0, fmt.Errorf("read BSD archive member name at offset %d: %w", off, err)
					}
					resolvedName = strings.TrimRight(string(bsdName), "\x00")
				}
				dataOff = payloadOff + nameLen
				payloadSize = size - nameLen
			case strings.HasPrefix(name, "/"):
				resolvedName = resolveGNULongName(name, longNames)
			default:
				resolvedName = strings.TrimSuffix(name, "/")
			}

			switch {
			case isArchiveSymbolTableName(resolvedName):
				member.keep = false
			case isSplitDWARFName(resolvedName):
				member.keep = false
				pruned++
			default:
				probeLen := min(payloadSize, int64(len(prefix)))
				if probeLen > 0 {
					if _, err := src.ReadAt(prefix[:probeLen], dataOff); err != nil {
						return nil, 0, fmt.Errorf("read archive member payload at offset %d: %w", off, err)
					}
				}
				if isZeroSectionCOFF(prefix[:probeLen], payloadSize) {
					member.keep = false
					pruned++
				}
			}
		}

		members = append(members, member)
		off = payloadOff + size + pad
	}

	return members, pruned, nil
}

func writeFilteredArchive(src, dst *os.File, members []arMember) error {
	bw := bufio.NewWriter(dst)
	if _, err := bw.WriteString(arMagic); err != nil {
		return err
	}

	var hdr [arHeaderSize]byte
	for _, m := range members {
		if !m.keep {
			continue
		}
		if _, err := src.ReadAt(hdr[:], m.headerOffset); err != nil {
			return err
		}
		if _, err := bw.Write(hdr[:]); err != nil {
			return err
		}
		if m.size > 0 {
			payload := io.NewSectionReader(src, m.headerOffset+arHeaderSize, m.size)
			if _, err := io.CopyN(bw, payload, m.size); err != nil {
				return err
			}
		}
		if m.size&1 != 0 {
			if err := bw.WriteByte('\n'); err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}

func isArchiveSymbolTableName(name string) bool {
	return name == "/" || name == "/SYM64/" || name == "__.SYMDEF" || name == "__.SYMDEF SORTED"
}

func resolveGNULongName(name string, longNames []byte) string {
	if len(name) < 2 || name[0] != '/' || name[1] < '0' || name[1] > '9' {
		return ""
	}
	off, err := strconv.Atoi(name[1:])
	if err != nil || off < 0 || off >= len(longNames) {
		return ""
	}
	rest := longNames[off:]
	if end := bytes.IndexAny(rest, "/\r\n\x00"); end >= 0 {
		rest = rest[:end]
	}
	return string(rest)
}

func isSplitDWARFName(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".dwo")
}

func isCOFFMachine(mach uint16) bool {
	switch mach {
	case pe.IMAGE_FILE_MACHINE_I386,
		pe.IMAGE_FILE_MACHINE_AMD64,
		pe.IMAGE_FILE_MACHINE_ARM,
		pe.IMAGE_FILE_MACHINE_ARMNT,
		pe.IMAGE_FILE_MACHINE_ARM64,
		pe.IMAGE_FILE_MACHINE_LOONGARCH32,
		pe.IMAGE_FILE_MACHINE_LOONGARCH64,
		pe.IMAGE_FILE_MACHINE_RISCV32,
		pe.IMAGE_FILE_MACHINE_RISCV64,
		pe.IMAGE_FILE_MACHINE_RISCV128:
		return true
	default:
		return false
	}
}

func isZeroSectionCOFF(prefix []byte, payloadSize int64) bool {
	if payloadSize >= coffHeaderSize && len(prefix) >= coffHeaderSize {
		mach := binary.LittleEndian.Uint16(prefix[0:2])
		if isCOFFMachine(mach) {
			return binary.LittleEndian.Uint16(prefix[2:4]) == 0
		}
	}
	if payloadSize >= coffBigObjHeaderSize && len(prefix) >= coffBigObjHeaderSize {
		sig1 := binary.LittleEndian.Uint16(prefix[0:2])
		sig2 := binary.LittleEndian.Uint16(prefix[2:4])
		version := binary.LittleEndian.Uint16(prefix[4:6])
		mach := binary.LittleEndian.Uint16(prefix[6:8])
		if sig1 == 0x0000 && sig2 == 0xffff && version >= 2 &&
			isCOFFMachine(mach) &&
			bytes.Equal(prefix[12:28], coffBigObjClassID[:]) {
			return binary.LittleEndian.Uint32(prefix[44:48]) == 0
		}
	}
	return false
}
