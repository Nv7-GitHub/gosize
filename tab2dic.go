package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func store_rec(values *map[string]interface{}, path []string, sz int) {
	if len(path) == 0 {
		var v int

		vraw, exists := (*values)["value"]
		if !exists {
			v = 0
		} else {
			v = vraw.(int)
		}
		v += sz

		(*values)["value"] = v
	} else {
		n := string(path[0])

		var d1 map[string]interface{}
		v, exists := (*values)["children"]
		if !exists {
			d1 = make(map[string]interface{})
		} else {
			d1 = v.(map[string]interface{})
		}

		var d map[string]interface{}
		v, exists = d1["children"]
		if !exists {
			d = make(map[string]interface{})
		} else {
			d = v.(map[string]interface{})
		}

		d1[n] = d
		(*values)["children"] = d
		store_rec(values, path[1:], sz)
	}
}

func tab2dic(sizes string) map[string]interface{} {
	values := make(map[string]interface{})
	vtables_sz := 0
	typeinfo_sz := 0
	init_sz := 0
	rest_sz := 0
	gotyp_sz := 0

	for _, line := range strings.Split(sizes, "\n") {
		if undefre.Match([]byte(line)) {
			continue
		}

		if line[len(line)-1] == '\n' {
			line = line[:len(line)-2]
		}

		matches := entriesre.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 {
			fmt.Println("unknown format:", line)
			continue
		}

		typ := strings.TrimSpace(matches[0][3])
		if typ == "U" {
			continue
		}

		szraw := strings.TrimSpace(matches[0][2])
		sz, err := strconv.Atoi(szraw)
		if err != nil {
			fmt.Println("unknown size format:", szraw)
			continue
		}
		if sz == 0 {
			continue
		}

		sym := strings.TrimSpace(matches[0][4])
		if sym == "" {
			continue
		}

		if strings.HasPrefix(sym, "construction vtable ") || strings.HasPrefix(sym, "vtable for") {
			vtables_sz += sz
			continue
		}

		if strings.HasPrefix(sym, "__static_initialization_and_destruction") {
			init_sz += sz
			continue
		}

		if strings.HasPrefix(sym, "typeinfo ") {
			typeinfo_sz += sz
			continue
		}

		if strings.HasPrefix(sym, "type..") {
			gotyp_sz += sz
			continue
		}

		parts := cppsymre.FindAllStringSubmatch(sym, -1)

		var prefix []string
		var partsre *regexp.Regexp

		if len(parts) > 0 {
			prefix = []string{"c/c++ · "}
			partsre = cpppathre
		} else {
			parts = gosymre.FindAllStringSubmatch(sym, -1)
			if len(parts) > 0 {
				prefix = []string{"go · "}
				partsre = gopathpartsre
			}
		}

		if len(parts) == 0 {
			fmt.Println("unknown", typ, "sym format:", sym)
			rest_sz += sz
			continue
		}

		path := parts[0][2]
		path = partsre.FindAllString(path, -1)[0]
		name := parts[0][1] + parts[0][3]
		fullpath := append(prefix, path, name)
		store_rec(&values, fullpath, sz)
	}
	store_rec(&values, []string{"c/c++ · ", "VTABLES"}, vtables_sz)
	store_rec(&values, []string{"c/c++ · ", "TYPEDATA"}, typeinfo_sz)
	store_rec(&values, []string{"c/c++ · ", "INITIALIZERS"}, init_sz)
	store_rec(&values, []string{"go · ", "TYPEDATA"}, gotyp_sz)
	store_rec(&values, []string{"UNKNOWN"}, rest_sz)

	return values
}
