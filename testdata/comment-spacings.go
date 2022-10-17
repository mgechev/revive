package commentspacings

//go:embed
//go:linkname

// This is a comment. 
type hello struct{
	random `json:random` //This is invalid  MATCH /no space between comment symbol and comment text/
}
