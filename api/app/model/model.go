package model

type ModificationItem struct {
    ID          uint64 `json:"id"`
    Title       string `json:"title"`
    Status      string `json:"status"`
    Deadline    string `json:"deadline"`
    LinkStatus  string `json:"link_status"`
    IssueNumber string `json:"issue_number,omitempty"`
    Details     string `json:"details,omitempty"`
}

type CreateModificationRequest struct {
    Title    string `json:"title"`
    Status   string `json:"status"`
    Deadline string `json:"deadline"`
    Details  string `json:"details,omitempty"`
}

type CreateModificationResponse struct {
    ModificationItem ModificationItem `json:"modification_item"`
}

type ListModificationsResponse struct {
    Modifications []ModificationItem `json:"modifications"`
}

type UpdateModificationRequest struct {
    Title    string `json:"title"`
    Status   string `json:"status"`
    Deadline string `json:"deadline"`
    Details  string `json:"details,omitempty"`
}
