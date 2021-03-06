package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var lastCached time.Time

func init() {
	getData()
}

// cache the database every 3 seconds.
func beginCache() {
	// tick := time.NewTicker(3 * time.Second)
	// go func() {
	// 	for range tick.C {
	// 					getData()
	// 	}
	// }()
	if time.Now().Sub(lastCached).Milliseconds() > 3000 {
		fmt.Println("caching")
		lastCached = time.Now()
		time.AfterFunc(3500*time.Millisecond, func() { getData() })
	}

}

// getChildren loads the child replies recursively
func getChildren(ID string) (childs []*postData) {
	children, err := rdb.ZRevRange(ctx, ID+":CHILDREN", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}

	for _, child := range children {
		data, err := rdb.HGetAll(ctx, "OBJECT:"+child).Result()
		if err != nil {
			fmt.Println(err)
		}

		childs = append(childs, makePost(data))
	}
	return
}

// getData gets the board data from redis
func getData() {
	posts = make(map[string][]*postData)
	tagmem, err := rdb.ZRevRange(ctx, "TAGS", 0, -1).Result()
	handleErr(err)

	tags = []string{}
	for _, tag := range tagmem {
		if !isDefaultTag(tag) {
			tags = append(tags, tag)
		}
		dbPosts, err := rdb.ZRevRange(ctx, tag, 0, -1).Result()
		handleErr(err)

		for _, post := range dbPosts {
			data, err := rdb.HGetAll(ctx, "OBJECT:"+post).Result()
			handleErr(err)

			posts[tag] = append(posts[tag], makePost(data))
		}
	}

	frontpage = make(map[string][]*postData)
	dbPosts, err := rdb.ZRevRange(ctx, "ALLPOSTS", 0, -1).Result()
	handleErr(err)

	for _, dbPost := range dbPosts {
		data, err := rdb.HGetAll(ctx, "OBJECT:"+dbPost).Result()
		handleErr(err)

		frontpage["all"] = append(frontpage["all"], makePost(data))
	}

}

func makeZmem(st string) *redis.Z {
	return &redis.Z{
		Member: st,
		Score:  0,
	}
}

func isDefaultTag(tag string) bool {
	for _, dtag := range defaultTags {
		if dtag == tag {
			return true

		}
	}
	return false
}
