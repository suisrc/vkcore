package corev

import (
	"reflect"
)

// 通过字段名获取标签列表
func ForFieldNameToTag(s interface{}, tag string, fields ...string) []string {
	tagMap := ForStructToTag(s, true, tag)
	tagArr := []string{}
	for _, fname := range fields {
		if tags, ok := tagMap[fname]; ok {
			tag := tags.Find(tag)
			if tag.Ignored() || tag.Empty() {
				continue
			}
			tagArr = append(tagArr, tag.Value)
		}
	}
	return tagArr
}

// 通过字段名或者标签+值
func ForFieldNameToTagValue(s interface{}, tag string, fields ...string) map[string]interface{} {
	tagMap := ForStructToTag(s, true, tag)
	tagVal := map[string]interface{}{}

	rValue := ForStructToValue(s)
	for _, fname := range fields {
		if tags, ok := tagMap[fname]; ok {
			tag := tags.Find(tag)
			if tag.Ignored() || tag.Empty() {
				continue
			}
			tagVal[tag.Value] = rValue.Field(tag.Index).Interface()
		}
	}
	return tagVal
}

// 通过字段名或者标签+值
func ForFieldNameToTagValueAll(s interface{}, tag string) map[string]interface{} {
	tagMap := ForStructToTag(s, true, tag)
	tagVal := map[string]interface{}{}

	rValue := ForStructToValue(s)
	for _, tags := range tagMap {
		tag := tags.Find(tag)
		if tag.Ignored() || tag.Empty() {
			continue
		}
		tagVal[tag.Value] = rValue.Field(tag.Index).Interface()
	}
	return tagVal
}

// 获取结构体实体
func ForStructToValue(s interface{}) reflect.Value {
	sv := reflect.ValueOf(s)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() == reflect.Slice {
		sv = sv.Elem()
		if sv.Kind() == reflect.Ptr {
			sv = sv.Elem()
		}
	}
	return sv
}

// ForStructToTag 返回结构体的标签
func ForStructToTag(s interface{}, mst bool, tags ...string) map[string]Tags {
	fields := make(map[string]Tags)

	st := reflect.TypeOf(s)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() == reflect.Slice {
		st = st.Elem()
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
	}

	cnt := st.NumField()
	for i := 0; i < cnt; i++ {
		field := st.Field(i)
		fields[field.Name] = TagsFor(field, i, mst, tags...)
	}
	return fields
}

//======================================================================
// var tags = "db"
// strings.Fields(tags)

type Tag struct {
	Index int
	Value string
	Name  string
}

func (t Tag) Empty() bool {
	return t.Value == ""
}

func (t Tag) Ignored() bool {
	return t.Value == "-"
}

type Tags []Tag

func (t Tags) Find(name string) Tag {
	for _, pTag := range t {
		if pTag.Name == name {
			return pTag
		}
	}
	return Tag{}
}

// if must = true, 保留第一标签使用属性名给出
func TagsFor(field reflect.StructField, idx int, mst bool, tags ...string) Tags {
	pTags := Tags{}
	for _, tag := range tags {
		if valTag := field.Tag.Get(tag); valTag != "" {
			pTags = append(pTags, Tag{idx, valTag, tag})
		}
	}

	if len(pTags) == 0 && mst {
		pTags = append(pTags, Tag{idx, field.Name, tags[0]})
	}
	return pTags
}
