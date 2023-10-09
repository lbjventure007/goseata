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
	"context"
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/seata/seata-go/pkg/constant"
	"github.com/seata/seata-go/pkg/tm"
	"net/http"
	"time"

	"github.com/seata/seata-go/pkg/client"
)

var (
	serverIpPort  = "http://localhost:8082"
	serverIpPort2 = "http://localhost:8081"
)

func main() {
	flag.Parse()
	client.InitPath("./conf/seatago.yml")

	bgCtx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	// sample update
	sampleUpdate(bgCtx)

	// sample insert on update
	// sampleInsertOnUpdate(bgCtx)
}

func sampleUpdate(ctx context.Context) {
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_Update",
		Timeout: time.Second * 30,
	}, updateData); err != nil {

		fmt.Printf(fmt.Sprintf("tm update data err, %v", err))
	}

}

func updateData(ctx context.Context) (re error) {
	request := gorequest.New()
	fmt.Println("branch transaction begin")

	// global transaction will roll back,because updateDataFail
	request.Post(serverIpPort+"/updateDataSuccess").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				fmt.Println("update updateDataSuccess fail 00")
				re = fmt.Errorf("update updateDataSuccess fail  01")
			}
		})

	request.Post(serverIpPort2+"/updateDataFail").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs1 []error) {
			if response.StatusCode != http.StatusOK {
				fmt.Println("update updateDataFail fail 02")
				re = fmt.Errorf("update data fail 02")

			}
		})

	if re != nil {
		fmt.Println(re)
	}
	return
}

func insertOnUpdateData(ctx context.Context) (re error) {
	request := gorequest.New()
	fmt.Println("branch transaction begin")

	// global transaction will roll back,because insertOnUpdateDataFail
	request.Post(serverIpPort+"/insertOnUpdateDataSuccess").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				re = fmt.Errorf("insert on update data success")
			}
		})

	request.Post(serverIpPort2+"/insertOnUpdateDataFail").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				re = fmt.Errorf("insert on update data fail")
			}
		})
	return
}

func sampleInsertOnUpdate(ctx context.Context) {
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_InsertOnUpdate",
		Timeout: time.Second * 30,
	}, insertOnUpdateData); err != nil {
		panic(fmt.Sprintf("tm insert on update data err, %v", err))
	}
}
