# integrating JWT verification

This is here just for ref only for now.
Below would be the "cooked" way, requiring supertokens.init in this server:

```ts
import { verifySession } from 'supertokens-node/recipe/session/framework/express'
import { SessionRequest } from 'supertokens-node/framework/express'
import UserRoles from 'supertokens-node/recipe/userroles'
import { SessionContainerInterface } from 'supertokens-node/lib/build/recipe/session/types'

interface ReqWithPGSettings extends SessionRequest {
  jwt: {
    claims: Record<string, string>
  }
}

function sessionVerifier() {
  return asyncHandlerWrapper(
    verifySession({
      sessionRequired: false, //allow un-authed access
      overrideGlobalClaimValidators: async globalValidators => [
        ...globalValidators,
        UserRoles.UserRoleClaim.validators.includes('member'),
        // UserRoles.PermissionClaim.validators.includes("read:own")
        // UserRoles.PermissionClaim.validators.includes("write:own")
      ],
    })
  )
}

function claimsHandler(req: ReqWithPGSettings, _res: Response, next: NextFunction) {
  // because we allow `sessionRequired: false`, we must check for a valid session first
  const { session } = req as ReqWithPGSettings
  if (session) {
    const userId = session.getUserId()
    // set authed claims on req...
  } else {
    // set guest claims on req...
  }

  console.log('LLDEBUG remove this log line below!!')
  console.dir(req.jwt)
  return next()
}
export const authMiddleware: RequestHandler[] = [sessionVerifier(), claimsHandler as RequestHandler]
```
