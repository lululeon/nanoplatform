import { Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import Add from './pages/Add'
import Unbuilts from './pages/Unbuilts'
import './App.css'

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/add" element={<Add />} />
        <Route path="/unbuilts" element={<Unbuilts />} />
      </Routes>
    </div>
  )
}

export default App
