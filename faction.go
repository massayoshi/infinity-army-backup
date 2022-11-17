package main

type Faction struct {
	Version string   `json:"version"`
	Units   []Unit   `json:"units"`
	Filters []Filter `json:"filters"`
	Resume  []struct {
		ID       int    `json:"id"`
		Isc      string `json:"isc"`
		IDArmy   int    `json:"idArmy"`
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Logo     string `json:"logo"`
		Type     int    `json:"type"`
		Category int    `json:"category"`
	} `json:"resume"`
	Fireteams []interface{} `json:"fireteams"`
	Relations []struct {
		Units []struct {
			Unit    int `json:"unit"`
			Profile int `json:"profile"`
		} `json:"units"`
		Min   int  `json:"min"`
		Max   int  `json:"max"`
		Group bool `json:"group"`
	} `json:"relations"`
	Specops struct {
		Equip []struct {
			Exp   int `json:"exp"`
			ID    int `json:"id"`
			Equip []struct {
				ID int `json:"id"`
			} `json:"equip,omitempty"`
		} `json:"equip"`
		Skills []struct {
			Exp     int   `json:"exp"`
			ID      int   `json:"id"`
			Extras  []int `json:"extras,omitempty"`
			Weapons []struct {
				ID int `json:"id"`
			} `json:"weapons,omitempty"`
			Equip []struct {
				ID int `json:"id"`
			} `json:"equip,omitempty"`
		} `json:"skills"`
		Weapons []struct {
			Exp int `json:"exp"`
			ID  int `json:"id"`
		} `json:"weapons"`
		Units []SpecopsUnit `json:"units"`
	} `json:"specops"`
	FireteamChart struct {
		Spec struct {
			CORE  int `json:"CORE"`
			HARIS int `json:"HARIS"`
			DUO   int `json:"DUO"`
		} `json:"spec"`
		Desc  string `json:"desc"`
		Teams []struct {
			Name  string   `json:"name"`
			Obs   string   `json:"obs"`
			Type  []string `json:"type"`
			Units []struct {
				Min      int    `json:"min"`
				Max      int    `json:"max"`
				Name     string `json:"name"`
				Comment  string `json:"comment"`
				Required bool   `json:"required"`
				Slug     string `json:"slug"`
			} `json:"units"`
		} `json:"teams"`
	} `json:"fireteamChart"`
}

type Unit struct {
	ID            int         `json:"id"`
	IDArmy        int         `json:"idArmy"`
	Canonical     int         `json:"canonical"`
	Isc           string      `json:"isc"`
	IscAbbr       interface{} `json:"iscAbbr"`
	Name          string      `json:"name"`
	ProfileGroups []struct {
		ID       int         `json:"id"`
		Category int         `json:"category"`
		Isc      string      `json:"isc"`
		Notes    interface{} `json:"notes"`
		Profiles []struct {
			ID       int           `json:"id"`
			Arm      int           `json:"arm"`
			Ava      int           `json:"ava"`
			Bs       int           `json:"bs"`
			Bts      int           `json:"bts"`
			Cc       int           `json:"cc"`
			Chars    []int         `json:"chars"`
			Equip    []interface{} `json:"equip"`
			Logo     string        `json:"logo"`
			Weapons  []interface{} `json:"weapons"`
			Includes []interface{} `json:"includes"`
			Move     []int         `json:"move"`
			Ph       int           `json:"ph"`
			S        int           `json:"s"`
			Str      bool          `json:"str"`
			Type     int           `json:"type"`
			W        int           `json:"w"`
			Wip      int           `json:"wip"`
			Name     string        `json:"name"`
			Notes    interface{}   `json:"notes"`
			Skills   []struct {
				ID    int   `json:"id"`
				Order int   `json:"order"`
				Extra []int `json:"extra,omitempty"`
				Q     int   `json:"q,omitempty"`
			} `json:"skills"`
			Peripheral []interface{} `json:"peripheral"`
		} `json:"profiles"`
		Options []struct {
			ID       int           `json:"id"`
			Chars    []interface{} `json:"chars"`
			Disabled bool          `json:"disabled"`
			Equip    []interface{} `json:"equip"`
			Minis    int           `json:"minis"`
			Orders   []struct {
				Type  string `json:"type"`
				List  int    `json:"list"`
				Total int    `json:"total"`
			} `json:"orders"`
			Includes []interface{} `json:"includes"`
			Points   int           `json:"points"`
			Swc      string        `json:"swc"`
			Weapons  []struct {
				ID    int `json:"id"`
				Order int `json:"order"`
			} `json:"weapons"`
			Name       string        `json:"name"`
			Skills     []interface{} `json:"skills"`
			Peripheral []interface{} `json:"peripheral"`
		} `json:"options"`
	} `json:"profileGroups"`
	Options []interface{} `json:"options"`
	Slug    string        `json:"slug"`
	Filters struct {
		Categories []int         `json:"categories"`
		Skills     []int         `json:"skills"`
		Equip      []interface{} `json:"equip"`
		Chars      []int         `json:"chars"`
		Types      []int         `json:"types"`
		Weapons    []int         `json:"weapons"`
		Ammunition []int         `json:"ammunition"`
	} `json:"filters"`
	Factions []int  `json:"factions"`
	Notes    string `json:"notes,omitempty"`
}

