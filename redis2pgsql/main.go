package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ctx = context.Background()

type Role struct {
	Description string
	Color       string
	Hoist       bool
	Joinable    bool
	Managed     bool
	Mentionable bool
	Name        string
	Permissions int
	Position    int
	ShortName   string `db:"short_name"`
	Sig         bool
	Sync        bool
	Type        string
}

func main() {
	var (
		err  error
		keys *redis.StringSliceCmd
		rdb  = redis.NewClient(&redis.Options{
			Addr:     "10.42.1.30:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		//namespace = "net.4amlunch.dev"
		namespace = "com.aba-eve"
		prefix    = fmt.Sprintf("%s.srv.perms", namespace)
	)

	db, err := sqlx.Connect("postgres", "host=10.42.1.30 user=chremoas_dev dbname=chremoas_dev_roles sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// Get Permission description
	//keys = rdb.Keys(ctx, fmt.Sprintf("%s:description:*", prefix))
	//tx := db.MustBegin()
	//for _, key := range keys.Val() {
	//	name := strings.Split(key, ":")
	//	description, err := rdb.Get(ctx, key).Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(name[len(name)-1], description)
	//	res := tx.MustExec("INSERT INTO permissions (name, description)"+
	//		"VALUES ($1, $2)", name[len(name)-1], description)
	//	fmt.Printf("%+v\n", res)
	//}
	//// This doesn't exist in redis but needs to exist here
	//tx.MustExec("INSERT INTO permissions (name, description) VALUES ('no_sync', 'IDs to not sync')")
	//err = tx.Commit()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// Get Permission Members
	//keys = rdb.Keys(ctx, fmt.Sprintf("%s:members:*", prefix))
	//tx = db.MustBegin()
	//for _, key := range keys.Val() {
	//	name := strings.Split(key, ":")
	//	val, err := rdb.SMembers(ctx, key).Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//	for _, v := range val {
	//		fmt.Println(key, v)
	//		res := tx.MustExec("INSERT INTO permission_membership (permission, user_id)"+
	//			"VALUES ((SELECT id FROM permissions WHERE name = $1), $2)", name[len(name)-1], v)
	//		fmt.Printf("%+v\n", res)
	//	}
	//}
	//err = tx.Commit()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// Get Filter Descriptions
	//keys = rdb.Keys(ctx, fmt.Sprintf("%s:filter_description:*", prefix))
	//tx = db.MustBegin()
	//for _, key := range keys.Val() {
	//	name := strings.Split(key, ":")
	//	description, err := rdb.Get(ctx, key).Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//	res := tx.MustExec("INSERT INTO filters (name, description)"+
	//		"VALUES ($1, $2)", name[len(name)-1], description)
	//	fmt.Printf("%+v\n", res)
	//}
	//err = tx.Commit()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// Get Filter Members
	keys = rdb.Keys(ctx, fmt.Sprintf("%s:filter_members:*", prefix))
	tx := db.MustBegin()
	for _, key := range keys.Val() {
		name := strings.Split(key, ":")
		val, err := rdb.SMembers(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		for _, v := range val {
			fmt.Println(key, v)
			if v == "" {
				continue
			}
			_, err := tx.Exec("INSERT INTO filter_membership (filter, user_id)"+
				"VALUES ((SELECT id FROM filters WHERE name = $1), $2)", name[len(name)-1], v)
			if err != nil {
				fmt.Printf("Error inserting filter membership: %s", err)
			}
			fmt.Printf("Inserting filter membership: %s %s\n", name[len(name)-1], v)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}

	// Get Roles
	//keys = rdb.Keys(ctx, fmt.Sprintf("%s:role:*", prefix))
	//tx = db.MustBegin()
	//for _, key := range keys.Val() {
	//	val := rdb.HGetAll(ctx, key).Val()
	//	hoist, _ := strconv.ParseBool(val["FilterA"])
	//	joinable, _ := strconv.ParseBool(val["Joinable"])
	//	managed, _ := strconv.ParseBool(val["Managed"])
	//	mentionable, _ := strconv.ParseBool(val["Mentionable"])
	//	permissions, _ := strconv.Atoi(val["Permissions"])
	//	position, _ := strconv.Atoi(val["Position"])
	//	sig, _ := strconv.ParseBool(val["Sig"])
	//	sync, _ := strconv.ParseBool(val["Sync"])
	//
	//	role := Role{
	//		Color:       val["Color"],
	//		Hoist:       hoist,
	//		Joinable:    joinable,
	//		Managed:     managed,
	//		Mentionable: mentionable,
	//		Name:        val["Name"],
	//		Permissions: permissions,
	//		Position:    position,
	//		ShortName:   val["ShortName"],
	//		Sig:         sig,
	//		Sync:        sync,
	//		Type:        val["Type"],
	//	}
	//	fmt.Printf("%s %s\n", val["Name"], key)
	//	rows, err := db.NamedQuery("INSERT INTO roles (color, hoist, joinable, managed, mentionable, name, permissions, position, role_nick, sig, sync, chat_type)"+
	//		"VALUES (:color, :hoist, :joinable, :managed, :mentionable, :name, :permissions, :position, :short_name, :sig, :sync, :type) RETURNING id", role)
	//	if err != nil {
	//		log.Fatalln(val["Name"], err)
	//	}
	//	var id int
	//	if rows.Next() {
	//		rows.Scan(&id)
	//	}
	//	rows.Close()
	//	log.Printf("id: %d", id)
	//
	//	if val["FilterA"] != "wildcard" {
	//		tx.MustExec("INSERT INTO role_filters (role, filter)"+
	//			"VALUES ($1, (SELECT id FROM filters WHERE name = $2))", id, val["FilterA"])
	//	}
	//
	//	if val["FilterB"] != "wildcard" {
	//		tx.MustExec("INSERT INTO role_filters (role, filter)"+
	//			"VALUES ($1, (SELECT id FROM filters WHERE name = $2))", id, val["FilterB"])
	//	}
	//}
	//err = tx.Commit()
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
