package fixtures

import (
	"fmt"
	"net/http"
)

var url = "http://example.com" // MATCH /prefer secure protocol https over http in "http://example.com"/

const prefix = "http://"
const urlPattern = "http://%s:%d" // MATCH /prefer secure protocol https over http in "http://%s:%d"/

var wsURL = "ws://example.com" // MATCH /prefer secure protocol wss over ws in "ws://example.com"/

const wsPrefix = "ws://"
const wsURLPattern = "ws://%s" // MATCH /prefer secure protocol wss over ws in "ws://%s"/

func unsecureURLScheme() {
	_ = fmt.Sprintf("http://%s", ipPort)                            // MATCH /prefer secure protocol https over http in "http://%s"/
	_ = fmt.Sprintf("http://%s/echo?msg=%s", ipPort, msg)           // MATCH /prefer secure protocol https over http in "http://%s/echo?msg=%s"/
	_ = "http:/::1"                                                 // MATCH /prefer secure protocol https over http in "http:/::1"/
	_ = "http://::/"                                                // MATCH /prefer secure protocol https over http in "http://::/"/
	http.Get("http://json-schema.org/draft-04/schema#/properties/") // MATCH /prefer secure protocol https over http in "http://json-schema.org/draft-04/schema#/properties/"/

	_ = fmt.Sprintf("ws://%s", ipPort)                            // MATCH /prefer secure protocol wss over ws in "ws://%s"/
	_ = fmt.Sprintf("ws://%s/echo?msg=%s", ipPort, msg)           // MATCH /prefer secure protocol wss over ws in "ws://%s/echo?msg=%s"/
	_ = "ws:/::1"                                                 // MATCH /prefer secure protocol wss over ws in "ws:/::1"/
	_ = "ws://::/"                                                // MATCH /prefer secure protocol wss over ws in "ws://::/"/
	http.Get("ws://json-schema.org/draft-04/schema#/properties/") // MATCH /prefer secure protocol wss over ws in "ws://json-schema.org/draft-04/schema#/properties/"/

	// Must not fail
	println("http://localhost:8080", "http://0.0.0.0:8080")
	if "http://127.0.0.1:80" == url {
	}

	println("ws://localhost", "ws://0.0.0.0")
	if "ws://127.0.0.1" == url {
	}

	_ = fmt.Sprintf("wss://%s", ipPort)
	_ = fmt.Sprintf("wss://%s/echo?msg=%s", ipPort, msg)
	_ = "wss:/::1"
	_ = "wss://::/"

	_ = fmt.Sprintf("https://%s", ipPort)
	_ = fmt.Sprintf("https://%s/echo?msg=%s", ipPort, msg)
	_ = "https:/::1"
	_ = "https://::/"
	http.Get("https://json-schema.org/draft-04/schema#/properties/")
	_ = "http://" + "http:/" + "http"
	_ = "ws://" + "ws:/" + "ws"
}
