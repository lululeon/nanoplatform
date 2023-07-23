import { gql } from '@apollo/client'

export const GET_UNBUILTS = gql`
  query {
    allUnbuilts {
      nodes {
        title
        currentStatus
        elevatorPitch
      }
    }
  }
`
