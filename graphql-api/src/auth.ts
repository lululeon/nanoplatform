import JsonWebToken, { JwtHeader, SigningKeyCallback } from 'jsonwebtoken'
import jwksClient from 'jwks-rsa'
import { NextFunction, Request, RequestHandler, Response } from 'express'

import dotenv from 'dotenv'

interface ReqWithPGSettings extends Request {
  ['jwt_userid']: string
  ['jwt_role']: string
}

type AsyncHandlerFn = (req: Request, res: Response, next: NextFunction) => Promise<void>
interface DecodedJWT {
  sub: string // uuid
  // 'st-perm': { t: nnn, v: [] }, // timestamp and value
  // 'st-role': { t: nnn, v: [] }, // timestamp and value
}

const asyncHandlerWrapper = (handler: AsyncHandlerFn) => (req: Request, res: Response, next: NextFunction) => {
  return handler(req, res, next)
    .then(() => {
      return next()
    })
    .catch((e: Error) => next(e))
}

dotenv.config()

const gqlApi = process.env.AUTH_SERVER_URL

const client = jwksClient({
  jwksUri: `${gqlApi}/jwt/jwks.json`,
  requestHeaders: {}, // Optional
  timeout: 5000, // Defaults to 30s otherwise!
  jwksRequestsPerMinute: 5,
})

function getKey(header: JwtHeader, callback: SigningKeyCallback) {
  client.getSigningKey(header.kid, function (err, key) {
    let signingKey = key ? key.getPublicKey() : undefined
    callback(err, signingKey)
  })
}

function verifyToken(token: string, req: ReqWithPGSettings, next: NextFunction) {
  if (!token) {
    console.log('verifyToken - no token received')
    return next()
  }

  JsonWebToken.verify(token, getKey, {}, function (err, decoded) {
    if (err) {
      console.log('verifyToken - will disregard token due to error:', err)
      req['jwt_userid'] = ''
      req['jwt_role'] = 'guest'
    } else {
      console.log('LLDEBUG - remove this log line!!')
      console.dir(decoded)
      req['jwt_userid'] = (decoded as DecodedJWT).sub
      req['jwt_role'] = 'member'
    }

    return next()
  })
}

function rawJwtVerifier(req: Request, _res: Response, next: NextFunction) {
  const token = (req.headers.authorization || '').trim().split(' ')[1]

  // we need to fetch key async, so we're delegating the next() middleware call to be invoked by the function's callback.
  // we're also delegating `req` so we can store some jwt-related metadata which the graphile middleware will pick up via pgSettings(req)
  // else middleware confusion followed by crash.
  verifyToken(token, req as ReqWithPGSettings, next)
}

export const authMiddleware: RequestHandler[] = [rawJwtVerifier]
