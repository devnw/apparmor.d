// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

const (
	tokLINK   = "link"
	tokFILE   = "file"
	tokOWNER  = "owner"
	tokSUBSET = "subset"
)

func init() {
	requirements[tokFILE] = requirement{
		"access": {"m", "r", "w", "l", "k"},
		"transition": {
			"ix", "ux", "Ux", "px", "Px", "cx", "Cx", "pix", "Pix", "cix",
			"Cix", "pux", "PUx", "cux", "CUx", "x",
		},
	}
}

func isOwner(log map[string]string) bool {
	fsuid, hasFsUID := log["fsuid"]
	ouid, hasOuUID := log["ouid"]
	isDbus := strings.Contains(log["operation"], "dbus")
	if hasFsUID && hasOuUID && fsuid == ouid && ouid != "0" && !isDbus {
		return true
	}
	return false
}

// cmpFileAccess compares two access strings for file rules.
// It is aimed to be used in slices.SortFunc.
func cmpFileAccess(i, j string) int {
	if slices.Contains(requirements[tokFILE]["access"], i) &&
		slices.Contains(requirements[tokFILE]["access"], j) {
		return requirementsWeights[tokFILE]["access"][i] - requirementsWeights[tokFILE]["access"][j]
	}
	if slices.Contains(requirements[tokFILE]["transition"], i) &&
		slices.Contains(requirements[tokFILE]["transition"], j) {
		return requirementsWeights[tokFILE]["transition"][i] - requirementsWeights[tokFILE]["transition"][j]
	}
	if slices.Contains(requirements[tokFILE]["access"], i) {
		return -1
	}
	return 1
}

type File struct {
	RuleBase
	Qualifier
	Owner  bool
	Path   string
	Access []string
	Target string
}

func newFileFromLog(log map[string]string) Rule {
	accesses, err := toAccess("file-log", log["requested_mask"])
	if err != nil {
		panic(fmt.Errorf("newFileFromLog(%v): %w", log, err))
	}
	if slices.Compare(accesses, []string{"l"}) == 0 {
		return newLinkFromLog(log)
	}
	return &File{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Access:    accesses,
		Target:    log["target"],
	}
}

func (r *File) Less(other any) bool {
	o, _ := other.(*File)
	letterR := getLetterIn(fileAlphabet, r.Path)
	letterO := getLetterIn(fileAlphabet, o.Path)
	if fileWeights[letterR] != fileWeights[letterO] && letterR != "" && letterO != "" {
		return fileWeights[letterR] < fileWeights[letterO]
	}
	if r.Path != o.Path {
		return r.Path < o.Path
	}
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Target != o.Target {
		return r.Target < o.Target
	}
	if o.Owner != r.Owner {
		return r.Owner
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *File) Equals(other any) bool {
	o, _ := other.(*File)
	return r.Path == o.Path && slices.Equal(r.Access, o.Access) && r.Owner == o.Owner &&
		r.Target == o.Target && r.Qualifier.Equals(o.Qualifier)
}

func (r *File) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *File) Constraint() constraint {
	return blockKind
}

func (r *File) Kind() string {
	return tokFILE
}

type Link struct {
	RuleBase
	Qualifier
	Owner  bool
	Subset bool
	Path   string
	Target string
}

func newLinkFromLog(log map[string]string) Rule {
	return &Link{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Target:    log["target"],
	}
}

func (r *Link) Less(other any) bool {
	o, _ := other.(*Link)
	if r.Path != o.Path {
		return r.Path < o.Path
	}
	if r.Target != o.Target {
		return r.Target < o.Target
	}
	if o.Owner != r.Owner {
		return r.Owner
	}
	if r.Subset != o.Subset {
		return r.Subset
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Link) Equals(other any) bool {
	o, _ := other.(*Link)
	return r.Subset == o.Subset && r.Owner == o.Owner && r.Path == o.Path &&
		r.Target == o.Target && r.Qualifier.Equals(o.Qualifier)
}

func (r *Link) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Link) Constraint() constraint {
	return blockKind
}

func (r *Link) Kind() string {
	return tokLINK
}
