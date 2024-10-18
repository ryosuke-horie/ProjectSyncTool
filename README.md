# GitHub Issue連携自動化ツール【サンプル】

## 動作イメージ

![動作イメージ](/assets/operation_image.gif)
![連携されたIssue](/assets/issue.jpg)

## 概要

このプロジェクトは、GitHub Issuesと連携して修正対応項目を自動的に管理するためのツールです。
GitHubアカウントを持たないユーザーとの協力を効率化し、タスクの漏れを防ぎつつ、Issueベースの開発に統合することを目的としています。
フロントエンドにはNext.js、バックエンドにはGo（TinyGo）を採用し、Cloudflareのサービスを利用してホスティングおよびデータベース管理を行っています。

## 課題・背景

- 修正依頼の管理が煩雑で、漏れや重複が発生しやすい
- チャットツールを介した情報共有では、タスクの追跡が困難
- 既存のプロジェクト管理ツールとの連携が不十分で自動化が難しい

## 目的

- 修正対応項目を自動的にGitHubに連携し、タスクの漏れを防ぐ
- GitHubアカウントを持たないユーザーとの協力を支援し、Issue管理を利用できるようにする
- GitHub Issueベースの開発スタイルに合わせた効率的なプロジェクト管理を実現する

## 主な機能

- タスク入力フォーム
  - 対応項目名、詳細、期限を入力可能
- 自動Issue連携
  - 入力内容を基にGitHub Actionsを利用してIssueを自動生成
- 進捗管理
  - UI上でタスクの進捗状況を「対応済み」に更新可能

## 想定するワークフロー

![想定するワークフロー](/assets/workflow.jpg)

## システム構成

- フロントエンド
  - フレームワーク: Next.js
  - ホスティング: Cloudflare Pages
- バックエンド
  - 言語: Go (TinyGo)
  - サーバーレス環境: Cloudflare Workers
データベース
  - サービス: Cloudflare D1
- その他
  - GitHub
  - GitHub Actions

## コスト

本プロジェクトは無料で利用可能です。
インフラはCloudflareの無料枠内で運用しており、GitHub Actionsも無料枠内で十分に対応できます。

## 使用方法（概要）

1. NextjsをCloudflare Pagesにデプロイ
2. バックエンドをCloudflare Workersにデプロイ
3. wranglerを利用してD1を作成しwrangler.tomlに設定
4. GitHubを利用しないユーザーとの共同作業が発生するプロジェクトのリポジトリに.github/workflows/issue-fetch.ymlを配置(定時実行の設定を追加する)
5. フロントエンドのURLを共有し、ユーザーが対応項目を入力
