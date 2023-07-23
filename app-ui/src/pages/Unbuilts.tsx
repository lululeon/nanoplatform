import { useQuery } from '@apollo/client'
import { GET_UNBUILTS } from '../gql/queries'

export default function () {
  const { loading, error, data } = useQuery(GET_UNBUILTS)

  if (loading) return <p>Loading...</p>
  if (error) return <p>Error : {error.message}</p>

  return (
    <>
      <h1>Unbuilts</h1>
      <div className="card">
        {data.unbuilts.nodes.map((ub: any) => (
          <div>
            <p style={{ border: '1px solid grey', margin: '1em' }}>{ub.title}</p>
            <p>{ub.elevatorPitch}</p>
            <small>{ub.currentStatus}</small>
          </div>
        ))}
      </div>

      <a href="/">go home</a>
    </>
  )
}
