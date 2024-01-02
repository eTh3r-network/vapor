This file aims to detail the eTh3r protocol to its latest version.

`current version: 0.1-beta (0x0001)`

eTh3r is a project aiming at providing end-to-end encryption using open source hardware, eliminating any risk for a 3rd party reading on the conversation. The protocol used is called by the same name.

Each client is unique and identified by it's public key id that it transfers to the server during the handshake.

## Initialisation

A client first initiate a connexion to a server then sends a "hey" package, with its version:
    `c->s: 0x0531b00b 0001`, first 4 bytes are constant and last two are for the version

Packet to which the server acknowledges: (see error handling at the end)
    `s->c: 0xa0`

The client will then proceeds to send its public key:
    `c->s: 0x0e1f 0800xxxxxxx`, first 2 bytes are constant then two bytes for the key length, following by the key itself

The server acknowledges again:
    `s->c: 0xa0`


This covers pretty much everything about initialisation.

## Knocking

In order for a client c1 to initiate a communication with a client c2, it needs to ask its consent. We call that a knock.

c1 knocks c2:
    `c1->s: 0xee 08bbbb`, first byte is constant, second byte corresponds to the id of c2's key id length, followed by the said key id

The server acknowledges the knock and transmits to c2:
    `s->c1: 0xa0ee`, constant
    `s->c2: 0xae 08aaaa`, first byte is constant, second byte corresponds to the id of c1's key id length, followed by the said key id

c2 sends its response and the server forwards it:
    `c2->s: 0xab 01 08aaaa`, first byte is constant, second byte corresponds to the answer (anything else than 1 means no, 1 means yes), then followed by the usual c1's key id length, c1's key id
    `s->c1: 0xab 01 08bbbb`, same

The server then sends to both clients a room id they will be able to use to communicate:
    `s->c1: 0xac 04dddd 08bbbb`, first byte is constant, second byte is the room id's length followed by the said room id, then followed by the same length, value but for c2's key
    `s->c2: 0xac 04dddd 08aaaa`, same but with c1's key

On the server, a connexion is now established between both client.

## Key retrieval

In order for the communication to be secure, each client must know the other's public key. Since we consider that the key id is easily exchangeable, we will once again use this to identify the client:

The client asks for a key:
    `c->s: 0xba 08aaaa`, 1st byte constant, then (1 byte length, value) of the wanted key's id

The server answers with the key:
    `s->c: 0xa0ba 0800xxxxxxxx`, 2 first bytes are constant then 2 bytes for the key length followed by the key

If the server doesn't know the key:
    `s->c: 0xca ba08aaaa`, 1st byte constant followed by a copy of the request for traceability

## Message

The client simply sends to the server:
    `c1->s: 0xda 04dddd 01 00000100mmmm`, first byte is constant, then (1 byte length, value) for the room id, 1 byte if the message is encrypted or not (01 is so, anything else means plain text), 4 bytes for the payload length followed by the payload

And the server forwards it and acknowledges:
    `s->c2: 0xda04dddd0100000100mmm`, a copy of the message
    `s->c1: 0xa0da`

## Room termination

A client can ask for the closing of a room:
    `c1->s: 0xaf 04dddd`, 1st byte constant then (1 byte length, value) of the room id

The server broadcasts the termination and acknowledges:
    `s->c2: 0xaf04dddd`, a copy
    `s->c1: 0xa0af`, constant


# Error codes

| Code | Description |
|----|----|
|`0xa1`| Wrong packet length (handshake) |
|`0xa2`| Wrong packet identifier (handshake) |
|`0xa4`| Unsupported version |
|`0xaa`| Error while receiving the key (minimal length not satisfied) |
|`0xab`| Key packet malformed (wrong constant, not following path) |
|`0xac`| Key payload malformed (wrong length) |
|`0xba`| General packer handler, malformation (null length) |
|`0xe0`| Out of path |
|`0xfd`| General packet handler, unkown packed id |
|`0xfe`| Not implemented |
|`0xff`| Error while reading the packet |
