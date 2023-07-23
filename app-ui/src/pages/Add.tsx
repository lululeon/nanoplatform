import React, { useState } from 'react'
import '../App.css'
import { useThingsStore } from '../store/things'
import type { Color } from '../store/things'

function Add() {
  const [name, setName] = useState('')
  const [color, setColor] = useState<Color>('red')

  const handleChange = (evt: React.FormEvent<HTMLInputElement>): void => {
    setName(evt.currentTarget.value)
  }

  const handleSelect = (evt: React.FormEvent<HTMLSelectElement>): void => {
    setColor(evt.currentTarget.value as Color)
  }
  const addFn = useThingsStore(state => state.addThing)

  return (
    <>
      <h1>Add a thing!</h1>
      <div className="card">
        <label htmlFor={name}>Name:</label>
        <input id={name} name="name" type="text" value={name} onChange={handleChange} />
        <label htmlFor="color">Color:</label>
        <select id={color} value={color} onChange={handleSelect}>
          <option value="red" selected={color === 'red'}>
            Red
          </option>
          <option value="yellow" selected={color === 'yellow'}>
            Yellow
          </option>
          <option value="blue" selected={color === 'blue'}>
            Blue
          </option>
        </select>
        <button onClick={() => addFn({ name, color })}>Save</button>
      </div>
      <a href="/">go home</a>
    </>
  )
}

export default Add
