package ajson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var (
	jsonExample = []byte(`{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`)
)

type storeExample struct {
	Store struct {
		Book []struct {
			Category string  `json:"category"`
			Author   string  `json:"author"`
			Title    string  `json:"title"`
			Price    float64 `json:"price"`
			Isbn     string  `json:"isbn,omitempty"`
		} `json:"book"`
		Bicycle struct {
			Color string  `json:"color"`
			Price float64 `json:"price"`
		} `json:"bicycle"`
	} `json:"store"`
}

type testCase struct {
	name  string
	input []byte
	_type NodeType
	value []byte
}

func simpleCorrupted(name string) *testCase {
	return &testCase{name: name, input: []byte(name)}
}

func simpleValid(test *testCase, t *testing.T) {
	root, err := Unmarshal(test.input)
	if err != nil {
		t.Errorf("Error on Unmarshal(%s): %s", test.name, err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal(%s): root is nil", test.name)
	} else if root.Type() != test._type {
		t.Errorf("Error on Unmarshal(%s): wrong type", test.name)
	} else if !bytes.Equal(root.Source(), test.value) {
		t.Errorf("Error on Unmarshal(%s): %s != %s", test.name, root.Source(), test.value)
	}
}

func simpleInvalid(test *testCase, t *testing.T) {
	root, err := Unmarshal(test.input)
	if err == nil {
		t.Errorf("Error on Unmarshal(%s): error expected, got '%s'", test.name, root.Source())
	} else if root != nil {
		t.Errorf("Error on Unmarshal(%s): root is not nil", test.name)
	}
}

func TestUnmarshal_NumericSimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "1", input: []byte("1"), _type: Numeric, value: []byte("1")},
		{name: "-1", input: []byte("-1"), _type: Numeric, value: []byte("-1")},

		{name: "1234567890", input: []byte("1234567890"), _type: Numeric, value: []byte("1234567890")},
		{name: "-123", input: []byte("-123"), _type: Numeric, value: []byte("-123")},

		{name: "123.456", input: []byte("123.456"), _type: Numeric, value: []byte("123.456")},
		{name: "-123.456", input: []byte("-123.456"), _type: Numeric, value: []byte("-123.456")},

		{name: "1e3", input: []byte("1e3"), _type: Numeric, value: []byte("1e3")},
		{name: "1e+3", input: []byte("1e+3"), _type: Numeric, value: []byte("1e+3")},
		{name: "1e-3", input: []byte("1e-3"), _type: Numeric, value: []byte("1e-3")},
		{name: "-1e3", input: []byte("-1e3"), _type: Numeric, value: []byte("-1e3")},
		{name: "-1e-3", input: []byte("-1e-3"), _type: Numeric, value: []byte("-1e-3")},

		{name: "1.123e3456", input: []byte("1.123e3456"), _type: Numeric, value: []byte("1.123e3456")},
		{name: "1.123e-3456", input: []byte("1.123e-3456"), _type: Numeric, value: []byte("1.123e-3456")},
		{name: "-1.123e3456", input: []byte("-1.123e3456"), _type: Numeric, value: []byte("-1.123e3456")},
		{name: "-1.123e-3456", input: []byte("-1.123e-3456"), _type: Numeric, value: []byte("-1.123e-3456")},

		{name: "1E3", input: []byte("1E3"), _type: Numeric, value: []byte("1E3")},
		{name: "1E-3", input: []byte("1E-3"), _type: Numeric, value: []byte("1E-3")},
		{name: "-1E3", input: []byte("-1E3"), _type: Numeric, value: []byte("-1E3")},
		{name: "-1E-3", input: []byte("-1E-3"), _type: Numeric, value: []byte("-1E-3")},

		{name: "1.123E3456", input: []byte("1.123E3456"), _type: Numeric, value: []byte("1.123E3456")},
		{name: "1.123E-3456", input: []byte("1.123E-3456"), _type: Numeric, value: []byte("1.123E-3456")},
		{name: "-1.123E3456", input: []byte("-1.123E3456"), _type: Numeric, value: []byte("-1.123E3456")},
		{name: "-1.123E-3456", input: []byte("-1.123E-3456"), _type: Numeric, value: []byte("-1.123E-3456")},

		{name: "-1.123E-3456 with spaces", input: []byte(" \r -1.123E-3456 \t\n"), _type: Numeric, value: []byte("-1.123E-3456")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_NumericSimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		simpleCorrupted("+1"),
		simpleCorrupted("+1.1"),
		simpleCorrupted("+1e1"),
		simpleCorrupted("+1E1"),
		simpleCorrupted("x1"),
		simpleCorrupted("1+1"),
		simpleCorrupted("-1+"),
		simpleCorrupted("."),
		simpleCorrupted("-"),
		simpleCorrupted("+"),
		simpleCorrupted("-."),
		simpleCorrupted("+."),
		simpleCorrupted("e"),
		simpleCorrupted("e+"),
		simpleCorrupted("e1"), // exp without base part
		simpleCorrupted("e+1-"),
		simpleCorrupted("1null"),
		simpleCorrupted("1.123e3.456"), // exp part must be integer type
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_StringSimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "blank", input: []byte("\"\""), _type: String, value: []byte("\"\"")},
		{name: "char", input: []byte("\"c\""), _type: String, value: []byte("\"c\"")},
		{name: "word", input: []byte("\"cat\""), _type: String, value: []byte("\"cat\"")},
		{name: "spaces", input: []byte("  \"good cat or dog\"\r\n "), _type: String, value: []byte("\"good cat or dog\"")},
		{name: "backslash", input: []byte("\"good \\\"cat\\\"\""), _type: String, value: []byte("\"good \\\"cat\\\"\"")},
		{name: "backslash 2", input: []byte("\"good \\\\\\\"cat\\\"\""), _type: String, value: []byte("\"good \\\\\\\"cat\\\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_StringSimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		{name: "one quote", input: []byte("\"")},
		{name: "white NL", input: []byte("\"foo\nbar\"")},
		{name: "white R", input: []byte("\"foo\rbar\"")},
		{name: "white Tab", input: []byte("\"foo\tbar\"")},
		{name: "one quote char", input: []byte("\"c")},
		{name: "wrong quotes", input: []byte("'cat'")},
		{name: "double string", input: []byte("\"Hello\" \"World\"")},
		{name: "quotes in quotes", input: []byte("\"good \"cat\"\"")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_NullSimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "lower", input: []byte("null"), _type: Null, value: []byte("null")},
		{name: "spaces", input: []byte("  null\r\n "), _type: Null, value: []byte("null")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_NullSimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		{name: "nul", input: []byte("nul")},
		{name: "NILL", input: []byte("NILL")},
		{name: "Null", input: []byte("Null")},
		{name: "NULL", input: []byte("NULL")},
		{name: "spaces", input: []byte("Nu ll")},
		{name: "null1", input: []byte("null1")},
		{name: "double", input: []byte("null null")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_BoolSimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "lower true", input: []byte("true"), _type: Bool, value: []byte("true")},
		{name: "lower false", input: []byte("false"), _type: Bool, value: []byte("false")},
		{name: "spaces true", input: []byte("  true\r\n "), _type: Bool, value: []byte("true")},
		{name: "spaces false", input: []byte("  false\r\n "), _type: Bool, value: []byte("false")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_BoolSimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		simpleCorrupted("tru"),
		simpleCorrupted("fals"),
		simpleCorrupted("tre"),
		simpleCorrupted("fal se"),
		simpleCorrupted("true false"),
		simpleCorrupted("True"),
		simpleCorrupted("TRUE"),
		simpleCorrupted("False"),
		simpleCorrupted("FALSE"),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_ArraySimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "[]", input: []byte("[]"), _type: Array, value: []byte("[]")},
		{name: "[1]", input: []byte("[1]"), _type: Array, value: []byte("[1]")},
		{name: "[1,2,3]", input: []byte("[1,2,3]"), _type: Array, value: []byte("[1,2,3]")},
		{name: "[1, 2, 3]", input: []byte("[1, 2, 3]"), _type: Array, value: []byte("[1, 2, 3]")},
		{name: "[1,[2],3]", input: []byte("[1,[2],3]"), _type: Array, value: []byte("[1,[2],3]")},
		{name: "[[],[],[]]", input: []byte("[[],[],[]]"), _type: Array, value: []byte("[[],[],[]]")},
		{name: "[[[[[]]]]]", input: []byte("[[[[[]]]]]"), _type: Array, value: []byte("[[[[[]]]]]")},
		{name: "[true,null,1,\"foo\",[]]", input: []byte("[true,null,1,\"foo\",[]]"), _type: Array, value: []byte("[true,null,1,\"foo\",[]]")},
		{name: "spaces", input: []byte("\n\r [\n1\n ]\r\n"), _type: Array, value: []byte("[\n1\n ]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_ArraySimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		simpleCorrupted("[,]"),
		simpleCorrupted("[]\\"),
		simpleCorrupted("[1,]"),
		simpleCorrupted("[[]"),
		simpleCorrupted("[]]"),
		simpleCorrupted("1[]"),
		simpleCorrupted("[]1"),
		simpleCorrupted("[[]1]"),
		simpleCorrupted("‌[],[]"),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_ObjectSimpleSuccess(t *testing.T) {
	tests := []*testCase{
		{name: "{}", input: []byte("{}"), _type: Object, value: []byte("{}")},
		{name: `{ \r\n }`, input: []byte("{ \r\n }"), _type: Object, value: []byte("{ \r\n }")},
		{name: `{"key":1}`, input: []byte(`{"key":1}`), _type: Object, value: []byte(`{"key":1}`)},
		{name: `{"key":true}`, input: []byte(`{"key":true}`), _type: Object, value: []byte(`{"key":true}`)},
		{name: `{"key":"value"}`, input: []byte(`{"key":"value"}`), _type: Object, value: []byte(`{"key":"value"}`)},
		{name: `{"foo":"bar","baz":"foo"}`, input: []byte(`{"foo":"bar", "baz":"foo"}`), _type: Object, value: []byte(`{"foo":"bar", "baz":"foo"}`)},
		{name: "spaces", input: []byte(`  {  "foo"  :  "bar"  , "baz"   :   "foo"   }    `), _type: Object, value: []byte(`{  "foo"  :  "bar"  , "baz"   :   "foo"   }`)},
		{name: "nested", input: []byte(`{"foo":{"bar":{"baz":{}}}}`), _type: Object, value: []byte(`{"foo":{"bar":{"baz":{}}}}`)},
		{name: "array", input: []byte(`{"array":[{},{},{"foo":[{"bar":["baz"]}]}]}`), _type: Object, value: []byte(`{"array":[{},{},{"foo":[{"bar":["baz"]}]}]}`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleValid(test, t)
		})
	}
}

func TestUnmarshal_ObjectSimpleCorrupted(t *testing.T) {
	tests := []*testCase{
		simpleCorrupted("{,}"),
		simpleCorrupted("{:}"),
		simpleCorrupted(`{"foo"}`),
		simpleCorrupted(`{"foo":}`),
		simpleCorrupted(`{:"foo"}`),
		simpleCorrupted(`{"foo":bar}`),
		simpleCorrupted(`{"foo":"bar",}`),
		simpleCorrupted(`{}{}`),
		simpleCorrupted(`{},{}`),
		simpleCorrupted(`{[},{]}`),
		simpleCorrupted(`{[,]}`),
		simpleCorrupted(`{[]}`),
		simpleCorrupted(`{}1`),
		simpleCorrupted(`1{}`),
		simpleCorrupted(`{"x"::1}`),
		simpleCorrupted(`{null:null}`),
		simpleCorrupted(`{"foo:"bar"}`),
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			simpleInvalid(test, t)
		})
	}
}

func TestUnmarshal_Array(t *testing.T) {
	root, err := Unmarshal([]byte(" [1,[\"1\",[1,[1,2,3]]]]\r\n"))
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal: root is nil")
	} else if root.Type() != Array {
		t.Errorf("Error on Unmarshal: wrong type")
	} else {
		array, err := root.GetArray()
		if err != nil {
			t.Errorf("Error on root.GetArray(): %s", err.Error())
		} else if len(array) != 2 {
			t.Errorf("Error on root.GetArray(): expected 2 elements")
		} else if val, err := array[0].GetNumeric(); err != nil {
			t.Errorf("Error on array[0].GetNumeric(): %s", err.Error())
		} else if val != 1 {
			t.Errorf("Error on array[0].GetNumeric(): expected to be '1'")
		} else if val, err := array[1].GetArray(); err != nil {
			t.Errorf("Error on array[1].GetArray(): %s", err.Error())
		} else if len(val) != 2 {
			t.Errorf("Error on array[1].GetArray(): expected 2 elements")
		} else if el, err := val[0].GetString(); err != nil {
			t.Errorf("Error on val[0].GetString(): %s", err.Error())
		} else if el != "1" {
			t.Errorf("Error on val[0].GetString(): expected to be '\"1\"'")
		}
	}
}

