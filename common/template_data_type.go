package common

type TemplateDataType interface {
	AddDataType(string, string)
	DelDataType(string)
	GetDataType(string) string
	NewTemplate(string) Template
}

type TemplateDataTypeImpl struct {
	TemplateDataType
	DataTypeList map[string]string
}

func NewTemplateDataType() TemplateDataType {
	return &TemplateDataTypeImpl{
		DataTypeList: make(map[string]string),
	}
}

func (this *TemplateDataTypeImpl) AddDataType(name string, pattern string) {
	this.DataTypeList[name] = pattern
}

func (this *TemplateDataTypeImpl) GetDataType(name string) string {
	return this.DataTypeList[name]
}

func (this *TemplateDataTypeImpl) DelDataType(name string) {
	delete(this.DataTypeList, name)
}

func (this *TemplateDataTypeImpl) NewTemplate(pattern string) Template {
	tpl := &TemplateImpl{
		RawString:    pattern,
		DataTypeList: this,
		Args:         []string{},
	}
	tpl.parse()

	return tpl
}
