package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	_ "github.com/go-sql-driver/mysql"
	// "math/rand"
)

type Message struct {
	chat      string `json:"chat"`
	text      string `json:"text"`
	sender    string `json:"sender"`
	send_time int64  `json:"send_time"`
}

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/tiktok")
	if err != nil {
		panic("Here is the error "+err.Error())
	}
	defer db.Close()
	insert, err := db.Query("INSERT INTO test (chat, text, sender, send_time) VALUES (?, ?, ?, ?)", req.Message.Chat, req.Message.Text, req.Message.Sender, req.Message.SendTime)
	if err != nil {
		resp := rpc.NewSendResponse()
		resp.Code, resp.Msg = 500, "failure in sending message..."
		return resp, nil
	}
	defer insert.Close()

	resp := rpc.NewSendResponse()
	fmt.Println("receive message: ", req.String())

	resp.Code, resp.Msg = 0, "Message sent successfully"
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/tiktok")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	get, err := db.Query("select chat, sender, text , send_time from test where chat = ? order by send_time limit ? offset ?", req.Chat, req.Limit, req.Cursor)
	//  order by send_time limit ? offset ?  , req.Limit, req.Cursor
	if err != nil {
		resp := rpc.NewPullResponse()
		resp.Code, resp.Msg = 500, "failure in pulling message..."
		return resp, nil
	}
	defer get.Close()
	fmt.Println("Pulled messages from databases")

	start := req.GetCursor()
	retr := int64(req.GetLimit())

	messages := make([]Message, 0)
	for get.Next() {
		var curMsg Message
		err := get.Scan(&curMsg.chat, &curMsg.sender, &curMsg.text, &curMsg.send_time)
		if err != nil {
			resp := rpc.NewPullResponse()
			resp.Code, resp.Msg = 500, "failure in pulling message..."
			return resp, nil
		}
		messages = append(messages, curMsg)
		fmt.Println(curMsg)
	}

	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = 0, "success"
	resp.Messages = make([]*rpc.Message, 0, len(messages))
	for _, msg := range messages {
		resp.Messages = append(resp.Messages, &rpc.Message{
			Chat:     msg.chat,
			Text:     msg.text,
			Sender:   msg.sender,
			SendTime: msg.send_time,
		})
	}
	if *req.Reverse {
		for i, j := 0, len(resp.Messages)-1; i < j; i, j = i+1, j-1 {
			resp.Messages[i], resp.Messages[j] = resp.Messages[j], resp.Messages[i]
		}
	}

	// declare nextCursor as i64
	nextCursor := int64(0)
	if len(messages) == 0 {
		nextCursor = start
	} else {
		nextCursor = start + retr
	}

	queryStatement2 := fmt.Sprintf("SELECT COUNT(id) FROM test WHERE chat = '%v'", req.Chat)

	getCount, err := db.Query(queryStatement2)
	if err != nil {
		resp := rpc.NewPullResponse()
		resp.Code, resp.Msg = 500, "failure in pulling message..."
		return resp, nil
	}
	defer getCount.Close()

	var count int64
	getCount.Next()
	getCount.Scan(&count)

	hasMore := false
	if count > nextCursor {
		hasMore = true
	}
	if !hasMore {
		nextCursor = 0
	}
	resp.Code, resp.Msg = 0, "success"
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor
	return resp, nil
}
