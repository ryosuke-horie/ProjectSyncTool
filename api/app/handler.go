package app

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "time"

    "ryosuke-horie/project-management-homemade-api/app/model"
)

type itemHandler struct {
    db *sql.DB
}

var _ http.Handler = (*itemHandler)(nil)

func NewItemHandler(db *sql.DB) http.Handler {
    return &itemHandler{
        db: db,
    }
}

// HTTPリクエストの処理 
func (h *itemHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    switch req.URL.Path {
    case "/items":
        switch req.Method {
        case http.MethodGet:
            h.listItems(w, req)
        case http.MethodPost:
            h.createItem(w, req)
        case http.MethodPut:
            h.updateItem(w, req)
        case http.MethodDelete:
            h.deleteItem(w, req)
        default:
            h.handleErr(w, http.StatusMethodNotAllowed, "Method not allowed")
        }
    case "/items/for-sync":
        if req.Method == http.MethodGet {
            h.listItemsForSync(w, req)
        } else {
            h.handleErr(w, http.StatusMethodNotAllowed, "Method not allowed")
        }
    case "/items/mark-linked":
        if req.Method == http.MethodPost {
            h.markLinkedHandler(w, req)
        } else {
            h.handleErr(w, http.StatusMethodNotAllowed, "Method not allowed")
        }
    default:
        h.handleErr(w, http.StatusNotFound, "Not found")
    }
}

// エラーレスポンスの送信
func (h *itemHandler) handleErr(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{
        "error": msg,
    })
}

// 修正項目の作成
func (h *itemHandler) createItem(w http.ResponseWriter, req *http.Request) {
    var createReq model.CreateModificationRequest
    if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // バリデーション
    if createReq.Title == "" || createReq.Status == "" || createReq.Deadline == "" {
        h.handleErr(w, http.StatusBadRequest, "Missing required fields")
        return
    }

    // ステータスの検証
    validStatuses := map[string]bool{
        "未着手": true,
    }
    if !validStatuses[createReq.Status] {
        h.handleErr(w, http.StatusBadRequest, "Invalid status value")
        return
    }

    // deadline の検証とフォーマット
    deadlineTime, err := time.Parse("2006-01-02", createReq.Deadline)
    if err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid deadline format. Use YYYY-MM-DD.")
        return
    }

    // データベースへの挿入
    stmt, err := h.db.Prepare(`
        INSERT INTO modifications (title, status, deadline, details)
        VALUES (?, ?, ?, ?)
    `)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to prepare statement")
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(createReq.Title, createReq.Status, deadlineTime.Format("2006-01-02"), createReq.Details)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to insert data")
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to retrieve item ID")
        return
    }

    // 作成された項目を返す
    createdItem := model.ModificationItem{
        ID:          uint64(id),
        Title:       createReq.Title,
        Status:      createReq.Status,
        Deadline:    deadlineTime.Format("2006-01-02"),
        LinkStatus:  "未連携", // デフォルト値
        Details:     createReq.Details,
        IssueNumber: "", // 初期状態では空
    }

    res := model.CreateModificationResponse{
        ModificationItem: createdItem,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
    }
}

// listItems は修正項目の一覧を取得します。
func (h *itemHandler) listItems(w http.ResponseWriter, req *http.Request) {
    rows, err := h.db.Query(`
        SELECT id, title, status, deadline, details
        FROM modifications
        ORDER BY id DESC
    `)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to retrieve items")
        return
    }
    defer rows.Close()

    items := []model.ModificationItem{}
    for rows.Next() {
        var item model.ModificationItem
        var deadlineStr string
        var details sql.NullString

        if err := rows.Scan(&item.ID, &item.Title, &item.Status, &deadlineStr, &details); err != nil {
            h.handleErr(w, http.StatusInternalServerError, "Failed to parse items")
            return
        }

        item.Deadline = deadlineStr
        if details.Valid {
            item.Details = details.String
        } else {
            item.Details = ""
        }

        items = append(items, item)
    }

    res := model.ListModificationsResponse{
        Modifications: items,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
    }
}

