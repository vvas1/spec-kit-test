import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Typography } from '@mui/material'
import IssueForm from '../components/IssueForm'
import { createIssue, getUsers } from '../services/api'
import type { User } from '../types/issue'

export default function CreateIssuePage() {
  const navigate = useNavigate()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getUsers()
      .then(setUsers)
      .catch(() => setUsers([]))
      .finally(() => setLoading(false))
  }, [])

  const handleSubmit = async (data: {
    title: string
    description: string
    status: string
    assigneeId: string
  }) => {
    await createIssue({
      title: data.title,
      description: data.description || undefined,
      status: data.status || undefined,
      assigneeId: data.assigneeId || undefined,
    })
    navigate('/')
  }

  if (loading) return <Typography>Loading...</Typography>

  return (
    <>
      <Typography variant="h5" sx={{ mb: 2 }}>Create issue</Typography>
      <IssueForm users={users} onSubmit={handleSubmit} submitLabel="Create" />
    </>
  )
}
