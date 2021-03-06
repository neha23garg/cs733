package main

import (
	"bufio"
	"bytes"
	"errors"
	//"fmt"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	//"sync"
	"testing"
	"time"
	//raft "github.com/kedareab/cs733/assignment4/raftnode"
	cluster "github.com/cs733-iitb/cluster"
	mock "github.com/cs733-iitb/cluster/mock"
	//raft "github.com/neha23garg/CS733/assignment4"
)

var chandle []*ClientHandler
var rafts []*RaftNode

func TestBasicStep(t *testing.T) {

	cluster_Config := cluster.Config{
		Peers: []cluster.PeerConfig{
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
			{Id: 5},
		}}

	//create cluster
	cl, err := mock.NewCluster(cluster_Config)
	if err != nil {
		fmt.Println("error", err)
	}
	for i := 1; i <= 5; i++ {
		configuration := Config{Id: i, LogDir: "./NodeLog" + strconv.Itoa(i), ElectionTimeout: 800, HeartbeatTimeout: 150, config_MockCluster: cl, Ports: 9090 + i}

		chandle = append(chandle, InitializeClientHandler(configuration))
		chandle[i-1].serverMain()
	}

	//config.MCluster = mcluster

	time.Sleep(1 * time.Second)
}

//leader test cases
func TestRPC_BasicSequential(t *testing.T) {
	cl := mkClient("localhost:9091")
	defer cl.close()

	// Read non-existent file cs733net
	/*m, err := cl.read("cs733net")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)

	// Read non-existent file cs733net
	m, err = cl.delete("cs733net")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)*/

	// Write file cs733net
	data := "Cloud fun"
	m, err := cl.write("cs733net", data, 0)
	expect(t, m, &Msg{Kind: 'O'}, "write success", err)

	// Expect to read it back
	m, err = cl.read("cs733net")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)

	// CAS in new value
	version1 := m.Version
	data2 := "Cloud fun 2"
	// Cas new value
	m, err = cl.cas("cs733net", version1, data2, 0)
	expect(t, m, &Msg{Kind: 'O'}, "cas success", err)

	// Expect to read it back
	m, err = cl.read("cs733net")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(data2)}, "read my cas", err)

	// Expect Cas to fail with old version
	m, err = cl.cas("cs733net", version1, data, 0)
	expect(t, m, &Msg{Kind: 'V'}, "cas version mismatch", err)

	// Expect a failed cas to not have succeeded. Read should return data2.
	m, err = cl.read("cs733net")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(data2)}, "failed cas to not have succeeded", err)

	// delete
	m, err = cl.delete("cs733net")
	expect(t, m, &Msg{Kind: 'O'}, "delete success", err)

	// Expect to not find the file
	m, err = cl.read("cs733net")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)
}

/*var leaderPort int

func TestRPC_Shutdown(t *testing.T) {
	for i := 1; i <= 5; i++ {
		//fmt.Println("rafts chandle", chandle[i-1].rnode)
		rafts = append(rafts, &chandle[i-1].rnode)
	}
	//fmt.Println("rafts", rafts)
	ldr := getLeader(rafts)
	if ldr != nil {
		fmt.Println("leader********", ldr.ID())
		ldr.Shutdown()
		time.Sleep(4000 * time.Millisecond)
		newldr := getNewLeader(rafts, ldr.ID())
		//fmt.Println("new Leader", newldr.ID())
		if newldr != nil {
			fmt.Println("new Leader*", newldr.ID())

			//new leader can be selected from the remaining raft nodes
			expectNotInt(t, ldr.ID(), newldr.ID())

			leaderPort = FetchLeaderPort(newldr.ID())
			//fmt.Printf("leaderport***", leaderPort)
			fmt.Println("leade port info****", strconv.Itoa(leaderPort))
			cl := mkClient("localhost:" + strconv.Itoa(leaderPort))
			defer cl.close()
			// Write file cs733net
			data := "Cloud funz"
			m, err := cl.write("cs733netz", data, 0)
			expect(t, m, &Msg{Kind: 'O'}, "write success", err)

			// Expect to read it back
			m, err = cl.read("cs733netz")
			expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)

			data = "Cloud funzz"
			m, err = cl.write("cs733netzz", data, 0)
			expect(t, m, &Msg{Kind: 'O'}, "write success", err)

			// Expect to read it back
			m, err = cl.read("cs733netzz")
			expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)

			ldr.nodeLog.Close()
			newNode := makeSingleRafts(ldr)
			time.Sleep(800 * time.Millisecond)
			fmt.Println(newNode.ID())
			/*time.Sleep(3000 * time.Millisecond)
			data = "Cloud funzz"
			m, err = cl.write("cs733netzz", data, 0)
			expect(t, m, &Msg{Kind: 'O'}, "write success", err)

			// Expect to read it back
			m, err = cl.read("cs733netzz")
			expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)
			// CAS in new value
			/*	version1 := m.Version
				data2 := "Cloud fun 2"
				// Cas new value
				m, err = cl.cas("cs733netzz", version1, data2, 0)
				expect(t, m, &Msg{Kind: 'O'}, "cas success", err)

				// Expect to read it back
				m, err = cl.read("cs733netzz")
				expect(t, m, &Msg{Kind: 'C', Contents: []byte(data2)}, "read my cas", err)
		}

	}

}

func FetchLeaderPort(id int) int {
	fmt.Println("inside FetchLeaderPortDetails", id)
	for i := 0; i <= 4; i++ {
		if chandle[i].rnode.ID() == id {
			return chandle[i].ClientPort
		}
	}
	return 0
}*/

