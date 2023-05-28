# FoodApp_Golang_LINEBOT_Lambda_RDS
1食分のカロリーを返答するLINEBot


## 目次

- [イメージ](#イメージ)
- [バックグラウンド(開発背景)](#バックグラウンド(開発背景))
- [サービス概要(使用法)](#サービス概要(使用法))
- [使用技術](#使用技術)


▽▽リンクアクセスには、許可を得てください。▽▽
- [システム構成図](#システム構成図)
- [データベース構成](#データベース構成)
- [LINE_Bot](#LINE_Bot)

△△リンクアクセスには、許可を得てください。△△

## イメージ
 <img width="220" alt="image" src="https://github.com/Yuta-Haruna/FoodApp_Golang_LINEBOT_Lambda_RDS/assets/50592688/f06331cb-a14b-4b8f-80cf-ef80b1de85d9">

## バックグラウンド(開発背景)
　ダイエットなどで摂取カロリーの計算を行っていますが、1食分をg単位でカロリー変換しなければならず使い勝手が悪く計算しづらかった。
 特に、重量(g)がわからないとカロリー計算が概算になってしまう。
 あまり市場に1食分(概算)単位のカロリー表示アプリがなかったため、本アプリの開発に至った。

## サービス概要(使用法)
1. 本アプリのQRコードを読み取る。
2. トークを開始する。
3. 調べたい料理名を入力する。

## 使用技術

### 環境
- AWS
- AWS(Lambda)
- AWS VPC
- AWS InterNetGateway
- RDS(MySQL)

### 言語
- Golang(go ver1.19)

### ライブラリ
- LINE SDK

### ツール
- GitHub(Git)
- VSCODE
- Windows
- TablePlus
- diagrams.net (設計用)
- Googleスプレッドシート

## システム構成図
https://app.diagrams.net/#G1gEIdjbOthDrG6d6J4AZlAKsZMo79B32u

![image](https://user-images.githubusercontent.com/50592688/232229040-f586921b-a1fa-4700-8a59-b2aa04a5037f.png)


## データベース構成
https://app.diagrams.net/#G1UuAvHwzMAnGrqwA0_Rl3_dnjPJuxgWML

・DB
![image](https://user-images.githubusercontent.com/50592688/232228472-bbe5ebc9-0470-4a3c-9364-b03066f33189.png)

・テーブル詳細
![image](https://user-images.githubusercontent.com/50592688/232228498-e328f25f-38d0-4f47-8725-d040725d126f.png)

## LINE_Bot
※必ず開発者へ連絡してください。
　(連絡ない場合は、機能を停止しています。)
 
・QRコード

![image](https://user-images.githubusercontent.com/50592688/232229454-2e67344c-db7b-4ee0-8770-e4e5cedc72fb.png)


