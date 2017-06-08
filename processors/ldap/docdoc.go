// Code generated by "bitfanDoc "; DO NOT EDIT
package ldapprocessor

import "github.com/vjeantet/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
  Name:       "ldapprocessor",
  ImportPath: "/Users/sodadi/go/src/github.com/vjeantet/bitfan/processors/ldap",
  Doc:        "Performs a search for a specified filter on the directory and fire events with results",
  DocShort:   "",
  Options:    &doc.ProcessorOptions{
    Doc:     "",
    Options: []*doc.ProcessorOption{
      &doc.ProcessorOption{
        Name:           "Add_field",
        Alias:          "",
        Doc:            "If this filter is successful, add any arbitrary fields to this event.",
        Required:       false,
        Type:           "hash",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Tags",
        Alias:          "",
        Doc:            "If this filter is successful, add arbitrary tags to the event. Tags can be dynamic\nand include parts of the event using the %{field} syntax.",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Type",
        Alias:          "",
        Doc:            "Add a type field to all events handled by this input",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Host",
        Alias:          "host",
        Doc:            "ldap hostname",
        Required:       true,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "host => \"ldap.forumsys.com\"",
      },
      &doc.ProcessorOption{
        Name:           "Port",
        Alias:          "port",
        Doc:            "ldap port",
        Required:       true,
        Type:           "int",
        DefaultValue:   "389",
        PossibleValues: []string{},
        ExampleLS:      "port => 389",
      },
      &doc.ProcessorOption{
        Name:           "BindDn",
        Alias:          "bind_dn",
        Doc:            "Bind dn",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "bind_dn => \"cn=read-only-admin,dc=example,dc=com\"",
      },
      &doc.ProcessorOption{
        Name:           "BindPassword",
        Alias:          "bind_password",
        Doc:            "Bind password",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "bind_password => \"password\"",
      },
      &doc.ProcessorOption{
        Name:           "BaseDn",
        Alias:          "base_dn",
        Doc:            "Base DN\nIf bind_dn is not specified or is empty, an anonymous bind is attempted.\nThis is defined in https://tools.ietf.org/html/rfc2251#section-4.2.2",
        Required:       true,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "base_dn => \"dc=example,dc=com\"",
      },
      &doc.ProcessorOption{
        Name:           "SearchBase",
        Alias:          "search_base",
        Doc:            "A search base (the distinguished name of the search base object) defines the\nlocation in the directory from which the LDAP search begins.",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "SearchFilter",
        Alias:          "search_filter",
        Doc:            "The search filter can be simple or advanced, using boolean operators in the format\ndescribed in the LDAP documentation (see [RFC4515](http://www.faqs.org/rfcs/rfc4515) for full information on filters).",
        Required:       true,
        Type:           "string",
        DefaultValue:   "\"(objectClass=*)\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "SearchAttributes",
        Alias:          "search_attributes",
        Doc:            "An array of the required attributes, e.g. [\"mail\", \"sn\", \"cn\"].\n\nNote that the \"dn\" is always returned irrespective of which attributes types are requested.\n\nUsing this parameter is much more efficient than the default action (which is to return all attributes and their associated values).\n\nThe use of this parameter should therefore be considered good practice.",
        Required:       false,
        Type:           "array",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "search_attributes => [\"mail\", \"sn\", \"cn\"]",
      },
      &doc.ProcessorOption{
        Name:           "SearchScope",
        Alias:          "search_scope",
        Doc:            "The SCOPE setting is the starting point of an LDAP search and the depth from the\nbase DN to which the search should occur.\n\nThere are three options (values) that can be assigned to the SCOPE parameter:\n\n* **base** : indicate searching only the entry at the base DN, resulting in only that entry being returned\n* **one** : indicate searching all entries one level under the base DN - but not including the base DN and not including any entries under that one level under the base DN.\n* **subtree** : indicate searching of all entries at all levels under and including the specified base DN\n\n![scope](../ldapscope.gif)",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"subtree\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "SizeLimit",
        Alias:          "size_limit",
        Doc:            "Maximum entries to return (leave empty to let the server decide)",
        Required:       false,
        Type:           "int",
        DefaultValue:   "0",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "PagingSize",
        Alias:          "paging_size",
        Doc:            "Desired page size in order to execute LDAP queries to fulfill the\nsearch request.\n\nSet 0 to not use Paging",
        Required:       false,
        Type:           "int",
        DefaultValue:   "1000",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "EventBy",
        Alias:          "event_by",
        Doc:            "Send an event row by row or one event with all results\npossible values \"entry\", \"result\"",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"entry\"",
        PossibleValues: []string{},
        ExampleLS:      "",
      },
      &doc.ProcessorOption{
        Name:           "Interval",
        Alias:          "interval",
        Doc:            "Set an interval when this processor is used as a input",
        Required:       false,
        Type:           "string",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "interval => \"10\"",
      },
      &doc.ProcessorOption{
        Name:           "Var",
        Alias:          "var",
        Doc:            "You can set variable to be used in Search Query by using ${var}.\neach reference will be replaced by the value of the variable found in search query content\nThe replacement is case-sensitive.",
        Required:       false,
        Type:           "hash",
        DefaultValue:   nil,
        PossibleValues: []string{},
        ExampleLS:      "var => {\"hostname\"=>\"myhost\",\"varname\"=>\"varvalue\"}",
      },
      &doc.ProcessorOption{
        Name:           "Target",
        Alias:          "target",
        Doc:            "Define the target field for placing the retrieved data. If this setting is omitted,\nthe data will be stored in the \"data\" field\nSet the value to \".\" to store value to the root (top level) of the event",
        Required:       false,
        Type:           "string",
        DefaultValue:   "\"data\"",
        PossibleValues: []string{},
        ExampleLS:      "target => \"data\"",
      },
    },
  },
  Ports: []*doc.ProcessorPort{},
}
}