func TestRPC_Binary(t *testing.T) {
	cl := mkClient("localhost:9093")
	defer cl.close()

	// Write binary contents
	data := "\x00\x01\r\n\x03" // some non-ascii, some crlf chars
	m, err := cl.write("binfile", data, 0)
	expect(t, m, &Msg{Kind: 'O'}, "write success", err)

	// Expect to read it back
	m, err = cl.read("binfile")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)

}

func TestRPC_Chunks(t *testing.T) {
	// Should be able to accept a few bytes at a time
	cl := mkClient("localhost:9091")
	defer cl.close()
	var err error
	snd := func(chunk string) {
		if err == nil {
			err = cl.send(chunk)
		}
	}

	// Send the command "write teststream 10\r\nabcdefghij\r\n" in multiple chunks
	// Nagle's algorithm is disabled on a write, so the server should get these in separate TCP packets.
	snd("wr")
	time.Sleep(10 * time.Millisecond)
	snd("ite test")
	time.Sleep(10 * time.Millisecond)
	snd("stream 1")
	time.Sleep(10 * time.Millisecond)
	snd("0\r\nabcdefghij\r")
	time.Sleep(10 * time.Millisecond)
	snd("\n")
	var m *Msg
	m, err = cl.rcv()
	expect(t, m, &Msg{Kind: 'O'}, "writing in chunks should work", err)
}

func TestRPC_Batch(t *testing.T) {
	// Send multiple commands in one batch, expect multiple responses
	cl := mkClient("localhost:9091")
	defer cl.close()
	m, err := cl.read("cs733net78")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)
	cmds := "write batch1 3\r\nabc\r\n" +
		"write batch2 4\r\ndefg\r\n"
		//"read batch1\r\n"

	cl.send(cmds)
	//time.Sleep(2 * time.Second)
	m, err = cl.rcv()
	expect(t, m, &Msg{Kind: 'O'}, "write batch1 success", err)
	m, err = cl.rcv()
	expect(t, m, &Msg{Kind: 'O'}, "write batch2 success", err)
	//m, err = cl.rcv()
	//expect(t, m, &Msg{Kind: 'C', Contents: []byte("abc")}, "read batch1", err)
}

func TestRPC_BasicTimer(t *testing.T) {
	cl := mkClient("localhost:9091")
	defer cl.close()

	// Write file cs733, with expiry time of 2 seconds
	str := "Cloud fun"
	m, err := cl.write("cs733", str, 2)
	expect(t, m, &Msg{Kind: 'O'}, "write success", err)

	// Expect to read it back immediately.
	m, err = cl.read("cs733")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(str)}, "read my cas", err)

	time.Sleep(3 * time.Second)

	// Expect to not find the file after expiry
	m, err = cl.read("cs733")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)

	// Recreate the file with expiry time of 1 second
	m, err = cl.write("cs733", str, 1)
	expect(t, m, &Msg{Kind: 'O'}, "file recreated", err)

	// Overwrite the file with expiry time of 4. This should be the new time.
	m, err = cl.write("cs733", str, 3)
	expect(t, m, &Msg{Kind: 'O'}, "file overwriten with exptime=4", err)

	// The last expiry time was 3 seconds. We should expect the file to still be around 2 seconds later
	time.Sleep(2 * time.Second)

	// Expect the file to not have expired.
	m, err = cl.read("cs733")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(str)}, "file to not expire until 4 sec", err)

	time.Sleep(3 * time.Second)
	// 5 seconds since the last write. Expect the file to have expired
	m, err = cl.read("cs733")
	expect(t, m, &Msg{Kind: 'F'}, "file not found after 4 sec", err)

	// Create the file with an expiry time of 1 sec. We're going to delete it
	// then immediately create it. The new file better not get deleted.
	m, err = cl.write("cs733", str, 1)
	expect(t, m, &Msg{Kind: 'O'}, "file created for delete", err)

	m, err = cl.delete("cs733")
	expect(t, m, &Msg{Kind: 'O'}, "deleted ok", err)

	m, err = cl.write("cs733", str, 0) // No expiry
	expect(t, m, &Msg{Kind: 'O'}, "file recreated", err)

	time.Sleep(1100 * time.Millisecond) // A little more than 1 sec
	m, err = cl.read("cs733")
	expect(t, m, &Msg{Kind: 'C'}, "file should not be deleted", err)

}

