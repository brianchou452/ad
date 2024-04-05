# Dcard Backend Intern Assignment

## Introduction

使用 Golang 設計並且實作一個簡化的廣告投放服務。該服務應該有兩個 API，一個用於產生廣告;一個用於列出廣告。詳細作業要求請參考 [這裡](https://drive.google.com/file/d/1dnDiBDen7FrzOAJdKZMDJg479IC77_zT/view)。
實作成果達到 [15000+ requests per second](#壓力測試結果)。

## Prerequisites

- Golang 1.22.0
- Docker

## 如何進行開發

在專案根目錄建立一個 `.env` 檔案，將 `.env.example` 的內容複製到 `.env` 中，視需求修改裡面的內容。

在專案根目錄執行以下指令：

```bash
./scripts/cicd.sh DEV_UP
```

這個指令會使用Docker 建立 MongoDB、Redis 和 API 的容器。執行指令前請確認 Docker 已經啟動。 Windows 使用者可能需要使用 WSL 來執行這個指令。
如果想要在本機端執行，可以跑完上面指令後使用以下指令：

```bash
docker stop dcard-ad-backend-api-1
go run main.go
```

建議搭配 VSCode [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 套件 和 `scripts/APITesting.http` 來進行 API 測試。

### 如何進行壓力測試
  
```bash
./scripts/cicd.sh PRESSURE_TEST
```

## 專案架構

```text
.
├── .devcontainer                    VSCode DevContainer 設定檔
├── api/
│   ├── get_ad_handler.go            GetAd API 主要程式碼所在
│   ├── post_create_ad_handler.go    PostAd API 主要程式碼所在
│   └── ...
├── assets/                          說明文件的圖片附件
│   └── ...
├── database/                        資料庫操作＆設定
│   ├── ad.go                        Mongo DB 針對廣告的相關操作
│   ├── database.go                  Mongo DB 設定
│   ├── redis.go                     Redis 針對廣告的相關操作
│   └── ...
├── docker/                          docker 相關文件
│   └── ...
├── model/                           整份程式共用的 Struct
│   └── ...
├── scripts/
│   ├── k6/
│   │   ├── main.js                  壓力測試的主要文件
│   │   └── ...
│   ├── APITesting.http              REST Client 測試檔案
│   ├── cicd.sh                      開發常用指令集
│   └── ...
├── main.go                          主程式
├── .env                             環境變數檔案
└── ...
```

## 壓力測試結果

17271.78 requests per second

設備:  
MacBook Air (M1, 2020, 16GB RAM, 512GB SSD)  
macOS Sonoma 14.4.1（23E224）  

![image](assets/pressure-test.png)
