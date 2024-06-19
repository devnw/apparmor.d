// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const PIVOTROOT Kind = "pivot_root"

type PivotRoot struct {
	RuleBase
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

func newPivotRootFromLog(log map[string]string) Rule {
	return &PivotRoot{
		RuleBase:      newBaseFromLog(log),
		Qualifier:     newQualifierFromLog(log),
		OldRoot:       log["srcname"],
		NewRoot:       log["name"],
		TargetProfile: "",
	}
}

func (r *PivotRoot) Validate() error {
	return nil
}

func (r *PivotRoot) Compare(other Rule) int {
	o, _ := other.(*PivotRoot)
	if res := compare(r.OldRoot, o.OldRoot); res != 0 {
		return res
	}
	if res := compare(r.NewRoot, o.NewRoot); res != 0 {
		return res
	}
	if res := compare(r.TargetProfile, o.TargetProfile); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *PivotRoot) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *PivotRoot) Constraint() constraint {
	return blockKind
}

func (r *PivotRoot) Kind() Kind {
	return PIVOTROOT
}
