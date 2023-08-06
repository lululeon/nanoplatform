import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'

const gqlEndpoint = process.env.GQL_ENDPOINT

const link = new HttpLink({
  uri: gqlEndpoint,
  credentials: 'include',
  fetchOptions: {
    mode: 'no-cors',
  },
})

// ref: https://www.apollographql.com/docs/react/networking/authentication/
const authLink = setContext((_, { headers }) => {
  // >>>
  // TEMP: get the token from frontend st client
  // ...
  const token = process.env.TEMP_ACCESS_TOKEN
  // <<<

  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  }
})

const gqlClient = new ApolloClient({
  link: authLink.concat(link),
  cache: new InMemoryCache(),
})

export { gqlClient }
