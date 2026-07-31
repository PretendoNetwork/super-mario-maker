package constants

const (
	// STARS_1ST_MEDAL represents the number of stars needed to acquire the 1st medal
	STARS_1ST_MEDAL uint32 = 1

	// STARS_2ND_MEDAL represents the number of stars needed to acquire the 2nd medal
	STARS_2ND_MEDAL uint32 = 50

	// STARS_3RD_MEDAL represents the number of stars needed to acquire the 3rd medal
	STARS_3RD_MEDAL uint32 = 150

	// STARS_4TH_MEDAL represents the number of stars needed to acquire the 4th medal
	STARS_4TH_MEDAL uint32 = 300

	// STARS_5TH_MEDAL represents the number of stars needed to acquire the 5th medal
	STARS_5TH_MEDAL uint32 = 500

	// STARS_6TH_MEDAL represents the number of stars needed to acquire the 6th medal
	STARS_6TH_MEDAL uint32 = 800

	// STARS_7TH_MEDAL represents the number of stars needed to acquire the 7th medal
	STARS_7TH_MEDAL uint32 = 1300

	// STARS_8TH_MEDAL represents the number of stars needed to acquire the 8th medal
	STARS_8TH_MEDAL uint32 = 2000

	// STARS_9TH_MEDAL represents the number of stars needed to acquire the 9th medal
	STARS_9TH_MEDAL uint32 = 3000

	// STARS_10TH_MEDAL represents the number of stars needed to acquire the 10th medal
	STARS_10TH_MEDAL uint32 = 5000

	// MAX_COURSE_UPLOADS_0TH_1ST_MEDAL represents the number of courses a user can upload when you have 0 or 1 medal
	// Nintendo sets this to 10 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_0TH_1ST_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_2ND_MEDAL represents the number of courses a user can upload after getting the 2nd medal
	// Nintendo sets this to 20 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_2ND_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_3RD_MEDAL represents the number of courses a user can upload after getting the 3rd medal
	// Nintendo sets this to 30 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_3RD_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_4TH_MEDAL represents the number of courses a user can upload after getting the 4th medal
	// Nintendo sets this to 40 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_4TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_5TH_MEDAL represents the number of courses a user can upload after getting the 5th medal
	// Nintendo sets this to 50 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_5TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_6TH_MEDAL represents the number of courses a user can upload after getting the 6th medal
	// Nintendo sets this to 60 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_6TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_7TH_MEDAL represents the number of courses a user can upload after getting the 7th medal
	// Nintendo sets this to 70 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_7TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_8TH_MEDAL represents the number of courses a user can upload after getting the 8th medal
	// Nintendo sets this to 80 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_8TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_9TH_MEDAL represents the number of courses a user can upload after getting the 9th medal
	// Nintendo sets this to 90 by default and users earn more upload slots up to 100.
	// This is a stupid, unfun, mechanic so everyone gets 100 by default. Can be more, but 100 is fine tbh
	MAX_COURSE_UPLOADS_9TH_MEDAL uint32 = 100

	// MAX_COURSE_UPLOADS_10TH_MEDAL represents the number of courses you can upload after getting the 10th medal
	// Nintendo sets this to 100 by default.
	MAX_COURSE_UPLOADS_10TH_MEDAL uint32 = 100

	// COURSE_WORLD_EASY_FAILURE_RATE_MIN is the minimum failure rate for Easy levels on the Course World
	COURSE_WORLD_EASY_FAILURE_RATE_MIN uint32 = 0

	// COURSE_WORLD_EASY_FAILURE_RATE_MAX is the maximum failure rate for Easy levels on the Course World
	COURSE_WORLD_EASY_FAILURE_RATE_MAX uint32 = 34

	// COURSE_WORLD_NORMAL_FAILURE_RATE_MIN is the minimum failure rate for Normal levels on the Course World
	COURSE_WORLD_NORMAL_FAILURE_RATE_MIN uint32 = 35

	// COURSE_WORLD_NORMAL_FAILURE_RATE_MAX is the maximum failure rate for Normal levels on the Course World
	COURSE_WORLD_NORMAL_FAILURE_RATE_MAX uint32 = 74

	// COURSE_WORLD_EXPERT_FAILURE_RATE_MIN is the minimum failure rate for Expert levels on the Course World
	COURSE_WORLD_EXPERT_FAILURE_RATE_MIN uint32 = 75

	// COURSE_WORLD_EXPERT_FAILURE_RATE_MAX is the maximum failure rate for Expert levels on the Course World
	COURSE_WORLD_EXPERT_FAILURE_RATE_MAX uint32 = 95

	// COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MIN is the minimum failure rate for Super Expert levels on the Course World
	COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MIN uint32 = 96

	// COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MAX is the maximum failure rate for Super Expert levels on the Course World
	COURSE_WORLD_SUPER_EXPERT_FAILURE_RATE_MAX uint32 = 100

	// COURSE_DOWNLOAD_WIIU represents the number of courses downloaded by the Wii U at the same time
	COURSE_DOWNLOAD_WIIU uint32 = 100

	// COURSE_DOWNLOAD_3DS represents the number of courses downloaded by the 3DS at the same time
	COURSE_DOWNLOAD_3DS uint32 = 40

	// CHANGED_DATE_YEAR contains the year of an unknown date, possibly when the config was last changed (2020-01-01 12:00)
	CHANGED_DATE_YEAR uint32 = 2020

	// CHANGED_DATE_MONTH contains the month of an unknown date, possibly when the config was last changed (2020-01-01 12:00)
	CHANGED_DATE_MONTH uint32 = 1

	// CHANGED_DATE_DAY contains the day of an unknown date, possibly when the config was last changed (2020-01-01 12:00)
	CHANGED_DATE_DAY uint32 = 1

	// CHANGED_DATE_HOUR contains the hour of an unknown date, possibly when the config was last changed (2020-01-01 12:00)
	CHANGED_DATE_HOUR uint32 = 12

	// CHANGED_DATE_MINUTE contains the minute of an unknown date, possibly when the config was last changed (2020-01-01 12:00)
	CHANGED_DATE_MINUTE uint32 = 0

	// BOOKMARK_DATE_YEAR contains the year of when the SMM bookmark was released (2015-12-22 5:00)
	BOOKMARK_DATE_YEAR uint32 = 2015

	// BOOKMARK_DATE_MONTH contains the month of when the SMM bookmark was released (2015-12-22 5:00)
	BOOKMARK_DATE_MONTH uint32 = 12

	// BOOKMARK_DATE_DAY contains the day of when the SMM bookmark was released (2015-12-22 5:00)
	BOOKMARK_DATE_DAY uint32 = 22

	// BOOKMARK_DATE_HOUR contains the hour of when the SMM bookmark was released (2015-12-22 5:00)
	BOOKMARK_DATE_HOUR uint32 = 5

	// BOOKMARK_DATE_MINUTE contains the minute of when the SMM bookmark was released (2015-12-22 5:00)
	BOOKMARK_DATE_MINUTE uint32 = 0
)
