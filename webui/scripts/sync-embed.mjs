import { readdir, readFile, mkdir, writeFile, unlink } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const source = fileURLToPath(new URL('../dist/', import.meta.url))
const target = fileURLToPath(new URL('../../internal/server/admin-ui/', import.meta.url))
const check = process.argv.includes('--check')

async function files(root, sub = '') {
  const result = []
  for (const entry of await readdir(path.join(root, sub), { withFileTypes: true })) {
    const relative = path.join(sub, entry.name)
    if (entry.isSymbolicLink()) throw new Error(`symlink is not a generated asset: ${relative}`)
    if (entry.isDirectory()) result.push(...await files(root, relative))
    else if (entry.isFile()) result.push(relative)
  }
  return result.sort()
}

const built = await files(source)
if (!built.includes('index.html')) throw new Error('index.html was not built')
const embedded = await files(target).catch((error) => {
  if (error.code === 'ENOENT') return []
  throw error
})
const differences = []
for (const file of built) {
  const content = await readFile(path.join(source, file))
  const old = await readFile(path.join(target, file)).catch((error) => {
    if (error.code !== 'ENOENT') throw error
    return null
  })
  if (!old || !content.equals(old)) {
    differences.push(file)
    if (!check) {
      await mkdir(path.dirname(path.join(target, file)), { recursive: true })
      await writeFile(path.join(target, file), content)
    }
  }
}
for (const file of embedded.filter((file) => !built.includes(file))) {
  differences.push(file)
  if (!check) await unlink(path.join(target, file))
}
if (check && differences.length) throw new Error(`embedded UI differs from build: ${differences.join(', ')}`)
console.log(`${check ? 'Verified' : 'Synchronized'} ${built.length} embedded assets`)
