package main

import (
	"database/sql"
	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	_ "github.com/go-sql-driver/mysql"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
)


func main() {
	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/tiktok")
	if err != nil {
		panic(err.Error())
	}

	// Create the tables in the database.
	drop, _ := db.Query("DROP TABLE IF EXISTS messages")
	defer drop.Close()
	create, err := db.Query("CREATE TABLE messages (id int PRIMARY KEY AUTO_INCREMENT, chat varchar(255) NOT NULL, text varchar(255) NOT NULL ,sender varchar(255) NOT NULL, send_time BIGINT NOT NULL )")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	defer create.Close()

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
