package service

const (
	TEST_USER_ONE   = "1"
	TEST_USER_TWO   = "2"
	TEST_USER_THREE = "3"

	INITIAL_SKIP = 0

	TEST_TAG_ONE   = "tag1"
	TEST_TAG_TWO   = "tag2"
	TEST_TAG_THREE = "tag3"

	TEST_NOT_FOLLOWING_ID = "4"


	TEST_USER_ONE_USERNAME = "username1"
	TEST_USER_TWO_USERNAME = "username2"
	TEST_USER_THREE_USERNAME = "username3"

	TEST_NOT_FOLLOWING_USERNAME = "username_not_following"
)

func getTestUsername(userID string) string {
	switch userID {
	case TEST_USER_ONE:
		return TEST_USER_ONE_USERNAME
	case TEST_USER_TWO:
		return TEST_USER_TWO_USERNAME
	case TEST_USER_THREE:
		return TEST_USER_THREE_USERNAME
	case TEST_NOT_FOLLOWING_ID:
		return TEST_NOT_FOLLOWING_USERNAME
	default:
		return ""
	}
}