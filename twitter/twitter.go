package twitter

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

type Client interface {
	GetFollowers() ([]int64, error)
	GetFollowings() ([]int64, error)
	GetUsersByIds(userIds []int64) ([]interface{}, error)
}

type TwitterCLient struct {
	api *anaconda.TwitterApi
}

func NewClient(apiKey, apiSecret, token, secret string) Client {
	anaconda.SetConsumerKey(apiKey)
	anaconda.SetConsumerSecret(apiSecret)
	api := anaconda.NewTwitterApi(token, secret)

	return &TwitterCLient{
		api: api,
	}
}

func (c *TwitterCLient) GetFollowers() ([]int64, error) {
	v := url.Values{}
	v.Set("count", "5000")

	var cursor int64 = -1

	followerIds := make([]int64, 0, 5000)

	for cursor != 0 {
		if cursor > 0 {
			v.Set("cursor", strconv.FormatInt(cursor, 10))
		}

		followers, err := c.api.GetFollowersIds(v)
		if err != nil {
			return nil, fmt.Errorf("Error getting followers, cursor:%d. %s", cursor, err)
		}

		if len(followers.Ids) > 0 {
			followerIds = append(followerIds, followers.Ids...)
		}

		cursor = followers.Next_cursor
	}

	return followerIds, nil
}

func (c *TwitterCLient) GetFollowings() ([]int64, error) {
	v := url.Values{}
	v.Set("count", "5000")

	var cursor int64 = -1

	followingIds := make([]int64, 0, 5000)

	for cursor != 0 {
		if cursor > 0 {
			v.Set("cursor", strconv.FormatInt(cursor, 10))
		}

		followings, err := c.api.GetFriendsIds(v)
		if err != nil {
			return nil, fmt.Errorf("Error getting followings, cursor:%d. %s", cursor, err)
		}

		if len(followings.Ids) > 0 {
			followingIds = append(followingIds, followings.Ids...)
		}

		cursor = followings.Next_cursor
	}

	return followingIds, nil
}

func (c *TwitterCLient) GetUsersByIds(userIds []int64) ([]interface{}, error) {
	v := url.Values{}
	v.Set("include_entities", "true")

	users, err := c.api.GetUsersLookupByIds(userIds, v)
	if err != nil {
		return nil, fmt.Errorf("error getting users by ids, %s", err)
	}

	outUsers := make([]interface{}, 0)
	for _, user := range users {
		outUsers = append(outUsers, user)
	}

	return outUsers, nil
}