type SpecopsUnit struct {
	ID            int         `json:"id"`
	IDArmy        int         `json:"idArmy"`
	Canonical     int         `json:"canonical"`
	Isc           string      `json:"isc"`
	IscAbbr       interface{} `json:"iscAbbr"`
	Notes         interface{} `json:"notes"`
	Name          string      `json:"name"`
	ProfileGroups []struct {
		Notes    interface{} `json:"notes"`
		Isc      string      `json:"isc"`
		Profiles []struct {
			Bts      int           `json:"bts"`
			Cc       int           `json:"cc"`
			Move     []int         `json:"move"`
			Notes    interface{}   `json:"notes"`
			Includes []interface{} `json:"includes"`
			Type     int           `json:"type"`
			Ava      int           `json:"ava"`
			Str      bool          `json:"str"`
			Bs       int           `json:"bs"`
			S        int           `json:"s"`
			Equip    []interface{} `json:"equip"`
			W        int           `json:"w"`
			Ph       int           `json:"ph"`
			Name     string        `json:"name"`
			Logo     string        `json:"logo"`
			ID       int           `json:"id"`
			Arm      int           `json:"arm"`
			Weapons  []interface{} `json:"weapons"`
			Chars    []int         `json:"chars"`
			Wip      int           `json:"wip"`
			Skills   []struct {
				Q     int   `json:"q"`
				Extra []int `json:"extra"`
				ID    int   `json:"id"`
				Order int   `json:"order"`
			} `json:"skills"`
			Peripheral []interface{} `json:"peripheral"`
		} `json:"profiles"`
		Options []struct {
			Includes []interface{} `json:"includes"`
			Minis    int           `json:"minis"`
			Points   int           `json:"points"`
			Equip    []interface{} `json:"equip"`
			Name     string        `json:"name"`
			Disabled bool          `json:"disabled"`
			Orders   []struct {
				Type  string `json:"type"`
				List  int    `json:"list"`
				Total int    `json:"total"`
			} `json:"orders"`
			ID      int `json:"id"`
			Weapons []struct {
				ID    int `json:"id"`
				Order int `json:"order"`
			} `json:"weapons"`
			Chars  []interface{} `json:"chars"`
			Swc    string        `json:"swc"`
			Skills []struct {
				Q     int `json:"q"`
				ID    int `json:"id"`
				Order int `json:"order"`
			} `json:"skills"`
			Peripheral []interface{} `json:"peripheral"`
		} `json:"options"`
		ID       int `json:"id"`
		Category int `json:"category"`
	} `json:"profileGroups"`
	Options []interface{} `json:"options"`
	Slug    string        `json:"slug"`
	Filters struct {
		Categories []int         `json:"categories"`
		Skills     []int         `json:"skills"`
		Equip      []interface{} `json:"equip"`
		Chars      []int         `json:"chars"`
		Types      []int         `json:"types"`
		Weapons    []int         `json:"weapons"`
		Ammunition []int         `json:"ammunition"`
	} `json:"filters"`
	Factions []interface{} `json:"factions"`
}

type Filter struct {
	Peripheral []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Mercs bool   `json:"mercs,omitempty"`
	} `json:"peripheral"`
	Attrs []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Abbr string `json:"abbr"`
		Min  int    `json:"min"`
		Max  int    `json:"max"`
	} `json:"attrs"`
	Points []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"points"`
	Swc      []interface{} `json:"swc"`
	Category []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Ammunition []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Wiki  string `json:"wiki,omitempty"`
		Mercs bool   `json:"mercs,omitempty"`
	} `json:"ammunition"`
	Chars []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Wiki string `json:"wiki,omitempty"`
	} `json:"chars"`
	Type []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	Equip []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Wiki  string `json:"wiki"`
		Mercs bool   `json:"mercs,omitempty"`
	} `json:"equip"`
	Skills []struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Wiki    string `json:"wiki,omitempty"`
		Mercs   bool   `json:"mercs,omitempty"`
		Specops bool   `json:"specops,omitempty"`
	} `json:"skills"`
	Weapons []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Mercs bool   `json:"mercs,omitempty"`
		Wiki  string `json:"wiki,omitempty"`
	} `json:"weapons"`
	Extras []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Mercs bool   `json:"mercs,omitempty"`
	} `json:"extras"`
}
