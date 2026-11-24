// Copyright 2026 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libc // import "modernc.org/libc"

import (
	"syscall"
	"unsafe"

	"modernc.org/libc/errno"
)

// On Windows errno holds the C runtime's values, those of errno.h, while
// GetLastError returns the system error code of the call that failed. setErrno
// keeps the two apart: errno gets the system error mapped as the UCRT maps it,
// and t.lastError, which GetLastError returns, the system error itself.

// setErrnoValue sets errno to the C errno value e. GetLastError returns e too,
// as it did when it read errno.
func (t *TLS) setErrnoValue(e int32) {
	*(*int32)(unsafe.Pointer(t.errnop)) = e
	t.lastError = uint32(e)
}

// setErrnoSyscall sets errno for the system error e, and GetLastError to e.
func (t *TLS) setErrnoSyscall(e syscallErrno) {
	*(*int32)(unsafe.Pointer(t.errnop)) = winErrno(e)
	t.lastError = uint32(e)
}

// winErrno returns the errno value for the system error e. For Windows error
// codes it is the UCRT's mapping (__acrt_errno_from_os_error), for the values
// package syscall invents for package os, the errno of the same name.
func winErrno(e syscall.Errno) int32 {
	if e == 0 {
		return 0
	}

	if e >= syscall.APPLICATION_ERROR {
		if v, ok := goErrnos[e]; ok {
			return v
		}

		return errno.EINVAL
	}

	if v, ok := winErrnos[uint32(e)]; ok {
		return v
	}

	switch {
	case e >= 19 && e <= 36: // ERROR_WRITE_PROTECT .. ERROR_SHARING_BUFFER_EXCEEDED
		return errno.EACCES
	case e >= 188 && e <= 202: // ERROR_INVALID_STARTING_CODESEG .. ERROR_INFLOOP_IN_RELOC_CHAIN
		return errno.ENOEXEC
	}
	return errno.EINVAL
}

// winErrnos is the UCRT's table of Windows error codes and their errno values.
var winErrnos = map[uint32]int32{
	1:    errno.EINVAL,    // ERROR_INVALID_FUNCTION
	2:    errno.ENOENT,    // ERROR_FILE_NOT_FOUND
	3:    errno.ENOENT,    // ERROR_PATH_NOT_FOUND
	4:    errno.EMFILE,    // ERROR_TOO_MANY_OPEN_FILES
	5:    errno.EACCES,    // ERROR_ACCESS_DENIED
	6:    errno.EBADF,     // ERROR_INVALID_HANDLE
	7:    errno.ENOMEM,    // ERROR_ARENA_TRASHED
	8:    errno.ENOMEM,    // ERROR_NOT_ENOUGH_MEMORY
	9:    errno.ENOMEM,    // ERROR_INVALID_BLOCK
	10:   errno.E2BIG,     // ERROR_BAD_ENVIRONMENT
	11:   errno.ENOEXEC,   // ERROR_BAD_FORMAT
	12:   errno.EINVAL,    // ERROR_INVALID_ACCESS
	13:   errno.EINVAL,    // ERROR_INVALID_DATA
	15:   errno.ENOENT,    // ERROR_INVALID_DRIVE
	16:   errno.EACCES,    // ERROR_CURRENT_DIRECTORY
	17:   errno.EXDEV,     // ERROR_NOT_SAME_DEVICE
	18:   errno.ENOENT,    // ERROR_NO_MORE_FILES
	33:   errno.EACCES,    // ERROR_LOCK_VIOLATION
	53:   errno.ENOENT,    // ERROR_BAD_NETPATH
	65:   errno.EACCES,    // ERROR_NETWORK_ACCESS_DENIED
	67:   errno.ENOENT,    // ERROR_BAD_NET_NAME
	80:   errno.EEXIST,    // ERROR_FILE_EXISTS
	82:   errno.EACCES,    // ERROR_CANNOT_MAKE
	83:   errno.EACCES,    // ERROR_FAIL_I24
	87:   errno.EINVAL,    // ERROR_INVALID_PARAMETER
	89:   errno.EAGAIN,    // ERROR_NO_PROC_SLOTS
	108:  errno.EACCES,    // ERROR_DRIVE_LOCKED
	109:  errno.EPIPE,     // ERROR_BROKEN_PIPE
	112:  errno.ENOSPC,    // ERROR_DISK_FULL
	114:  errno.EBADF,     // ERROR_INVALID_TARGET_HANDLE
	128:  errno.ECHILD,    // ERROR_WAIT_NO_CHILDREN
	129:  errno.ECHILD,    // ERROR_CHILD_NOT_COMPLETE
	130:  errno.EBADF,     // ERROR_DIRECT_ACCESS_HANDLE
	131:  errno.EINVAL,    // ERROR_NEGATIVE_SEEK
	132:  errno.EACCES,    // ERROR_SEEK_ON_DEVICE
	145:  errno.ENOTEMPTY, // ERROR_DIR_NOT_EMPTY
	158:  errno.EACCES,    // ERROR_NOT_LOCKED
	161:  errno.ENOENT,    // ERROR_BAD_PATHNAME
	164:  errno.EAGAIN,    // ERROR_MAX_THRDS_REACHED
	167:  errno.EACCES,    // ERROR_LOCK_FAILED
	183:  errno.EEXIST,    // ERROR_ALREADY_EXISTS
	206:  errno.ENOENT,    // ERROR_FILENAME_EXCED_RANGE
	215:  errno.EAGAIN,    // ERROR_NESTING_NOT_ALLOWED
	1113: errno.EILSEQ,    // ERROR_NO_UNICODE_TRANSLATION
	1816: errno.ENOMEM,    // ERROR_NOT_ENOUGH_QUOTA
}

