package validate

var states = map[string]bool{
	"AL": true,
	"AK": true,
	"AS": true,
	"AZ": true,
	"AR": true,
	"CA": true,
	"CO": true,
	"CT": true,
	"DE": true,
	"DC": true,
	"FM": true,
	"FL": true,
	"GA": true,
	"GU": true,
	"HI": true,
	"ID": true,
	"IL": true,
	"IN": true,
	"IA": true,
	"KS": true,
	"KY": true,
	"LA": true,
	"ME": true,
	"MH": true,
	"MD": true,
	"MA": true,
	"MI": true,
	"MN": true,
	"MS": true,
	"MO": true,
	"MT": true,
	"NE": true,
	"NV": true,
	"NH": true,
	"NJ": true,
	"NM": true,
	"NY": true,
	"NC": true,
	"ND": true,
	"MP": true,
	"OH": true,
	"OK": true,
	"OR": true,
	"PW": true,
	"PA": true,
	"PR": true,
	"RI": true,
	"SC": true,
	"SD": true,
	"TN": true,
	"TX": true,
	"UT": true,
	"VT": true,
	"VI": true,
	"VA": true,
	"WA": true,
	"WV": true,
	"WI": true,
	"WY": true,
}

// State validation
func State(stateAbbr string) (errs []error) {
	if len(stateAbbr) != 2 {
		errs = append(errs, validationError("US state must be of length 2"))
	} else if states[stateAbbr] != true {
		errs = append(errs, validationError("Invalid US state"))
	}

	return errs
}

func ZipCode(zipCode string) (errs []error) {
	if len(zipCode) != 5 {
		errs = append(errs, validationError("Invalid zip code"))
	}

	return errs
}
