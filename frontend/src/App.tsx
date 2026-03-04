import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import { Link as MuiLink } from '@mui/material'
import CreateIssuePage from './pages/CreateIssuePage'

const theme = createTheme()

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <nav style={{ padding: 8 }}>
          <MuiLink component={Link} to="/" sx={{ mr: 2 }}>
            Home
          </MuiLink>
          <MuiLink component={Link} to="/create">
            Create issue
          </MuiLink>
        </nav>
        <Routes>
          <Route path="/" element={<div>Issue Tracker</div>} />
          <Route path="/create" element={<CreateIssuePage />} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
