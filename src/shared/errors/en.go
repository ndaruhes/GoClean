package errors

var en = map[string]string{
	// ERR 400 LISTS
	"ERR-400001": "User has been activated",
	"ERR-400002": "You have completed your profile",
	"ERR-400003": "Credit package is invalid",
	"ERR-400004": "Job is invalid",
	"ERR-400005": "Application has been made",
	"ERR-400006": "Application requires talent's CV",
	"ERR-400007": "Invalid application status",
	"ERR-400008": "Admin activation is in invalid state",
	"ERR-400009": "Invalid offer status",
	"ERR-400010": "Company is in invalid state",
	"ERR-400011": "Invalid attachment source status",
	"ERR-400012": "Talent activate in invalid state",
	"ERR-400013": "CV already purchased",
	"ERR-400014": "Insufficient credit",
	"ERR-400015": "Your old password doesn't match",

	// ERR 401 LISTS
	"ERR-40101": "User has been deactivated",

	// ERR 403 LISTS
	"ERR-403001": "User is not allowed to update talent's data",
	"ERR-403002": "User is not allowed to update job's data",
	"ERR-403003": "Job was published, you can't edit the job",
	"ERR-403004": "Job was closed, you can't edit the job",
	"ERR-403005": "Forbidden action",
	"ERR-403006": "Mismatch company ID",
	"ERR-403007": "User is not allowed to update employer's data",

	// ERR 404 LISTS
	"ERR-404001": "Job not found",
	"ERR-404002": "Data not found",
	"ERR-404003": "Invalid URL",
	"ERR-404004": "Experience is not found",
	"ERR-404005": "Certification is not found",
	"ERR-404006": "Project is not found",
	"ERR-404007": "Education is not found",

	// ERR 410
	"ERR-410001": "Token is expired",
	"ERR-410002": "Token is invalid",

	// ERR 429 LISTS
	"ERR-429001": "Too many request in the last hour",

	// ERR 500 LISTS
	"ERR-500001": "An unexpected error occurred. Please contact the support team for assistance. TraceID: {}",
}
