package rule

import (
	"fmt"
	"go/ast"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/mgechev/revive/lint"
)

// StructTagRule lints struct tags.
type StructTagRule struct {
	userDefined map[string][]string // map: key -> []option
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

type tagChecker func(ctx *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeded bool)

// populate tag checkers map
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
	userDefined    map[string][]string // map: key -> []option
	usedTagNbr     map[int]bool        // list of used tag numbers
	usedTagName    map[string]bool     // list of used tag keys
	isAtLeastGo124 bool
}

func (ctx checkContext) isUserDefined(key tagKey, opt string) bool {
	if ctx.userDefined == nil {
		return false
	}

	options := ctx.userDefined[string(key)]
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

	r.userDefined = make(map[string][]string, len(arguments))
	for _, arg := range arguments {
		item, ok := arg.(string)
		if !ok {
			return fmt.Errorf("invalid argument to the %s rule. Expecting a string, got %v (of type %T)", r.Name(), arg, arg)
		}
		parts := strings.Split(item, ",")
		if len(parts) < 2 {
			return fmt.Errorf("invalid argument to the %s rule. Expecting a string of the form key[,option]+, got %s", r.Name(), item)
		}
		key := strings.TrimSpace(parts[0])
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
	userDefined    map[string][]string // map: key -> []option
	isAtLeastGo124 bool
	tagCheckers    map[tagKey]tagChecker
}

func (w lintStructTagRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.StructType:
		isEmptyStruct := n.Fields == nil || n.Fields.NumFields() < 1
		if isEmptyStruct {
			return nil // skip empty structs
		}

		ctx := &checkContext{
			userDefined:    w.userDefined,
			usedTagNbr:     map[int]bool{},
			usedTagName:    map[string]bool{},
			isAtLeastGo124: w.isAtLeastGo124,
		}

		for _, f := range n.Fields.List {
			if f.Tag != nil {
				w.checkTaggedField(ctx, f)
			}
		}
	}

	return w
}

// checkTaggedField checks the tag of the given field.
// precondition: the field has a tag
func (w lintStructTagRule) checkTaggedField(ctx *checkContext, f *ast.Field) {
	if len(f.Names) > 0 && !f.Names[0].IsExported() {
		w.addFailure(f, "tag on not-exported field "+f.Names[0].Name)
	}

	tags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
	if err != nil || tags == nil {
		w.addFailure(f.Tag, "malformed tag")
		return
	}

	for _, tag := range tags.Tags() {
		if msg, ok := w.checkTagNameIfNeed(ctx, tag); !ok {
			w.addFailure(f.Tag, msg)
		}

		checker, ok := w.tagCheckers[tagKey(tag.Key)]
		if !ok {
			continue // we don't have a checker for the tag
		}

		msg, ok := checker(ctx, tag, f.Type)
		if !ok {
			w.addFailure(f.Tag, msg)
		}
	}
}

func (w lintStructTagRule) checkTagNameIfNeed(ctx *checkContext, tag *structtag.Tag) (message string, succeded bool) {
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
	if _, ok := ctx.usedTagName[mapKey]; ok {
		return fmt.Sprintf("duplicate tag name: %q", tagName), false
	}

	ctx.usedTagName[mapKey] = true

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

func checkASN1Tag(ctx *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeded bool) {
	checkList := append(tag.Options, tag.Name)
	for _, opt := range checkList {
		switch opt {
		case "application", "explicit", "generalized", "ia5", "omitempty", "optional", "set", "utf8":
			// do nothing
		default:
			msg, ok := checkCompoundANS1Option(ctx, opt, fieldType)
			if !ok {
				return msg, false
			}
		}
	}

	return "", true
}

func checkCompoundANS1Option(ctx *checkContext, opt string, fieldType ast.Expr) (message string, succeded bool) {
	parts := strings.Split(opt, ":")
	switch parts[0] {
	case "tag":
		tagNumber := strings.TrimLeft(opt, "tag:")
		number, err := strconv.Atoi(tagNumber)
		if err != nil {
			return fmt.Sprintf("ASN1 tag must be a number, got %q", tagNumber), false
		}
		if ctx.usedTagNbr[number] {
			return fmt.Sprintf("duplicated tag number %v", number), false
		}
		ctx.usedTagNbr[number] = true
	case "default":
		if len(parts) < 2 {
			return "malformed default for ASN1 tag", false
		}
		if !typeValueMatch(fieldType, parts[1]) {
			return "field type and default value type mismatch", false
		}
	default:
		if !ctx.isUserDefined(keyASN1, opt) {
			return fmt.Sprintf("unknown option %q in ASN1 tag", opt), false
		}
	}
	return "", true
}

func checkDatastoreTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "flatten", "noindex", "omitempty":
		default:
			if ctx.isUserDefined(keyDatastore, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in Datastore tag", opt), false
		}
	}

	return "", true
}

