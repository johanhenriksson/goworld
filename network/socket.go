package network

import (
    "fmt"
    "net"
    "time"
    "encoding/binary"
)

type OutMessage struct {
    Data    []byte
    Sent    time.Time
}

type Socket struct {
    Timeout     time.Duration
    Address     *net.UDPAddr /* remote address */
    outbox      map[uint16]OutMessage /* outgoing packets not yet ack'd */
    buffer      []byte /* receieve buffer */
    acked       uint32 /* previous ack bitfield */
    ack         uint16 /* latest ack sent */
    seqnum      uint16 /* next sequence number */
    conn        *net.UDPConn
    server      bool
}

func NewSocket(connection *net.UDPConn, addr *net.UDPAddr, buffSize int, server bool) *Socket {
    return &Socket {
        Timeout: 3 * time.Second,
        Address: addr,

        outbox: make(map[uint16]OutMessage),
        buffer: make([]byte, buffSize),
        seqnum: 1,
        conn: connection,
        server: server,
    }
}

func (s *Socket) Send(buffer []byte) {
    sq := s.seqnum
    s.seqnum++

    /* Store outgoing message */
    s.outbox[sq] = OutMessage {
        Sent: time.Now(),
        Data: buffer,
    }

    /* Write header */
    header := make([]byte, HeaderLength)
    binary.BigEndian.PutUint16(header[0:2], sq)
    binary.BigEndian.PutUint16(header[2:4], s.ack)
    binary.BigEndian.PutUint32(header[4:8], s.acked)

    /* Prepend packet header & send */
    packet := append(header, buffer...)
    fmt.Printf("Sent %d -> %q to %s\n", sq, packet, s.Address)
    s.Write(packet)
}

func (s *Socket) Write(buffer []byte) {
    if s.server {
        s.conn.WriteToUDP(buffer, s.Address)
    } else {
        s.conn.Write(buffer)
    }
}

func (s *Socket) Read(buffer []byte) (int, *net.UDPAddr, error) {
    if s.server {
        return s.conn.ReadFromUDP(s.buffer)
    } else {
        n, err := s.conn.Read(s.buffer)
        return n, s.Address, err
    }
}

func (s *Socket) Recv(buffer []byte) {
    header := buffer[0:HeaderLength]
    packet := buffer[HeaderLength:]

    seq   := binary.BigEndian.Uint16(header[0:2])

    fmt.Println("Recieved", seq, "<-", string(packet), "from", s.Address)

    /* Update ack local fields with the incoming sequence number */
    if seq < s.ack {
        i := s.ack - seq
        if i >= 32 {
            /* Packet is too old - drop it */
            fmt.Println("Dropping old packet", seq)
            return
        }
        mask := uint32(1 << i)
        s.acked = s.acked | mask
        fmt.Println("Ack old packet", seq)
    }
    if seq == s.ack {
        /* Duplicate packet? Drop it */
        fmt.Println("Dropping duplicate packet", seq)
        return
    }
    if seq > s.ack {
        i := seq - s.ack
        s.acked = (s.acked << i) | 1
        s.ack = seq
        fmt.Println("Ack new packet", seq)
    }

    /* Mark outgoing messages as delivered */
    ack   := binary.BigEndian.Uint16(header[2:4])
    acked := binary.BigEndian.Uint32(header[4:8])
    if _, ok := s.outbox[ack]; ok {
        delete(s.outbox, ack)
    }
    for i := 0; i < 32; i++ {
        x := uint32(1 << uint32(i))
        if acked & x > 0 {
            p_ack := ack - uint16(i)
            delete(s.outbox, p_ack)
        }
    }
}

func (s *Socket) Close() {
}

func (s *Socket) Update(dt float64) {
    now := time.Now()
    lost := make([]uint16, 0, 4)
    for sq, msg := range s.outbox {
        age := now.Sub(msg.Sent)
        if age > s.Timeout {
            /* Packet lost */
            lost = append(lost, sq)
            fmt.Println("Lost packet", sq, string(msg.Data))
        }
    }

    for _, sq := range lost {
        delete(s.outbox, sq)
    }
}
