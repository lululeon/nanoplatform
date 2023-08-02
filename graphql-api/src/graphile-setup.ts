import type { Pool } from 'pg'
import { PostGraphileOptions } from 'postgraphile'
import PgSimplifyInflectorPlugin from '@graphile-contrib/pg-simplify-inflector'

const isProd = process.env.NODE_ENV === 'prod' || process.env.NODE_ENV === 'production'

// Connection string (or pg.Pool) for PostGraphile to use
export const database: string | Pool = process.env.DATABASE_URL || ''

// PostGraphile options; see https://www.graphile.org/postgraphile/usage-library/#api-postgraphilepgconfig-schemaname-options
export const options: PostGraphileOptions = {
  appendPlugins: [PgSimplifyInflectorPlugin],

  // Everything returned by pgSettings is applied to the current session with set_config($key, $value, true);
  // note that set_config only supports string values so it is best to only feed pgSettings string values
  pgSettings(req) {
    // extract headers for use by pg if needed
    return {
      'headers.x-src': req.headers['x-src'],
    }
  },

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

export const port: number = process.env.PORT ? parseInt(process.env.PORT, 10) : 5000
