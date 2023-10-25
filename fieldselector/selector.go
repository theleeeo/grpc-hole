package fieldselector

import "fmt"

const (
	Delimiter           = "."
	Root      Selection = Delimiter
)

type Selection string

func (s Selection) String() string {
	return string(s)
}

func (s Selection) IsRoot() bool {
	return s == Root
}

func (s Selection) IsEmpty() bool {
	return s == ""
}

func (s Selection) AppendField(field string) Selection {
	if s.IsRoot() {
		return Selection(s.String() + field)
	}
	return Selection(s.String() + Delimiter + field)
}

func (s Selection) GetField() (string, error) {
	var field string
	_, err := fmt.Sscanf(s.String(), "%*s.%s", &field)
	if err != nil {
		return "", err
	}
	return field, nil
}

func (s Selection) SetIndex(index int) Selection {
	return Selection(fmt.Sprintf("%s[%d]", s.String(), index))
}

func (s Selection) SetKey(index string) Selection {
	return Selection(fmt.Sprintf("%s[%s]", s.String(), index))
}

func (s Selection) GetIndex() (int, error) {
	var index int
	_, err := fmt.Sscanf(s.String(), "%*s[%d]", &index)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func (s Selection) GetKey() (string, error) {
	var key string
	_, err := fmt.Sscanf(s.String(), "%*s[%s]", &key)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s Selection) GetParent() Selection {
	if s.IsRoot() {
		return s
	}
	var parent string
	_, err := fmt.Sscanf(s.String(), "%s.%*s", &parent)
	if err != nil {
		return s
	}
	return Selection(parent)
}