// goErrnos maps the errno values package syscall invents on Windows, from
// syscall.APPLICATION_ERROR on, to the C errno values of the same name.
var goErrnos = map[syscall.Errno]int32{
	syscall.E2BIG:           errno.E2BIG,
	syscall.EACCES:          errno.EACCES,
	syscall.EADDRINUSE:      errno.EADDRINUSE,
	syscall.EADDRNOTAVAIL:   errno.EADDRNOTAVAIL,
	syscall.EAFNOSUPPORT:    errno.EAFNOSUPPORT,
	syscall.EAGAIN:          errno.EAGAIN,
	syscall.EALREADY:        errno.EALREADY,
	syscall.EBADF:           errno.EBADF,
	syscall.EBADMSG:         errno.EBADMSG,
	syscall.EBUSY:           errno.EBUSY,
	syscall.ECANCELED:       errno.ECANCELED,
	syscall.ECHILD:          errno.ECHILD,
	syscall.ECONNABORTED:    errno.ECONNABORTED,
	syscall.ECONNREFUSED:    errno.ECONNREFUSED,
	syscall.ECONNRESET:      errno.ECONNRESET,
	syscall.EDEADLK:         errno.EDEADLK,
	syscall.EDEADLOCK:       errno.EDEADLOCK,
	syscall.EDESTADDRREQ:    errno.EDESTADDRREQ,
	syscall.EDOM:            errno.EDOM,
	syscall.EEXIST:          errno.EEXIST,
	syscall.EFAULT:          errno.EFAULT,
	syscall.EFBIG:           errno.EFBIG,
	syscall.EHOSTUNREACH:    errno.EHOSTUNREACH,
	syscall.EIDRM:           errno.EIDRM,
	syscall.EILSEQ:          errno.EILSEQ,
	syscall.EINPROGRESS:     errno.EINPROGRESS,
	syscall.EINTR:           errno.EINTR,
	syscall.EINVAL:          errno.EINVAL,
	syscall.EIO:             errno.EIO,
	syscall.EISCONN:         errno.EISCONN,
	syscall.EISDIR:          errno.EISDIR,
	syscall.ELOOP:           errno.ELOOP,
	syscall.EMFILE:          errno.EMFILE,
	syscall.EMLINK:          errno.EMLINK,
	syscall.EMSGSIZE:        errno.EMSGSIZE,
	syscall.ENAMETOOLONG:    errno.ENAMETOOLONG,
	syscall.ENETDOWN:        errno.ENETDOWN,
	syscall.ENETRESET:       errno.ENETRESET,
	syscall.ENETUNREACH:     errno.ENETUNREACH,
	syscall.ENFILE:          errno.ENFILE,
	syscall.ENOBUFS:         errno.ENOBUFS,
	syscall.ENODATA:         errno.ENODATA,
	syscall.ENODEV:          errno.ENODEV,
	syscall.ENOEXEC:         errno.ENOEXEC,
	syscall.ENOLCK:          errno.ENOLCK,
	syscall.ENOLINK:         errno.ENOLINK,
	syscall.ENOMEM:          errno.ENOMEM,
	syscall.ENOMSG:          errno.ENOMSG,
	syscall.ENOPROTOOPT:     errno.ENOPROTOOPT,
	syscall.ENOSPC:          errno.ENOSPC,
	syscall.ENOSR:           errno.ENOSR,
	syscall.ENOSTR:          errno.ENOSTR,
	syscall.ENOSYS:          errno.ENOSYS,
	syscall.ENOTCONN:        errno.ENOTCONN,
	syscall.ENOTEMPTY:       errno.ENOTEMPTY,
	syscall.ENOTRECOVERABLE: errno.ENOTRECOVERABLE,
	syscall.ENOTSOCK:        errno.ENOTSOCK,
	syscall.ENOTSUP:         errno.ENOTSUP,
	syscall.ENOTTY:          errno.ENOTTY,
	syscall.ENXIO:           errno.ENXIO,
	syscall.EOPNOTSUPP:      errno.EOPNOTSUPP,
	syscall.EOVERFLOW:       errno.EOVERFLOW,
	syscall.EOWNERDEAD:      errno.EOWNERDEAD,
	syscall.EPERM:           errno.EPERM,
	syscall.EPIPE:           errno.EPIPE,
	syscall.EPROTO:          errno.EPROTO,
	syscall.EPROTONOSUPPORT: errno.EPROTONOSUPPORT,
	syscall.EPROTOTYPE:      errno.EPROTOTYPE,
	syscall.ERANGE:          errno.ERANGE,
	syscall.EROFS:           errno.EROFS,
	syscall.ESPIPE:          errno.ESPIPE,
	syscall.ESRCH:           errno.ESRCH,
	syscall.ETIME:           errno.ETIME,
	syscall.ETIMEDOUT:       errno.ETIMEDOUT,
	syscall.ETXTBSY:         errno.ETXTBSY,
	syscall.EWOULDBLOCK:     errno.EWOULDBLOCK,
	syscall.EXDEV:           errno.EXDEV,
	syscall.EWINDOWS:        errno.ENOSYS,
}

