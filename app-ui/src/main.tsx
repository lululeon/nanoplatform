import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { ApolloProvider } from '@apollo/client'
import { gqlClient } from './providers/apollo-client.ts'
import App from './App.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <ApolloProvider client={gqlClient}>
        <App />
      </ApolloProvider>
    </BrowserRouter>
  </React.StrictMode>
)
