import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client'

const gqlEndpoint = process.env.GQL_ENDPOINT
const link = new HttpLink({
  uri: gqlEndpoint,
  fetchOptions: {
    mode: 'no-cors',
  },
})

const gqlClient = new ApolloClient({
  link: link,
  cache: new InMemoryCache(),
})

export { gqlClient }
