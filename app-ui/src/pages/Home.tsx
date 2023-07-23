import { useState } from 'react'
import reactLogo from '../assets/react.svg'
import viteLogo from '/vite.svg'
import '../App.css'
import { useThingsStore } from '../store/things'

function Home() {
  const [count, setCount] = useState(0)
  const { things, removeThing, removeAll } = useThingsStore(state => state)

  return (
    <>
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount(count => count + 1)}>count is {count}</button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">Click on the Vite and React logos to learn more</p>
      <h2>There are {things.length} things</h2>
      <button onClick={removeAll}>Click here to delete everything</button>
      {things.map(thing => (
        <div className="thing">
          <p style={{ color: thing.color }}>{thing.name}</p>
          <small>{thing.color}</small>
          <br />
          <button onClick={() => removeThing(thing)}>Remove</button>
        </div>
      ))}
      <a href="/add">go here to add more things (test zustand state mgt)</a>
      <br />
      <a href="/unbuilts">go here to list unbuilts (test apollo client and grapqhql backend)</a>
      <br />
    </>
  )
}

export default Home
