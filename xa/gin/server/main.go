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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seata/seata-go-samples/util"
	"github.com/seata/seata-go/pkg/client"
	ginmiddleware "github.com/seata/seata-go/pkg/integration/gin"
	"log"
)

var db *sql.DB

func main() {
	client.InitPath("./conf/seatago.yml")
	db = util.GetXAMySqlDb()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// NOTE: when use gin，must set ContextWithFallback true when gin version >= 1.8.1
	//r.ContextWithFallback = true

	r.Use(ginmiddleware.TransactionMiddleware())
	r.POST("/updateDataFail", updateDataFailHandler)
	r.POST("/updateDataSuccess", updateDataSuccessHandler)
	r.POST("/selectForUpdateSuccess", selectForUpdateSuccHandler)

	r.POST("/xa", func(c *gin.Context) {
		log.Println("get tm insertOnUpdateData")
		if err := UpdateDataSuccess(c); err != nil {
			c.JSON(http.StatusBadRequest, "insertOnUpdateData failure")
			return
		}
		if err := UpdateDataFail(c); err != nil {
			c.JSON(http.StatusBadRequest, "insertOnUpdateData failure")
			return
		}
		c.JSON(http.StatusOK, "insertOnUpdateData ok")
	})

	if err := r.Run(":8080"); err != nil {
		log.Println("start tcc server fatal: %v", err)
	}
}

func updateDataSuccessHandler(c *gin.Context) {
	fmt.Println("get tm updateData")
	if err := UpdateDataSuccess(c); err != nil {
		c.JSON(http.StatusBadRequest, "updateData failure")
		fmt.Println(" updateData fail")
		return
	}
	fmt.Println(" updateData ok")
	c.JSON(http.StatusOK, "updateData ok")
}

func updateDataFailHandler(c *gin.Context) {
	fmt.Println("get tm updateData")
	if err := UpdateDataFail(c); err != nil {
		c.JSON(http.StatusBadRequest, "updateData failure")
		fmt.Println(" updateData fail")
		return
	}
	fmt.Println(" updateData ok")
	c.JSON(http.StatusOK, "updateData ok")
}

func selectForUpdateSuccHandler(c *gin.Context) {
	log.Println("execute select for update")
	if err := SelectForUpdateSucc(c); err != nil {
		c.JSON(http.StatusBadRequest, "select for update failed")
		return
	}
	c.JSON(http.StatusOK, "select for update success")
}

func SelectForUpdateSucc(ctx *gin.Context) error {
	sql := "select id, user_id from order_tbl where id=? for update"
	ret, err := db.ExecContext(ctx, sql, 333)
	if err != nil {
		fmt.Printf("select for udpate failed, err:%v\n", err)
		return err
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("select for udpate failed, err:%v\n", err)
		return err
	}
	fmt.Printf("select for udpate success： %d.\n", rows)
	return nil
}

func InsertOnUpdateDataSuccess(ctx *gin.Context) error {
	sql := "insert into order_tbl (id, user_id, commodity_code, count, money, descs) values (?, ?, ?, ?, ?, ?) " +
		"on duplicate key update descs=?"
	ret, err := db.ExecContext(ctx, sql, 1, "NO-100001", "C100000", 100, nil, "init desc", fmt.Sprintf("insert on update descs %d", time.Now().Unix()))
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

func UpdateDataFail(ctx *gin.Context) error {
	sql := "update order_tbl set descs=? where id=? and count=?"
	ret, err := db.ExecContext(ctx, sql, "test02", 2, 10)
	if err != nil {
		fmt.Println("update-- failed, err:%v\n", err)
		return nil
	}

	rows, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("update--- failed, err:%v\n", err)
		return nil
	}
	if rows == 0 {
		return errors.New("修改失败")
	}
	fmt.Println("update0--- success： %d.\n", rows)
	return nil
}
func UpdateDataSuccess(ctx *gin.Context) error {
	sql := "update order_tbl set descs=? where id=?"
	ret, err := db.ExecContext(ctx, sql, "test01", 1)
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
