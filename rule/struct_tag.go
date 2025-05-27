package rule

import (
	"fmt"
	"go/ast"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// StructTagRule lints struct tags.
type StructTagRule struct {
	userDefined map[tagKey][]string // map: key -> []option
}

type tagKey string

const (
	keyASN1         tagKey = "asn1"
	keyBSON         tagKey = "bson"
	keyDatastore    tagKey = "datastore"
	keyDefault      tagKey = "default"
	keyJSON         tagKey = "json"
	keyMapstructure tagKey = "mapstructure"
	keyProperties   tagKey = "properties"
	keyProtobuf     tagKey = "protobuf"
	keyRequired     tagKey = "required"
	keyTOML         tagKey = "toml"
	keyURL          tagKey = "url"
	keyValidate     tagKey = "validate"
	keyXML          tagKey = "xml"
	keyYAML         tagKey = "yaml"
)

type tagChecker func(checkCtx *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeeded bool)

var tagCheckers = map[tagKey]tagChecker{
	keyASN1:         checkASN1Tag,
	keyBSON:         checkBSONTag,
	keyDatastore:    checkDatastoreTag,
	keyDefault:      checkDefaultTag,
	keyJSON:         checkJSONTag,
	keyMapstructure: checkMapstructureTag,
	keyProperties:   checkPropertiesTag,
	keyProtobuf:     checkProtobufTag,
	keyRequired:     checkRequiredTag,
	keyTOML:         checkTOMLTag,
	keyURL:          checkURLTag,
	keyValidate:     checkValidateTag,
	keyXML:          checkXMLTag,
	keyYAML:         checkYAMLTag,
}

type checkContext struct {
	userDefined    map[tagKey][]string // map: key -> []option
	usedTagNbr     map[int]bool        // list of used tag numbers
	usedTagName    map[string]bool     // list of used tag keys
	isAtLeastGo124 bool
}

func (checkCtx checkContext) isUserDefined(key tagKey, opt string) bool {
	if checkCtx.userDefined == nil {
		return false
	}

	options := checkCtx.userDefined[key]
	return slices.Contains(options, opt)
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *StructTagRule) Configure(arguments lint.Arguments) error {
	if len(arguments) == 0 {
		return nil
	}

	err := checkNumberOfArguments(1, arguments, r.Name())
	if err != nil {
		return err
	}

	r.userDefined = make(map[tagKey][]string, len(arguments))
	for _, arg := range arguments {
		item, ok := arg.(string)
		if !ok {
			return fmt.Errorf("invalid argument to the %s rule. Expecting a string, got %v (of type %T)", r.Name(), arg, arg)
		}
		parts := strings.Split(item, ",")
		if len(parts) < 2 {
			return fmt.Errorf("invalid argument to the %s rule. Expecting a string of the form key[,option]+, got %s", r.Name(), item)
		}
		key := tagKey(strings.TrimSpace(parts[0]))
		for i := 1; i < len(parts); i++ {
			option := strings.TrimSpace(parts[i])
			r.userDefined[key] = append(r.userDefined[key], option)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *StructTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintStructTagRule{
		onFailure:      onFailure,
		userDefined:    r.userDefined,
		isAtLeastGo124: file.Pkg.IsAtLeastGoVersion(lint.Go124),
		tagCheckers:    tagCheckers,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*StructTagRule) Name() string {
	return "struct-tag"
}

type lintStructTagRule struct {
	onFailure      func(lint.Failure)
	userDefined    map[tagKey][]string // map: key -> []option
	isAtLeastGo124 bool
	tagCheckers    map[tagKey]tagChecker
}

func (w lintStructTagRule) Visit(node ast.Node) ast.Visitor {
	if n, ok := node.(*ast.StructType); ok {
		isEmptyStruct := n.Fields == nil || n.Fields.NumFields() < 1
		if isEmptyStruct {
			return nil // skip empty structs
		}

		checkCtx := &checkContext{
			userDefined:    w.userDefined,
			usedTagNbr:     map[int]bool{},
			usedTagName:    map[string]bool{},
			isAtLeastGo124: w.isAtLeastGo124,
		}

		for _, f := range n.Fields.List {
			if f.Tag != nil {
				w.checkTaggedField(checkCtx, f)
			}
		}
	}

	return w
}

// checkTaggedField checks the tag of the given field.
// precondition: the field has a tag
func (w lintStructTagRule) checkTaggedField(checkCtx *checkContext, f *ast.Field) {
	if len(f.Names) > 0 && !f.Names[0].IsExported() {
		w.addFailuref(f, "tag on not-exported field %s", f.Names[0].Name)
	}

	tags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
	if err != nil || tags == nil {
		w.addFailuref(f.Tag, "malformed tag")
		return
	}

	for _, tag := range tags.Tags() {
		if msg, ok := w.checkTagNameIfNeed(checkCtx, tag); !ok {
			w.addFailureWithTagKey(f.Tag, msg, tag.Key)
		}

		checker, ok := w.tagCheckers[tagKey(tag.Key)]
		if !ok {
			continue // we don't have a checker for the tag
		}

		msg, ok := checker(checkCtx, tag, f.Type)
		if !ok {
			w.addFailureWithTagKey(f.Tag, msg, tag.Key)
		}
	}
}

func (w lintStructTagRule) checkTagNameIfNeed(checkCtx *checkContext, tag *structtag.Tag) (message string, succeeded bool) {
	isUnnamedTag := tag.Name == "" || tag.Name == "-"
	if isUnnamedTag {
		return "", true
	}

	key := tagKey(tag.Key)
	switch key {
	case keyBSON, keyJSON, keyXML, keyYAML, keyProtobuf:
	default:
		return "", true
	}

	tagName := w.getTagName(tag)
	if tagName == "" {
		return "", true // No tag name found
	}

	// We concat the key and name as the mapping key here
	// to allow the same tag name in different tag type.
	mapKey := tag.Key + ":" + tagName
	if _, ok := checkCtx.usedTagName[mapKey]; ok {
		return fmt.Sprintf("duplicated tag name %q", tagName), false
	}

	checkCtx.usedTagName[mapKey] = true

	return "", true
}

func (lintStructTagRule) getTagName(tag *structtag.Tag) string {
	key := tagKey(tag.Key)
	switch key {
	case keyProtobuf:
		for _, option := range tag.Options {
			if tagKey, found := strings.CutPrefix(option, "name="); found {
				return tagKey
			}
		}
		return "" // protobuf tag lacks 'name' option
	default:
		return tag.Name
	}
}

func checkASN1Tag(checkCtx *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeeded bool) {
	checkList := slices.Concat(tag.Options, []string{tag.Name})
	for _, opt := range checkList {
		switch opt {
		case "application", "explicit", "generalized", "ia5", "omitempty", "optional", "set", "utf8":
			// do nothing
		default:
			msg, ok := checkCompoundANS1Option(checkCtx, opt, fieldType)
			if !ok {
				return msg, false
			}
		}
	}

	return "", true
}

func checkCompoundANS1Option(checkCtx *checkContext, opt string, fieldType ast.Expr) (message string, succeeded bool) {
	key, value, _ := strings.Cut(opt, ":")
	switch key {
	case "tag":
		number, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Sprintf("tag must be a number but is %q", value), false
		}
		if checkCtx.usedTagNbr[number] {
			return fmt.Sprintf(msgDuplicatedTagNumber, number), false
		}
		checkCtx.usedTagNbr[number] = true
	case "default":
		if !typeValueMatch(fieldType, value) {
			return msgTypeMismatch, false
		}
	default:
		if !checkCtx.isUserDefined(keyASN1, opt) {
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}
	return "", true
}

func checkDatastoreTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "flatten", "noindex", "omitempty":
		default:
			if checkCtx.isUserDefined(keyDatastore, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkDefaultTag(_ *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeeded bool) {
	if !typeValueMatch(fieldType, tag.Name) {
		return msgTypeMismatch, false
	}

	return "", true
}

func checkBSONTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "inline", "minsize", "omitempty":
		default:
			if checkCtx.isUserDefined(keyBSON, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkJSONTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty", "string":
		case "":
			// special case for JSON key "-"
			if tag.Name != "-" {
				return "option can not be empty", false
			}
		case "omitzero":
			if checkCtx.isAtLeastGo124 {
				continue
			}
			return `prior Go 1.24, option "omitzero" is unsupported`, false
		default:
			if checkCtx.isUserDefined(keyJSON, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkMapstructureTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty", "reminder", "squash":
		default:
			if checkCtx.isUserDefined(keyMapstructure, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkPropertiesTag(_ *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeeded bool) {
	options := tag.Options
	if len(options) == 0 {
		return "", true
	}

	seenOptions := map[string]bool{}
	for _, opt := range options {
		msg, ok := fmt.Sprintf("unknown or malformed option %q", opt), false
		if key, value, found := strings.Cut(opt, "="); found {
			msg, ok = checkCompoundPropertiesOption(key, value, fieldType, seenOptions)
		}

		if !ok {
			return msg, false
		}
	}

	return "", true
}

func checkCompoundPropertiesOption(key, value string, fieldType ast.Expr, seenOptions map[string]bool) (message string, succeeded bool) {
	if _, ok := seenOptions[key]; ok {
		return fmt.Sprintf(msgDuplicatedOption, key), false
	}
	seenOptions[key] = true

	if strings.TrimSpace(value) == "" {
		return fmt.Sprintf("option %q not of the form %s=value", key, key), false
	}

	switch key {
	case "default":
		if !typeValueMatch(fieldType, value) {
			return msgTypeMismatch, false
		}
	case "layout":
		if astutils.GoFmt(fieldType) != "time.Time" {
			return "layout option is only applicable to fields of type time.Time", false
		}
	}

	return "", true
}

func checkProtobufTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	// check name
	switch tag.Name {
	case "bytes", "fixed32", "fixed64", "group", "varint", "zigzag32", "zigzag64":
		// do nothing
	default:
		return fmt.Sprintf("invalid tag name %q", tag.Name), false
	}

	return checkProtobufOptions(checkCtx, tag.Options)
}

func checkProtobufOptions(checkCtx *checkContext, options []string) (message string, succeeded bool) {
	seenOptions := map[string]bool{}
	hasName := false
	for _, opt := range options {
		opt := strings.Split(opt, "=")[0]

		if number, err := strconv.Atoi(opt); err == nil {
			_, alreadySeen := checkCtx.usedTagNbr[number]
			if alreadySeen {
				return fmt.Sprintf(msgDuplicatedTagNumber, number), false
			}
			checkCtx.usedTagNbr[number] = true
			continue // option is an integer
		}

		switch opt {
		case "json", "opt", "proto3", "rep", "req":
			// do nothing
		case "name":
			hasName = true
		default:
			if checkCtx.isUserDefined(keyProtobuf, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}

		_, alreadySeen := seenOptions[opt]
		if alreadySeen {
			return fmt.Sprintf(msgDuplicatedOption, opt), false
		}
		seenOptions[opt] = true
	}

	if !hasName {
		return `mandatory option "name" not found`, false
	}

	return "", true
}

func checkRequiredTag(_ *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	switch tag.Name {
	case "true", "false":
		return "", true
	default:
		return `required should be "true" or "false"`, false
	}
}

func checkTOMLTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty":
		default:
			if checkCtx.isUserDefined(keyTOML, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkURLTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	var delimiter = ""
	for _, opt := range tag.Options {
		switch opt {
		case "int", "omitempty", "numbered", "brackets":
		case "unix", "unixmilli", "unixnano": // TODO : check that the field is of type time.Time
		case "comma", "semicolon", "space":
			if delimiter == "" {
				delimiter = opt
				continue
			}
			return fmt.Sprintf("can not set both %q and %q as delimiters", opt, delimiter), false
		default:
			if checkCtx.isUserDefined(keyURL, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkValidateTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	previousOption := ""
	seenKeysOption := false
	options := append([]string{tag.Name}, tag.Options...)
	for _, opt := range options {
		switch opt {
		case "keys":
			if previousOption != "dive" {
				return `option "keys" must follow a "dive" option`, false
			}
			seenKeysOption = true
		case "endkeys":
			if !seenKeysOption {
				return `option "endkeys" without a previous "keys" option`, false
			}
			seenKeysOption = false
		default:
			parts := strings.Split(opt, "|")
			errMsg, ok := checkValidateOptionsAlternatives(checkCtx, parts)
			if !ok {
				return errMsg, false
			}
		}
		previousOption = opt
	}

	return "", true
}

func checkXMLTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "any", "attr", "cdata", "chardata", "comment", "innerxml", "omitempty", "typeattr":
		default:
			if checkCtx.isUserDefined(keyXML, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkYAMLTag(checkCtx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "flow", "inline", "omitempty":
		default:
			if checkCtx.isUserDefined(keyYAML, opt) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, opt), false
		}
	}

	return "", true
}

func checkValidateOptionsAlternatives(checkCtx *checkContext, alternatives []string) (message string, succeeded bool) {
	for _, alternative := range alternatives {
		alternative := strings.TrimSpace(alternative)
		lhs, _, found := strings.Cut(alternative, "=")
		if found {
			_, ok := validateLHS[lhs]
			if ok || checkCtx.isUserDefined(keyValidate, lhs) {
				continue
			}
			return fmt.Sprintf(msgUnknownOption, lhs), false
		}

		badOpt, ok := areValidateOpts(alternative)
		if ok || checkCtx.isUserDefined(keyValidate, badOpt) {
			continue
		}

		return fmt.Sprintf(msgUnknownOption, badOpt), false
	}

	return "", true
}

func typeValueMatch(t ast.Expr, val string) bool {
	tID, ok := t.(*ast.Ident)
	if !ok {
		return true
	}

	typeMatches := true
	switch tID.Name {
	case "bool":
		typeMatches = val == "true" || val == "false"
	case "float64":
		_, err := strconv.ParseFloat(val, 64)
		typeMatches = err == nil
	case "int":
		_, err := strconv.ParseInt(val, 10, 64)
		typeMatches = err == nil
	case "string":
	case "nil":
	default:
		// unchecked type
	}

	return typeMatches
}

func (w lintStructTagRule) addFailureWithTagKey(n ast.Node, msg, tagKey string) {
	w.addFailuref(n, "%s in %s tag", msg, tagKey)
}

func (w lintStructTagRule) addFailuref(n ast.Node, msg string, args ...any) {
	w.onFailure(lint.Failure{
		Node:       n,
		Failure:    fmt.Sprintf(msg, args...),
		Confidence: 1,
	})
}

func areValidateOpts(opts string) (string, bool) {
	parts := strings.Split(opts, "|")
	for _, opt := range parts {
		_, ok := validateSingleOptions[opt]
		if !ok {
			return opt, false
		}
	}

	return "", true
}

const (
	msgDuplicatedOption    = "duplicated option %q"
	msgDuplicatedTagNumber = "duplicated tag number %v"
	msgUnknownOption       = "unknown option %q"
	msgTypeMismatch        = "type mismatch between field type and default value type"
)

var validateSingleOptions = map[string]struct{}{
	"alpha":                         {},
	"alphanum":                      {},
	"alphanumunicode":               {},
	"alphaunicode":                  {},
	"ascii":                         {},
	"base32":                        {},
	"base64":                        {},
	"base64rawurl":                  {},
	"base64url":                     {},
	"bcp47_language_tag":            {},
	"bic":                           {},
	"boolean":                       {},
	"btc_addr":                      {},
	"btc_addr_bech32":               {},
	"cidr":                          {},
	"cidrv4":                        {},
	"cidrv6":                        {},
	"credit_card":                   {},
	"cron":                          {},
	"cve":                           {},
	"datauri":                       {},
	"dir":                           {},
	"dirpath":                       {},
	"dive":                          {},
	"dns_rfc1035_label":             {},
	"e164":                          {},
	"ein":                           {},
	"email":                         {},
	"eth_addr":                      {},
	"eth_addr_checksum":             {},
	"file":                          {},
	"filepath":                      {},
	"fqdn":                          {},
	"hexadecimal":                   {},
	"hexcolor":                      {},
	"hostname":                      {},
	"hostname_port":                 {},
	"hostname_rfc1123":              {},
	"hsl":                           {},
	"hsla":                          {},
	"html":                          {},
	"html_encoded":                  {},
	"http_url":                      {},
	"image":                         {},
	"ip":                            {},
	"ip_addr":                       {},
	"ip4_addr":                      {},
	"ip6_addr":                      {},
	"ipv4":                          {},
	"ipv6":                          {},
	"isbn":                          {},
	"isbn10":                        {},
	"isbn13":                        {},
	"isdefault":                     {},
	"iso3166_1_alpha_numeric":       {},
	"iso3166_1_alpha_numeric_eu":    {},
	"iso3166_1_alpha2":              {},
	"iso3166_1_alpha2_eu":           {},
	"iso3166_1_alpha3":              {},
	"iso3166_1_alpha3_eu":           {},
	"iso3166_2":                     {},
	"iso4217":                       {},
	"iso4217_numeric":               {},
	"issn":                          {},
	"json":                          {},
	"jwt":                           {},
	"latitude":                      {},
	"longitude":                     {},
	"lowercase":                     {},
	"luhn_checksum":                 {},
	"mac":                           {},
	"md4":                           {},
	"md5":                           {},
	"mongodb":                       {},
	"mongodb_connection_string":     {},
	"multibyte":                     {},
	"number":                        {},
	"numeric":                       {},
	"port":                          {},
	"postcode_iso3166_alpha2":       {},
	"postcode_iso3166_alpha2_field": {},
	"printascii":                    {},
	"required":                      {},
	"rgb":                           {},
	"rgba":                          {},
	"ripemd128":                     {},
	"ripemd160":                     {},
	"semver":                        {},
	"sha256":                        {},
	"sha384":                        {},
	"sha512":                        {},
	"ssn":                           {},
	"tcp_addr":                      {},
	"tcp4_addr":                     {},
	"tcp6_addr":                     {},
	"tiger128":                      {},
	"tiger160":                      {},
	"tiger192":                      {},
	"timezone":                      {},
	"udp_addr":                      {},
	"udp4_addr":                     {},
	"udp6_addr":                     {},
	"ulid":                          {},
	"unix_addr":                     {},
	"uppercase":                     {},
	"uri":                           {},
	"url":                           {},
	"url_encoded":                   {},
	"urn_rfc2141":                   {},
	"uuid":                          {},
	"uuid_rfc4122":                  {},
	"uuid3":                         {},
	"uuid3_rfc4122":                 {},
	"uuid4":                         {},
	"uuid4_rfc4122":                 {},
	"uuid5":                         {},
	"uuid5_rfc4122":                 {},
}

// These are options that are used in expressions of the form:
//
//	<option> = <RHS>
var validateLHS = map[string]struct{}{
	"contains":             {},
	"containsany":          {},
	"containsfield":        {},
	"containsrune":         {},
	"datetime":             {},
	"endsnotwith":          {},
	"endswith":             {},
	"eq":                   {},
	"eq_ignore_case":       {},
	"eqcsfield":            {},
	"eqfield":              {},
	"excluded_if":          {},
	"excluded_unless":      {},
	"excluded_with":        {},
	"excluded_with_all":    {},
	"excluded_without":     {},
	"excluded_without_all": {},
	"excludes":             {},
	"excludesall":          {},
	"excludesfield":        {},
	"excludesrune":         {},
	"fieldcontains":        {},
	"fieldexcludes":        {},
	"gt":                   {},
	"gtcsfield":            {},
	"gte":                  {},
	"gtecsfield":           {},
	"gtefield":             {},
	"gtfield":              {},
	"len":                  {},
	"lt":                   {},
	"ltcsfield":            {},
	"lte":                  {},
	"ltecsfield":           {},
	"ltefield":             {},
	"ltfield":              {},
	"max":                  {},
	"min":                  {},
	"ne":                   {},
	"ne_ignore_case":       {},
	"necsfield":            {},
	"nefield":              {},
	"oneof":                {},
	"oneofci":              {},
	"required_if":          {},
	"required_unless":      {},
	"required_with":        {},
	"required_with_all":    {},
	"required_without":     {},
	"required_without_all": {},
	"skip_unless":          {},
	"spicedb":              {},
	"startsnotwith":        {},
	"startswith":           {},
	"unique":               {},
	"validateFn":           {},
}
