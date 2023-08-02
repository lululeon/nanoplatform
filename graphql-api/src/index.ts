import express, { Express } from 'express'
import cors from 'cors'
import { postgraphile } from 'postgraphile'
import { database, options, port } from './graphile-setup'
import dotenv from 'dotenv'
import { verifyToken } from './jwt-verify'

dotenv.config()

const app: Express = express()

const origins = (process.env.ALLOWED_ORIGINS || '').split(',')
const schemas = (process.env.GQL_EXPOSED_SCHEMAS || '').split(',')

const corsOptions = {
  origin: origins,
  methods: ['GET', 'POST'],
}

app.use(cors(corsOptions))
app.use('/graphql', verifyToken)
app.use(postgraphile(database, schemas, options))

const server = app.listen(port, () => {
  console.log(`GrapqhQL Server is running on port: ${port}`)

  const address = server.address()
  console.log(`GrapqhQL Server up and running: ${address}:${port}`)
})
