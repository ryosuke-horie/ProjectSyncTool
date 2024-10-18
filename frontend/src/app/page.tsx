"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

type ModificationItem = {
  id: number
  title: string
  status: "未着手" | "進行中" | "対応済み" | "完了"
  deadline: string
  linkStatus: string
  details?: string
  issueNumber?: number
}

export default function ModificationManagement() {
  const [items, setItems] = useState<ModificationItem[]>([])
  const [newItem, setNewItem] = useState<Omit<ModificationItem, "id" | "linkStatus" | "issueNumber">>({
    title: "",
    status: "未着手",
    deadline: "",
    details: "",
  })
  const [editItem, setEditItem] = useState<ModificationItem | null>(null)
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string>("")

  // ダイアログの開閉状態を管理するための状態変数
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)

  // 環境変数からベースURLを取得
  const workerUrl = process.env.NEXT_PUBLIC_WORKER_URL

  // コンポーネントのマウント時に修正項目を取得
  useEffect(() => {
    fetchItems()
  }, [])

  const fetchItems = async () => {
    setLoading(true)
    setError("")
    try {
      const response = await fetch(`${workerUrl}/items`)
      if (!response.ok) {
        throw new Error(`Error fetching items: ${response.statusText}`)
      }
      const data = await response.json()

      // レスポンスがネストされている場合とされていない場合に対応
      const modifications = data.modification_items || data.modifications || data.items || []

      const mappedItems: ModificationItem[] = modifications.map((item: any) => ({
        id: item.id,
        title: item.title,
        status: item.status,
        deadline: item.deadline,
        linkStatus: item.link_status || "未設定",
        details: item.details || "",
        issueNumber: item.issue_number || undefined,
      }))

      setItems(mappedItems)
    } catch (err: any) {
      setError(err.message || "Failed to fetch items")
    } finally {
      setLoading(false)
    }
  }

  const addItem = async () => {
    // 入力バリデーション
    if (!newItem.title || !newItem.deadline) {
      setError("タイトルと期限は必須です。")
      return
    }

    setLoading(true)
    setError("")
    try {
      const response = await fetch(`${workerUrl}/items`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(newItem),
      })

      if (!response.ok) {
        throw new Error(`Error adding item: ${response.statusText}`)
      }

      const data = await response.json()

      // レスポンスがネストされている場合とされていない場合に対応
      const modificationItem = data.modification_item || data

      const newModificationItem: ModificationItem = {
        id: modificationItem.id,
        title: modificationItem.title,
        status: modificationItem.status,
        deadline: modificationItem.deadline,
        linkStatus: modificationItem.link_status || "未設定",
        details: modificationItem.details || "",
        issueNumber: modificationItem.issue_number || undefined,
      }

      setItems(prevItems => [...prevItems, newModificationItem])
      setNewItem({ title: "", status: "未着手", deadline: "", details: "" })

      // ダイアログを閉じる
      setIsAddDialogOpen(false)
    } catch (err: any) {
      setError(err.message || "Failed to add item")
    } finally {
      setLoading(false)
    }
  }

  const updateItem = async () => {
    if (!editItem) return

    setLoading(true)
    setError("")
    try {
      const response = await fetch(`${workerUrl}/items?id=${editItem.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          title: editItem.title,
          status: editItem.status,
          deadline: editItem.deadline,
          details: editItem.details,
        }),
      })

      if (!response.ok) {
        throw new Error(`Error updating item: ${response.statusText}`)
      }

      const updatedData = await response.json()

      const modificationItem = updatedData.modification_item || updatedData

      const updatedItem: ModificationItem = {
        id: modificationItem.id,
        title: modificationItem.title,
        status: modificationItem.status,
        deadline: modificationItem.deadline,
        linkStatus: modificationItem.link_status || "未設定",
        details: modificationItem.details || "",
        issueNumber: modificationItem.issue_number || undefined,
      }

      setItems(prevItems => prevItems.map(item => (item.id === updatedItem.id ? updatedItem : item)))

      // ダイアログを閉じる
      setIsEditDialogOpen(false)
      setEditItem(null)
    } catch (err: any) {
      setError(err.message || "Failed to update item")
    } finally {
      setLoading(false)
    }
  }

  const deleteItem = async (id: number) => {
    setLoading(true)
    setError("")
    try {
      const response = await fetch(`${workerUrl}/items?id=${id}`, {
        method: "DELETE",
      })

      if (!response.ok) {
        throw new Error(`Error deleting item: ${response.statusText}`)
      }

      setItems(prevItems => prevItems.filter(item => item.id !== id))
    } catch (err: any) {
      setError(err.message || "Failed to delete item")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto p-4">
      {/* ...（以下は変更なし）... */}
    </div>
  )
}
