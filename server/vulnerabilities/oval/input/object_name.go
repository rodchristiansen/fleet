package oval_input

type objectNameXML struct {
	VarRef   string `xml:"var_ref,attr"`
	VarCheck string `xml:"var_check,attr"`
	Value    string `xml:",chardata"`
}
