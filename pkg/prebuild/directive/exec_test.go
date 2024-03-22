// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/arduino/go-paths-helper"
)

func TestExec_Apply(t *testing.T) {
	tests := []struct {
		name          string
		rootApparmord *paths.Path
		opt           *Option
		profile       string
		want          string
	}{
		{
			name:          "exec",
			rootApparmord: paths.New("../../../apparmor.d/groups/kde/"),
			opt: &Option{
				Name: "exec",
				Args: map[string]string{"DiscoverNotifier": ""},
				File: nil,
				Raw:  "  #aa:exec DiscoverNotifier",
			},
			profile: `  #aa:exec DiscoverNotifier`,
			want: `  @{lib}/@{multiarch}/{,libexec/}DiscoverNotifier Px,
  @{lib}/DiscoverNotifier Px,`,
		},
		{
			name:          "exec-unconfined",
			rootApparmord: paths.New("../../../apparmor.d/groups/freedesktop/"),
			opt: &Option{
				Name: "exec",
				Args: map[string]string{"U": "", "polkit-agent-helper": ""},
				File: nil,
				Raw:  "  #aa:exec U polkit-agent-helper",
			},
			profile: `  #aa:exec U polkit-agent-helper`,
			want: `  @{lib}/polkit-[0-9]/polkit-agent-helper-[0-9] Ux,
  @{lib}/polkit-agent-helper-[0-9] Ux,`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootApparmord = tt.rootApparmord
			if got := Directives["exec"].Apply(tt.opt, tt.profile); got != tt.want {
				t.Errorf("Exec.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}