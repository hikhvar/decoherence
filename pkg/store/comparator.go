package store

type compareResult struct {
	foundOther bool
	equalOther bool
	info       FileInfo
}

type diff struct {
	expected FileInfo
	got      FileInfo
}

type Result struct {
	missingInNew []diff
	missingInOld []diff
	notEqual     []diff
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
