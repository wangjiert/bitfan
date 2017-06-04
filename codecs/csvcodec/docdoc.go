// Code generated by "bitfanDoc -codec csvDecoder"; DO NOT EDIT
package csvcodec

import "github.com/vjeantet/bitfan/processors/doc"

func (p *csvDecoder) Doc() *doc.Codec {
	return &doc.Codec{
  Name:       "csvcodec",
  ImportPath: "/Users/sodadi/go/src/github.com/vjeantet/bitfan/codecs/csvcodec",
  Doc:        "Parses comma-separated value data into individual fields",
  DocShort:   "",
  Options:    &doc.CodecOptions{
    Doc:     "Parses comma-separated value data into individual fields",
    Options: []*doc.CodecOption{
      &doc.CodecOption{
        Name:           "Separator",
        Alias:          "",
        Doc:            "Define the column separator value. If this is not specified, the default is a comma ,. Optional",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\",\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.CodecOption{
        Name:           "AutogenerateColumnNames",
        Alias:          "autogenerate_column_names",
        Doc:            "Define whether column names should autogenerated or not. Defaults to true.\nIf set to false, columns not having a header specified will not be parsed.",
        Required:       false,
        Type:           "bool",
        DefaultValue:   "true",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.CodecOption{
        Name:           "QuoteChar",
        Alias:          "quote_char",
        Doc:            "Define the character used to quote CSV fields. If this is not specified the default is a double quote \". Optional.",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"\\\"\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.CodecOption{
        Name:           "Columns",
        Alias:          "",
        Doc:            "Define a list of column names (in the order they appear in the CSV, as\nif it were a header line).\n\nIf columns is not configured, or there are not enough columns specified,\nthe default column names are \"column1\", \"column2\", etc.\n\nIn the case that there are more columns in the data than specified in this column\nlist, extra columns will be auto-numbered:\n(e.g. \"user_defined_1\", \"user_defined_2\", \"column3\", \"column4\", etc.)",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
    },
  },
}
}