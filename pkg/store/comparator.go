package store

import (
	"github.com/jedib0t/go-pretty/table"
)

type compareResult struct {
	foundOther bool
	equalOther bool
	info       FileInfo
}

type diff struct {
	expected FileInfo
	got      FileInfo
}

func (d diff) Rows() (rows []table.Row) {
	fields := d.got.Fields()
	path := d.got.RelativePath
	if path == "" {
		path = d.expected.RelativePath
	}
	for _, name := range fields {
		if d.expected.ValueOf(name) != d.got.ValueOf(name) {
			rows = append(rows, table.Row{path, name, d.expected.ValueOf(name), d.got.ValueOf(name)})
		}
	}
	return rows
}

type Result struct {
	missingInNew []diff
	missingInOld []diff
	notEqual     []diff
}

func (r Result) render(csv bool) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Path", "Attribute", "Expected", "Got"})
	for _, d := range r.notEqual {
		t.AppendRows(d.Rows())
	}
	for _, d := range r.missingInNew {
		t.AppendRow(table.Row{d.expected.RelativePath, "present", true, false})
	}
	for _, d := range r.missingInOld {
		t.AppendRow(table.Row{d.expected.RelativePath, "present", false, true})
	}
	if csv {
		return t.RenderCSV()
	}
	return t.Render()

}

func (r Result) Render(csv bool) string {
	return r.render(csv)
}

func ComputeDiffs(old, new []FileInfo) Result {
	oldMap := make(map[string]compareResult)
	newMap := make(map[string]compareResult)
	r := Result{}
	for _, finfo := range old {
		oldMap[finfo.RelativePath] = compareResult{
			foundOther: false,
			equalOther: false,
			info:       finfo,
		}
	}
	for _, finfo := range new {
		oldItem, found := oldMap[finfo.RelativePath]
		newMap[finfo.RelativePath] = compareResult{
			foundOther: found,
			equalOther: finfo.Equals(&oldItem.info),
			info:       finfo,
		}
		if !found {
			r.missingInOld = append(r.missingInOld, diff{
				expected: finfo,
			})
		}
	}
	for path, value := range oldMap {
		newItem, found := newMap[path]
		oldMap[path] = compareResult{
			foundOther: found,
			equalOther: value.info.Equals(&newItem.info),
			info:       value.info,
		}
		if !found {
			r.missingInNew = append(r.missingInNew, diff{
				expected: value.info,
			})
		}
	}
	for _, value := range oldMap {
		if !value.equalOther && value.foundOther {
			r.notEqual = append(r.notEqual, diff{
				expected: value.info,
				got:      newMap[value.info.RelativePath].info,
			})
		}
	}
	return r
}
