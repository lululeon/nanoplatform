import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

type Color = 'red' | 'yellow' | 'blue'
interface Thing {
  name: string
  color: Color
}
interface ThingsStore {
  things: Thing[]
  addThing: (thing: Thing) => void
  removeThing: (thing: Thing) => void
  removeAll: () => void
}

const useThingsStore = create<ThingsStore>()(
  devtools(
    persist(
      set => ({
        things: [],
        addThing: (newThing: Thing) => set((state: ThingsStore) => ({ things: [...state.things, newThing] })),
        removeThing: (thing: Thing) =>
          set((state: ThingsStore) => ({ things: state.things.filter(_ => _.name != thing.name) })),
        removeAll: () => set(() => ({ things: [] })),
      }),
      {
        name: 'thing-store',
      }
    )
  )
)

export { useThingsStore }

export type { Color, Thing, ThingsStore }