func TestUnmarshal_Object(t *testing.T) {
	root, err := Unmarshal([]byte(`{"foo":{"bar":[null]}, "baz":true}`))
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal: root is nil")
	} else if !root.IsObject() {
		t.Errorf("Error on Unmarshal: wrong type")
	} else {
		object, err := root.GetObject()
		if err != nil {
			t.Errorf("Error on root.GetObject(): %s", err.Error())
		} else if foo, ok := object["foo"]; !ok {
			t.Errorf("Error on getting foo from map")
		} else if !foo.IsObject() {
			t.Errorf("Child element type error [foo]")
		} else if obj, err := foo.GetObject(); err != nil {
			t.Errorf("Error on foo.GetObject(): %s", err.Error())
		} else if bar, ok := obj["bar"]; !ok {
			t.Errorf("Error on getting bar from map")
		} else if !bar.IsArray() {
			t.Errorf("Child element type error [bar]")
		} else if baz, ok := object["baz"]; !ok {
			t.Errorf("Error on getting baz from map")
		} else if !baz.IsBool() {
			t.Errorf("Child element type error [baz]")
		} else if val, err := baz.GetBool(); err != nil {
			t.Errorf("Error on baz.GetBool(): %s", err.Error())
		} else if !val {
			t.Errorf("Error on getting boolean")
		}
	}
}

func TestUnmarshalSafe(t *testing.T) {
	safe, err := UnmarshalSafe(jsonExample)
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if safe == nil {
		t.Errorf("Error on Unmarshal: safe is nil")
	} else {
		root, err := Unmarshal(jsonExample)
		if err != nil {
			t.Errorf("Error on Unmarshal: %s", err.Error())
		} else if root == nil {
			t.Errorf("Error on Unmarshal: root is nil")
		} else if !bytes.Equal(root.Source(), safe.Source()) {
			t.Errorf("Error on UnmarshalSafe: values not same")
		}
	}
}

func TestUnmarshal_Must(t *testing.T) {
	root, err := Unmarshal(jsonExample)
	if err != nil {
		t.Errorf("Error on Unmarshal: %s", err.Error())
	} else if root == nil {
		t.Errorf("Error on Unmarshal: root is nil")
	} else {
		category := root.MustObject()["store"].MustObject()["book"].MustArray()[2].MustObject()["category"].MustString()
		if category != "fiction" {
			t.Errorf("Error on Unmarshal: data corrupted")
		}
	}
}

