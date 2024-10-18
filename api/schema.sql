CREATE TABLE IF NOT EXISTS modifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    status TEXT CHECK(status IN ('未着手', '進行中', '対応済み', '完了')) NOT NULL,
    deadline TEXT NOT NULL,
    link_status TEXT CHECK(link_status IN ('未連携', '連携済み')) NOT NULL DEFAULT '未連携',
    issue_number TEXT,
    details TEXT
);