func checkDefaultTag(_ *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeded bool) {
	if !typeValueMatch(fieldType, tag.Name) {
		return "field type and default value type mismatch", false
	}

	return "", true
}

func checkBSONTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "inline", "minsize", "omitempty":
		default:
			if ctx.isUserDefined(keyBSON, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in BSON tag", opt), false
		}
	}

	return "", true
}

func checkJSONTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty", "string":
		case "":
			// special case for JSON key "-"
			if tag.Name != "-" {
				return "option can not be empty in JSON tag", false
			}
		case "omitzero":
			if ctx.isAtLeastGo124 {
				continue
			}
			fallthrough
		default:
			if ctx.isUserDefined(keyJSON, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in JSON tag", opt), false
		}
	}

	return "", true
}

func checkMapstructureTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty", "reminder", "squash":
		default:
			if ctx.isUserDefined(keyMapstructure, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in Mapstructure tag", opt), false
		}
	}

	return "", true
}

func checkPropertiesTag(ctx *checkContext, tag *structtag.Tag, fieldType ast.Expr) (message string, succeded bool) {
	options := tag.Options
	if len(options) == 0 {
		return "", true
	}

	seenOptions := map[string]bool{}
	for _, opt := range options {
		var msg string
		var ok bool
		parts := strings.Split(opt, "=")
		switch len(parts) {
		case 2:
			msg, ok = checkCompoundPropertiesOption(parts[0], parts[1], fieldType, seenOptions)
		default:
			msg, ok = fmt.Sprintf("unknown or malformed option %q in properties tag", opt), false
		}
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func checkCompoundPropertiesOption(key, value string, fieldType ast.Expr, seenOptions map[string]bool) (message string, succeded bool) {
	if _, ok := seenOptions[key]; ok {
		return fmt.Sprintf("duplicated option %q in properties tag", key), false
	}
	seenOptions[key] = true

	if strings.TrimSpace(value) == "" {
		return fmt.Sprintf("expected option %q to be of the form %s=value in properties tag", key, key), false
	}

	switch key {
	case "default":
		if !typeValueMatch(fieldType, value) {
			return "field type and default value type mismatch", false
		}
	case "layout":
		if gofmt(fieldType) != "time.Time" {
			return "layout option is only applicable to fields of type time.Time", false
		}
	}

	return "", true
}

func checkProtobufTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	// check name
	switch tag.Name {
	case "bytes", "fixed32", "fixed64", "group", "varint", "zigzag32", "zigzag64":
		// do nothing
	default:
		return fmt.Sprintf("invalid protobuf tag name %q", tag.Name), false
	}

	// check options
	seenOptions := map[string]bool{}
	for _, opt := range tag.Options {
		if number, err := strconv.Atoi(opt); err == nil {
			_, alreadySeen := ctx.usedTagNbr[number]
			if alreadySeen {
				return fmt.Sprintf("duplicated tag number %v", number), false
			}
			ctx.usedTagNbr[number] = true
			continue // option is an integer
		}

		switch {
		case opt == "opt" || opt == "proto3" || opt == "rep" || opt == "req":
			// do nothing
		case strings.Contains(opt, "="):
			o := strings.Split(opt, "=")[0]
			_, alreadySeen := seenOptions[o]
			if alreadySeen {
				return fmt.Sprintf("protobuf tag has duplicated option %q", o), false
			}
			seenOptions[o] = true
			continue
		}
	}
	_, hasName := seenOptions["name"]
	if !hasName {
		return `protobuf tag lacks mandatory option "name"`, false
	}

	for k := range seenOptions {
		switch k {
		case "name", "json":
			// do nothing
		default:
			if ctx.isUserDefined(keyProtobuf, k) {
				continue
			}
			return fmt.Sprintf("unknown option %q in protobuf tag", k), false
		}
	}

	return "", true
}

func checkRequiredTag(_ *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	if tag.Name != "true" && tag.Name != "false" {
		return `required should be "true" or "false"`, false
	}

	return "", true
}

func checkTOMLTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "omitempty":
		default:
			if ctx.isUserDefined(keyTOML, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in TOML tag", opt), false
		}
	}

	return "", true
}

func checkURLTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
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
			return fmt.Sprintf("can not set both %q and %q as delimiters in URL tag", opt, delimiter), false
		default:
			if ctx.isUserDefined(keyURL, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in URL tag", opt), false
		}
	}

	return "", true
}

func checkValidateTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	previousOption := ""
	seenKeysOption := false
	options := append([]string{tag.Name}, tag.Options...)
	for _, opt := range options {
		switch opt {
		case "keys":
			if previousOption != "dive" {
				return `option "keys" must follow a "dive" option in validate tag`, false
			}
			seenKeysOption = true
		case "endkeys":
			if !seenKeysOption {
				return `option "endkeys" without a previous "keys" option in validate tag`, false
			}
			seenKeysOption = false
		default:
			parts := strings.Split(opt, "|")
			errMsg, ok := checkValidateOptionsAlternatives(ctx, parts)
			if !ok {
				return errMsg, false
			}
		}
		previousOption = opt
	}

	return "", true
}

func checkXMLTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "any", "attr", "cdata", "chardata", "comment", "innerxml", "omitempty", "typeattr":
		default:
			if ctx.isUserDefined(keyXML, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in XML tag", opt), false
		}
	}

	return "", true
}

func checkYAMLTag(ctx *checkContext, tag *structtag.Tag, _ ast.Expr) (message string, succeded bool) {
	for _, opt := range tag.Options {
		switch opt {
		case "flow", "inline", "omitempty":
		default:
			if ctx.isUserDefined(keyYAML, opt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in YAML tag", opt), false
		}
	}

	return "", true
}

func checkValidateOptionsAlternatives(ctx *checkContext, alternatives []string) (message string, succeded bool) {
	for _, alternative := range alternatives {
		alternative := strings.TrimSpace(alternative)
		parts := strings.Split(alternative, "=")
		switch len(parts) {
		case 1:
			badOpt, ok := areValidateOpts(parts[0])
			if ok || ctx.isUserDefined(keyValidate, badOpt) {
				continue
			}
			return fmt.Sprintf("unknown option %q in validate tag", badOpt), false
		case 2:
			lhs := parts[0]
			_, ok := validateLHS[lhs]
			if ok || ctx.isUserDefined(keyValidate, lhs) {
				continue
			}
			return fmt.Sprintf("unknown option %q in validate tag", lhs), false
		default:
			return fmt.Sprintf("malformed options %q in validate tag, not expected more than one '='", alternative), false
		}
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

func (w lintStructTagRule) addFailure(n ast.Node, msg string) {
	w.onFailure(lint.Failure{
		Node:       n,
		Failure:    msg,
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

var validateSingleOptions = map[string]struct{}{
	"alpha":                     {},
	"alphanum":                  {},
	"alphanumunicode":           {},
	"alphaunicode":              {},
	"ascii":                     {},
	"base32":                    {},
	"base64":                    {},
	"base64url":                 {},
	"bcp47_language_tag":        {},
	"boolean":                   {},
	"bic":                       {},
	"btc_addr":                  {},
	"btc_addr_bech32":           {},
	"cidr":                      {},
	"cidrv4":                    {},
	"cidrv6":                    {},
	"country_code":              {},
	"credit_card":               {},
	"cron":                      {},
	"cve":                       {},
	"datauri":                   {},
	"dir":                       {},
	"dirpath":                   {},
	"dive":                      {},
	"dns_rfc1035_label":         {},
	"e164":                      {},
	"email":                     {},
	"eth_addr":                  {},
	"file":                      {},
	"filepath":                  {},
	"fqdn":                      {},
	"hexadecimal":               {},
	"hexcolor":                  {},
	"hostname":                  {},
	"hostname_port":             {},
	"hostname_rfc1123":          {},
	"hsl":                       {},
	"hsla":                      {},
	"html":                      {},
	"html_encoded":              {},
	"image":                     {},
	"ip":                        {},
	"ip4_addr":                  {},
	"ip6_addr":                  {},
	"ip_addr":                   {},
	"ipv4":                      {},
	"ipv6":                      {},
	"isbn":                      {},
	"isbn10":                    {},
	"isbn13":                    {},
	"isdefault":                 {},
	"iso3166_1_alpha2":          {},
	"iso3166_1_alpha3":          {},
	"iscolor":                   {},
	"json":                      {},
	"jwt":                       {},
	"latitude":                  {},
	"longitude":                 {},
	"lowercase":                 {},
	"luhn_checksum":             {},
	"mac":                       {},
	"mongodb":                   {},
	"mongodb_connection_string": {},
	"multibyte":                 {},
	"nostructlevel":             {},
	"number":                    {},
	"numeric":                   {},
	"omitempty":                 {},
	"printascii":                {},
	"required":                  {},
	"rgb":                       {},
	"rgba":                      {},
	"semver":                    {},
	"ssn":                       {},
	"structonly":                {},
	"tcp_addr":                  {},
	"tcp4_addr":                 {},
	"tcp6_addr":                 {},
	"timezone":                  {},
	"udp4_addr":                 {},
	"udp6_addr":                 {},
	"ulid":                      {},
	"unique":                    {},
	"unix_addr":                 {},
	"uppercase":                 {},
	"uri":                       {},
	"url":                       {},
	"url_encoded":               {},
	"urn_rfc2141":               {},
	"uuid":                      {},
	"uuid3":                     {},
	"uuid4":                     {},
	"uuid5":                     {},
}

var validateLHS = map[string]struct{}{
	"contains":             {},
	"containsany":          {},
	"containsfield":        {},
	"containsrune":         {},
	"datetime":             {},
	"endsnotwith":          {},
	"endswith":             {},
	"eq":                   {},
	"eqfield":              {},
	"eqcsfield":            {},
	"excluded_if":          {},
	"excluded_unless":      {},
	"excludes":             {},
	"excludesall":          {},
	"excludesfield":        {},
	"excludesrune":         {},
	"gt":                   {},
	"gtcsfield":            {},
	"gtecsfield":           {},
	"len":                  {},
	"lt":                   {},
	"lte":                  {},
	"ltcsfield":            {},
	"ltecsfield":           {},
	"max":                  {},
	"min":                  {},
	"ne":                   {},
	"necsfield":            {},
	"oneof":                {},
	"oneofci":              {},
	"required_if":          {},
	"required_unless":      {},
	"required_with":        {},
	"required_with_all":    {},
	"required_without":     {},
	"required_without_all": {},
	"spicedb":              {},
	"startsnotwith":        {},
	"startswith":           {},
	"unique":               {},
}