// Examples from https://json.org/example.html
func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name: "glossary",
			value: `{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}`,
		},
		{
			name: "menu",
			value: `{"menu": {
  "id": "file",
  "value": "File",
  "popup": {
    "menuitem": [
      {"value": "New", "onclick": "CreateNewDoc()"},
      {"value": "Open", "onclick": "OpenDoc()"},
      {"value": "Close", "onclick": "CloseDoc()"}
    ]
  }
}}`,
		},
		{
			name: "widget",
			value: `{"widget": {
    "debug": "on",
    "window": {
        "title": "Sample Konfabulator Widget",
        "name": "main_window",
        "width": 500,
        "height": 500
    },
    "image": { 
        "src": "Images/Sun.png",
        "name": "sun1",
        "hOffset": 250,
        "vOffset": 250,
        "alignment": "center"
    },
    "text": {
        "data": "Click Here",
        "size": 36,
        "style": "bold",
        "name": "text1",
        "hOffset": 250,
        "vOffset": 100,
        "alignment": "center",
        "onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
    }
}}    `,
		},
		{
			name: "web-app",
			value: `{"web-app": {
  "servlet": [   
    {
      "servlet-name": "cofaxCDS",
      "servlet-class": "org.cofax.cds.CDSServlet",
      "init-param": {
        "configGlossary:installationAt": "Philadelphia, PA",
        "configGlossary:adminEmail": "ksm@pobox.com",
        "configGlossary:poweredBy": "Cofax",
        "configGlossary:poweredByIcon": "/images/cofax.gif",
        "configGlossary:staticPath": "/content/static",
        "templateProcessorClass": "org.cofax.WysiwygTemplate",
        "templateLoaderClass": "org.cofax.FilesTemplateLoader",
        "templatePath": "templates",
        "templateOverridePath": "",
        "defaultListTemplate": "listTemplate.htm",
        "defaultFileTemplate": "articleTemplate.htm",
        "useJSP": false,
        "jspListTemplate": "listTemplate.jsp",
        "jspFileTemplate": "articleTemplate.jsp",
        "cachePackageTagsTrack": 200,
        "cachePackageTagsStore": 200,
        "cachePackageTagsRefresh": 60,
        "cacheTemplatesTrack": 100,
        "cacheTemplatesStore": 50,
        "cacheTemplatesRefresh": 15,
        "cachePagesTrack": 200,
        "cachePagesStore": 100,
        "cachePagesRefresh": 10,
        "cachePagesDirtyRead": 10,
        "searchEngineListTemplate": "forSearchEnginesList.htm",
        "searchEngineFileTemplate": "forSearchEngines.htm",
        "searchEngineRobotsDb": "WEB-INF/robots.db",
        "useDataStore": true,
        "dataStoreClass": "org.cofax.SqlDataStore",
        "redirectionClass": "org.cofax.SqlRedirection",
        "dataStoreName": "cofax",
        "dataStoreDriver": "com.microsoft.jdbc.sqlserver.SQLServerDriver",
        "dataStoreUrl": "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
        "dataStoreUser": "sa",
        "dataStorePassword": "dataStoreTestQuery",
        "dataStoreTestQuery": "SET NOCOUNT ON;select test='test';",
        "dataStoreLogFile": "/usr/local/tomcat/logs/datastore.log",
        "dataStoreInitConns": 10,
        "dataStoreMaxConns": 100,
        "dataStoreConnUsageLimit": 100,
        "dataStoreLogLevel": "debug",
        "maxUrlLength": 500}},
    {
      "servlet-name": "cofaxEmail",
      "servlet-class": "org.cofax.cds.EmailServlet",
      "init-param": {
      "mailHost": "mail1",
      "mailHostOverride": "mail2"}},
    {
      "servlet-name": "cofaxAdmin",
      "servlet-class": "org.cofax.cds.AdminServlet"},
 
    {
      "servlet-name": "fileServlet",
      "servlet-class": "org.cofax.cds.FileServlet"},
    {
      "servlet-name": "cofaxTools",
      "servlet-class": "org.cofax.cms.CofaxToolsServlet",
      "init-param": {
        "templatePath": "toolstemplates/",
        "log": 1,
        "logLocation": "/usr/local/tomcat/logs/CofaxTools.log",
        "logMaxSize": "",
        "dataLog": 1,
        "dataLogLocation": "/usr/local/tomcat/logs/dataLog.log",
        "dataLogMaxSize": "",
        "removePageCache": "/content/admin/remove?cache=pages&id=",
        "removeTemplateCache": "/content/admin/remove?cache=templates&id=",
        "fileTransferFolder": "/usr/local/tomcat/webapps/content/fileTransferFolder",
        "lookInContext": 1,
        "adminGroupID": 4,
        "betaServer": true}}],
  "servlet-mapping": {
    "cofaxCDS": "/",
    "cofaxEmail": "/cofaxutil/aemail/*",
    "cofaxAdmin": "/admin/*",
    "fileServlet": "/static/*",
    "cofaxTools": "/tools/*"},
 
  "taglib": {
    "taglib-uri": "cofax.tld",
    "taglib-location": "/WEB-INF/tlds/cofax.tld"}}}`,
		},
		{
			name: "SVG Viewer",
			value: `{"menu": {
    "header": "SVG Viewer",
    "items": [
        {"id": "Open"},
        {"id": "OpenNew", "label": "Open New"},
        null,
        {"id": "ZoomIn", "label": "Zoom In"},
        {"id": "ZoomOut", "label": "Zoom Out"},
        {"id": "OriginalView", "label": "Original View"},
        null,
        {"id": "Quality"},
        {"id": "Pause"},
        {"id": "Mute"},
        null,
        {"id": "Find", "label": "Find..."},
        {"id": "FindAgain", "label": "Find Again"},
        {"id": "Copy"},
        {"id": "CopyAgain", "label": "Copy Again"},
        {"id": "CopySVG", "label": "Copy SVG"},
        {"id": "ViewSVG", "label": "View SVG"},
        {"id": "ViewSource", "label": "View Source"},
        {"id": "SaveAs", "label": "Save As"},
        null,
        {"id": "Help"},
        {"id": "About", "label": "About Adobe CVG Viewer..."}
    ]
}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Unmarshal([]byte(test.value))
			if err != nil {
				t.Errorf("Error on Unmarshal: %s", err.Error())
			}
		})
	}
}

func BenchmarkUnmarshal_AJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root, err := Unmarshal(jsonExample)
		if err != nil || root == nil {
			b.Errorf("Error on Unmarshal")
		}
	}
}

func BenchmarkUnmarshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		root := new(storeExample)
		err := json.Unmarshal(jsonExample, &root)
		if err != nil || root == nil {
			b.Errorf("Error on Unmarshal")
		}
	}
}

// Calculate AVG price from different types of objects, JSON from: https://goessner.net/articles/JsonPath/index.html#e3
func ExampleUnmarshal() {
	data := []byte(`{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": { "color": "red",
      "price": 19.95
    },
    "tools": null
  }
}`)

	root, err := Unmarshal(data)
	if err != nil {
		panic(err)
	}

	store := root.MustKey("store").MustObject()

	var prices float64
	size := 0
	for _, objects := range store {
		if objects.IsArray() && objects.Size() > 0 {
			size += objects.Size()
			for _, object := range objects.MustArray() {
				prices += object.MustKey("price").MustNumeric()
			}
		} else if objects.IsObject() && objects.HasKey("price") {
			size++
			prices += objects.MustKey("price").MustNumeric()
		}
	}

	if size > 0 {
		fmt.Println("AVG price:", prices/float64(size))
	} else {
		fmt.Println("AVG price:", 0)
	}
}

// Interface object interface and render html link. JSON from: https://tools.ietf.org/html/rfc7159#section-13
func ExampleUnmarshal_unpack() {
	data := []byte(`{
        "Image": {
            "Width":  800,
            "Height": 600,
            "Title":  "View from 15th Floor",
            "Thumbnail": {
                "Url":    "http://www.example.com/image/481989943",
                "Height": 125,
                "Width":  100
            },
            "Animated" : false,
            "IDs": [116, 943, 234, 38793]
          }
      }`)

	root, err := Unmarshal(data)
	if err != nil {
		panic(err)
	}
	object, err := root.Interface()
	if err != nil {
		panic(err)
	}
	image := object.(map[string]interface{})["Image"].(map[string]interface{})
	thumbnail := image["Thumbnail"].(map[string]interface{})
	fmt.Printf(
		`<a href="%s?width=%.0f&height=%.0f" title="%s"><img src="%s?width=%.0f&height=%.0f" /></a>`,
		thumbnail["Url"],
		image["Width"],
		image["Height"],
		image["Title"],
		thumbnail["Url"],
		thumbnail["Width"],
		thumbnail["Height"],
	)
}

func ExampleMust() {
	data := []byte(`{
        "Image": {
            "Width":  800,
            "Height": 600,
            "Title":  "View from 15th Floor",
            "Thumbnail": {
                "Url":    "http://www.example.com/image/481989943",
                "Height": 125,
                "Width":  100
            },
            "Animated" : false,
            "IDs": [116, 943, 234, 38793]
          }
      }`)

	root := Must(Unmarshal(data))
	fmt.Printf("Object has %d inheritors inside", root.Size())
	// Output:
	// Object has 1 inheritors inside
}

func ExampleMust_panic() {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Unmarshal(): %s", rec)
		}
	}()
	data := []byte(`{]`)

	root := Must(Unmarshal(data))
	fmt.Printf("Object has %d inheritors inside", root.Size())
	// Output:
	// Unmarshal(): wrong symbol ']' at 1
}

func TestUnmarshal_main(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// region Custom
		{
			name:    `[]`,
			args:    args{[]byte(`[]`)},
			want:    []interface{}{},
			wantErr: false,
		},
		{
			name:    `{}`,
			args:    args{[]byte(`{}`)},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    `[1,2,3]`,
			args:    args{[]byte(`[1,2,3]`)},
			want:    []interface{}{float64(1), float64(2), float64(3)},
			wantErr: false,
		},
		{
			name:    `"string"`,
			args:    args{[]byte(`"string"`)},
			want:    "string",
			wantErr: false,
		},
		{
			name:    `"string\""`,
			args:    args{[]byte(`"string\""`)},
			want:    "string\"",
			wantErr: false,
		},
		{
			name:    `"UTF-8 ßtring"`,
			args:    args{[]byte(`"UTF-8 \u00dftring"`)},
			want:    "UTF-8 ßtring",
			wantErr: false,
		},
		{
			name:    `[{"null":null}]`,
			args:    args{[]byte(`[{"null":null}]`)},
			want:    []interface{}{map[string]interface{}{"null": nil}},
			wantErr: false,
		},
		{
			name:    `{"key"}`,
			args:    args{[]byte(`{"key"}`)},
			wantErr: true,
		},
		{
			name:    `1e`,
			args:    args{[]byte(`1e`)},
			wantErr: true,
		},
		{
			name:    `1e+`,
			args:    args{[]byte(`1e+`)},
			wantErr: true,
		},
		{
			name:    `1e-`,
			args:    args{[]byte(`1e-`)},
			wantErr: true,
		},
		{
			name:    `-1.3e2`,
			args:    args{[]byte(`-1.3e2`)},
			want:    float64(-130),
			wantErr: false,
		},
		{
			name:    `y_object_key_unicode`,
			args:    args{[]byte(`{"\u041f\u043e\u043b\u0442\u043e\u0440\u0430\n\u0417\u0435\u043c\u043b\u0435\u043a\u043e\u043f\u0430": true}`)},
			want:    map[string]interface{}{"Полтора\nЗемлекопа": true},
			wantErr: false,
		},
		// endregion
		// region From https://json.org/example.html
		{
			name: "example#glossary",
			args: args{[]byte(`{
						"glossary": {
							"title": "example glossary",
							"GlossDiv": {
								"title": "S",
								"GlossList": {
									"GlossEntry": {
										"ID": "SGML",
										"SortAs": "SGML",
										"GlossTerm": "Standard Generalized Markup Language",
										"Acronym": "SGML",
										"Abbrev": "ISO 8879:1986",
										"GlossDef": {
											"para": "A meta-markup language, used to create markup languages such as DocBook.",
											"GlossSeeAlso": ["GML", "XML"]
										},
										"GlossSee": "markup"
									}
								}
							}
						}
					}`)},
			want: map[string]interface{}{
				"glossary": map[string]interface{}{
					"title": "example glossary",
					"GlossDiv": map[string]interface{}{
						"title": "S",
						"GlossList": map[string]interface{}{
							"GlossEntry": map[string]interface{}{
								"ID":        "SGML",
								"SortAs":    "SGML",
								"GlossTerm": "Standard Generalized Markup Language",
								"Acronym":   "SGML",
								"Abbrev":    "ISO 8879:1986",
								"GlossDef": map[string]interface{}{
									"para":         "A meta-markup language, used to create markup languages such as DocBook.",
									"GlossSeeAlso": []interface{}{"GML", "XML"},
								},
								"GlossSee": "markup",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: `example#menu`,
			args: args{[]byte(`{"menu": {
					  "id": "file",
					  "value": "File",
					  "popup": {
						"menuitem": [
						  {"value": "New", "onclick": "CreateNewDoc()"},
						  {"value": "Open", "onclick": "OpenDoc()"},
						  {"value": "Close", "onclick": "CloseDoc()"}
						]
					  }
					}}`)},
			want: map[string]interface{}{"menu": map[string]interface{}{
				"id":    "file",
				"value": "File",
				"popup": map[string]interface{}{
					"menuitem": []interface{}{
						map[string]interface{}{"value": "New", "onclick": "CreateNewDoc()"},
						map[string]interface{}{"value": "Open", "onclick": "OpenDoc()"},
						map[string]interface{}{"value": "Close", "onclick": "CloseDoc()"},
					},
				},
			}},
			wantErr: false,
		},
		{
			name: `example#widget`,
			args: args{[]byte(`{"widget": {
						"debug": "on",
						"window": {
							"title": "Sample Konfabulator Widget",
							"name": "main_window",
							"width": 500,
							"height": 500
						},
						"image": {
							"src": "Images/Sun.png",
							"name": "sun1",
							"hOffset": 250,
							"vOffset": 250,
							"alignment": "center"
						},
						"text": {
							"data": "Click Here",
							"size": 36,
							"style": "bold",
							"name": "text1",
							"hOffset": 250,
							"vOffset": 100,
							"alignment": "center",
							"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
						}
					}}`)},
			want: map[string]interface{}{"widget": map[string]interface{}{
				"debug": "on",
				"window": map[string]interface{}{
					"title":  "Sample Konfabulator Widget",
					"name":   "main_window",
					"width":  float64(500),
					"height": float64(500),
				},
				"image": map[string]interface{}{
					"src":       "Images/Sun.png",
					"name":      "sun1",
					"hOffset":   float64(250),
					"vOffset":   float64(250),
					"alignment": "center",
				},
				"text": map[string]interface{}{
					"data":      "Click Here",
					"size":      float64(36),
					"style":     "bold",
					"name":      "text1",
					"hOffset":   float64(250),
					"vOffset":   float64(100),
					"alignment": "center",
					"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;",
				},
			}},
			wantErr: false,
		},
		{
			name: `example#web-app`,
			args: args{[]byte(`{"web-app": {
					  "servlet": [
						{
						  "servlet-name": "cofaxCDS",
						  "servlet-class": "org.cofax.cds.CDSServlet",
						  "init-param": {
							"configGlossary:installationAt": "Philadelphia, PA",
							"configGlossary:adminEmail": "ksm@pobox.com",
							"configGlossary:poweredBy": "Cofax",
							"configGlossary:poweredByIcon": "/images/cofax.gif",
							"configGlossary:staticPath": "/content/static",
							"templateProcessorClass": "org.cofax.WysiwygTemplate",
							"templateLoaderClass": "org.cofax.FilesTemplateLoader",
							"templatePath": "templates",
							"templateOverridePath": "",
							"defaultListTemplate": "listTemplate.htm",
							"defaultFileTemplate": "articleTemplate.htm",
							"useJSP": false,
							"jspListTemplate": "listTemplate.jsp",
							"jspFileTemplate": "articleTemplate.jsp",
							"cachePackageTagsTrack": 200,
							"cachePackageTagsStore": 200,
							"cachePackageTagsRefresh": 60,
							"cacheTemplatesTrack": 100,
							"cacheTemplatesStore": 50,
							"cacheTemplatesRefresh": 15,
							"cachePagesTrack": 200,
							"cachePagesStore": 100,
							"cachePagesRefresh": 10,
							"cachePagesDirtyRead": 10,
							"searchEngineListTemplate": "forSearchEnginesList.htm",
							"searchEngineFileTemplate": "forSearchEngines.htm",
							"searchEngineRobotsDb": "WEB-INF/robots.db",
							"useDataStore": true,
							"dataStoreClass": "org.cofax.SqlDataStore",
							"redirectionClass": "org.cofax.SqlRedirection",
							"dataStoreName": "cofax",
							"dataStoreDriver": "com.microsoft.jdbc.sqlserver.SQLServerDriver",
							"dataStoreUrl": "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
							"dataStoreUser": "sa",
							"dataStorePassword": "dataStoreTestQuery",
							"dataStoreTestQuery": "SET NOCOUNT ON;select test='test';",
							"dataStoreLogFile": "/usr/local/tomcat/logs/datastore.log",
							"dataStoreInitConns": 10,
							"dataStoreMaxConns": 100,
							"dataStoreConnUsageLimit": 100,
							"dataStoreLogLevel": "debug",
							"maxUrlLength": 500}},
						{
						  "servlet-name": "cofaxEmail",
						  "servlet-class": "org.cofax.cds.EmailServlet",
						  "init-param": {
						  "mailHost": "mail1",
						  "mailHostOverride": "mail2"}},
						{
						  "servlet-name": "cofaxAdmin",
						  "servlet-class": "org.cofax.cds.AdminServlet"},

						{
						  "servlet-name": "fileServlet",
						  "servlet-class": "org.cofax.cds.FileServlet"},
						{
						  "servlet-name": "cofaxTools",
						  "servlet-class": "org.cofax.cms.CofaxToolsServlet",
						  "init-param": {
							"templatePath": "toolstemplates/",
							"log": 1,
							"logLocation": "/usr/local/tomcat/logs/CofaxTools.log",
							"logMaxSize": "",
							"dataLog": 1,
							"dataLogLocation": "/usr/local/tomcat/logs/dataLog.log",
							"dataLogMaxSize": "",
							"removePageCache": "/content/admin/remove?cache=pages&id=",
							"removeTemplateCache": "/content/admin/remove?cache=templates&id=",
							"fileTransferFolder": "/usr/local/tomcat/webapps/content/fileTransferFolder",
							"lookInContext": 1,
							"adminGroupID": 4,
							"betaServer": true}}],
					  "servlet-mapping": {
						"cofaxCDS": "/",
						"cofaxEmail": "/cofaxutil/aemail/*",
						"cofaxAdmin": "/admin/*",
						"fileServlet": "/static/*",
						"cofaxTools": "/tools/*"},

					  "taglib": {
						"taglib-uri": "cofax.tld",
						"taglib-location": "/WEB-INF/tlds/cofax.tld"}}}`)},
			want: map[string]interface{}{"web-app": map[string]interface{}{
				"servlet": []interface{}{
					map[string]interface{}{
						"servlet-name":  "cofaxCDS",
						"servlet-class": "org.cofax.cds.CDSServlet",
						"init-param": map[string]interface{}{
							"configGlossary:installationAt": "Philadelphia, PA",
							"configGlossary:adminEmail":     "ksm@pobox.com",
							"configGlossary:poweredBy":      "Cofax",
							"configGlossary:poweredByIcon":  "/images/cofax.gif",
							"configGlossary:staticPath":     "/content/static",
							"templateProcessorClass":        "org.cofax.WysiwygTemplate",
							"templateLoaderClass":           "org.cofax.FilesTemplateLoader",
							"templatePath":                  "templates",
							"templateOverridePath":          "",
							"defaultListTemplate":           "listTemplate.htm",
							"defaultFileTemplate":           "articleTemplate.htm",
							"useJSP":                        false,
							"jspListTemplate":               "listTemplate.jsp",
							"jspFileTemplate":               "articleTemplate.jsp",
							"cachePackageTagsTrack":         float64(200),
							"cachePackageTagsStore":         float64(200),
							"cachePackageTagsRefresh":       float64(60),
							"cacheTemplatesTrack":           float64(100),
							"cacheTemplatesStore":           float64(50),
							"cacheTemplatesRefresh":         float64(15),
							"cachePagesTrack":               float64(200),
							"cachePagesStore":               float64(100),
							"cachePagesRefresh":             float64(10),
							"cachePagesDirtyRead":           float64(10),
							"searchEngineListTemplate":      "forSearchEnginesList.htm",
							"searchEngineFileTemplate":      "forSearchEngines.htm",
							"searchEngineRobotsDb":          "WEB-INF/robots.db",
							"useDataStore":                  true,
							"dataStoreClass":                "org.cofax.SqlDataStore",
							"redirectionClass":              "org.cofax.SqlRedirection",
							"dataStoreName":                 "cofax",
							"dataStoreDriver":               "com.microsoft.jdbc.sqlserver.SQLServerDriver",
							"dataStoreUrl":                  "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
							"dataStoreUser":                 "sa",
							"dataStorePassword":             "dataStoreTestQuery",
							"dataStoreTestQuery":            "SET NOCOUNT ON;select test='test';",
							"dataStoreLogFile":              "/usr/local/tomcat/logs/datastore.log",
							"dataStoreInitConns":            float64(10),
							"dataStoreMaxConns":             float64(100),
							"dataStoreConnUsageLimit":       float64(100),
							"dataStoreLogLevel":             "debug",
							"maxUrlLength":                  float64(500)}},
					map[string]interface{}{
						"servlet-name":  "cofaxEmail",
						"servlet-class": "org.cofax.cds.EmailServlet",
						"init-param": map[string]interface{}{
							"mailHost":         "mail1",
							"mailHostOverride": "mail2"}},
					map[string]interface{}{
						"servlet-name":  "cofaxAdmin",
						"servlet-class": "org.cofax.cds.AdminServlet"},

					map[string]interface{}{
						"servlet-name":  "fileServlet",
						"servlet-class": "org.cofax.cds.FileServlet"},
					map[string]interface{}{
						"servlet-name":  "cofaxTools",
						"servlet-class": "org.cofax.cms.CofaxToolsServlet",
						"init-param": map[string]interface{}{
							"templatePath":        "toolstemplates/",
							"log":                 float64(1),
							"logLocation":         "/usr/local/tomcat/logs/CofaxTools.log",
							"logMaxSize":          "",
							"dataLog":             float64(1),
							"dataLogLocation":     "/usr/local/tomcat/logs/dataLog.log",
							"dataLogMaxSize":      "",
							"removePageCache":     "/content/admin/remove?cache=pages&id=",
							"removeTemplateCache": "/content/admin/remove?cache=templates&id=",
							"fileTransferFolder":  "/usr/local/tomcat/webapps/content/fileTransferFolder",
							"lookInContext":       float64(1),
							"adminGroupID":        float64(4),
							"betaServer":          true}}},
				"servlet-mapping": map[string]interface{}{
					"cofaxCDS":    "/",
					"cofaxEmail":  "/cofaxutil/aemail/*",
					"cofaxAdmin":  "/admin/*",
					"fileServlet": "/static/*",
					"cofaxTools":  "/tools/*"},

				"taglib": map[string]interface{}{
					"taglib-uri":      "cofax.tld",
					"taglib-location": "/WEB-INF/tlds/cofax.tld"}}},
			wantErr: false,
		},
		{
			name: "example#menu2",
			args: args{[]byte(`{"menu": {
						"header": "SVG Viewer",
						"items": [
							{"id": "Open"},
							{"id": "OpenNew", "label": "Open New"},
							null,
							{"id": "ZoomIn", "label": "Zoom In"},
							{"id": "ZoomOut", "label": "Zoom Out"},
							{"id": "OriginalView", "label": "Original View"},
							null,
							{"id": "Quality"},
							{"id": "Pause"},
							{"id": "Mute"},
							null,
							{"id": "Find", "label": "Find..."},
							{"id": "FindAgain", "label": "Find Again"},
							{"id": "Copy"},
							{"id": "CopyAgain", "label": "Copy Again"},
							{"id": "CopySVG", "label": "Copy SVG"},
							{"id": "ViewSVG", "label": "View SVG"},
							{"id": "ViewSource", "label": "View Source"},
							{"id": "SaveAs", "label": "Save As"},
							null,
							{"id": "Help"},
							{"id": "About", "label": "About Adobe CVG Viewer..."}
						]
					}}`)},
			want: map[string]interface{}{"menu": map[string]interface{}{
				"header": "SVG Viewer",
				"items": []interface{}{
					map[string]interface{}{"id": "Open"},
					map[string]interface{}{"id": "OpenNew", "label": "Open New"},
					nil,
					map[string]interface{}{"id": "ZoomIn", "label": "Zoom In"},
					map[string]interface{}{"id": "ZoomOut", "label": "Zoom Out"},
					map[string]interface{}{"id": "OriginalView", "label": "Original View"},
					nil,
					map[string]interface{}{"id": "Quality"},
					map[string]interface{}{"id": "Pause"},
					map[string]interface{}{"id": "Mute"},
					nil,
					map[string]interface{}{"id": "Find", "label": "Find..."},
					map[string]interface{}{"id": "FindAgain", "label": "Find Again"},
					map[string]interface{}{"id": "Copy"},
					map[string]interface{}{"id": "CopyAgain", "label": "Copy Again"},
					map[string]interface{}{"id": "CopySVG", "label": "Copy SVG"},
					map[string]interface{}{"id": "ViewSVG", "label": "View SVG"},
					map[string]interface{}{"id": "ViewSource", "label": "View Source"},
					map[string]interface{}{"id": "SaveAs", "label": "Save As"},
					nil,
					map[string]interface{}{"id": "Help"},
					map[string]interface{}{"id": "About", "label": "About Adobe CVG Viewer..."},
				},
			}},
			wantErr: false,
		},
		// endregion
		// region TestSuite from https://github.com/nst/JSONTestSuite/blob/master/test_parsing/
		{
			name:    `i_number_double_huge_neg_exp.json`,
			args:    args{[]byte(`[123.456e-789]`)},
			want:    []interface{}{float64(123.456e-789)},
			wantErr: false,
		},
		{
			name:    `n_number_.-1`,
			args:    args{[]byte(`[.-1]`)},
			wantErr: true,
		},
		{
			name:    `n_array_double_extra_comma`,
			args:    args{[]byte(`["x",,]`)},
			wantErr: true,
		},
		{
			name:    `y_number_simple_real`,
			args:    args{[]byte(`[123.456789]`)},
			want:    []interface{}{float64(123.456789)},
			wantErr: false,
		},
		{
			name:    `n_object_non_string_key_but_huge_number_instead`,
			args:    args{[]byte(`{9999E9999:1}`)},
			wantErr: true,
		},
		{
			name:    `y_array_empty-string`,
			args:    args{[]byte(`[""]`)},
			want:    []interface{}{""},
			wantErr: false,
		},
		{
			name:    `y_number_0e1`,
			args:    args{[]byte(`[0e1]`)},
			want:    []interface{}{float64(0)},
			wantErr: false,
		},
		{
			name:    `n_number_expression`,
			args:    args{[]byte(`[1+2]`)},
			wantErr: true,
		},
		{
			name:    `y_structure_string_empty`,
			args:    args{[]byte(`""`)},
			want:    "",
			wantErr: false,
		},
		{
			name:    `y_string_unicode_U+200B_ZERO_WIDTH_SPACE`,
			args:    args{[]byte(`["\u200B"]`)},
			want:    []interface{}{"​"},
			wantErr: false,
		},
		{
			name: `y_number_double_close_to_zero`,
			args: args{[]byte(`[-0.000000000000000000000000000000000000000000000000000000000000000000000000000001]
		`)},
			want:    []interface{}{float64(-0.000000000000000000000000000000000000000000000000000000000000000000000000000001)},
			wantErr: false,
		},
		{
			name:    `n_number_real_without_fractional_part`,
			args:    args{[]byte(`[1.]`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicode_U+FFFE_nonchar`,
			args:    args{[]byte(`["\uFFFE"]`)},
			want:    []interface{}{"￾"},
			wantErr: false,
		},
		{
			name:    `n_number_0.e1`,
			args:    args{[]byte(`[0.e1]`)},
			wantErr: true,
		},
		{
			name:    `n_number_1_000`,
			args:    args{[]byte(`[1 000.0]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_object_close_array`,
			args:    args{[]byte(`{]`)},
			wantErr: true,
		},
		{
			name:    `n_number_1.0e+`,
			args:    args{[]byte(`[1.0e+]`)},
			wantErr: true,
		},
		{
			name:    `n_number_1.0e-`,
			args:    args{[]byte(`[1.0e-]`)},
			wantErr: true,
		},
		{
			name:    `y_string_double_escape_a`,
			args:    args{[]byte(`["\\a"]`)},
			want:    []interface{}{`\a`},
			wantErr: false,
		},
		{
			name:    `n_object_missing_value`,
			args:    args{[]byte(`{"a":`)},
			wantErr: true,
		},
		{
			name:    `y_number_minus_zero`,
			args:    args{[]byte(`[-0]`)},
			want:    []interface{}{float64(0)},
			wantErr: false,
		},
		{
			name:    `n_object_key_with_single_quotes`,
			args:    args{[]byte(`{key: 'value'}`)},
			wantErr: true,
		},
		{
			name:    `n_structure_object_unclosed_no_value`,
			args:    args{[]byte(`{"":`)},
			wantErr: true,
		},
		{
			name:    `n_object_non_string_key`,
			args:    args{[]byte(`{1:1}`)},
			wantErr: true,
		},
		{
			name:    `n_number_2.e+3`,
			args:    args{[]byte(`[2.e+3]`)},
			wantErr: true,
		},
		{
			name:    `y_structure_lonely_negative_real`,
			args:    args{[]byte(`-0.1`)},
			want:    float64(-0.1),
			wantErr: false,
		},
		{
			name:    `n_structure_open_array_comma`,
			args:    args{[]byte(`[,`)},
			wantErr: true,
		},
		{
			name:    `y_string_simple_ascii`,
			args:    args{[]byte(`["asd "]`)},
			want:    []interface{}{`asd `},
			wantErr: false,
		},
		{
			name:    `n_array_missing_value`,
			args:    args{[]byte(`[   , ""]`)},
			wantErr: true,
		},
		{
			name:    `y_string_backslash_doublequotes`,
			args:    args{[]byte(`["\""]`)},
			want:    []interface{}{`"`},
			wantErr: false,
		},
		{
			name:    `n_structure_lone-invalid-utf-8`,
			args:    args{[]byte(`�`)},
			wantErr: true,
		},
		{
			name:    `n_number_+Inf`,
			args:    args{[]byte(`[+Inf]`)},
			wantErr: true,
		},
		{
			name:    `n_number_2.e-3`,
			args:    args{[]byte(`[2.e-3]`)},
			wantErr: true,
		},
		{
			name:    `y_string_u+2028_line_sep`,
			args:    args{[]byte(`[" "]`)},
			want:    []interface{}{" "},
			wantErr: false,
		},
		{
			name:    `y_number_real_capital_e_pos_exp`,
			args:    args{[]byte(`[1E+2]`)},
			want:    []interface{}{float64(100)},
			wantErr: false,
		},
		{
			name:    `y_array_empty`,
			args:    args{[]byte(`[]`)},
			want:    []interface{}{},
			wantErr: false,
		},
		{
			name:    `y_string_unicode_2`,
			args:    args{[]byte(`["⍂㈴⍂"]`)},
			want:    []interface{}{"⍂㈴⍂"},
			wantErr: false,
		},
		{
			name:    `y_string_in_array`,
			args:    args{[]byte(`["asd"]`)},
			want:    []interface{}{"asd"},
			wantErr: false,
		},
		{
			name:    `n_number_real_garbage_after_e`,
			args:    args{[]byte(`[1ea]`)},
			wantErr: true,
		},
		{
			name:    `n_object_double_colon`,
			args:    args{[]byte(`{"x"::"b"}`)},
			wantErr: true,
		},
		{
			name:    `n_object_with_trailing_garbage`,
			args:    args{[]byte(`{"a":"b"}#`)},
			wantErr: true,
		},
		{
			name:    `y_string_allowed_escapes`,
			args:    args{[]byte(`["\"\\\/\b\f\n\r\t"]`)},
			want:    []interface{}{"\"\\/\b\f\n\r\t"},
			wantErr: false,
		},
		{
			name:    `n_object_missing_key`,
			args:    args{[]byte(`{:"b"}`)},
			wantErr: true,
		},
		{
			name:    `n_object_with_single_string`,
			args:    args{[]byte(`{ "foo" : "bar", "a" }`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_object`,
			args:    args{[]byte(`{`)},
			wantErr: true,
		},
		{
			name:    `n_array_inner_array_no_comma`,
			args:    args{[]byte(`[3[4]]`)},
			wantErr: true,
		},
		{
			name:    `n_object_emoji`,
			args:    args{[]byte(`{🇨🇭}`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_object_open_array`,
			args:    args{[]byte(`{[`)},
			wantErr: true,
		},
		{
			name:    `n_number_0e+`,
			args:    args{[]byte(`[0e+]`)},
			wantErr: true,
		},
		{
			name:    `n_object_missing_colon`,
			args:    args{[]byte(`{"a" b}`)},
			wantErr: true,
		},
		{
			name:    `n_structure_array_trailing_garbage`,
			args:    args{[]byte(`[1]x`)},
			wantErr: true,
		},
		{
			name:    `y_number_negative_one`,
			args:    args{[]byte(`[-1]`)},
			want:    []interface{}{float64(-1)},
			wantErr: false,
		},
		{
			name:    `y_string_nonCharacterInUTF-8_U+10FFFF`,
			args:    args{[]byte(`["􏿿"]`)},
			want:    []interface{}{"􏿿"},
			wantErr: false,
		},
		{
			name: `n_structure_open_array_object`,
			args: args{func() []byte {
				result := make([]byte, 50000)
				word := []byte(`[{"":`)
				for i := 0; i < 10000; i++ {
					result = append(result, word...)
				}
				return result
			}},
			wantErr: true,
		},
		{
			name:    `y_string_space`,
			args:    args{[]byte(`" "`)},
			want:    " ",
			wantErr: false,
		},
		{
			name:    `n_string_incomplete_surrogate_escape_invalid`,
			args:    args{[]byte(`["\uD800\uD800\x"]`)},
			wantErr: true,
		},
		{
			name:    `n_array_items_separated_by_semicolon`,
			args:    args{[]byte(`[1:2]`)},
			wantErr: true,
		},
		{
			name:    `n_string_single_string_no_double_quotes`,
			args:    args{[]byte(`abc`)},
			wantErr: true,
		},
		{
			name:    `n_structure_unclosed_array_partial_null`,
			args:    args{[]byte(`[ false, nul`)},
			wantErr: true,
		},
		{
			name: `n_object_bracket_key`,
			args: args{[]byte(`{[: "x"}
`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicode_U+FDD0_nonchar`,
			args:    args{[]byte(`["\uFDD0"]`)},
			want:    []interface{}{"﷐"},
			wantErr: false,
		},
		{
			name:    `y_string_uEscape`,
			args:    args{[]byte(`["\u0061\u30af\u30EA\u30b9"]`)},
			want:    []interface{}{"aクリス"},
			wantErr: false,
		},
		{
			name: `y_array_with_1_and_newline`,
			args: args{[]byte(`[1
]`)},
			want:    []interface{}{float64(1)},
			wantErr: false,
		},
		{
			name:    `n_structure_single_eacute`,
			args:    args{[]byte(`�`)},
			wantErr: true,
		},
		{
			name:    `n_multidigit_number_then_00`,
			args:    args{[]byte("123\000")},
			wantErr: true,
		},
		{
			name:    `n_structure_capitalized_True`,
			args:    args{[]byte(`[True]`)},
			wantErr: true,
		},
		{
			name:    `y_array_with_several_null`,
			args:    args{[]byte(`[1,null,null,null,2]`)},
			want:    []interface{}{float64(1), nil, nil, nil, float64(2)},
			wantErr: false,
		},
		{
			name:    `y_object_duplicated_key_and_value`,
			args:    args{[]byte(`{"a":"b","a":"b"}`)},
			want:    map[string]interface{}{"a": "b"},
			wantErr: false,
		},
		{
			name:    `y_number_negative_zero`,
			args:    args{[]byte(`[-0]`)},
			want:    []interface{}{float64(0)},
			wantErr: false,
		},
		{
			name:    `y_string_escaped_noncharacter`,
			args:    args{[]byte(`["\uFFFF"]`)},
			want:    []interface{}{"￿"},
			wantErr: false,
		},
		{
			name:    `n_number_-1.0.`,
			args:    args{[]byte(`[-1.0.]`)},
			wantErr: true,
		},
		{
			name:    `n_number_minus_infinity`,
			args:    args{[]byte(`[-Infinity]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_whitespace_U+2060_word_joiner`,
			args:    args{[]byte(`[⁠]`)},
			wantErr: true,
		},
		{
			name:    `n_number_invalid-negative-real`,
			args:    args{[]byte(`[-123.123foo]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_unclosed_array_unfinished_false`,
			args:    args{[]byte(`[ true, fals`)},
			wantErr: true,
		},
		{
			name:    `n_array_unclosed`,
			args:    args{[]byte(`[""`)},
			wantErr: true,
		},
		{
			name:    `y_number_real_capital_e`,
			args:    args{[]byte(`[1E22]`)},
			want:    []interface{}{float64(1e22)},
			wantErr: false,
		},
		{
			name:    `n_structure_comma_instead_of_closing_brace`,
			args:    args{[]byte(`{"x": true,`)},
			wantErr: true,
		},
		{
			name:    `n_object_no-colon`,
			args:    args{[]byte(`{"a"`)},
			wantErr: true,
		},
		{
			name:    `n_array_comma_after_close`,
			args:    args{[]byte(`[""],`)},
			wantErr: true,
		},
		{
			name:    `n_number_-2.`,
			args:    args{[]byte(`[-2.]`)},
			wantErr: true,
		},
		{
			name:    `y_object_duplicated_key`,
			args:    args{[]byte(`{"a":"b","a":"c"}`)},
			want:    map[string]interface{}{"a": "c"},
			wantErr: false,
		},
		{
			name:    `n_object_garbage_at_end`,
			args:    args{[]byte(`{"a":"a" 123}`)},
			wantErr: true,
		},
		{
			name:    `n_string_no_quotes_with_bad_escape`,
			args:    args{[]byte(`[\n]`)},
			wantErr: true,
		},
		{
			name:    `n_string_with_trailing_garbage`,
			args:    args{[]byte(`""x`)},
			wantErr: true,
		},
		{
			name:    `y_number_real_fraction_exponent`,
			args:    args{[]byte(`[123.456e78]`)},
			want:    []interface{}{float64(123.456e78)},
			wantErr: false,
		},
		{
			name:    `y_object`,
			args:    args{[]byte(`{"asd":"sdf", "dfg":"fgh"}`)},
			want:    map[string]interface{}{"asd": "sdf", "dfg": "fgh"},
			wantErr: false,
		},
		{
			name:    `n_structure_unclosed_array_unfinished_true`,
			args:    args{[]byte(`[ false, tru`)},
			wantErr: true,
		},
		{
			name:    `y_array_heterogeneous`,
			args:    args{[]byte(`[null, 1, "1", {}]`)},
			want:    []interface{}{nil, float64(1), "1", map[string]interface{}{}},
			wantErr: false,
		},
		{
			name:    `n_number_neg_real_without_int_part`,
			args:    args{[]byte(`[-.123]`)},
			wantErr: true,
		},
		{
			name:    `n_string_1_surrogate_then_escape`,
			args:    args{[]byte(`["\uD800\"]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_object_followed_by_closing_object`,
			args:    args{[]byte(`{}}`)},
			wantErr: true,
		},
		{
			name:    `y_array_with_trailing_space`,
			args:    args{[]byte(`[2] `)},
			want:    []interface{}{float64(2)},
			wantErr: false,
		},
		{
			name:    `n_object_unquoted_key`,
			args:    args{[]byte(`{a: "b"}`)},
			wantErr: true,
		},
		{
			name:    `n_number_hex_1_digit`,
			args:    args{[]byte(`[0x1]`)},
			wantErr: true,
		},
		{
			name:    `n_object_trailing_comment_open`,
			args:    args{[]byte(`{"a":"b"}/**//`)},
			wantErr: true,
		},
		{
			name:    `n_structure_angle_bracket_.`,
			args:    args{[]byte(`<.>`)},
			wantErr: true,
		},
		{
			name:    `y_structure_true_in_array`,
			args:    args{[]byte(`[true]`)},
			want:    []interface{}{true},
			wantErr: false,
		},
		{
			name:    `n_object_missing_semicolon`,
			args:    args{[]byte(`{"a" "b"}`)},
			wantErr: true,
		},
		{
			name:    `n_array_number_and_several_commas`,
			args:    args{[]byte(`[1,,]`)},
			wantErr: true,
		},
		{
			name:    `n_number_minus_sign_with_trailing_garbage`,
			args:    args{[]byte(`[-foo]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_ascii-unicode-identifier`,
			args:    args{[]byte(`aå`)},
			wantErr: true,
		},
		{
			name:    `n_number_minus_space_1`,
			args:    args{[]byte(`[- 1]`)},
			wantErr: true,
		},
		{
			name:    `y_array_false`,
			args:    args{[]byte(`[false]`)},
			want:    []interface{}{false},
			wantErr: false,
		},
		{
			name:    `n_number_1eE2`,
			args:    args{[]byte(`[1eE2]`)},
			wantErr: true,
		},
		{
			name:    `n_string_unescaped_crtl_char`,
			args:    args{[]byte("[\"a\000a\"]")},
			wantErr: true,
		},
		{
			name:    `n_number_invalid-utf-8-in-int`,
			args:    args{[]byte(`[0�]`)},
			wantErr: true,
		},
		{
			name:    `n_array_unclosed_trailing_comma`,
			args:    args{[]byte(`[1,`)},
			wantErr: true,
		},
		{
			name:    `n_structure_array_with_extra_array_close`,
			args:    args{[]byte(`[1]]`)},
			wantErr: true,
		},
		{
			name:    `n_number_invalid-utf-8-in-bigger-int`,
			args:    args{[]byte(`[123�]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_single_star`,
			args:    args{[]byte(`*`)},
			wantErr: true,
		},
		{
			name:    `n_structure_object_with_trailing_garbage`,
			args:    args{[]byte(`{"a": true} "x"`)},
			wantErr: true,
		},
		{
			name:    `y_object_string_unicode`,
			args:    args{[]byte(`{"title":"\u041f\u043e\u043b\u0442\u043e\u0440\u0430 \u0417\u0435\u043c\u043b\u0435\u043a\u043e\u043f\u0430" }`)},
			want:    map[string]interface{}{"title": "Полтора Землекопа"},
			wantErr: false,
		},
		{
			name:    `n_incomplete_true`,
			args:    args{[]byte(`[tru]`)},
			wantErr: true,
		},
		{
			name:    `y_array_with_leading_space`,
			args:    args{[]byte(` [1]`)},
			want:    []interface{}{float64(1)},
			wantErr: false,
		},
		{
			name:    `n_object_several_trailing_commas`,
			args:    args{[]byte(`{"id":0,,,,,}`)},
			wantErr: true,
		},
		{
			name:    `n_string_1_surrogate_then_escape_u1x`,
			args:    args{[]byte(`["\uD800\u1x"]`)},
			wantErr: true,
		},
		{
			name:    `y_string_surrogates_U+1D11E_MUSICAL_SYMBOL_G_CLEF`,
			args:    args{[]byte(`["\uD834\uDd1e"]`)},
			want:    []interface{}{"𝄞"},
			wantErr: false,
		},
		{
			name:    `n_structure_open_open`,
			args:    args{[]byte(`["\{["\{["\{["\{`)},
			wantErr: true,
		},
		{
			name:    `y_number`,
			args:    args{[]byte(`[123e65]`)},
			want:    []interface{}{float64(123e65)},
			wantErr: false,
		},
		{
			name:    `n_number_neg_with_garbage_at_end`,
			args:    args{[]byte(`[-1x]`)},
			wantErr: true,
		},
		{
			name:    `y_string_1_2_3_bytes_UTF-8_sequences`,
			args:    args{[]byte(`["\u0060\u012a\u12AB"]`)},
			want:    []interface{}{"`Īካ"},
			wantErr: false,
		},
		{
			name:    `y_string_unicode`,
			args:    args{[]byte(`["\uA66D"]`)},
			want:    []interface{}{"ꙭ"},
			wantErr: false,
		},
		{
			name:    `n_incomplete_false`,
			args:    args{[]byte(`[fals]`)},
			wantErr: true,
		},
		{
			name:    `n_number_U+FF11_fullwidth_digit_one`,
			args:    args{[]byte(`[１]`)},
			wantErr: true,
		},
		{
			name:    `n_incomplete_null`,
			args:    args{[]byte(`[nul]`)},
			wantErr: true,
		},
		{
			name: `y_array_ending_with_newline`,
			args: args{[]byte(`["a"
]`)},
			want:    []interface{}{"a"},
			wantErr: false,
		},
		{
			name:    `n_number_with_alpha`,
			args:    args{[]byte(`[1.2a-3]`)},
			wantErr: true,
		},
		{
			name:    `y_object_empty_key`,
			args:    args{[]byte(`{"":0}`)},
			want:    map[string]interface{}{"": float64(0)},
			wantErr: false,
		},
		{
			name:    `n_number_0.3e`,
			args:    args{[]byte(`[0.3e]`)},
			wantErr: true,
		},
		{
			name:    `y_object_basic`,
			args:    args{[]byte(`{"asd":"sdf"}`)},
			want:    map[string]interface{}{"asd": "sdf"},
			wantErr: false,
		},
		{
			name:    `y_number_real_pos_exponent`,
			args:    args{[]byte(`[1e+2]`)},
			want:    []interface{}{float64(100)},
			wantErr: false,
		},
		{
			name:    `n_string_invalid-utf-8-in-escape`,
			args:    args{[]byte(`["\u�"]`)},
			wantErr: true,
		},
		{
			name: `y_structure_trailing_newline`,
			args: args{[]byte(`["a"]
`)},
			want:    []interface{}{"a"},
			wantErr: false,
		},
		{
			name:    `y_object_empty`,
			args:    args{[]byte(`{}`)},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    `n_number_neg_int_starting_with_zero`,
			args:    args{[]byte(`[-012]`)},
			wantErr: true,
		},
		{
			name:    `y_number_negative_int`,
			args:    args{[]byte(`[-123]`)},
			want:    []interface{}{float64(-123)},
			wantErr: false,
		},
		{
			name:    `y_string_unescaped_char_delete`,
			args:    args{[]byte(`[""]`)},
			want:    []interface{}{""},
			wantErr: false,
		},
		{
			name:    `n_number_0_capital_E+`,
			args:    args{[]byte(`[0E+]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_uescaped_LF_before_string`,
			args:    args{[]byte(`[\u000A""]`)},
			wantErr: true,
		},
		{
			name:    `n_string_invalid_backslash_esc`,
			args:    args{[]byte(`["\a"]`)},
			wantErr: true,
		},
		{
			name:    `n_number_invalid+-`,
			args:    args{[]byte(`[0e+-1]`)},
			wantErr: true,
		},
		{
			name:    `n_number_++`,
			args:    args{[]byte(`[++1234]`)},
			wantErr: true,
		},
		{
			name:    `y_string_u+2029_par_sep`,
			args:    args{[]byte(`[" "]`)},
			want:    []interface{}{" "},
			wantErr: false,
		},
		{
			name:    `n_number_with_leading_zero`,
			args:    args{[]byte(`[012]`)},
			wantErr: true,
		},
		{
			name:    `n_number_Inf`,
			args:    args{[]byte(`[Inf]`)},
			wantErr: true,
		},
		{
			name:    `y_object_simple`,
			args:    args{[]byte(`{"a":[]}`)},
			want:    map[string]interface{}{"a": []interface{}{}},
			wantErr: false,
		},
		{
			name: `y_string_escaped_control_character`,
			args: args{[]byte(`["\u0012"]`)},
			want: []interface{}{""},
			wantErr: false,
		},
		{
			name:    `y_structure_lonely_null`,
			args:    args{[]byte(`null`)},
			want:    nil,
			wantErr: false,
		},
		{
			name:    `y_structure_lonely_true`,
			args:    args{[]byte(`true`)},
			want:    true,
			wantErr: false,
		},
		{
			name:    `y_array_null`,
			args:    args{[]byte(`[null]`)},
			want:    []interface{}{nil},
			wantErr: false,
		},
		{
			name:    `n_number_0_capital_E`,
			args:    args{[]byte(`[0E]`)},
			wantErr: true,
		},
		{
			name: `n_structure_100000_opening_arrays`,
			args: args{func() []byte {
				result := make([]byte, 100000)
				for i := 0; i < 100000; i++ {
					result[i] = '['
				}
				return result
			}},
			wantErr: true,
		},
		{
			name:    `n_string_single_doublequote`,
			args:    args{[]byte(`"`)},
			wantErr: true,
		},
		{
			name:    `y_string_with_del_character`,
			args:    args{[]byte(`["aa"]`)},
			want:    []interface{}{"aa"},
			wantErr: false,
		},
		{
			name:    `n_structure_trailing_#`,
			args:    args{[]byte(`{"a":"b"}#{}`)},
			wantErr: true,
		},
		{
			name:    `n_string_unicode_CapitalU`,
			args:    args{[]byte(`"\UA66D"`)},
			wantErr: true,
		},
		{
			name:    `n_structure_double_array`,
			args:    args{[]byte(`[][]`)},
			wantErr: true,
		},
		{
			name:    `n_number_hex_2_digits`,
			args:    args{[]byte(`[0x42]`)},
			wantErr: true,
		},
		{
			name:    `n_string_1_surrogate_then_escape_u`,
			args:    args{[]byte(`["\uD800\u"]`)},
			wantErr: true,
		},
		{
			name:    `y_string_in_array_with_leading_space`,
			args:    args{[]byte(`[ "asd"]`)},
			want:    []interface{}{"asd"},
			wantErr: false,
		},
		{
			name:    `n_structure_lone-open-bracket`,
			args:    args{[]byte(`[`)},
			wantErr: true,
		},
		{
			name:    `n_string_start_escape_unclosed`,
			args:    args{[]byte(`["\`)},
			wantErr: true,
		},
		{
			name: `y_object_with_newlines`,
			args: args{[]byte(`{
"a": "b"
}`)},
			want:    map[string]interface{}{"a": "b"},
			wantErr: false,
		},
		{
			name:    `y_string_unicode_U+10FFFE_nonchar`,
			args:    args{[]byte(`["\uDBFF\uDFFE"]`)},
			want:    []interface{}{"􏿾"},
			wantErr: false,
		},
		{
			name:    `n_number_invalid-utf-8-in-exponent`,
			args:    args{[]byte(`[1e1�]`)},
			wantErr: true,
		},
		{
			name:    `n_array_extra_comma`,
			args:    args{[]byte(`["",]`)},
			wantErr: true,
		},
		{
			name:    `y_string_utf8`,
			args:    args{[]byte(`["€𝄞"]`)},
			want:    []interface{}{"€𝄞"},
			wantErr: false,
		},
		{
			name:    `y_number_after_space`,
			args:    args{[]byte(`[ 4]`)},
			want:    []interface{}{float64(4)},
			wantErr: false,
		},
		{
			name:    `n_structure_angle_bracket_null`,
			args:    args{[]byte(`[<null>]`)},
			wantErr: true,
		},
		{
			name:    `n_array_1_true_without_comma`,
			args:    args{[]byte(`[1 true]`)},
			wantErr: true,
		},
		{
			name:    `n_number_with_alpha_char`,
			args:    args{[]byte(`[1.8011670033376514H-308]`)},
			wantErr: true,
		},
		{
			name:    `y_object_escaped_null_in_key`,
			args:    args{[]byte(`{"foo\u0000bar": 42}`)},
			want:    map[string]interface{}{"foo\u0000bar": float64(42)},
			wantErr: false,
		},
		{
			name:    `y_string_double_escape_n`,
			args:    args{[]byte(`["\\n"]`)},
			want:    []interface{}{"\\n"},
			wantErr: false,
		},
		{
			name:    `y_number_0e+1`,
			args:    args{[]byte(`[0e+1]`)},
			want:    []interface{}{float64(0)},
			wantErr: false,
		},
		{
			name:    `n_array_extra_close`,
			args:    args{[]byte(`["x"]]`)},
			wantErr: true,
		},
		{
			name:    `n_number_-01`,
			args:    args{[]byte(`[-01]`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicode_U+1FFFE_nonchar`,
			args:    args{[]byte(`["\uD83F\uDFFE"]`)},
			want:    []interface{}{"🿾"},
			wantErr: false,
		},
		{
			name:    `y_number_int_with_exp`,
			args:    args{[]byte(`[20e1]`)},
			want:    []interface{}{float64(200)},
			wantErr: false,
		},
		{
			name:    `y_string_three-byte-utf-8`,
			args:    args{[]byte(`["\u0821"]`)},
			want:    []interface{}{"ࠡ"},
			wantErr: false,
		},
		{
			name:    `n_array_colon_instead_of_comma`,
			args:    args{[]byte(`["": 1]`)},
			wantErr: true,
		},
		{
			name:    `n_array_invalid_utf8`,
			args:    args{[]byte(`[�]`)},
			wantErr: true,
		},
		{
			name:    `n_string_invalid_unicode_escape`,
			args:    args{[]byte(`["\uqqqq"]`)},
			wantErr: true,
		},
		{
			name:    `n_array_unclosed_with_object_inside`,
			args:    args{[]byte(`[{}`)},
			wantErr: true,
		},
		{
			name: `n_string_unescaped_newline`,
			args: args{[]byte(`["new
line"]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_number_with_trailing_garbage`,
			args:    args{[]byte(`2@`)},
			wantErr: true,
		},
		{
			name:    `y_string_one-byte-utf-8`,
			args:    args{[]byte(`["\u002c"]`)},
			want:    []interface{}{","},
			wantErr: false,
		},
		{
			name:    `n_structure_UTF8_BOM_no_data`,
			args:    args{[]byte("\uFEFF")},
			wantErr: true,
		},
		{
			name:    `y_structure_lonely_false`,
			args:    args{[]byte(`false`)},
			want:    false,
			wantErr: false,
		},
		{
			name:    `y_string_pi`,
			args:    args{[]byte(`["π"]`)},
			want:    []interface{}{"π"},
			wantErr: false,
		},
		{
			name:    `y_string_null_escape`,
			args:    args{[]byte(`["\u0000"]`)},
			want:    []interface{}{"\u0000"},
			wantErr: false,
		},
		{
			name:    `n_number_starting_with_dot`,
			args:    args{[]byte(`[.123]`)},
			wantErr: true,
		},
		{
			name:    `y_structure_lonely_int`,
			args:    args{[]byte(`42`)},
			want:    float64(42),
			wantErr: false,
		},
		{
			name:    `n_structure_array_with_unclosed_string`,
			args:    args{[]byte(`["asd]`)},
			wantErr: true,
		},
		{
			name:    `y_string_nbsp_uescaped`,
			args:    args{[]byte(`["new\u00A0line"]`)},
			want:    []interface{}{"new\u00A0line"},
			wantErr: false,
		},
		{
			name:    `n_object_repeated_null_null`,
			args:    args{[]byte(`{null:null,null:null}`)},
			wantErr: true,
		},
		{
			name:    `n_string_backslash_00`,
			args:    args{[]byte(`["\ "]`)},
			wantErr: true,
		},
		{
			name:    `n_string_single_quote`,
			args:    args{[]byte(`['single quote']`)},
			wantErr: true,
		},
		{
			name:    `n_single_space`,
			args:    args{[]byte(` `)},
			wantErr: true,
		},
		{
			name:    `n_array_star_inside`,
			args:    args{[]byte(`[*]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_array_apostrophe`,
			args:    args{[]byte(`['`)},
			wantErr: true,
		},
		{
			name:    `y_number_real_neg_exp`,
			args:    args{[]byte(`[1e-2]`)},
			want:    []interface{}{float64(1e-2)},
			wantErr: false,
		},
		{
			name:    `n_structure_open_array_string`,
			args:    args{[]byte(`["a"`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicode_U+2064_invisible_plus`,
			args:    args{[]byte(`["\u2064"]`)},
			want:    []interface{}{"⁤"},
			wantErr: false,
		},
		{
			name:    `n_object_trailing_comma`,
			args:    args{[]byte(`{"id":0,}`)},
			wantErr: true,
		},
		{
			name:    `y_array_arraysWithSpaces`,
			args:    args{[]byte(`[[]   ]`)},
			want:    []interface{}{[]interface{}{}},
			wantErr: false,
		},
		{
			name:    `n_number_0.1.2`,
			args:    args{[]byte(`[0.1.2]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_null-byte-outside-string`,
			args:    args{[]byte("[\000]")},
			wantErr: true,
		},
		{
			name:    `n_array_just_minus`,
			args:    args{[]byte(`[-]`)},
			wantErr: true,
		},
		{
			name:    `n_string_escaped_backslash_bad`,
			args:    args{[]byte(`["\\\"]`)},
			wantErr: true,
		},
		{
			name: `n_string_unescaped_tab`,
			args: args{[]byte(`["	"]`)},
			wantErr: true,
		},
		{
			name:    `n_number_.2e-3`,
			args:    args{[]byte(`[.2e-3]`)},
			wantErr: true,
		},
		{
			name:    `n_number_1.0e`,
			args:    args{[]byte(`[1.0e]`)},
			wantErr: true,
		},
		{
			name:    `n_array_a_invalid_utf8`,
			args:    args{[]byte(`[a�]`)},
			wantErr: true,
		},
		{
			name:    `n_number_0.3e+`,
			args:    args{[]byte(`[0.3e+]`)},
			wantErr: true,
		},
		{
			name:    `n_string_incomplete_escape`,
			args:    args{[]byte(`["\"]`)},
			wantErr: true,
		},
		{
			name:    `n_object_unterminated-value`,
			args:    args{[]byte(`{"a":"a`)},
			wantErr: true,
		},
		{
			name:    `n_array_incomplete_invalid_value`,
			args:    args{[]byte(`[x`)},
			wantErr: true,
		},
		{
			name:    `y_number_simple_int`,
			args:    args{[]byte(`[123]`)},
			want:    []interface{}{float64(123)},
			wantErr: false,
		},
		{
			name:    `n_object_comma_instead_of_colon`,
			args:    args{[]byte(`{"x", null}`)},
			wantErr: true,
		},
		{
			name:    `n_structure_U+2060_word_joined`,
			args:    args{[]byte(`[⁠]`)},
			wantErr: true,
		},
		{
			name:    `y_structure_whitespace_array`,
			args:    args{[]byte(` [] `)},
			want:    []interface{}{},
			wantErr: false,
		},
		{
			name:    `n_string_escape_x`,
			args:    args{[]byte(`["\x00"]`)},
			wantErr: true,
		},
		{
			name: `n_array_spaces_vertical_tab_formfeed`,
			args: args{[]byte(`["a"\f]`)},
			wantErr: true,
		},
		{
			name:    `n_number_real_with_invalid_utf8_after_e`,
			args:    args{[]byte(`[1e�]`)},
			wantErr: true,
		},
		{
			name: `n_array_unclosed_with_new_lines`,
			args: args{[]byte(`[1,
1
,1`)},
			wantErr: true,
		},
		{
			name: `n_string_escaped_ctrl_char_tab`,
			args: args{[]byte(`["\	"]`)},
			wantErr: true,
		},
		{
			name: `n_structure_whitespace_formfeed`,
			args: args{[]byte(`[]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_close_unopened_array`,
			args:    args{[]byte(`1]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_object_open_string`,
			args:    args{[]byte(`{"a`)},
			wantErr: true,
		},
		{
			name:    `n_number_infinity`,
			args:    args{[]byte(`[Infinity]`)},
			wantErr: true,
		},
		{
			name:    `n_string_leading_uescaped_thinspace`,
			args:    args{[]byte(`[\u0020"asd"]`)},
			wantErr: true,
		},
		{
			name:    `y_string_two-byte-utf-8`,
			args:    args{[]byte(`["\u0123"]`)},
			want:    []interface{}{"ģ"},
			wantErr: false,
		},
		{
			name:    `y_string_reservedCharacterInUTF-8_U+1BFFF`,
			args:    args{[]byte(`["𛿿"]`)},
			want:    []interface{}{"𛿿"},
			wantErr: false,
		},
		{
			name:    `y_string_nonCharacterInUTF-8_U+FFFF`,
			args:    args{[]byte(`["￿"]`)},
			want:    []interface{}{"￿"},
			wantErr: false,
		},
		{
			name:    `n_array_incomplete`,
			args:    args{[]byte(`["x"`)},
			wantErr: true,
		},
		{
			name:    `n_structure_no_data`,
			args:    args{[]byte(``)},
			wantErr: true,
		},
		{
			name:    `y_string_comments`,
			args:    args{[]byte(`["a/*b*/c/*d//e"]`)},
			want:    []interface{}{"a/*b*/c/*d//e"},
			wantErr: false,
		},
		{
			name:    `n_structure_unclosed_object`,
			args:    args{[]byte(`{"asd":"asd"`)},
			wantErr: true,
		},
		{
			name:    `n_string_incomplete_escaped_character`,
			args:    args{[]byte(`["\u00A"]`)},
			wantErr: true,
		},
		{
			name:    `n_object_trailing_comment`,
			args:    args{[]byte(`{"a":"b"}/**/`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_array_open_string`,
			args:    args{[]byte(`["a`)},
			wantErr: true,
		},
		{
			name:    `n_number_-NaN`,
			args:    args{[]byte(`[-NaN]`)},
			wantErr: true,
		},
		{
			name:    `y_string_accepted_surrogate_pairs`,
			args:    args{[]byte(`["\ud83d\ude39\ud83d\udc8d"]`)},
			want:    []interface{}{"😹💍"},
			wantErr: false,
		},
		{
			name:    `n_object_lone_continuation_byte_in_key_and_trailing_comma`,
			args:    args{[]byte(`{"�":"0",}`)},
			wantErr: true,
		},
		{
			name:    `n_string_escaped_emoji`,
			args:    args{[]byte(`["\🌀"]`)},
			wantErr: true,
		},
		{
			name:    `y_number_real_capital_e_neg_exp`,
			args:    args{[]byte(`[1E-2]`)},
			want:    []interface{}{float64(1e-2)},
			wantErr: false,
		},
		{
			name:    `n_number_+1`,
			args:    args{[]byte(`[+1]`)},
			wantErr: true,
		},
		{
			name:    `n_object_bad_value`,
			args:    args{[]byte(`["x", truth]`)},
			wantErr: true,
		},
		{
			name:    `n_number_2.e3`,
			args:    args{[]byte(`[2.e3]`)},
			wantErr: true,
		},
		{
			name:    `n_string_1_surrogate_then_escape_u1`,
			args:    args{[]byte(`["\uD800\u1"]`)},
			wantErr: true,
		},
		{
			name: `n_array_newlines_unclosed`,
			args: args{[]byte(`["a",
4
,1,`)},
			wantErr: true,
		},
		{
			name:    `n_number_NaN`,
			args:    args{[]byte(`[NaN]`)},
			wantErr: true,
		},
		{
			name:    `n_array_number_and_comma`,
			args:    args{[]byte(`[1,]`)},
			wantErr: true,
		},
		{
			name:    `n_array_comma_and_number`,
			args:    args{[]byte(`[,1]`)},
			wantErr: true,
		},
		{
			name:    `n_object_trailing_comment_slash_open`,
			args:    args{[]byte(`{"a":"b"}//`)},
			wantErr: true,
		},
		{
			name:    `n_string_incomplete_surrogate`,
			args:    args{[]byte(`["\uD834\uDd"]`)},
			wantErr: true,
		},
		{
			name:    `n_object_single_quote`,
			args:    args{[]byte(`{'a':0}`)},
			wantErr: true,
		},
		{
			name:    `y_object_extreme_numbers`,
			args:    args{[]byte(`{ "min": -1.0e+28, "max": 1.0e+28 }`)},
			want:    map[string]interface{}{"min": -1.0e+28, "max": 1.0e+28},
			wantErr: false,
		},
		{
			name:    `n_number_0e`,
			args:    args{[]byte(`[0e]`)},
			wantErr: true,
		},
		{
			name:    `y_string_accepted_surrogate_pair`,
			args:    args{[]byte(`["\uD801\udc37"]`)},
			want:    []interface{}{"𐐷"},
			wantErr: false,
		},
		{
			name:    `n_string_invalid_utf8_after_escape`,
			args:    args{[]byte(`["\�"]`)},
			wantErr: true,
		},
		{
			name:    `n_array_double_comma`,
			args:    args{[]byte(`[1,,2]`)},
			wantErr: true,
		},
		{
			name:    `n_number_9.e+`,
			args:    args{[]byte(`[9.e+]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_open_object_string_with_apostrophes`,
			args:    args{[]byte(`{'a'`)},
			wantErr: true,
		},
		{
			name:    `n_array_just_comma`,
			args:    args{[]byte(`[,]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_object_with_comment`,
			args:    args{[]byte(`{"a":/*comment*/"b"}`)},
			wantErr: true,
		},
		{
			name:    `y_structure_lonely_string`,
			args:    args{[]byte(`"asd"`)},
			want:    "asd",
			wantErr: false,
		},
		{
			name:    `n_structure_open_array_open_object`,
			args:    args{[]byte(`[{`)},
			wantErr: true,
		},
		{
			name: `y_object_long_strings`,
			args: args{[]byte(`{"x":[{"id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}], "id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}`)},
			want: map[string]interface{}{
				"x": []interface{}{
					map[string]interface{}{"id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
				},
				"id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			wantErr: false,
		},
		{
			name:    `n_object_two_commas_in_a_row`,
			args:    args{[]byte(`{"a":"b",,"c":"d"}`)},
			wantErr: true,
		},
		{
			name:    `y_string_uescaped_newline`,
			args:    args{[]byte(`["new\u000Aline"]`)},
			want:    []interface{}{"new\nline"},
			wantErr: false,
		},
		{
			name:    `n_structure_unclosed_array`,
			args:    args{[]byte(`[1`)},
			wantErr: true,
		},
		{
			name:    `n_string_accentuated_char_no_quotes`,
			args:    args{[]byte(`[é]`)},
			wantErr: true,
		},
		{
			name:    `n_structure_end_array`,
			args:    args{[]byte(`]`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicode_escaped_double_quote`,
			args:    args{[]byte(`["\u0022"]`)},
			want:    []interface{}{"\""},
			wantErr: false,
		},
		{
			name:    `y_string_last_surrogates_1_and_2`,
			args:    args{[]byte(`["\uDBFF\uDFFF"]`)},
			want:    []interface{}{"􏿿"},
			wantErr: false,
		},
		{
			name:    `n_structure_open_object_comma`,
			args:    args{[]byte(`{,`)},
			wantErr: true,
		},
		{
			name:    `y_string_unicodeEscapedBackslash`,
			args:    args{[]byte(`["\u005C"]`)},
			want:    []interface{}{"\\"},
			wantErr: false,
		},
		{
			name:    `y_string_backslash_and_u_escaped_zero`,
			args:    args{[]byte(`["\\u0000"]`)},
			want:    []interface{}{`\u0000`},
			wantErr: false,
		},
		{
			name:    `n_object_trailing_comment_slash_open_incomplete`,
			args:    args{[]byte(`{"a":"b"}/`)},
			wantErr: true,
		},
		{
			name:    `n_structure_incomplete_UTF8_BOM`,
			args:    args{[]byte(`�{}`)},
			wantErr: true,
		},
		{
			name:    `n_structure_unicode-identifier`,
			args:    args{[]byte(`å`)},
			wantErr: true,
		},
		{
			name:    `y_number_real_exponent`,
			args:    args{[]byte(`[123e45]`)},
			want:    []interface{}{123e45},
			wantErr: false,
		},
		// endregion
		// region Generated
		{
			name: "[[[...500...]]]",
			args: args{func() []byte {
				result := make([]byte, 1000)
				for i := 0; i < 500; i++ {
					result[i] = '['
					result[999-i] = ']'
				}
				return result
			}},
			want: func() interface{} {
				var prev interface{}
				var current []interface{}
				for i := 0; i < 500; i++ {
					current = make([]interface{}, 0)
					if prev != nil {
						current = append(current, prev)
					}
					prev = current
				}
				return current
			},
			wantErr: false,
		},
		{
			name: "[[[...5000...]]]",
			args: args{func() []byte {
				result := make([]byte, 10000)
				for i := 0; i < 5000; i++ {
					result[i] = '['
					result[9999-i] = ']'
				}
				return result
			}},
			want: func() interface{} {
				var prev interface{}
				var current []interface{}
				for i := 0; i < 5000; i++ {
					current = make([]interface{}, 0)
					if prev != nil {
						current = append(current, prev)
					}
					prev = current
				}
				return current
			},
			wantErr: false,
		},
		// endregion
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input []byte
			switch tt.args.data.(type) {
			case []byte:
				input = tt.args.data.([]byte)
			case func() []byte:
				input = tt.args.data.(func() []byte)()
			}
			root, err := Unmarshal(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v.", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if root == nil {
				t.Errorf("Unmarshal() return nil as result")
				return
			}

			var want interface{}
			switch tt.want.(type) {
			case func() interface{}:
				want = tt.want.(func() interface{})()
			default:
				want = tt.want
			}

			got, err := root.Interface()
			if err != nil {
				t.Errorf("Unmarshal() error = %v, wantErr %v. got = %v", err, tt.wantErr, got)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Unmarshal() got = %v, want %v", got, want)
			}
		})
	}
}
