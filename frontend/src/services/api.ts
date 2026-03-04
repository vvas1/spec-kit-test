import type { Issue, User } from '../types/issue'

const BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const url = path.startsWith('http') ? path : `${BASE}${path}`
  const res = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error((data as { error?: string }).error || res.statusText)
  }
  return data as T
}

export async function createIssue(input: {
  title: string
  description?: string
  status?: string
  assigneeId?: string
}) {
  return request<Issue>('/api/issues', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}

export async function getUsers() {
  const res = await request<{ items: User[] }>('/api/users')
  return res.items
}

export { request, BASE }
