import JsonWebToken, { JwtHeader, SigningKeyCallback } from 'jsonwebtoken'
import jwksClient from 'jwks-rsa'
import dotenv from 'dotenv'

dotenv.config()

const gqlApi = process.env.AUTH_SERVER_URL

const client = jwksClient({
  jwksUri: `${gqlApi}/auth/jwt/jwks.json`,
})

function getKey(header: JwtHeader, callback: SigningKeyCallback) {
  client.getSigningKey(header.kid, function (err, key) {
    var signingKey = key!.getPublicKey()
    callback(err, signingKey)
  })
}

export const verifyToken = (token: string) => {
  JsonWebToken.verify(token, getKey, {}, function (_err, decoded) {
    let decodedToken = decoded

    // use token here:
    console.log(`TEMP LOGGING: decoded token is: ${decodedToken}`)
  })
}
