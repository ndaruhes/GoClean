package locales

// === BASIC MESSAGES ===
var BasicEN = map[string]string{
	"BASIC-0001": "Upps, something went wrong",
}

// === SUCCESS MESSAGES ===
var SuccessEN = map[string]string{
	// AUTH LISTS
	"SUCCESS-0001": "Register successfully",
	"SUCCESS-0002": "Login successfully",

	//	DATABASE LIST
	"SUCCESS-0003": "Success migrate database",
}

// === ERROR MESSAGES ===
var ErrorEN = map[string]string{
	// ERROR 400 LISTS
	"ERROR-400001": "User has been activated",
	"ERROR-400002": "You have completed your profile",
	"ERROR-400003": "Please fill form correctly",

	// ERROR 401 LISTS
	"ERROR-401001": "Email has already registered",
	"ERROR-401002": "You are not registered",
	"ERROR-401003": "Credentials doesn't match",

	//	ERROR 500 LISTS
	"ERROR-50001": "Failed to migrate database",
	"ERROR-50002": "Migrate key doesn't valid",
}
