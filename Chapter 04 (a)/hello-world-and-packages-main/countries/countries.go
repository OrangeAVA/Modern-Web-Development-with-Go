package countries

func GetCountry(code string) (country string) {
	switch code {
	case "FR":
		country = "France"
	case "IT":
		country = "Italy"
	case "IN":
		country = "India"
	case "US":
		country = "United States"
	default:
		country = "Unknown"
	}
	return
}
