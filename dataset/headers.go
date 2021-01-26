package dataset

import "errors"

func checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID string) error {
	if userAccessToken == "" {
		return errors.New("no user access token header set")
	}
	if collectionID == "" {
		return errors.New("no collection ID header set")
	}
	return nil
}
