  import { customAlphabet } from 'nanoid'
export function createUniqueId(): string {
  const nanoid = customAlphabet('abcdefghijklmnopqrstuvwxyz012345679_-', 10)
  return nanoid(10)
}