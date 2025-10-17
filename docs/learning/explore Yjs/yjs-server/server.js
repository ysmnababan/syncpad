import http from 'http'
import WebSocket, { WebSocketServer } from 'ws'
import * as Y from 'yjs'
import * as syncProtocol from 'y-protocols/sync.js'
import * as awarenessProtocol from 'y-protocols/awareness.js'
import * as encoding from 'lib0/encoding.js'
import * as decoding from 'lib0/decoding.js'

// Port for WebSocket server
const PORT = 8081

// Store documents per room
const docs = new Map()

function getYDoc(roomName) {
  if (!docs.has(roomName)) {
    const ydoc = new Y.Doc()
    const awareness = new awarenessProtocol.Awareness(ydoc)

    // Listen for updates on this doc
    ydoc.on('update', (update, origin) => {
      // Log metadata and human-readable text
      console.log(`[room=${roomName}] Received update (origin=${origin || 'client'})`)
      console.log('Update bytes length:', update.length)

      // Apply to a temporary doc to see readable text
      const tmp = new Y.Doc()
      Y.applyUpdate(tmp, update)
      console.log('Current text after update:', tmp.getText('quill').toString())
      console.log('---------------------------')
    })

    docs.set(roomName, { ydoc, awareness })
  }
  return docs.get(roomName)
}

// HTTP server (required for WebSocket upgrade)
const server = http.createServer()

const wss = new WebSocketServer({ noServer: true })

wss.on('connection', (ws, request) => {
  const roomName = (request.url?.split('/')[1] || 'default')
  const { ydoc, awareness } = getYDoc(roomName)

  console.log(`Client connected to room=${roomName}`)

  // Send initial state to client
  const encoder = encoding.createEncoder()
  syncProtocol.writeSyncStep1(encoder, ydoc)
  ws.send(encoding.toUint8Array(encoder))

  ws.on('message', (data) => {
    const u8 = data instanceof Buffer ? new Uint8Array(data) : new Uint8Array(data)
    if (u8.length === 0) return

    const decoder = decoding.createDecoder(u8)
    const messageType = decoder.readUint8()

    if (messageType === 0) { // sync
      try {
        // Apply sync message safely
        syncProtocol.readSyncMessage(decoder, ws, ydoc, encoding.createEncoder())
      } catch (e) {
        console.error('Failed to read sync message:', e)
      }
    } else if (messageType === 1) { // awareness
      try {
        awarenessProtocol.applyAwarenessUpdate(awareness, u8.subarray(1), ws)
      } catch (e) {
        console.error('Failed to apply awareness update:', e)
      }
    } else {
      console.warn('Unknown message type:', messageType)
    }
  })

  ws.on('close', () => console.log(`Client disconnected from room=${roomName}`))
})

// Handle HTTP Upgrade
server.on('upgrade', (req, socket, head) => {
  wss.handleUpgrade(req, socket, head, (ws) => wss.emit('connection', ws, req))
})

server.listen(PORT, () => console.log(`Yjs WebSocket debug server running at ws://localhost:${PORT}`))
