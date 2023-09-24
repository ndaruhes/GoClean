package locales

// === SUCCESS MESSAGES ===
var SuccessEN = map[string]string{
	//BASIC LISTS
	"SUCCESS-BASIC-0001": "Good Job Bro",

	// AUTH LISTS
	"SUCCESS-AUTH-0001": "Register successfully",
	"SUCCESS-AUTH-0002": "Login successfully",

	//	DATABASE LISTS
	"SUCCESS-DB-0001": "Success migrate database",
	"SUCCESS-DB-0002": "Success seed data to database",

	//	BLOG LISTS
	"SUCCESS-BLOG-0001": "Blog successfully created",
	"SUCCESS-BLOG-0002": "Blog successfully updated",
	"SUCCESS-BLOG-0003": "Blog successfully deleted",
	"SUCCESS-BLOG-0004": "Blog successfully published",
	"SUCCESS-BLOG-0005": "Blog slug successfully updated",
	"SUCCESS-BLOG-0006": "Blog successfully updated to draft",
	"SUCCESS-BLOG-0007": "Get blog successfully",
}

// === ERROR MESSAGES ===
var ErrorEN = map[string]string{
	// ERROR 400 LISTS
	"ERROR-400001": "User has been activated",
	"ERROR-400002": "You have completed your profile",
	"ERROR-400003": "Please fill form correctly",
	"ERROR-400004": "Please upload blog cover",

	// ERROR 401 LISTS
	"ERROR-401001": "Email has already registered",
	"ERROR-401002": "You are not registered",
	"ERROR-401003": "Credentials doesn't match",

	// ERROR 403 LIST
	"ERROR-403001": "Can't edit blog, invalid status",

	// ERROR 404 LIST
	"ERROR-404001": "Data not found",

	//	ERROR 500 LISTS
	"ERROR-50001": "Failed to migrate infrastructures",
	"ERROR-50002": "Migrate key doesn't valid",
	"ERROR-50003": "Upps, something went wrong",
}