func TestRPC_BasicSequential_Follower(t *testing.T) {
	cl := mkClient("localhost:9092")
	defer cl.close()

	// Read non-existent file cs733net
	m, err := cl.read("cs733net1")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)

	// Read non-existent file cs733net
	m, err = cl.delete("cs733net1")
	expect(t, m, &Msg{Kind: 'F'}, "file not found", err)

	// Write file cs733net
	data := "Cloud fun"
	m, err = cl.write("cs733net1", data, 0)
	expect(t, m, &Msg{Kind: 'O'}, "write success", err)

	// Expect to read it back
	m, err = cl.read("cs733net1")
	expect(t, m, &Msg{Kind: 'C', Contents: []byte(data)}, "read my write", err)

}

// nclients write to the same file. At the end the file should be
// any one clients' last write

func TestRPC_ConcurrentWrites(t *testing.T) {
	nclients := 2
	niters := 10
	clients := make([]*Client, nclients)
	for i := 0; i < nclients; i++ {
		cl := mkClient("localhost:9091")
		if cl == nil {
			t.Fatalf("Unable to create client #%d", i)
		}
		defer cl.close()
		clients[i] = cl
	}

	errCh := make(chan error, nclients)
	var sem sync.WaitGroup // Used as a semaphore to coordinate goroutines to begin concurrently
	sem.Add(1)
	ch := make(chan *Msg, nclients*niters) // channel for all replies
	for i := 0; i < nclients; i++ {
		go func(i int, cl *Client) {
			sem.Wait()
			for j := 0; j < niters; j++ {
				str := fmt.Sprintf("cl %d %d", i, j)
				m, err := cl.write("concWrite", str, 0)
				if err != nil {
					errCh <- err
					break
				} else {
					ch <- m
				}
			}
		}(i, clients[i])
	}
	time.Sleep(100 * time.Millisecond) // give goroutines a chance
	sem.Done()                         // Go!

	// There should be no errors
	for i := 0; i < nclients*niters; i++ {
		select {
		case m := <-ch:
			if m.Kind != 'O' {
				t.Fatalf("Concurrent write failed with kind=%c", m.Kind)
			}
		case err := <-errCh:
			t.Fatal(err)
		}
	}
	m, _ := clients[0].read("concWrite")
	// Ensure the contents are of the form "cl <i> 9"
	// The last write of any client ends with " 9"
	if !(m.Kind == 'C' && strings.HasSuffix(string(m.Contents), " 9")) {
		t.Fatalf("Expected to be able to read after 1000 writes. Got msg = %v", m)
	}
}

//----------------------------------------------------------------------
// Utility functions

type Msg struct {
	// Kind = the first character of the command. For errors, it
	// is the first letter after "ERR_", ('V' for ERR_VERSION, for
	// example), except for "ERR_CMD_ERR", for which the kind is 'M'
	Kind        byte
	Filename    string
	Contents    []byte
	Numbytes    int
	Exptime     int // expiry time in seconds
	Version     int
	RedirectUrl string
}

func (cl *Client) read(filename string) (*Msg, error) {
	cmd := "read " + filename + "\r\n"
	return cl.sendRcv(cmd)
}

func (cl *Client) write(filename string, contents string, exptime int) (*Msg, error) {
	var cmd string
	if exptime == 0 {
		cmd = fmt.Sprintf("write %s %d\r\n", filename, len(contents))
	} else {
		cmd = fmt.Sprintf("write %s %d %d\r\n", filename, len(contents), exptime)
	}
	cmd += contents + "\r\n"
	return cl.sendRcv(cmd)
}

func (cl *Client) cas(filename string, version int, contents string, exptime int) (*Msg, error) {
	var cmd string
	if exptime == 0 {
		cmd = fmt.Sprintf("cas %s %d %d\r\n", filename, version, len(contents))
	} else {
		cmd = fmt.Sprintf("cas %s %d %d %d\r\n", filename, version, len(contents), exptime)
	}
	cmd += contents + "\r\n"
	return cl.sendRcv(cmd)
}

func (cl *Client) delete(filename string) (*Msg, error) {
	cmd := "delete " + filename + "\r\n"
	return cl.sendRcv(cmd)
}

var errNoConn = errors.New("Connection is closed")

