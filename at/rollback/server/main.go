/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seata/seata-go-samples/util"
	"github.com/seata/seata-go/pkg/client"
	ginmiddleware "github.com/seata/seata-go/pkg/integration/gin"
	"github.com/seata/seata-go/pkg/util/log"
)

var db *sql.DB

func main() {
	client.InitPath("./conf/seatago.yml")
	db = util.GetAtMySqlDb()

	r := gin.Default()

	// NOTE: when use gin，must set ContextWithFallback true when gin version >= 1.8.1
	// r.ContextWithFallback = true

	r.Use(ginmiddleware.TransactionMiddleware())

	r.POST("/updateDataSuccess", func(c *gin.Context) {
		fmt.Println("get tm updateData")
		if err := updateDataSuccess(c); err != nil {
			fmt.Println("get updateDataSuccess fail")
			c.JSON(http.StatusBadRequest, "updateData failure")
			return
		}

		c.JSON(http.StatusOK, "updateData ok")
	})

	r.POST("/insertOnUpdateDataSuccess", func(c *gin.Context) {
		log.Infof("get tm insertOnUpdateData")
		if err := insertOnUpdateDataSuccess(c); err != nil {
			c.JSON(http.StatusBadRequest, "insertOnUpdateData failure")
			return
		}
		c.JSON(http.StatusInternalServerError, "insertOnUpdateData failure")
	})

	if err := r.Run(":8082"); err != nil {
		log.Fatalf("start tcc server fatal: %v", err)
	}
}

func insertOnUpdateDataSuccess(ctx *gin.Context) error {
	sql := "insert into order_tbl (id, user_id, commodity_code, count, money, descs) values (?, ?, ?, ?, ?, ?) " +
		"on duplicate key update descs=?"
	ret, err := db.ExecContext(ctx, sql, 1, "NO-100001", "C100000", 100, nil, "init desc", fmt.Sprintf("insert on update success %d", time.Now().Unix()))
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return nil
	}

	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return nil
	}
	fmt.Printf("update success： %d.\n", rows)
	return nil
}

func updateDataSuccess(ctx *gin.Context) error {
	sql := "update order_tbl set descs=? where id=?"
	ret, err := db.ExecContext(ctx, sql, "test001", 4)
	if err != nil {
		fmt.Println("update 8080 fail：", err)
		return nil
	}

	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("update 8080 fail：", err)
		return nil
	}

	fmt.Println("update 8080 success：", rows)
	return nil
}
