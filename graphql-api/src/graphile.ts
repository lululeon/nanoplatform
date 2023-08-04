import { PostGraphileOptions } from 'postgraphile'
import PgSimplifyInflectorPlugin from '@graphile-contrib/pg-simplify-inflector'

// PostGraphile options; see https://www.graphile.org/postgraphile/usage-library/#api-postgraphilepgconfig-schemaname-options
export default function getGraphileOptions(isProd: boolean): PostGraphileOptions {
  return {
    appendPlugins: [PgSimplifyInflectorPlugin],

    // Everything returned by pgSettings is applied to the current session with set_config($key, $value, true);
    // note that set_config only supports string values so it is best to only feed pgSettings string values
    // graphile docs: keys must contain either one or two period (.) characters, and the prefix (the bit before the first period) must not be used by any Postgres extension.
    // Variables without periods will be interpreted as internal Postgres settings, such as 'role', and will be applied by Postgres.
    pgSettings: async (req: any) => ({
      'jwt.claims.user_id': req['jwt_userid'],
      role: req.pgrole,
    }),

    // we control schema changes and should restart when appropriate
    watchPg: false,

    graphiql: isProd ? false : true,
    enhanceGraphiql: isProd ? false : true,
    subscriptions: true,
    dynamicJson: true,
    setofFunctionsContainNulls: false,
    ignoreRBAC: false,
    showErrorStack: 'json',
    extendedErrors: ['hint', 'detail', 'errcode'],
    allowExplain: true,
    legacyRelations: 'omit',
    exportGqlSchemaPath: `${__dirname}/schema.graphql`,
    disableQueryLog: isProd ? true : false,
    sortExport: true,
  }
}
