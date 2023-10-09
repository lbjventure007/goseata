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
	"log"
	"net/http"
	"time"

	"github.com/seata/seata-go/pkg/client"
)

var serverIpPort = "http://localhost:8080"

func main() {
	flag.Parse()
	client.InitPath("./conf/seatago.yml")
	fmt.Println("----1---")
	bgCtx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	fmt.Println("----3---")
	sampleXa(bgCtx)
	//sampleTwoUpdate(bgCtx)
	//fmt.Println("----4---")
	//// sample update
	//sampleUpdate(bgCtx)
	//fmt.Println("----5---")
	// sample insert on update
	//sampleInsertOnUpdate(bgCtx)

	// sample select for update
	//sampleSelectForUpdate(bgCtx)
	fmt.Println("----6---")
}

func updateXA(ctx context.Context) (re error) {
	req := gorequest.New()

	log.Println("branch transaction begin")

	req.Post(serverIpPort+"/updateDataSuccess").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				fmt.Println("update xa 01 failed")
				re = fmt.Errorf("xa update 01 failed")
			}
			fmt.Println("update xa 01 ok")
		})
	req.Post(serverIpPort+"/updateDataFail").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				fmt.Println("update xa failed")
				re = fmt.Errorf("xa update 02 failed")
			}
			fmt.Println("update xa 02 ok")
		})
	return
}
func updateDataTwo(ctx context.Context) (re error) {
	req := gorequest.New()

	log.Println("branch transaction begin")

	req.Post(serverIpPort+"/updateDataFail").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				re = fmt.Errorf("select for update failed")
			}
		})
	return
}
func selectForUpdate(ctx context.Context) (re error) {
	req := gorequest.New()

	log.Println("branch transaction begin")

	req.Post(serverIpPort+"/selectForUpdateSuccess").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				re = fmt.Errorf("select for update failed")
			}
		})
	return
}

func sampleSelectForUpdate(ctx context.Context) {
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_SelectForUpdate",
		Timeout: time.Second * 30,
	}, selectForUpdate); err != nil {
		panic(fmt.Sprintf("tm select for update data err, %v", err))
	}
}

func insertOnUpdateData(ctx context.Context) (re error) {
	request := gorequest.New()
	log.Println("branch transaction begin")
	request.Post(serverIpPort+"/insertOnUpdateDataSuccess").
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

func updateData(ctx context.Context) (re error) {
	log.Println("branch transaction begin000")
	request := gorequest.New()
	log.Println("branch transaction begin111")
	request.Post(serverIpPort+"/updateDataSuccess").
		Set(constant.XidKey, tm.GetXID(ctx)).
		End(func(response gorequest.Response, body string, errs []error) {
			if response.StatusCode != http.StatusOK {
				re = fmt.Errorf("update data fail")
			}
		})
	return
}

func sampleUpdate(ctx context.Context) {
	fmt.Println("-------")
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_Update",
		Timeout: time.Second * 30,
	}, updateData); err != nil {
		panic(fmt.Sprintf("tm update data err, %v", err))
	}
}
func sampleXa(ctx context.Context) {
	fmt.Println("-------")
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_Update",
		Timeout: time.Second * 30,
	}, updateXA); err != nil {
		fmt.Println("tm update data err", err)
	}
}
func sampleTwoUpdate(ctx context.Context) {
	fmt.Println("-------")
	if err := tm.WithGlobalTx(ctx, &tm.GtxConfig{
		Name:    "ATSampleLocalGlobalTx_Update",
		Timeout: time.Second * 30,
	}, updateDataTwo); err != nil {
		panic(fmt.Sprintf("tm update data err, %v", err))
	}
}
