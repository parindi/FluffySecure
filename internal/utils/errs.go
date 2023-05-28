package utils

import (
	"fmt"
	"runtime"
)

// ErrSliceSortAlphabetical is a helper type that can be used with sort.Sort to sort a slice of errors in alphabetical
// order. Usage is simple just do sort.Sort(ErrSliceSortAlphabetical([]error{})).
type ErrSliceSortAlphabetical []error

func (s ErrSliceSortAlphabetical) Len() int { return len(s) }

func (s ErrSliceSortAlphabetical) Less(i, j int) bool { return s[i].Error() < s[j].Error() }

func (s ErrSliceSortAlphabetical) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetExpectedErrTxt returns error text for expected errs. THIS IS A TEST UTILITY FUNCTION.
func GetExpectedErrTxt(err string) string {
	switch runtime.GOOS {
	case windows:
		switch err {
		case strPathNotFound:
			return fmt.Sprintf(errFmtWindowsNotFound, strOpen, strPath)
		case strStat + strPathNotFound:
			return fmt.Sprintf(errFmtWindowsNotFound, strStat, strPath)
		case strFileNotFound:
			return fmt.Sprintf(errFmtWindowsNotFound, strOpen, strFile)
		case strStat + strFileNotFound:
			return fmt.Sprintf(errFmtWindowsNotFound, strStat, strFile)
		case strIsDir:
			return "read %s: The handle is invalid."
		}
	default:
		switch err {
		case strPathNotFound, strFileNotFound:
			return fmt.Sprintf(errFmtLinuxNotFound, strOpen)
		case strStat + strPathNotFound, strStat + strFileNotFound:
			return fmt.Sprintf(errFmtLinuxNotFound, strStat)
		case strIsDir:
			return "read %s: is a directory"
		}
	}

	return ""
}
