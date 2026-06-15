import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter, Route, Routes } from 'react-router'
import Login from './pages/Login.tsx'
import Register from './pages/Register.tsx'
import Chats from './pages/Chats.tsx'
import Chat from './pages/Chat.tsx'
import ProtectedRoute from './conn/protectedRoute.tsx'
import CreateChat from './pages/CreateChat.tsx'

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route element={<ProtectedRoute />}>
        <Route path="/createChat" element={<CreateChat />} />
        <Route path="/chat" element={<Chat />} />
        <Route path="/" element={<Chats />} />
      </Route>
    </Routes>
  </BrowserRouter>
)
