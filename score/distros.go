package score

func Top50Distros() []Distro {
	return []Distro{
		// Tier 1 - Principales (desktop-friendly)
		{ID: "cachyos", Name: "CachyOS", Rolling: 10, Easy: 6, DIY: 6, Performance: 10, DevFocus: 9, Popularity: 3698, Trend: TrendUp}, // ❌ no logo → caerá en generic
		{ID: "mint", Name: "Linux Mint", Rolling: 1, Easy: 10, DIY: 2, Performance: 6, DevFocus: 6, Popularity: 2714, Trend: TrendDown},
		{ID: "mx", Name: "MX Linux", Rolling: 3, Easy: 8, DIY: 4, Performance: 6, DevFocus: 6, Popularity: 1951, Trend: TrendStable},
		{ID: "debian", Name: "Debian", Rolling: 1, Easy: 7, DIY: 6, Performance: 7, DevFocus: 8, Popularity: 1589, Trend: TrendStable},
		{ID: "endeavour", Name: "EndeavourOS", Rolling: 10, Easy: 7, DIY: 7, Performance: 8, DevFocus: 9, Popularity: 1529, Trend: TrendUp},
		{ID: "pop", Name: "Pop!_OS", Rolling: 4, Easy: 9, DIY: 3, Performance: 8, DevFocus: 9, Popularity: 1346, Trend: TrendStable},
		{ID: "manjaro", Name: "Manjaro", Rolling: 10, Easy: 8, DIY: 5, Performance: 7, DevFocus: 8, Popularity: 1105, Trend: TrendStable},
		{ID: "ubuntu", Name: "Ubuntu", Rolling: 2, Easy: 10, DIY: 2, Performance: 7, DevFocus: 9, Popularity: 1072, Trend: TrendStable},
		{ID: "fedora_newlogo_newcolor", Name: "Fedora", Rolling: 5, Easy: 8, DIY: 4, Performance: 8, DevFocus: 9, Popularity: 1048, Trend: TrendStable},
		{ID: "zorin", Name: "Zorin OS", Rolling: 1, Easy: 10, DIY: 1, Performance: 6, DevFocus: 5, Popularity: 1004, Trend: TrendStable},
		{ID: "suse", Name: "openSUSE", Rolling: 6, Easy: 7, DIY: 6, Performance: 8, DevFocus: 8, Popularity: 789, Trend: TrendStable},
		{ID: "nobara", Name: "Nobara", Rolling: 8, Easy: 7, DIY: 4, Performance: 9, DevFocus: 6, Popularity: 721, Trend: TrendUp}, // ❌ no logo
		{ID: "elementary", Name: "elementary OS", Rolling: 1, Easy: 10, DIY: 1, Performance: 6, DevFocus: 5, Popularity: 593, Trend: TrendStable},
		{ID: "nixos", Name: "NixOS", Rolling: 10, Easy: 2, DIY: 10, Performance: 9, DevFocus: 10, Popularity: 555, Trend: TrendUp},
		{ID: "garuda", Name: "Garuda Linux", Rolling: 10, Easy: 7, DIY: 6, Performance: 9, DevFocus: 8, Popularity: 452, Trend: TrendStable},
		{ID: "kali", Name: "Kali Linux", Rolling: 5, Easy: 5, DIY: 6, Performance: 6, DevFocus: 7, Popularity: 419, Trend: TrendStable},
		{ID: "arch", Name: "Arch Linux", Rolling: 10, Easy: 3, DIY: 10, Performance: 9, DevFocus: 9, Popularity: 373, Trend: TrendStable},
		{ID: "alpine", Name: "Alpine Linux", Rolling: 5, Easy: 3, DIY: 8, Performance: 9, DevFocus: 8, Popularity: 353, Trend: TrendUp},

		// Tier 2
		{ID: "kubuntu", Name: "Kubuntu", Rolling: 2, Easy: 9, DIY: 3, Performance: 6, DevFocus: 7, Popularity: 313, Trend: TrendStable},
		{ID: "lite", Name: "Linux Lite", Rolling: 1, Easy: 10, DIY: 1, Performance: 6, DevFocus: 4, Popularity: 307, Trend: TrendStable},
		{ID: "tails", Name: "Tails", Rolling: 3, Easy: 6, DIY: 2, Performance: 4, DevFocus: 5, Popularity: 324, Trend: TrendStable},
		{ID: "parrot", Name: "Parrot OS", Rolling: 5, Easy: 6, DIY: 5, Performance: 6, DevFocus: 7, Popularity: 254, Trend: TrendStable},
		{ID: "void", Name: "Void Linux", Rolling: 10, Easy: 3, DIY: 9, Performance: 9, DevFocus: 8, Popularity: 203, Trend: TrendStable},
		{ID: "gentoo", Name: "Gentoo", Rolling: 10, Easy: 1, DIY: 10, Performance: 10, DevFocus: 9, Popularity: 218, Trend: TrendStable},
		{ID: "artix", Name: "Artix Linux", Rolling: 10, Easy: 4, DIY: 9, Performance: 8, DevFocus: 8, Popularity: 216, Trend: TrendStable},
		{ID: "solus", Name: "Solus", Rolling: 6, Easy: 9, DIY: 3, Performance: 7, DevFocus: 7, Popularity: 357, Trend: TrendStable},
		{ID: "qubes", Name: "Qubes OS", Rolling: 3, Easy: 3, DIY: 7, Performance: 6, DevFocus: 7, Popularity: 185, Trend: TrendStable},
		{ID: "rebornos", Name: "RebornOS", Rolling: 10, Easy: 7, DIY: 6, Performance: 8, DevFocus: 8, Popularity: 150, Trend: TrendStable},

		// Tier 3
		{ID: "antix", Name: "antiX", Rolling: 3, Easy: 6, DIY: 5, Performance: 8, DevFocus: 5, Popularity: 508, Trend: TrendStable},
		{ID: "lubuntu", Name: "Lubuntu", Rolling: 2, Easy: 9, DIY: 2, Performance: 7, DevFocus: 6, Popularity: 241, Trend: TrendStable},
		{ID: "xubuntu", Name: "Xubuntu", Rolling: 2, Easy: 9, DIY: 2, Performance: 7, DevFocus: 6, Popularity: 216, Trend: TrendStable},
		{ID: "openmandriva", Name: "OpenMandriva", Rolling: 6, Easy: 7, DIY: 4, Performance: 7, DevFocus: 7, Popularity: 262, Trend: TrendStable},
		{ID: "deepin", Name: "Deepin", Rolling: 3, Easy: 9, DIY: 2, Performance: 6, DevFocus: 6, Popularity: 238, Trend: TrendStable},
	}
}