type Client struct {
	conn   *net.TCPConn
	reader *bufio.Reader // a bufio Reader wrapper over conn
}

func mkClient(url string) *Client {
	var client *Client
	raddr, err := net.ResolveTCPAddr("tcp", url)
	if err == nil {
		conn, err := net.DialTCP("tcp", nil, raddr)
		if err == nil {
			client = &Client{conn: conn, reader: bufio.NewReader(conn)}
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func (cl *Client) send(str string) error {
	if cl.conn == nil {
		return errNoConn
	}
	_, err := cl.conn.Write([]byte(str))
	if err != nil {
		err = fmt.Errorf("Write error in SendRaw: %v", err)
		cl.conn.Close()
		cl.conn = nil
	}
	return err
}

func (cl *Client) sendRcv(str string) (msg *Msg, err error) {
	if cl.conn == nil {
		return nil, errNoConn
	}
	err = cl.send(str)
	if err == nil {
		msg, err = cl.rcv()
	}

	//fmt.Printf("\n Msg recd : %+v", msg)
	if msg != nil && msg.Kind == 'E' {
		msg, err = NewClientConnection(cl, msg.RedirectUrl, str)
	}
	return msg, err
}

func NewClientConnection(cl *Client, url string, str string) (msg *Msg, err error) {
	cl.close()
	client := mkClient(url)
	cl.conn = client.conn
	cl.reader = client.reader
	return cl.sendRcv(str)

}

func (cl *Client) close() {
	if cl != nil && cl.conn != nil {
		cl.conn.Close()
		cl.conn = nil
	}
}

func (cl *Client) rcv() (msg *Msg, err error) {
	// we will assume no errors in server side formatting
	line, err := cl.reader.ReadString('\n')
	if err == nil {
		msg, err = parseFirst(line)
		if err != nil {
			return nil, err
		}
		if msg.Kind == 'C' {
			contents := make([]byte, msg.Numbytes)
			var c byte
			for i := 0; i < msg.Numbytes; i++ {
				if c, err = cl.reader.ReadByte(); err != nil {
					break
				}
				contents[i] = c
			}
			if err == nil {
				msg.Contents = contents
				cl.reader.ReadByte() // \r
				cl.reader.ReadByte() // \n
			}
		}
	}
	if err != nil {
		cl.close()
	}
	return msg, err
}

func parseFirst(line string) (msg *Msg, err error) {
	fields := strings.Fields(line)
	msg = &Msg{}

	// Utility function fieldNum to int
	toInt := func(fieldNum int) int {
		var i int
		if err == nil {
			if fieldNum >= len(fields) {
				err = errors.New(fmt.Sprintf("Not enough fields. Expected field #%d in %s\n", fieldNum, line))
				return 0
			}
			i, err = strconv.Atoi(fields[fieldNum])
		}
		return i
	}

	if len(fields) == 0 {
		return nil, errors.New("Empty line. The previous command is likely at fault")
	}
	switch fields[0] {
	case "OK": // OK [version]
		msg.Kind = 'O'
		if len(fields) > 1 {
			msg.Version = toInt(1)
		}
	case "CONTENTS": // CONTENTS <version> <numbytes> <exptime> \r\n
		msg.Kind = 'C'
		msg.Version = toInt(1)
		msg.Numbytes = toInt(2)
		msg.Exptime = toInt(3)
	case "ERR_VERSION":
		msg.Kind = 'V'
		msg.Version = toInt(1)
	case "ERR_FILE_NOT_FOUND":
		msg.Kind = 'F'
	case "ERR_CMD_ERR":
		msg.Kind = 'M'
	case "ERR_INTERNAL":
		msg.Kind = 'I'
	case "ERR_REDIRECT":
		msg.Kind = 'E'
		msg.RedirectUrl = fields[1]
	default:
		err = errors.New("Unknown response " + fields[0])
	}
	if err != nil {
		return nil, err
	} else {
		return msg, nil
	}
}

func expect(t *testing.T, response *Msg, expected *Msg, errstr string, err error) {
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	ok := true
	if response.Kind != expected.Kind {
		ok = false
		errstr += fmt.Sprintf(" Got kind='%c', expected '%c'", response.Kind, expected.Kind)
	}
	if expected.Version > 0 && expected.Version != response.Version {
		ok = false
		errstr += " Version mismatch"
	}
	if response.Kind == 'C' {
		if expected.Contents != nil &&
			bytes.Compare(response.Contents, expected.Contents) != 0 {
			ok = false
		}
	}
	if !ok {
		t.Fatal("Expected " + errstr)
	}
}
func expectNotInt(t *testing.T, a int, b int) {
	if a == b {
		t.Error(fmt.Sprintf("Not Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}
