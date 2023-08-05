import express, { Express, NextFunction, Request, Response } from 'express'
import cors from 'cors'
import dotenv from 'dotenv'
import postgraphile from 'postgraphile'
import { authMiddleware } from './auth'
import getGraphileOptions from './graphile'

dotenv.config()

// params
const isProd = process.env.NODE_ENV === 'prod' || process.env.NODE_ENV === 'production'
const dbUrl: string = process.env.DATABASE_URL || ''
const schemas = (process.env.GQL_EXPOSED_SCHEMAS || '').split(',')
const origins = (process.env.ALLOWED_ORIGINS || '').split(',')

const corsOptions = {
  origin: origins,
  methods: ['GET', 'POST'],
}

// server setup
const app: Express = express()
app.use(cors(corsOptions))
app.get('/health', (_req: Request, res: Response, _next: NextFunction) => {
  res.sendStatus(200)
})
app.use('/graphql', authMiddleware)
app.use(postgraphile(dbUrl, schemas, getGraphileOptions(isProd)))

const port = process.env.PORT || ''
app.listen(port, () => {
  console.dir(`GrapqhQL Server up and running on port ${port}`)
})
