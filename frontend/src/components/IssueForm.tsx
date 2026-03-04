import { useState } from 'react'
import {
  TextField,
  Button,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  FormHelperText,
  Box,
} from '@mui/material'
import type { User } from '../types/issue'

const STATUS_OPTIONS = ['To Do', 'In Progress', 'Review', 'Done']
const MAX_TITLE = 200
const MAX_DESC = 10000

interface IssueFormProps {
  users: User[]
  initial?: { title: string; description: string; status: string; assigneeId: string }
  onSubmit: (data: { title: string; description: string; status: string; assigneeId: string }) => Promise<void>
  submitLabel?: string
}

export default function IssueForm({
  users,
  initial = { title: '', description: '', status: 'To Do', assigneeId: '' },
  onSubmit,
  submitLabel = 'Save',
}: IssueFormProps) {
  const [title, setTitle] = useState(initial.title)
  const [description, setDescription] = useState(initial.description)
  const [status, setStatus] = useState(initial.status)
  const [assigneeId, setAssigneeId] = useState(initial.assigneeId)
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    const t = title.trim()
    if (!t) {
      setError('Title is required')
      return
    }
    if (t.length > MAX_TITLE) {
      setError(`Title must be at most ${MAX_TITLE} characters`)
      return
    }
    if (description.length > MAX_DESC) {
      setError(`Description must be at most ${MAX_DESC} characters`)
      return
    }
    setSubmitting(true)
    try {
      await onSubmit({ title: t, description: description.trim(), status, assigneeId })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, maxWidth: 480 }}>
        <TextField
          label="Title"
          required
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          error={!!error && error.toLowerCase().includes('title')}
          helperText={title.length > MAX_TITLE ? `${title.length}/${MAX_TITLE}` : undefined}
          inputProps={{ maxLength: MAX_TITLE }}
        />
        <TextField
          label="Description"
          multiline
          rows={4}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          inputProps={{ maxLength: MAX_DESC }}
        />
        <FormControl>
          <InputLabel>Status</InputLabel>
          <Select
            value={status}
            label="Status"
            onChange={(e) => setStatus(e.target.value)}
          >
            {STATUS_OPTIONS.map((s) => (
              <MenuItem key={s} value={s}>{s}</MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControl>
          <InputLabel>Assignee</InputLabel>
          <Select
            value={assigneeId}
            label="Assignee"
            onChange={(e) => setAssigneeId(e.target.value)}
          >
            <MenuItem value="">Unassigned</MenuItem>
            {users.map((u) => (
              <MenuItem key={u.id} value={u.id}>{u.name}</MenuItem>
            ))}
          </Select>
        </FormControl>
        {error && <FormHelperText error>{error}</FormHelperText>}
        <Button type="submit" variant="contained" disabled={submitting}>
          {submitLabel}
        </Button>
      </Box>
    </form>
  )
}
