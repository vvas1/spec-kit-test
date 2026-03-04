export interface Issue {
  id: string
  title: string
  description: string
  status: string
  assigneeId: string
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  name: string
}

export interface CreateIssueInput {
  title: string
  description?: string
  status?: string
  assigneeId?: string
}
