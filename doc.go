// Package ajson implements decoding of JSON as defined in RFC 7159 without predefined mapping to a struct of golang, with support of JSONPath.
//
// All JSON structs reflects to a custom struct of Node, witch can be presented by it type and value.
//
// Method Unmarshal will scan all the byte slice to create a root node of JSON structure, with all it behaviors.
//
// Each Node has it's own type and calculated value, which will be calculated on demand.
// Calculated value saves in atomic.Value, so it's thread safe.
//
// Method JSONPath will returns slice of founded elements in current JSON data, by it's JSONPath.
//
// JSONPath selection described at http://goessner.net/articles/JsonPath/
package ajson
