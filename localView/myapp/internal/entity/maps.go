package entity

// Countries MAP
var Countries = map[string]string{
	"RU": "Russian Federation",
	"US": "United States of America",
	"GB": "Great Britain",
	"FR": "France",
	"BL": "Saint Barthélemy",
	"AT": "Austria",
	"BG": "Bulgaria",
	"DK": "Denmark",
	"CA": "Canada",
	"ES": "Spain",
	"CH": "Switzerland",
	"TR": "Turkey",
	"PE": "Peru",
	"NZ": "New Zealand",
	"MC": "Monaco",
}

// SMS/MMS providers MAP
var ProvidersMS = map[string]struct{}{
	"Topolo": {},
	"Rond":   {},
	"Kildy":  {},
}

// Voice Call providers MAP
var ProvidersCall = map[string]struct{}{
	"TransparentCalls": {},
	"E-Voice":          {},
	"JustPhone":        {},
}

// Email providers MAP
var ProvidersEmail = map[string]struct{}{
	"Gmail":      {},
	"Yahoo":      {},
	"Hotmail":    {},
	"MSN":        {},
	"Orange":     {},
	"Comcast":    {},
	"AOL":        {},
	"Live":       {},
	"RediffMail": {},
	"GMX":        {},
	"Protonmail": {},
	"Yandex":     {},
	"Mail.ru":    {},
}