// 更新処理
func (h *itemHandler) updateItem(w http.ResponseWriter, req *http.Request) {
    // URLからIDを取得（例: /items?id=1）
    idStr := req.URL.Query().Get("id")
    if idStr == "" {
        h.handleErr(w, http.StatusBadRequest, "Missing item ID")
        return
    }
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid item ID")
        return
    }

    var updateReq model.UpdateModificationRequest
    if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // バリデーション
    if updateReq.Title == "" || updateReq.Status == "" || updateReq.Deadline == "" {
        h.handleErr(w, http.StatusBadRequest, "Missing required fields")
        return
    }

    // ステータスの検証
    validStatuses := map[string]bool{
        "未着手": true,
        "進行中": true,
        "対応済み": true,
        "完了":  true,
    }
    if !validStatuses[updateReq.Status] {
        h.handleErr(w, http.StatusBadRequest, "Invalid status value")
        return
    }

    // deadline の検証とフォーマット
    deadlineTime, err := time.Parse("2006-01-02", updateReq.Deadline)
    if err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid deadline format. Use YYYY-MM-DD.")
        return
    }

    // データベースの更新（detailsを追加）
    stmt, err := h.db.Prepare(`
        UPDATE modifications
        SET title = ?, status = ?, deadline = ?, details = ?
        WHERE id = ?
    `)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to prepare statement")
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(updateReq.Title, updateReq.Status, deadlineTime.Format("2006-01-02"), updateReq.Details, id)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to update item")
        return
    }

    // 更新された項目を返す
    updatedItem := model.ModificationItem{
        ID:          id,
        Title:       updateReq.Title,
        Status:      updateReq.Status,
        Deadline:    deadlineTime.Format("2006-01-02"),
        LinkStatus:  "未連携", // 既存のLinkStatusを保持する場合は別途取得が必要
        Details:     updateReq.Details,
        IssueNumber: "", // 必要に応じて保持
    }

    res := model.CreateModificationResponse{
        ModificationItem: updatedItem,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
    }
}

// 修正項目の削除
func (h *itemHandler) deleteItem(w http.ResponseWriter, req *http.Request) {
    // URLからIDを取得（例: /items?id=1）
    idStr := req.URL.Query().Get("id")
    if idStr == "" {
        h.handleErr(w, http.StatusBadRequest, "Missing item ID")
        return
    }

    // idStr を uint64 に変換
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid item ID")
        return
    }

    // データベースからの削除
    stmt, err := h.db.Prepare(`
        DELETE FROM modifications
        WHERE id = ?
    `)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to prepare statement")
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to delete item")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// 連携前のアイテムを取得
func (h *itemHandler) listItemsForSync(w http.ResponseWriter, req *http.Request) {
    // 連携前のアイテムを取得
    rows, err := h.db.Query(`
        SELECT id, title, deadline, details
        FROM modifications
        WHERE link_status = '未連携'
    `)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to retrieve items")
        return
    }
    defer rows.Close()

    items := []model.ModificationItem{}
    for rows.Next() {
        var item model.ModificationItem
        var deadlineStr string
        var details sql.NullString

        if err := rows.Scan(&item.ID, &item.Title, &deadlineStr, &details); err != nil {
            h.handleErr(w, http.StatusInternalServerError, "Failed to parse items")
            return
        }

        item.Deadline = deadlineStr
        if details.Valid {
            item.Details = details.String
        } else {
            item.Details = ""
        }

        items = append(items, item)
    }

    res := model.ListModificationsResponse{
        Modifications: items,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode response: %v\n", err)
    }
}

type markLinkedRequest struct {
    ID          uint64 `json:"id"`
    IssueNumber string `json:"issue_number"`
}

// GitHubへ連携済みのデータに対し、ステータスを「連携済み」に変更する
func (h *itemHandler) markLinkedHandler(w http.ResponseWriter, req *http.Request) {
    var markReq markLinkedRequest
    if err := json.NewDecoder(req.Body).Decode(&markReq); err != nil {
        h.handleErr(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // アイテムの存在確認
    var exists bool
    err := h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM modifications WHERE id = ?)`, markReq.ID).Scan(&exists)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to check item existence")
        return
    }
    if !exists {
        h.handleErr(w, http.StatusNotFound, "Item not found")
        return
    }

    // ステータスの更新
    _, err = h.db.Exec(`
        UPDATE modifications
        SET link_status = '連携済み', issue_number = ?
        WHERE id = ?
    `, markReq.IssueNumber, markReq.ID)
    if err != nil {
        h.handleErr(w, http.StatusInternalServerError, "Failed to update item status")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Item marked as linked.",
    })
}
