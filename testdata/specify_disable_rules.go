package fixtures

//revive:disable // MATCH /rule name for lint disabling not found/
var nakedDisable = 1

//revive:disable:exported // no match - rules specified
var withRules = 2

//revive:disable-line // MATCH /rule name for lint disabling not found/
var nakedDisableLine = 3

//revive:disable-line:exported // no match - rules specified
var withRulesLine = 4

//revive:disable-next-line // MATCH /rule name for lint disabling not found/
var nakedDisableNextLine = 5

//revive:disable-next-line:exported // no match - rules specified
var withRulesNextLine = 6

//revive:enable // no match - enable directives are not checked
var enableDirective = 7