// strerrorTexts are the UCRT's messages. strerror gives "Unknown error" for the
// other values.
var strerrorTexts = map[int32]string{
	0:   "No error",
	1:   "Operation not permitted",
	2:   "No such file or directory",
	3:   "No such process",
	4:   "Interrupted function call",
	5:   "Input/output error",
	6:   "No such device or address",
	7:   "Arg list too long",
	8:   "Exec format error",
	9:   "Bad file descriptor",
	10:  "No child processes",
	11:  "Resource temporarily unavailable",
	12:  "Not enough space",
	13:  "Permission denied",
	14:  "Bad address",
	16:  "Resource device",
	17:  "File exists",
	18:  "Improper link",
	19:  "No such device",
	20:  "Not a directory",
	21:  "Is a directory",
	22:  "Invalid argument",
	23:  "Too many open files in system",
	24:  "Too many open files",
	25:  "Inappropriate I/O control operation",
	27:  "File too large",
	28:  "No space left on device",
	29:  "Invalid seek",
	30:  "Read-only file system",
	31:  "Too many links",
	32:  "Broken pipe",
	33:  "Domain error",
	34:  "Result too large",
	36:  "Resource deadlock avoided",
	38:  "Filename too long",
	39:  "No locks available",
	40:  "Function not implemented",
	41:  "Directory not empty",
	42:  "Illegal byte sequence",
	100: "address in use",
	101: "address not available",
	102: "address family not supported",
	103: "connection already in progress",
	104: "bad message",
	105: "operation canceled",
	106: "connection aborted",
	107: "connection refused",
	108: "connection reset",
	109: "destination address required",
	110: "host unreachable",
	111: "identifier removed",
	112: "operation in progress",
	113: "already connected",
	114: "too many symbolic link levels",
	115: "message size",
	116: "network down",
	117: "network reset",
	118: "network unreachable",
	119: "no buffer space",
	120: "no message available",
	121: "no link",
	122: "no message",
	123: "no protocol option",
	124: "no stream resources",
	125: "not a stream",
	126: "not connected",
	127: "state not recoverable",
	128: "not a socket",
	129: "not supported",
	130: "operation not supported",
	132: "value too large",
	133: "owner dead",
	134: "protocol error",
	135: "protocol not supported",
	136: "wrong protocol type",
	137: "stream timeout",
	138: "timed out",
	139: "text file busy",
	140: "operation would block",
}
