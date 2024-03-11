This file aims to detail the eTh3r protocol to its latest version.

`current version: 0.1-beta (0x0001)`

eTh3r is a project aiming at providing end-to-end encryption using open source hardware, eliminating any risk for a 3rd party reading on the conversation. The protocol used is called by the same name.

Each client is *unique* and identified by its *public key id* that links it to its *public key*.
*Public key* length can't be longer than `0xFFFF` bytes, thus their length will be expressed on 2 bytes.
*Public key id* lenght can't be longer than `0xFF` bytes, thus their length will be expressed on 1 byte.

## Initialisation

A client first initiates a connexion to a server then sends a "hey" package, with its version:
- `c->s: 0x0531b00b 0001`, first 4 bytes are constant (`0x0531b00b`) and last two are for the version (`0x0001`)

Packet to which the server acknowledges: (see error handling at the end)
- `s->c: 0xa0`

The client will then proceeds to send its public key:
- `c->s: 0x0e1f 0800xxxxxxx`, first 2 bytes are constant then (`length`+2) bytes for the key id: `[two bytes key length, following by the key itself]`

The server acknowledges again:
- `s->c: 0xa0`

This covers pretty much everything about initialisation.

## Knocking

In order for a client c1 to initiate a communication with a client c2, it needs to ask its consent. We call that a knock.

c1 knocks c2:
- `c1->s: 0xee 08bbbb`, first byte is constant, followed by (`lenght`+1) bytes for c2's key id: `[length, c2_uid]` 
The server acknowledges the knock and transmits to c2:
- `s->c1: 0xa0 ee`, constant
- `s->c2: 0xae 08aaaa`, first byte is constant, followed by (`lenght`+1) bytes for c1's key id: `[length, c1_uid]`

c2 sends its response and the server forwards it:
- `c2->s: 0xab 01 08aaaa`, first byte is constant, second byte corresponds to the answer (anything else than 1 means no, 1 means yes), then followed by (`lenght`+1) bytes for c1's key id: `[length, c1_uid]`
- `s->c1: 0xab 01 08bbbb`, same but for c2's key id

The server then sends to both clients a room id they will be able to use to communicate:
- `s->c1: 0xac 04dddd 08bbbb`, first byte is constant, second byte is (`room_lenght`+1) bytes representing the room id: `[room_length, rid]`, then followed by (`id_lenght`+1) bytes for c2's key id: `[id_length, c2_uid]`
- `s->c2: 0xac 04dddd 08aaaa`, same but with c1's key

On the server, a connexion is now established between both client.

## Key retrieval

In order for the communication to be secure, each client must know the other's public key. Since we consider that the key id is easily exchangeable, we will once again use this to identify the client:

The client asks for a key:
- `c->s: 0xba 08aaaa`, 1st byte constant, then (`lenght`+1) bytes for the wanted user's key id: `[length, uid]`

The server answers with the key:
- `s->c: 0xa0ba 0800xxxxxxxx`, 2 first bytes are constant then (`lenght`+2) bytes for the key: `[length, pub_key]`

If the server doesn't know the key:
- `s->c: 0xca ba 08aaaa`, 1st byte constant followed by a copy of the request for traceability

## Message

The client simply sends to the server:
- `c1->s: 0xda 04dddd 01 00000100mmmm`, first byte is constant, then (`rid_lenght`+1) bytes for the room id: `[rid_length, rid]`, 1 byte if the message is encrypted or not (01 if so, anything else means plain text), (`pl_lenght`+4) bytes for the payload: `[pl_length, payload]`

And the server forwards it and acknowledges:
- `s->c2: 0xda04dddd0100000100mmm`, a copy of the message
- `s->c1: 0xa0da`

## Room termination

A client can ask for the closing of a room:
- `c1->s: 0xaf 04dddd`, 1st byte constant then (`lenght`+1) bytes for the room id: `[length, rid]`

The server broadcasts the termination and acknowledges:
- `s->c2: 0xaf04dddd`, a copy
- `s->c1: 0xa0af`, constant


# Error codes

| Code | Description |
|----|----|
|`0xa1`| Wrong packet length (handshake) |
|`0xa2`| Wrong packet identifier (handshake) |
|`0xa4`| Unsupported version |
|`0xaa`| Error while receiving the key (minimal length not satisfied) |
|`0xab`| Key packet malformed (wrong constant, not following path) |
|`0xac`| Key payload malformed (wrong length) |
|`0xad`| Key id could not be generated |
|`0xba`| General packer handler, malformation (null length) |
|`0xca`| User not found |
|`0xe0`| Out of path |
|`0xfd`| General packet handler, unkown packed id |
|`0xfe`| Not implemented |
|`0xff`| Error while reading the packet |
